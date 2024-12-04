package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
	"github.com/nbnet/side-chain/core/types"
	tdTypes "github.com/tendermint/tendermint/abci/types"
	tmLog "github.com/tendermint/tendermint/libs/log"
	"math/big"
	"strings"
)

type Abci struct {
	Db      types.Db
	log     types.CustomLogger
	appHash []byte
}

func NewAbci(db types.Db, logger tmLog.Logger) *Abci {
	return &Abci{
		Db:      db,
		log:     types.CustomLogger{Logger: logger},
		appHash: types.DefaultHash.Bytes(),
	}
}

func (s *Abci) Info(info tdTypes.RequestInfo) tdTypes.ResponseInfo {
	return tdTypes.ResponseInfo{}
}

func (s *Abci) BeginBlock(block tdTypes.RequestBeginBlock) tdTypes.ResponseBeginBlock {
	return tdTypes.ResponseBeginBlock{}
}

func (s *Abci) CheckTx(tdTx tdTypes.RequestCheckTx) tdTypes.ResponseCheckTx {

	result := s.processCheckTx(tdTx.GetTx())

	return tdTypes.ResponseCheckTx{
		Code:      result.code,
		Log:       result.log,
		Info:      result.info,
		GasWanted: result.gas,
	}
}

func (s *Abci) DeliverTx(tdTx tdTypes.RequestDeliverTx) tdTypes.ResponseDeliverTx {

	result := s.processDeliverTx(tdTx.GetTx())

	return tdTypes.ResponseDeliverTx{
		Code:    result.code,
		Log:     result.log,
		Info:    result.info,
		GasUsed: result.gas,
	}
}

func (s *Abci) EndBlock(block tdTypes.RequestEndBlock) tdTypes.ResponseEndBlock {
	return tdTypes.ResponseEndBlock{}
}

func (s *Abci) Commit() tdTypes.ResponseCommit {
	s.log.Info(types.CommitTitle, "app_hash", fmt.Sprintf("%x", s.appHash))

	return tdTypes.ResponseCommit{
		Data: s.appHash,
	}
}

func (s *Abci) Query(query tdTypes.RequestQuery) tdTypes.ResponseQuery {
	return tdTypes.ResponseQuery{}
}

func (s *Abci) InitChain(chain tdTypes.RequestInitChain) tdTypes.ResponseInitChain {
	return tdTypes.ResponseInitChain{}
}

func (s *Abci) ListSnapshots(snapshots tdTypes.RequestListSnapshots) tdTypes.ResponseListSnapshots {
	return tdTypes.ResponseListSnapshots{}
}

func (s *Abci) OfferSnapshot(snapshot tdTypes.RequestOfferSnapshot) tdTypes.ResponseOfferSnapshot {
	return tdTypes.ResponseOfferSnapshot{}
}

func (s *Abci) LoadSnapshotChunk(chunk tdTypes.RequestLoadSnapshotChunk) tdTypes.ResponseLoadSnapshotChunk {
	return tdTypes.ResponseLoadSnapshotChunk{}
}

func (s *Abci) ApplySnapshotChunk(chunk tdTypes.RequestApplySnapshotChunk) tdTypes.ResponseApplySnapshotChunk {
	return tdTypes.ResponseApplySnapshotChunk{}
}

func (s *Abci) SetOption(option tdTypes.RequestSetOption) tdTypes.ResponseSetOption {
	return tdTypes.ResponseSetOption{}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// internalResult is a struct used to encapsulate the results of processing a transaction within the ABCI application.
// It includes details like log messages, additional info, gas usage, and an error code indicating the status of the transaction processing.
type internalResult struct {
	log     string
	info    string
	gas     int64
	code    uint32
	address common.Address
	ty      types.TxType
}

func (s *Abci) processCheckTx(txBytes []byte) internalResult {
	gas := int64(0)

	var tx types.Tx
	err := json.Unmarshal(txBytes, &tx)
	if err != nil {
		s.log.Error(types.ProcessTxTitle, types.ErrDecodeTx, err)
		return internalResult{
			code: 1,
			info: err.Error(),
			log:  types.ErrDecodeTx,
		}
	}

	switch tx.Ty {
	case types.Mint:
		var body types.MintBody
		err := mapstructure.Decode(tx.Body, &body)
		if err != nil {
			s.log.Error(types.ProcessTxTitle, types.ErrDecodeMintBody, err)
			return internalResult{
				code: 1,
				log:  types.ErrDecodeMintBody,
				info: err.Error(),
			}
		}

		// calculate hash
		digestHash, err := body.DigestHash()
		if err != nil {
			s.log.Error(types.ProcessTxTitle, types.ErrCalculateDigestHash, err)
			return internalResult{
				code: 1,
				log:  types.ErrCalculateDigestHash,
				info: err.Error(),
			}
		}
		address := common.HexToAddress(body.Address)

		// check signature
		err = tx.VerifySignature(address, digestHash)
		if err != nil {
			s.log.Error(types.ProcessTxTitle, types.ErrVerifySignature, err)
			return internalResult{
				code: 1,
				log:  types.ErrVerifySignature,
				info: err.Error(),
			}
		}

		nonce := new(big.Int).SetUint64(body.Nonce)
		// check nonce
		dbNonce, err := s.Db.GetAccountNonce(address)
		if err != nil {
			s.log.Error(types.ProcessTxTitle, types.ErrGetNonce, err)
			return internalResult{
				code: 1,
				log:  types.ErrGetNonce,
				info: err.Error(),
			}
		}

		if dbNonce.Cmp(nonce) != 0 {
			s.log.Debug(types.ProcessTxTitle, types.ErrNonceNotMatch, "", "expected", dbNonce, "get", nonce)
			return internalResult{
				code: 1,
				log:  types.ErrNonceNotMatch,
			}
		}

	case types.Blob:
		var body types.BlobBody
		err := mapstructure.Decode(tx.Body, &body)
		if err != nil {
			s.log.Error(types.ProcessTxTitle, types.ErrDecodeBlobBody, err)
			return internalResult{
				code: 1,
				log:  types.ErrDecodeBlobBody,
				info: err.Error(),
			}
		}

		if len(body.Address) == types.DefaultAddressSize {
			return internalResult{
				code: 1,
				log:  types.ErrInvalidAddress,
			}
		}

		address := common.HexToAddress(body.Address)

		if address.Cmp(types.DefaultAddress) == 0 {
			return internalResult{
				code: 1,
				log:  types.ErrInvalidAddress,
			}
		}

		g := body.Gas()
		gas = g.Int64()

		balance, err := s.Db.GetAccountBalance(address)
		if err != nil {
			return internalResult{
				code: 1,
				log:  types.ErrGetBalance,
				info: err.Error(),
			}
		}

		if balance.Cmp(g) < 0 {
			return internalResult{
				code: 1,
				log:  types.ErrInsufficientBalance,
				gas:  gas,
			}
		}

	case types.UnKnown:
		fallthrough
	default:
		s.log.Error(types.ProcessTxTitle, types.ErrUnknownTxBody, tx.Body)
		return internalResult{
			code: 1,
			log:  types.ErrUnknownTxBody,
		}
	}

	return internalResult{gas: gas}
}

func (s *Abci) processDeliverTx(txBytes []byte) internalResult {

	address := types.DefaultAddress
	var tx types.Tx
	// Success by default, only successful in checkTx will reach here
	_ = json.Unmarshal(txBytes, &tx)

	bytes := append(s.appHash, txBytes...)
	appHash := sha256.Sum256(bytes)
	s.appHash = appHash[:]

	switch tx.Ty {
	case types.Mint:

		var body types.MintBody
		// Success by default, only successful in checkTx will reach here
		_ = mapstructure.Decode(tx.Body, &body)
		address = common.HexToAddress(body.Address)

		err := s.Db.UpdateAccountNonce(address)
		if err != nil {
			s.log.Error(types.ProcessTxTitle, types.ErrUpdateNonce, err)
			return internalResult{
				code:    1,
				log:     types.ErrUpdateNonce,
				info:    err.Error(),
				address: address,
			}
		}

		a := body.Amount
		if strings.Contains(body.Amount, "0x") {
			a = body.Amount[2:]
		}

		if amount, ok := new(big.Int).SetString(a, 16); ok {
			err := s.Db.AddAccountBalance(address, amount)
			if err != nil {
				s.log.Error(types.ProcessTxTitle, types.ErrUpdateBalance, err)
				return internalResult{
					code:    1,
					log:     types.ErrUpdateBalance,
					info:    err.Error(),
					address: address,
				}
			}
		} else {
			s.log.Error(types.ProcessTxTitle, types.ErrDecodeAmount, body.Amount)
			return internalResult{
				code:    1,
				log:     types.ErrDecodeAmount,
				address: address,
			}
		}
	case types.Blob:
		var body types.BlobBody
		// Success by default, only successful in checkTx will reach here
		_ = mapstructure.Decode(tx.Body, &body)
		gas := body.Gas()
		address = common.HexToAddress(body.Address)
		if err := s.Db.SubAccountBalance(address, gas); err != nil {
			s.log.Error(types.ProcessTxTitle, types.ErrUpdateBalance, err)
			return internalResult{
				code:    1,
				log:     types.ErrUpdateBalance,
				info:    err.Error(),
				gas:     gas.Int64(),
				address: address,
			}
		}
	case types.UnKnown:
		fallthrough
	default:
		s.log.Error(types.ProcessTxTitle, types.ErrUnknownTxBody, tx.Body)
		return internalResult{
			code: 1,
			log:  types.ErrUnknownTxBody,
		}
	}

	return internalResult{
		address: address,
		ty:      tx.Ty,
	}
}
