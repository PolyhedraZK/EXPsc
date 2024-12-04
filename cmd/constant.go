package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

var (
	Validators        int
	DefaultValidators = 1

	RootDir        string
	DefaultRootDir = os.Getenv("HOME") + "/.side-chain"

	ValidatorDir        string
	DefaultValidatorDir = os.Getenv("HOME") + "/.side-chain/0"

	DefaultCompressTargetFolder = "compress"

	HostList        string
	DefaultHostList = "127.0.0.1"

	CreateEmptyBlocks bool
	DefaultNodeDbPath = os.Getenv("HOME") + "/.side-chain/0/node_db"

	MintTdRpc        string
	DefaultMintTdRpc = fmt.Sprintf("http://%s:%d", DefaultHostList, DefaultTdRpcPort)

	MintNodeRpc        string
	DefaultMintNodeRpc = fmt.Sprintf("http://%s:%d", DefaultHostList, DefaultNodePort)

	MintPrivateKeyPath        string
	DefaultMintPrivateKeyPath = ""

	QueryAddress string

	PortSpacingFactor = 100

	DefaultNodePort        = 7074
	DefaultBlockMaxTxBytes = 104857600 // 100MB
	DefaultRpcMaxBodyBytes = 209715200 // 200MB, max request body < 100MB

	TimeoutPropose            int
	CreateEmptyBlocksInterval int

	DefaultTimeoutPropose            = 9 // Use with DefaultMaxTxBytes=100MB
	DefaultCreateEmptyBlocksInterval = 20

	DefaultTdP2pPort   = 26656
	DefaultTdRpcPort   = 26657
	DefaultTdProxyPort = 26658

	DefaultAccountAddress    = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	DefaultAccountPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	DefaultAmount            = "0xde0b6b3a7640000" // 1ether
)

var (
	LogMaxSize    = 128 // mb
	LogMaxAge     = 30  // day
	LogMaxBackups = 100 //
	LogPath       = "log/node.log"
)
