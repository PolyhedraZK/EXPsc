package types

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/json-iterator/go"
	"github.com/nbnet/side-chain/core/utils"
	"math/big"
	"strings"
)

type TxType int

const (
	UnKnown TxType = iota + 1
	Mint
	Blob
)

func (t TxType) String() string {
	switch t {
	case Mint:
		return "mint"
	case Blob:
		return "blob"
	default:
		return "unknown"
	}
}

func (t TxType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TxType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case "mint":
		*t = Mint
	case "blob":
		*t = Blob
	default:
		*t = UnKnown
	}

	return nil
}

type Tx struct {
	Ty        TxType      `json:"type" mapstructure:"type"`
	Signature string      `json:"signature" mapstructure:"signature"`
	Body      interface{} `json:"body" mapstructure:"body"`
}

func (t *Tx) VerifySignature(address common.Address, digestHash []byte) error {

	pubkey, err := crypto.SigToPub(digestHash, common.Hex2Bytes(t.Signature))
	if err != nil {
		return err
	}

	recoverAddress := crypto.PubkeyToAddress(*pubkey)

	if address != recoverAddress {
		return fmt.Errorf("recover address %s not equal to address %s", recoverAddress, address)
	}

	return nil
}

func (t *Tx) GenGzipCompressBlobTx(body BlobBody) error {

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	utils.RemoveHexPrefix(body.Data)
	dataBytes, err := hex.DecodeString(body.Data)
	if err != nil {
		return err
	}

	if _, err := writer.Write(dataBytes); err != nil {
		return err
	}
	writer.Close()

	compressData := hex.EncodeToString(buf.Bytes())
	t.Ty = Blob
	t.Body = BlobBody{
		Data:    compressData,
		Address: body.Address,
	}

	return nil
}

func (t *Tx) ToBytes() ([]byte, error) {
	jsonType := jsoniter.ConfigCompatibleWithStandardLibrary

	result, _ := jsonType.Marshal(t)

	return result, nil
}

func JsonTxs(txs []*Tx) string {
	j, _ := json.Marshal(txs)
	return string(j)
}

type MintBody struct {
	Nonce   uint64 `json:"nonce" mapstructure:"nonce"`
	Amount  string `json:"amount" mapstructure:"amount"`
	Address string `json:"address" mapstructure:"address"`
}

func (m *MintBody) DigestHash() ([]byte, error) {

	jsonType := jsoniter.ConfigCompatibleWithStandardLibrary

	result, err := jsonType.Marshal(m)
	if err != nil {
		return nil, err
	}

	digestHash := sha256.Sum256(result)
	return digestHash[:], nil
}

type BlobBody struct {
	Data    string `json:"data" mapstructure:"data"`
	Address string `json:"address" mapstructure:"address"`
}

func (b *BlobBody) Gas() *big.Int {
	bytesLen := big.NewInt(int64(len(b.Data)))
	gas := bytesLen.Mul(bytesLen, DefaultPerByteFee)
	return gas
}
