/*
 * @Date: 2023-10-28 11:45:15
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-01 11:28:08
 * @FilePath: /1024/server/src/service/common.go
 */
package service

import (
	"context"
	"database/sql"
	"dev1024/src/storage/object"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func commitOrRollback(err error, tx *sql.Tx) error {
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			err = errors.Wrap(err, txErr.Error())
		}
	} else {
		if txErr := tx.Commit(); txErr != nil {
			err = errors.WithStack(txErr)
		}
	}
	return err
}

func getTxFromCtx(ctx context.Context) (context.Context, *sql.Tx, error) {
	txv := ctx.Value("tx")
	var tx *sql.Tx
	var err error
	if txv == nil {
		dbv := ctx.Value("db")
		if dbv == nil {
			return ctx, nil, errors.New("db==nil")
		}
		db := dbv.(*sql.DB)
		tx, err = db.Begin()
		if err != nil {
			return ctx, nil, errors.WithStack(err)
		}
		ctx = context.WithValue(ctx, "tx", tx)
	}
	return ctx, tx, nil
}

func getRedisConnFromCtx(ctx context.Context) (context.Context, *redis.Conn, error) {
	connValue := ctx.Value("redis.conn")
	var conn *redis.Conn
	if connValue == nil {
		redisValue := ctx.Value("redis")
		if redisValue == nil {
			return ctx, nil, errors.New("redis==nil")
		}
		redisClient := redisValue.(*redis.Client)
		conn = redisClient.Conn()
		ctx = context.WithValue(ctx, "redis.conn", conn)
	}
	return ctx, conn, nil
}

func getObjectStorageFromCtx(ctx context.Context) (context.Context, object.ObjectStorage, error) {
	objectStorageValue := ctx.Value("object_storage")
	var objectStorage object.ObjectStorage
	if objectStorageValue == nil {
		panic("objectStorageValue == nil")
	}
	objectStorage=objectStorageValue.(object.ObjectStorage)
	return ctx, objectStorage, nil
}
