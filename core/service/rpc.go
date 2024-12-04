package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/nbnet/side-chain/core/types"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	tmTypes "github.com/tendermint/tendermint/types"
	"io"
	"time"
)

type Rpc struct {
	rpcConfig *types.RpcConfig
	db        types.Db
	log       types.CustomLogger
	engine    *gin.Engine
	tdClient  *tmClient.HTTP
}

func NewRpc(rpcConfig *types.RpcConfig, db types.Db, logger tmLog.Logger, output io.Writer) *Rpc {

	tdClient, err := tmClient.New(rpcConfig.TdRpc, "/websocket")
	if err != nil {
		panic(err)
	}

	gin.DefaultWriter = output
	gin.DefaultErrorWriter = output

	engine := gin.Default()
	engine.Use(gin.LoggerWithWriter(output))
	engine.Use(gin.RecoveryWithWriter(output))

	return &Rpc{
		rpcConfig: rpcConfig,
		db:        db,
		log: types.CustomLogger{
			Logger: logger,
		},
		engine:   engine,
		tdClient: tdClient,
	}
}

func (rpc *Rpc) Start() {
	go func() {
		rpc.log.Info("Starting RPC service")

		rpc.engine.GET("/balance/:address", rpc.balanceHandler)
		rpc.engine.GET("/nonce/:address", rpc.nonceHandler)
		rpc.engine.POST("/blob", rpc.blobHandler)
		rpc.engine.Run(fmt.Sprintf("%s:%d", rpc.rpcConfig.Host, rpc.rpcConfig.Port))
	}()

	time.Sleep(time.Second * 3)
}

func (rpc *Rpc) balanceHandler(c *gin.Context) {
	addressStr := c.Param("address")
	address := common.HexToAddress(addressStr)
	balance, err := rpc.db.GetAccountBalance(address)
	if err != nil {
		rpc.log.Error(types.BalanceHandlerTitle, types.ErrGetBalance, err)
		c.JSON(500, types.NewRpcResp(err, types.NewRpcBalanceData(nil, 1)))
		return
	}

	c.JSON(200, types.NewRpcResp(err, types.NewRpcBalanceData(balance, 0)))
}

func (rpc *Rpc) nonceHandler(c *gin.Context) {
	addressStr := c.Param("address")
	address := common.HexToAddress(addressStr)
	nonce, err := rpc.db.GetAccountNonce(address)
	if err != nil {
		rpc.log.Error(types.NonceHandlerTitle, types.ErrGetBalance, err)
		c.JSON(500, types.NewRpcResp(err, types.NewRpcNonceData(nil, 1)))
		return
	}

	c.JSON(200, types.NewRpcResp(err, types.NewRpcNonceData(nonce, 0)))
}

func (rpc *Rpc) blobHandler(c *gin.Context) {
	var body types.BlobBody
	if err := c.BindJSON(&body); err != nil {
		rpc.log.Error(types.BlobHandlerTitle, types.ErrDecodeMintBody, err)
		c.JSON(400, types.NewRpcResp(err, types.NewRpcBlobData(1, nil)))
		return
	}

	tx := types.Tx{}
	if err := tx.GenGzipCompressBlobTx(body); err != nil {
		rpc.log.Error(types.BlobHandlerTitle, types.ErrGenGzipCompressBlobTx, err)
		c.JSON(500, types.NewRpcResp(err, types.NewRpcBlobData(1, nil)))
		return
	}

	j, err := json.Marshal(tx)
	if err != nil {
		rpc.log.Error(types.BlobHandlerTitle, types.ErrEncodeTx, err)
		c.JSON(500, types.NewRpcResp(err, types.NewRpcBlobData(1, nil)))
		return
	}

	tdTx := tmTypes.Tx(j)
	result, err := rpc.tdClient.BroadcastTxSync(context.Background(), tdTx)
	if err != nil {
		rpc.log.Error(types.BlobHandlerTitle, types.ErrBroadcastTxSync, err)
		c.JSON(500, types.NewRpcResp(err, types.NewRpcBlobData(1, nil)))
		return
	}

	c.JSON(200, types.NewRpcResp(nil, types.NewRpcBlobData(0, result)))
}
