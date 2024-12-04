package main

import (
	"github.com/natefinch/lumberjack"
	"github.com/nbnet/side-chain/core/service"
	coreTypes "github.com/nbnet/side-chain/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	tmLog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start side chain node",
	RunE:  start,
}

func init() {
	StartCmd.Flags().StringVarP(&ValidatorDir, "validator-dir", "v", DefaultValidatorDir, "Node directory")
}

// start initializes and starts a side chain node based on provided configurations.
// It reads configuration from TOML files, generates a node instance with the configurations,
// starts the node, and sets up a signal handler to gracefully stop the node upon receiving an interrupt signal.
// Returns an error if configuration parsing fails, otherwise, it exits the program after a signal is received.
func start(cmd *cobra.Command, args []string) error {

	tdConfig := &cfg.Config{}
	nodeConfig := &coreTypes.Config{}

	configPath := filepath.Join(ValidatorDir, "config")

	l, output := genLogger()

	tdConfigPath := filepath.Join(configPath, "config.toml")
	viper.SetConfigFile(tdConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&tdConfig); err != nil {
		panic(err)
	}

	nodeConfigPath := filepath.Join(configPath, "node.toml")
	viper.SetConfigFile(nodeConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&nodeConfig); err != nil {
		panic(err)
	}

	db := service.NewDbService(nodeConfig.Db, l)

	rpc := service.NewRpc(nodeConfig.Rpc, db, l, output)
	rpc.Start()

	n := genNode(tdConfig, db, l)

	n.Start()

	defer func() {
		n.Stop()
		n.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
	return nil
}

func genNode(config *cfg.Config, db coreTypes.Db, l tmLog.Logger) *node.Node {
	abci := service.NewAbci(db, l)
	return coreTypes.NewTd(abci, config, l)
}

func genLogger() (tmLog.Logger, io.Writer) {
	logPath := filepath.Join(ValidatorDir, LogPath)
	logWriter := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    LogMaxSize,
		MaxAge:     LogMaxAge,
		MaxBackups: LogMaxBackups,
		LocalTime:  false,
		Compress:   false,
	}

	l := tmLog.NewTMLogger(logWriter)
	return l, logWriter
}
