package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
)

var QueryCmd = &cobra.Command{
	Use:   "query [balance] [nonce]",
	Short: "Query balance/nonce",
	Args:  cobra.ExactArgs(1),
	RunE:  query,
}

func init() {
	QueryCmd.Flags().StringVarP(&QueryAddress, "address", "a", DefaultAccountAddress.String(), "Query by address")
	QueryCmd.Flags().StringVarP(&MintNodeRpc, "node-rpc", "n", DefaultMintNodeRpc, "RPC server address")
}

func query(cmd *cobra.Command, args []string) error {
	arg := args[0]
	baseUrl := fmt.Sprintf("http://%s:%d", DefaultHostList, DefaultNodePort)
	if len(MintNodeRpc) != 0 || MintNodeRpc != DefaultMintNodeRpc {
		baseUrl = MintNodeRpc
	}

	switch strings.ToLower(arg) {
	case "balance":
		balance, err := getBalance(QueryAddress, baseUrl)
		if err != nil {
			logger.Error("query balance", "err", err)
		}
		logger.Info("account balance", "balance", balance)
	case "nonce":
		nonce, err := getNonce(QueryAddress, baseUrl)
		if err != nil {
			logger.Error("query balance", "err", err)
		}
		logger.Info("account nonce", "nonce", nonce)
	default:
		logger.Error("invalid query type")
		return nil
	}
	return nil
}

func getBalance(address, baseUrl string) (string, error) {
	url := fmt.Sprintf("%s/balance/%s", baseUrl, address)
	response, err := http.Get(url)
	if err != nil {
		logger.Error("get account balance error", "err", err)
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return "", err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		logger.Error("parse bytes to map[string]interface{}:", err)
		return "", err
	}
	data, ok := m["data"].(map[string]interface{})
	if !ok {
		logger.Error("parse data to map[string]interface{}:", m)
		return "", nil
	}
	balance, ok := data["balance"].(string)
	if !ok {
		logger.Error("parse balance to string:", data)
		return "", nil
	}

	return balance, nil
}

func getNonce(address, baseUrl string) (float64, error) {
	url := fmt.Sprintf("%s/nonce/%s", baseUrl, address)
	response, err := http.Get(url)
	if err != nil {
		logger.Error("get account balance error", "err", err)
		return 0, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return 0, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		logger.Error("parse bytes to map[string]interface{}:", err)
		return 0, err
	}
	data, ok := m["data"].(map[string]interface{})
	if !ok {
		logger.Error("parse data to map[string]interface{}:", m)
		return 0, fmt.Errorf("parse data to map[string]interface{}: %s", m)
	}
	// uint64 to float64.........
	nonce, ok := data["nonce"].(float64)
	if !ok {
		logger.Error("parse nonce to float64:", data)
		return 0, fmt.Errorf("parse nonce to float64: %s", data)
	}

	return nonce, nil
}
