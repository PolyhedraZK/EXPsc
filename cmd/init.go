package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/naoina/toml"
	coreCfg "github.com/nbnet/side-chain/core/types"
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

var InitFilesCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize side chain",
	RunE:  initFiles,
}

// init sets up flags for the 'init' command to configure root directory, number of validators, and host list for node initialization.
func init() {
	InitFilesCmd.Flags().StringVarP(&RootDir, "root-dir", "r", DefaultRootDir, "Root directory, '"+DefaultRootDir+
		"' will be generated in the directory you specified, like $HOME/.side-chain")

	InitFilesCmd.Flags().IntVarP(&Validators, "validators", "v", DefaultValidators,
		"Number of Validators")

	InitFilesCmd.Flags().StringVarP(&HostList, "host-list", "l", DefaultHostList,
		"Host list, specify hosts for different nodes, separated by semicolons. like 192.168.31.64;192.168.73.2")

	InitFilesCmd.Flags().BoolVar(&CreateEmptyBlocks, "ceb", true, "Create empty blocks")

	InitFilesCmd.Flags().IntVar(&TimeoutPropose, "time-propose", DefaultTimeoutPropose, "Timeout Propose")
	InitFilesCmd.Flags().IntVar(&CreateEmptyBlocksInterval, "cebi", DefaultCreateEmptyBlocksInterval, "Create Empty Blocks Interval")
}

func initFiles(cmd *cobra.Command, args []string) error {
	return initFilesWithConfig()
}

// Used to create genesis files and modify peers in config
type tempConfig struct {
	nodeConfigPath    string
	nodeId            string
	actualGenesisPath string
}

// initFilesWithConfig initializes the necessary files and configurations for a side chain setup based on provided or default settings.
// It removes existing side-chain data, constructs file PVs, sets up directories, and generates configuration files for each validator.
// Additionally, it creates a GenesisDoc and modifies peer configurations across nodes. Returns an error if any step fails.
func initFilesWithConfig() error {

	sideChainPath := filepath.Join(RootDir)
	if err := os.RemoveAll(sideChainPath); err != nil {
		logger.Error("remove side-chain data fail", "err", err)
		return err
	}

	compressPath := filepath.Join(sideChainPath, DefaultCompressTargetFolder)
	if err := os.MkdirAll(compressPath, 0755); err != nil {
		logger.Error("create compress target folder fail", "err", err)
		return err
	}

	hostList := make([]string, 0)

	if strings.Contains(HostList, ";") && len(HostList) != 0 {
		hostList = strings.Split(HostList, ";")
	}

	// priv_validator_key list
	pvs := make([]*privval.FilePV, 0)
	// nodeId(8583a4e44cbff3ade6fff723a421a602c235e9ff) => p2pAddress(0.0.0.0:26656)
	p2pAddressMap := make(map[string]string)
	tempConfigs := make([]tempConfig, 0)

	for idx := range Validators {

		config := cfg.DefaultConfig()

		var p2pAddress string
		var nodePort int
		var tdPort int

		// set p2p rpc proxy address
		if len(hostList) == 0 { // local
			logger.Info("local mode")
			config.P2P.ListenAddress = fmt.Sprintf("tcp://127.0.0.1:%d", DefaultTdP2pPort+idx*PortSpacingFactor)
			config.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", DefaultTdRpcPort+idx*PortSpacingFactor)
			config.BaseConfig.ProxyApp = fmt.Sprintf("tcp://127.0.0.1:%d", DefaultTdProxyPort+idx*PortSpacingFactor)

			p2pAddress = fmt.Sprintf("127.0.0.1:%d", DefaultTdP2pPort+idx*PortSpacingFactor)
			nodePort = DefaultNodePort + idx*PortSpacingFactor
			tdPort = DefaultTdRpcPort + idx*PortSpacingFactor

		} else { // cluster
			logger.Info("cluster mode")
			config.P2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", DefaultTdP2pPort)
			config.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", DefaultTdRpcPort)
			config.BaseConfig.ProxyApp = fmt.Sprintf("tcp://0.0.0.0:%d", DefaultTdProxyPort)

			p2pAddress = fmt.Sprintf("%s:26656", hostList[idx])
			nodePort = DefaultNodePort
			tdPort = DefaultTdRpcPort
		}

		if !CreateEmptyBlocks {
			config.Consensus.CreateEmptyBlocks = false
		} else {
			config.Consensus.CreateEmptyBlocksInterval = time.Duration(CreateEmptyBlocksInterval) * time.Second
		}

		nodePath := filepath.Join(sideChainPath, fmt.Sprintf("%d", idx))
		if err := tmos.EnsureDir(nodePath, 0700); err != nil {
			logger.Error("ensure node dir fail", "err", err)
			return err
		}

		nodeConfigPath := filepath.Join(nodePath, "config")
		if err := tmos.EnsureDir(nodeConfigPath, 0700); err != nil {
			logger.Error("ensure config dir fail", "err", err)
			return err
		}

		nodeDataPath := filepath.Join(nodePath, "data")
		if err := tmos.EnsureDir(nodeDataPath, 0700); err != nil {
			logger.Error("ensure data dir fail", "err", err)
			return err
		}

		config.P2P.AllowDuplicateIP = true

		// set home
		config.RootDir = nodePath
		// set addr book, use full path. like: $HOME/.side-chain/0/config/addrbook.json
		config.P2P.AddrBook = filepath.Join(nodeConfigPath, "addrbook.json")
		// set priv validator, use partial paths. RootDir+PrivValidatorKey/PrivValidatorState will be used when creating td
		config.PrivValidatorKey = "config/priv_validator_key.json"
		config.PrivValidatorState = "data/priv_validator_state.json"
		// set wal path, use full path. like: $HOME/.side-chain/0/data/cs.wal/wal
		config.Consensus.WalPath = filepath.Join(nodePath, config.Consensus.WalPath)

		config.Consensus.TimeoutPropose = time.Duration(TimeoutPropose) * time.Second
		config.Consensus.TimeoutProposeDelta = time.Duration(TimeoutPropose+3) * time.Second

		// set the transaction size to 10mb
		config.Mempool.MaxTxBytes = DefaultBlockMaxTxBytes
		config.RPC.MaxBodyBytes = int64(DefaultRpcMaxBodyBytes)
		// set genesis path, use partial paths. RootDir+Genesis will be used when creating td
		config.Genesis = "config/genesis.json"

		// create file
		{
			actualP2PAddrBook := filepath.Join(nodeConfigPath, "addrbook.json")
			actualPrivValidatorKey := filepath.Join(nodeConfigPath, "priv_validator_key.json")
			actualPrivValidatorState := filepath.Join(nodeDataPath, "priv_validator_state.json")
			var pv *privval.FilePV

			if tmos.FileExists(actualP2PAddrBook) {
				pv = privval.LoadFilePV(actualPrivValidatorKey, actualPrivValidatorState)
				logger.Info("Found private validator", "keyFile", actualPrivValidatorKey,
					"stateFile", actualPrivValidatorState)
			} else {
				pv = privval.GenFilePV(actualPrivValidatorKey, actualPrivValidatorState)
				pv.Save()
				logger.Info("Generated private validator", "keyFile", actualPrivValidatorKey,
					"stateFile", actualPrivValidatorState)
			}
			pvs = append(pvs, pv)
		}

		// set node key, use partial paths. RootDir+NodeKey will be used when creating td
		config.NodeKey = "config/node_key.json"
		{
			actualNodeKey := filepath.Join(nodeConfigPath, "node_key.json")
			if tmos.FileExists(actualNodeKey) {
				logger.Info("Found node key", "path", actualNodeKey)
			} else {
				nodekey, err := p2p.LoadOrGenNodeKey(actualNodeKey)
				if err != nil {
					logger.Error("create node key fail", "err", err)
					return err
				}

				// set node map p2p address
				p2pAddressMap[string(nodekey.ID())] = p2pAddress

				tempConfigs = append(tempConfigs, tempConfig{
					nodeConfigPath:    nodeConfigPath,
					nodeId:            string(nodekey.ID()),
					actualGenesisPath: filepath.Join(nodeConfigPath, "genesis.json"),
				})

				logger.Info("Generated node key", "path", config.NodeKey)
			}
		}

		// create node config
		{
			nodeConfig := coreCfg.DefaultConfig(idx, nodePort, tdPort)
			nodeConfigFile := filepath.Join(nodeConfigPath, "node.toml")

			content, err := toml.Marshal(nodeConfig)
			if err != nil {
				return err
			}

			if err = tmos.WriteFile(nodeConfigFile, content, 0644); err != nil {
				logger.Error("write config fail", "err", err)
				return err
			}
		}

		configs = append(configs, config)

	}

	// gen genesis
	genDoc := types.GenesisDoc{
		ChainID:         fmt.Sprintf("side-chain-%v", tmrand.Str(6)),
		GenesisTime:     tmtime.Now(),
		ConsensusParams: types.DefaultConsensusParams(),
	}
	// set block max size 10mb
	genDoc.ConsensusParams.Block.MaxBytes = int64(DefaultBlockMaxTxBytes)

	for _, pv := range pvs {
		pubKey, err := pv.GetPubKey()
		if err != nil {
			logger.Error("get pubkey fail", "err", err)
			return err
		}

		genDoc.Validators = append(genDoc.Validators, types.GenesisValidator{
			Address: pubKey.Address(),
			PubKey:  pubKey,
			Power:   10,
		})
	}

	// Traverse to create genesis files and modify peers in config
	for idx, config := range configs {

		tempConfig := tempConfigs[idx]

		if err := genDoc.SaveAs(tempConfig.actualGenesisPath); err != nil {
			logger.Error("genesis save fail", "err", err)
			return err
		}

		seeds := genSeedsString(p2pAddressMap, tempConfig.nodeId)

		//config.P2P.Seeds = seeds
		config.P2P.PersistentPeers = seeds

		{
			configFile := filepath.Join(tempConfig.nodeConfigPath, "config.toml")

			var configMap map[string]interface{}
			if err := mapstructure.Decode(config, &configMap); err != nil {
				logger.Error("decode config fail", "err", err)
				return err
			}

			var configContent bytes.Buffer
			if err := toml.NewEncoder(&configContent).Encode(configMap); err != nil {
				logger.Error("encode config fail", "err", err)
				return err
			}

			if err := tmos.WriteFile(configFile, configContent.Bytes(), 0644); err != nil {
				logger.Error("write config fail", "err", err)
				return err
			}
		}

		compressFolderPath := fmt.Sprintf("%s/%d", sideChainPath, idx)
		tarPath := fmt.Sprintf("%s/%d.tar.gz", compressPath, idx)

		err := compressWithTar(compressFolderPath, tarPath)
		if err != nil {
			logger.Error("compress fail", "err", err)
			return err
		}

	}

	return nil
}

// genSeedsString generates a comma-separated string of key-value pairs from a map,
// excluding the entry with a key matching the provided filter.
// Each pair is formatted as "key@value". The resulting string is stripped of trailing commas.
func genSeedsString(m map[string]string, filter string) string {
	var str string

	for key, val := range m {
		if key == filter {
			continue
		}
		str += fmt.Sprintf("%s@%s", key, val) + ","
	}

	str = strings.TrimSuffix(str, ",")
	return str
}

func compressWithTar(folderPath, tarGzPath string) error {
	cmd := exec.Command("tar", "-czf", tarGzPath, "-C", folderPath, ".")
	return cmd.Run()
}
