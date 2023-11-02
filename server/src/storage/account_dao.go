/*
 * @Date: 2023-10-28 09:38:27
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-02 12:22:41
 * @FilePath: /1024-dev/1024/server/src/storage/account_dao.go
 */
package storage

import (
	"context"
	"dev1024/src/entities"
)

type AccountDao interface {
	GetJwtSignedKey(ctx context.Context) (string, error)
	SetJwtSignedKey(ctx context.Context, key string) error
	DelJwtSignedKey(ctx context.Context) error
	Create(ctx context.Context, account *entities.Account) (uint64, error)
	GetAccountByUsername(ctx context.Context, username string) (entities.Account, error)
	GetAccountByID(ctx context.Context, uid uint64) (entities.Account, error)
}

type AccountDaoBase struct {
	instance AccountDao
}

// GetAccountByID implements AccountDao.
func (t*AccountDaoBase) GetAccountByID(ctx context.Context, uid uint64) (entities.Account, error) {
	return t.instance.GetAccountByID(ctx,uid)
}

// DelJwtSignedKey implements AccountDao.
func (t *AccountDaoBase) DelJwtSignedKey(ctx context.Context) error {
	return t.instance.DelJwtSignedKey(ctx)
}

// GetJwtSignedKey implements AccountDao.
func (t *AccountDaoBase) GetJwtSignedKey(ctx context.Context) (string, error) {
	return t.instance.GetJwtSignedKey(ctx)
}

// GetUidAndPwdByUsername implements AccountDao.
func (t *AccountDaoBase) GetAccountByUsername(ctx context.Context, username string) (entities.Account, error) {
	return t.instance.GetAccountByUsername(ctx, username)
}

// Create implements AccountDao.
func (t *AccountDaoBase) Create(ctx context.Context, account *entities.Account) (uint64, error) {
	return t.instance.Create(ctx, account)
}

// SetJwtSignedKey implements AccountDao.
func (t *AccountDaoBase) SetJwtSignedKey(ctx context.Context, key string) error {
	return t.instance.SetJwtSignedKey(ctx, key)
}

func NewAccountDao(instance AccountDao) AccountDao {
	return &AccountDaoBase{
		instance: instance,
	}
}

var _ AccountDao = (*AccountDaoBase)(nil)
