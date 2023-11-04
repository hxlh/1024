/*
 * @Date: 2023-10-28 09:38:27
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-02 14:52:56
 * @FilePath: /1024-dev/1024/server/src/storage/account_dao_impl.go
 */
package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"dev1024/src/entities"
	"encoding/hex"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const SIGNKEY_EXPIRE_TIME = 1 * time.Hour

type AccountDaoImpl struct {
}

// GetAccountByID implements AccountDao.
func (t*AccountDaoImpl) GetAccountByID(ctx context.Context, uid uint64) (entities.Account, error) {
	var account entities.Account
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return account, errors.WithStack(err)
	}
	var avatar sql.NullString
	err = SelectTableBy(tx, "video1024.account", []string{
		"username",
		"nickname",
		"pwd",
		"avatar",
	}, []any{
		&account.Username,
		&account.NickName,
		&account.Pwd,
		&avatar,
	}, "uid", uid)
	if err != nil {
		return account, errors.WithStack(err)
	}
	account.Uid=uid
	account.Avatar=avatar.String
	return account,nil
}

// GetUidAndPwdByUsername implements AccountDao.
func (t *AccountDaoImpl) GetAccountByUsername(ctx context.Context, username string) (entities.Account, error) {
	var account entities.Account
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return account, errors.WithStack(err)
	}
	return t.getAccountByUsername(tx, username)
}

func (t *AccountDaoImpl) getAccountByUsername(tx *sql.Tx, username string) (entities.Account, error) {
	var account entities.Account
	var avatar sql.NullString
	err := SelectTableBy(tx, "video1024.account", []string{
		"uid",
		"nickname",
		"pwd",
		"avatar",
	}, []any{
		&account.Uid,
		&account.NickName,
		&account.Pwd,
		&avatar,
	}, "username", username)
	if err != nil {
		return account, errors.WithStack(err)
	}
	account.Avatar = avatar.String
	return account, nil
}

// Create implements AccountDao.
func (t *AccountDaoImpl) Create(ctx context.Context, account *entities.Account) (uint64, error) {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	uid, err := t.create(tx, account)
	return uid, err
}

// RegisterAccount implements AccountDao.
func (t *AccountDaoImpl) create(tx *sql.Tx, account *entities.Account) (uint64, error) {
	stmt, err := tx.Prepare("INSERT INTO video1024.account(username,nickname,pwd,register_time) VALUES(?,?,?,?)")
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer stmt.Close()

	hash := sha256.New()
	hash.Write([]byte(account.Username))
	hash.Write([]byte(account.Pwd))
	hashByte := hash.Sum(nil)
	hashStr := hex.EncodeToString(hashByte)

	res, err := stmt.Exec(account.Username, account.NickName, hashStr,time.Now().UnixMilli())
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
func (t *AccountDaoImpl) DelJwtSignedKey(ctx context.Context) error {
	ctx, conn, err := getRedisConnFromCtx(ctx)
	if err != nil {
		return err
	}
	err = conn.Del(ctx, "jwt.sign_key").Err()
	if err != nil {
		return err
	}
	return nil
}

// GetJwtSignedKey implements AccountDao.
func (t *AccountDaoImpl) GetJwtSignedKey(ctx context.Context) (string, error) {
	ctx, conn, err := getRedisConnFromCtx(ctx)
	if err != nil {
		return "", err
	}

	ans, err := conn.Get(ctx, "jwt.sign_key").Result()
	if err != nil {
		if err == redis.Nil {
			return "", redis.Nil
		}
		return "", errors.WithStack(err)
	}
	return ans, nil
}

// SetJwtSignedKey implements AccountDao.
func (t *AccountDaoImpl) SetJwtSignedKey(ctx context.Context, key string) error {
	ctx, conn, err := getRedisConnFromCtx(ctx)
	if err != nil {
		return err
	}
	err = conn.SetNX(ctx, "jwt.sign_key", key, SIGNKEY_EXPIRE_TIME).Err()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func NewAccountDaoImpl() *AccountDaoImpl {
	return &AccountDaoImpl{}
}

var _ AccountDao = (*AccountDaoImpl)(nil)
