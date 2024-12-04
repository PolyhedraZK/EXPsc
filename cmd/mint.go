package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nbnet/side-chain/core/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/os"
	"io/ioutil"
	"net/http"
	"strings"
)

var MintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint Token",
	RunE:  mint,
}

func init() {
	MintCmd.Flags().StringVarP(&MintTdRpc, "td-rpc", "r", DefaultMintTdRpc, "RPC server address")
	MintCmd.Flags().StringVarP(&MintNodeRpc, "node-rpc", "n", DefaultMintNodeRpc, "RPC server address")
	MintCmd.Flags().StringVarP(&MintPrivateKeyPath, "privatekey-path", "k", DefaultMintPrivateKeyPath, "Mint private key path")
}

func mint(cmd *cobra.Command, args []string) error {

	var privateKey *ecdsa.PrivateKey
	var address common.Address
	var err error

	if len(MintPrivateKeyPath) == 0 {
		privateKey, err = crypto.HexToECDSA(DefaultAccountPrivateKey)
		if err != nil {
			logger.Error("load private key error", "err", err)
			return err
		}
		address = crypto.PubkeyToAddress(privateKey.PublicKey)
	} else {
		b, err := os.ReadFile(MintPrivateKeyPath)
		if err != nil {
			logger.Error("read private key error", "err", err)
			return err
		}
		privateKey, err = crypto.HexToECDSA(strings.TrimSpace(string(b)))
		if err != nil {
			logger.Error("load private key error", "err", err)
			return err
		}
		address = crypto.PubkeyToAddress(privateKey.PublicKey)
	}

	nonce, err := getNonce(address.String(), MintNodeRpc)

	body := types.MintBody{
		Nonce:   uint64(nonce),
		Amount:  DefaultAmount,
		Address: address.String(),
	}

	mintTx := types.Tx{
		Ty:   types.Mint,
		Body: body,
	}

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err = encoder.Encode(mintTx.Body)
	if err != nil {
		logger.Error("RLP encode error", err)
		return err
	}

	digestHash, err := body.DigestHash()
	if err != nil {
		logger.Error("Digest hash error", err)
		return err
	}

	signature, err := crypto.Sign(digestHash, privateKey)
	if err != nil {
		logger.Error("Sign error", err)
		return err
	}

	mintTx.Signature = common.Bytes2Hex(signature)

	j, err := json.Marshal(&mintTx)
	if err != nil {
		logger.Error("JSON marshal error", err)
		return err
	}
	tx := base64.StdEncoding.EncodeToString(j)
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "broadcast_tx_sync",
		"params": map[string]interface{}{
			"tx": tx,
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		logger.Error("Error marshalling JSON:", err)
		return err
	}

	//url := fmt.Sprintf("http://%s:%d", DefaultHostList, DefaultTdRpcPort)

	response, err := http.Post(MintTdRpc, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Error sending POST request:", err)
		return err
	}
	defer response.Body.Close()

	// 读取响应内容
	context, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return err
	}

	logger.Info("Response", "Response Status:", response.Status)
	logger.Info("Response", "Response Body:", string(context))
	logger.Info("Mint Account", "Address", address, "Amount", "1ether")

	return nil
}
