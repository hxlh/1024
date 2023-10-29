/*
 * @Date: 2023-10-28 09:38:27
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-28 12:57:28
 * @FilePath: /1024/server/src/storage/account_dao.go
 */
package storage

import "dev1024/src/entities"

type AccountDao interface {
	GetJwtSignedKey() (string, error)
	SetJwtSignedKey(key string) error
	DelJwtSignedKey() error
	Save(account *entities.Account) (uint64, error)
	GetUidAndPwdByUsername(username string) (uint64,string, error)
}

type AccountDaoBase struct {
	instance AccountDao
}

// GetUidAndPwdByUsername implements AccountDao.
func (t *AccountDaoBase) GetUidAndPwdByUsername(username string) (uint64,string, error) {
	return t.instance.GetUidAndPwdByUsername(username)
}

// Save implements AccountDao.
func (t *AccountDaoBase) Save(account *entities.Account) (uint64, error) {
	return t.instance.Save(account)
}

// DelJwtSignedKey implements AccountDao.
func (t *AccountDaoBase) DelJwtSignedKey() error {
	return t.instance.DelJwtSignedKey()
}

// GetJwtSignedKey implements AccountDao.
func (t *AccountDaoBase) GetJwtSignedKey() (string, error) {
	return t.instance.GetJwtSignedKey()
}

// SetJwtSignedKey implements AccountDao.
func (t *AccountDaoBase) SetJwtSignedKey(key string) error {
	return t.instance.SetJwtSignedKey(key)
}

func NewAccountDao(instance AccountDao) AccountDao {
	return &AccountDaoBase{
		instance: instance,
	}
}

var _ AccountDao = (*AccountDaoBase)(nil)
