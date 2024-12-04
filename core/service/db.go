package service

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nbnet/side-chain/core/types"
	tmLog "github.com/tendermint/tendermint/libs/log"
	"math/big"
)

type DbService struct {
	config *types.DbConfig
	db     *badger.DB
	log    *types.CustomLogger
}

func NewDbService(config *types.DbConfig, logger tmLog.Logger) *DbService {

	log := &types.CustomLogger{
		Logger: logger,
	}

	db, err := badger.Open(
		badger.DefaultOptions(config.Path).WithLogger(log),
	)
	if err != nil {
		panic(err)
	}

	return &DbService{
		config: config,
		db:     db,
		log:    log,
	}
}

func (d *DbService) AddAccountBalance(address common.Address, amount *big.Int) error {
	result := d.db.Update(func(txn *badger.Txn) error {

		balance := big.NewInt(0)

		item, err := txn.Get(types.BalanceKey(address))

		if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
			d.log.Error(types.UpdateAccountBalanceTitle, types.ErrUpdateBalance, err)
			return err
		}

		if item != nil {
			err := item.Value(func(val []byte) error {
				balance.SetBytes(val)
				return nil
			})
			if err != nil {
				d.log.Error(types.UpdateAccountBalanceTitle, types.ErrUpdateBalance, err)
				return err
			}
		}

		err = txn.Set(types.BalanceKey(address), balance.Add(balance, amount).Bytes())
		if err != nil {
			d.log.Error(types.UpdateAccountBalanceTitle, types.ErrUpdateBalance, err)
			return err
		}
		d.log.Debug(types.UpdateAccountBalanceTitle, "Address", address, "Balance", balance)
		return nil
	})

	return result
}

func (d *DbService) SubAccountBalance(address common.Address, amount *big.Int) error {
	result := d.db.Update(func(txn *badger.Txn) error {

		balance := big.NewInt(0)

		item, err := txn.Get(types.BalanceKey(address))

		if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
			d.log.Error(types.UpdateAccountBalanceTitle, types.ErrUpdateBalance, err)
			return err
		}

		if item != nil {
			err := item.Value(func(val []byte) error {
				balance.SetBytes(val)
				return nil
			})
			if err != nil {
				d.log.Error(types.UpdateAccountBalanceTitle, types.ErrUpdateBalance, err)
				return err
			}
		}

		err = txn.Set(types.BalanceKey(address), balance.Sub(balance, amount).Bytes())
		if err != nil {
			d.log.Error(types.UpdateAccountBalanceTitle, types.ErrUpdateBalance, err)
			return err
		}
		d.log.Debug(types.UpdateAccountBalanceTitle, "Address", address, "Balance", balance)
		return nil
	})

	return result
}

func (d *DbService) UpdateAccountNonce(address common.Address) error {
	result := d.db.Update(func(txn *badger.Txn) error {

		nonce := big.NewInt(0)

		item, err := txn.Get(types.NonceKey(address))
		if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
			d.log.Error(types.UpdateAccountNonceTitle, types.ErrUpdateNonce, err)
			return err
		}

		if item != nil {
			err := item.Value(func(val []byte) error {
				nonce.SetBytes(val)
				return nil
			})
			if err != nil {
				d.log.Error(types.UpdateAccountNonceTitle, types.ErrUpdateNonce, err)
				return err
			}
		}

		err = txn.Set(types.NonceKey(address), nonce.Add(nonce, big.NewInt(1)).Bytes())
		if err != nil {
			d.log.Error(types.UpdateAccountNonceTitle, types.ErrUpdateNonce, err)
			return err
		}
		d.log.Debug(types.UpdateAccountNonceTitle, "Address", address, "Nonce", nonce)
		return nil
	})

	return result
}

func (d *DbService) GetAccountBalance(address common.Address) (*big.Int, error) {
	result := big.NewInt(0)
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(types.BalanceKey(address))

		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil
		}

		if err != nil {
			d.log.Error(types.GetAccountBalanceTitle, types.ErrGetBalance, err)
			return err
		}
		err = item.Value(func(val []byte) error {
			result.SetBytes(val)
			return nil
		})
		return err
	})

	return result, err
}

func (d *DbService) GetAccountNonce(address common.Address) (*big.Int, error) {
	result := big.NewInt(0)
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(types.NonceKey(address))

		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil
		}

		if err != nil {
			d.log.Error(types.GetAccountNonceTitle, types.ErrGetNonce, err)
			return err
		}
		err = item.Value(func(val []byte) error {
			result.SetBytes(val)
			return nil
		})
		return err
	})

	return result, err
}
