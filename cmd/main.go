package main

import (
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

var (
	configs = make([]*cfg.Config, 0)
	logger  = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "side chain",
		Short: "This is a side chain application",
	}

	rootCmd.AddCommand(
		InitFilesCmd,
		StartCmd,
		MintCmd,
		QueryCmd,
	)

	rootCmd.Execute()
}
