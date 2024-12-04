package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Db interface {
	AddAccountBalance(address common.Address, amount *big.Int) error
	SubAccountBalance(address common.Address, amount *big.Int) error
	UpdateAccountNonce(address common.Address) error
	GetAccountBalance(address common.Address) (*big.Int, error)
	GetAccountNonce(address common.Address) (*big.Int, error)
}
