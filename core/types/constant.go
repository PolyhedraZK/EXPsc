package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	BalanceKeyPrefix = []byte("balance")
	NonceKeyPrefix   = []byte("nonce")

	// wei
	DefaultPerByteFee  = new(big.Int).SetUint64(10)
	DefaultAddress     = common.HexToAddress("0x0000000000000000000000000000000000000000")
	DefaultHash        = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	DefaultAddressSize = 40
)

var (
	DeliverTxTitle            = "DeliverTx"
	EndBlockTitle             = "EndBlock"
	CommitTitle               = "Commit"
	UpdateAccountBalanceTitle = "UpdateAccountBalance"
	UpdateAccountNonceTitle   = "UpdateAccountNonce"
	GetAccountBalanceTitle    = "GetAccountBalance"
	GetAccountNonceTitle      = "GetAccountNonce"
	ProcessTxTitle            = "ProcessTx"
	BalanceHandlerTitle       = "BalanceHandler"
	NonceHandlerTitle         = "NonceHandler"
	BlobHandlerTitle          = "BlobHandler"
)

var (
	ErrDecodeTx              = "DecodeTxError"
	ErrEncodeTx              = "EncodeTxError"
	ErrDecodeMintBody        = "DecodeMintBodyError"
	ErrDecodeBlobBody        = "DecodeBlobBodyError"
	ErrCalculateDigestHash   = "CalculateDigestHashError"
	ErrVerifySignature       = "VerifySignatureError"
	ErrUpdateNonce           = "UpdateNonceError"
	ErrUpdateBalance         = "UpdateBalanceError"
	ErrDecodeAmount          = "DecodeAmountError"
	ErrUnknownTxBody         = "UnknownTxBody"
	ErrGetBalance            = "GetBalanceError"
	ErrGetNonce              = "GetNonceError"
	ErrInsufficientBalance   = "InsufficientBalance"
	ErrNonceNotMatch         = "NonceNotMatch"
	ErrInvalidAddress        = "InvalidAddress"
	ErrGenGzipCompressBlobTx = "GenGzipCompressBlobTxError"
	ErrBroadcastTxSync       = "BroadcastTxSyncError"
	ErrProcessCommit         = "ProcessCommitError"
	ErrTxToBytes             = "TxToBytesErr"
)

func BalanceKey(address common.Address) []byte {
	return append(BalanceKeyPrefix, address.Bytes()...)
}

func NonceKey(address common.Address) []byte {
	return append(NonceKeyPrefix, address.Bytes()...)
}
