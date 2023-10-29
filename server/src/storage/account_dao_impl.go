/*
 * @Date: 2023-10-28 09:38:27
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-28 13:00:28
 * @FilePath: /1024/server/src/storage/account_dao_impl.go
 */
package storage

import (
	"context"
	"database/sql"
	"dev1024/src/entities"

	"github.com/pkg/errors"
)

type AccountDaoImpl struct {
	ctx    context.Context
	baseDB *BaseDB
}

// GetUidAndPwdByUsername implements AccountDao.
func (t *AccountDaoImpl) GetUidAndPwdByUsername(username string) (uint64, string, error) {
	tx, err := t.baseDB.DB.Begin()
	if err != nil {
		return 0, "", errors.WithStack(err)
	}
	uid, pwd, err := t.getUidAndPwdByUsername(tx, username)
	return uid, pwd, commitOrRollback(err, tx)
}

func (t *AccountDaoImpl) getUidAndPwdByUsername(tx *sql.Tx, username string) (uint64, string, error) {
	stmt, err := tx.Prepare("SELECT uid,pwd from video1024.account WHERE username = ?")
	if err != nil {
		return 0, "", errors.WithStack(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		return 0, "", errors.WithStack(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, "", errors.New("username does not exist")
	}
	var uid uint64
	var pwd string
	err = rows.Scan(&uid, &pwd)
	if err != nil {
		return 0, "", errors.WithStack(err)
	}
	return uid, pwd, nil
}

// Save implements AccountDao.
func (t *AccountDaoImpl) Save(account *entities.Account) (uint64, error) {
	tx, err := t.baseDB.DB.Begin()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	uid, err := t.save(tx, account)
	return uid, commitOrRollback(err, tx)
}

// RegisterAccount implements AccountDao.
func (t *AccountDaoImpl) save(tx *sql.Tx, account *entities.Account) (uint64, error) {
	stmt, err := tx.Prepare("INSERT INTO video1024.account(username,nickname,pwd) VALUES(?,?,?)")
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(account.Username, account.NickName, account.Pwd)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	// 获取主键
	uid, err := res.LastInsertId()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return uint64(uid), nil
}

// DelJwtSignedKey implements AccountDao.
func (t *AccountDaoImpl) DelJwtSignedKey() error {
	err := t.baseDB.RC.Del(t.ctx, "jwt.sign_key").Err()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetJwtSignedKey implements AccountDao.
func (t *AccountDaoImpl) GetJwtSignedKey() (string, error) {
	ans, err := t.baseDB.RC.Get(t.ctx, "jwt.sign_key").Result()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return ans, nil
}

// SetJwtSignedKey implements AccountDao.
func (t *AccountDaoImpl) SetJwtSignedKey(key string) error {
	err := t.baseDB.RC.Set(t.ctx, "jwt.sign_key", key, 0).Err()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func NewAccountDaoImpl(base *BaseDB, ctx context.Context) *AccountDaoImpl {
	return &AccountDaoImpl{
		ctx:    ctx,
		baseDB: base,
	}
}

var _ AccountDao = (*AccountDaoImpl)(nil)
