package types

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tmTypes "github.com/tendermint/tendermint/rpc/core/types"
	"math/big"
)

func NewRpcResp(err error, data gin.H) gin.H {
	return gin.H{
		"jsonrpc": "2.0",
		"id":      0,
		"error":   err,
		"data":    data,
	}
}

func NewRpcBalanceData(balance *big.Int, code int) gin.H {

	if balance == nil {
		balance = big.NewInt(0)
	}

	return gin.H{
		"code":    code,
		"balance": fmt.Sprintf("0x%x", balance.Bytes()),
	}
}

func NewRpcNonceData(nonce *big.Int, code int) gin.H {

	if nonce == nil {
		nonce = big.NewInt(0)
	}

	return gin.H{
		"code":  code,
		"nonce": nonce.Uint64(),
	}
}

func NewRpcBlobData(code int, data *tmTypes.ResultBroadcastTx) gin.H {

	result := gin.H{
		"code": code,
		"data": "",
		"log":  "",
		"hash": "",
	}

	if data != nil {
		result["code"] = data.Code
		result["data"] = data.Data.String()
		result["log"] = data.Log
		result["hash"] = data.Hash.String()
	}

	return result
}
