/*
 * @Date: 2023-10-27 07:58:50
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-01 16:45:44
 * @FilePath: /1024/server/src/service/account_service_impl.go
 */
package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"dev1024/src/entities"
	"dev1024/src/storage"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const JWT_EXPIRE_DUR = time.Minute * 30

type MyCustomClaims struct {
	jwt.RegisteredClaims
	Uid       uint64 `json:"uid"`
	Username  string `json:"username"`
	LoginTime int64  `json:"deadline"`
}

type AccountServiceImpl struct {
	accountDao storage.AccountDao
}

// Register implements AccountService.
func (t *AccountServiceImpl) Register(ctx context.Context, account *entities.Account) error {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return err
	}
	_, err = t.accountDao.Create(ctx, account)
	return commitOrRollback(err, tx)
}

// GenJwtSignKey implements AccountService.
func (t *AccountServiceImpl) GenJwtSignKey(ctx context.Context) string {
	var randNum int64
	if err := binary.Read(rand.Reader, binary.BigEndian, &randNum); err != nil {
		panic(err)
	}
	randNumStr := strconv.FormatInt(randNum, 10)
	// randNum:=
	key := "video1024" + randNumStr
	hash := sha256.New()
	_, err := hash.Write([]byte(key))
	if err != nil {
		panic(err)
	}
	s := hash.Sum(nil)
	return hex.EncodeToString(s)
}

// Login implements AccountService.
func (t *AccountServiceImpl) Login(ctx context.Context, username string, password string) (entities.Account, string, error) {
	var account entities.Account
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return account, "", err
	}
	account,err= t.accountDao.GetAccountByUsername(ctx, username)
	account.Username=username
	err = commitOrRollback(err, tx)
	if err != nil {
		return account, "", errors.WithStack(err)
	}

	hash := sha256.New()
	hash.Write([]byte(username))
	hash.Write([]byte(password))
	hashByte := hash.Sum(nil)
	hashStr := hex.EncodeToString(hashByte)

	if account.Pwd != hashStr {
		return account, "", errors.New("Incorrect password")
	}

	claim := MyCustomClaims{
		Uid:       account.Uid,
		Username:  username,
		LoginTime: time.Now().Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(JWT_EXPIRE_DUR)},
		},
	}

	ctx, conn, err := getRedisConnFromCtx(ctx)
	if err != nil {
		return account, "", err
	}
	defer conn.Close()
	signKey, err := t.accountDao.GetJwtSignedKey(ctx)
	if err != nil {
		if err != redis.Nil {
			return account, "", err
		}
		// key 不存在
		signKey = t.GenJwtSignKey(ctx)
		// 写回缓存
		err = t.accountDao.SetJwtSignedKey(ctx, signKey)
		if err != nil {
			return account, "", err
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenSigned, err := token.SignedString([]byte(signKey))

	if err != nil {
		return account, "", err
	}
	return account, tokenSigned, nil
}

// 验证token是否有效
func (t *AccountServiceImpl) JwtAuth(ctx context.Context, token string) (*MyCustomClaims, error) {
	ctx, conn, err := getRedisConnFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	signKey, err := t.accountDao.GetJwtSignedKey(ctx)
	if err != nil {
		return nil, err
	}
	claim := &MyCustomClaims{}
	_, err = jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return claim, nil
}

func NewAccountServiceImpl(accountDao storage.AccountDao) *AccountServiceImpl {
	return &AccountServiceImpl{
		accountDao: accountDao,
	}
}

var _ AccountService = (*AccountServiceImpl)(nil)
