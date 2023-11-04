/*
 * @Date: 2023-10-27 07:58:50
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-29 08:56:35
 * @FilePath: /1024/server/src/service/account_service_impl.go
 */
package service

import (
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
func (t *AccountServiceImpl) Register(account *entities.Account) error {
	_, err := t.accountDao.Save(account)
	return err
}

// GenJwtSignKey implements AccountService.
func (t *AccountServiceImpl) GenJwtSignKey() string {
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
func (t *AccountServiceImpl) Login(username string, password string) (entities.LoginInfo, error) {
	var loginInfo entities.LoginInfo
	uid, pwd, err := t.accountDao.GetUidAndPwdByUsername(username)
	if err != nil {
		return loginInfo, err
	}
	if pwd != password {
		return loginInfo, errors.New("Incorrect password")
	}

	claim := MyCustomClaims{
		Uid:       uid,
		Username:  username,
		LoginTime: time.Now().Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(JWT_EXPIRE_DUR)},
		},
	}

	signKey, err := t.accountDao.GetJwtSignedKey()
	if err != nil {
		if err != redis.Nil {
			return loginInfo, err
		}
		// key 不存在
		signKey = t.GenJwtSignKey()
		// 写回缓存
		err = t.accountDao.SetJwtSignedKey(signKey)
		if err != nil {
			return loginInfo, err
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenSigned, err := token.SignedString([]byte(signKey))

	if err != nil {
		return loginInfo, err
	}
	loginInfo.Uid = uid
	loginInfo.UserName = username
	loginInfo.Token = tokenSigned
	return loginInfo, nil
}

// 验证token是否有效
func (t *AccountServiceImpl) JwtAuth(token string) error {
	signKey, err := t.accountDao.GetJwtSignedKey()
	if err != nil {
		return err
	}
	parseClaim := MyCustomClaims{}
	_, err = jwt.ParseWithClaims(token, &parseClaim, func(t *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

func NewAccountServiceImpl(accountDao storage.AccountDao) *AccountServiceImpl {
	return &AccountServiceImpl{
		accountDao: accountDao,
	}
}

var _ AccountService = (*AccountServiceImpl)(nil)
