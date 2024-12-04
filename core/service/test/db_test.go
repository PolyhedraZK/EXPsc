package test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nbnet/side-chain/core/service"
	"github.com/nbnet/side-chain/core/types"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"os"
	"testing"
)

// TestDbBalance tests the database service's ability to update and retrieve an account's balance.
// It creates a new database service instance, sets an initial balance for a given Ethereum address,
// then fetches and logs the balance to ensure it was correctly stored and retrieved.
func TestDbBalance(t *testing.T) {

	address := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	path := "./.test_db/node_db"
	amount := new(big.Int).SetUint64(1000000)

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	config := types.DbConfig{
		Path: path,
	}

	db := service.NewDbService(&config, logger)

	// set
	err := db.AddAccountBalance(address, amount)
	if err != nil {
		panic(err)
	}

	// get
	balance, err := db.GetAccountBalance(address)
	if err != nil {
		panic(err)
	}
	logger.Info("TestDbBalance", "Balance", balance)
}

// TestDbNonce tests the functionality of updating and retrieving a nonce for an Ethereum account using a database service.
// It sets a nonce value for a given address, then fetches and logs the retrieved nonce to ensure it matches the set value.
// Parameters:
// t *testing.T - The testing object used for assertions and logging within the test.
func TestDbNonce(t *testing.T) {

	address := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	path := "./.test_db/.node_db"

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	config := types.DbConfig{
		Path: path,
	}

	db := service.NewDbService(&config, logger)

	// set
	err := db.UpdateAccountNonce(address)
	if err != nil {
		panic(err)
	}

	// get
	nonce, err := db.GetAccountNonce(address)
	if err != nil {
		panic(err)
	}
	logger.Info("TestDbNonce", "Nonce", nonce)
}
