/*
 * @Date: 2023-10-28 11:45:15
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-01 18:09:15
 * @FilePath: /1024/server/src/storage/common.go
 */
package storage

import (
	"context"
	"database/sql"
	"dev1024/src/storage/object"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func getTxFromCtx(ctx context.Context) (context.Context, *sql.Tx, error) {
	txv := ctx.Value("tx")
	var tx *sql.Tx
	if txv == nil {
		panic("txv == nil ")
	}
	tx = txv.(*sql.Tx)
	return ctx, tx, nil
}

func getRedisConnFromCtx(ctx context.Context) (context.Context, *redis.Conn, error) {
	connValue := ctx.Value("redis.conn")
	var conn *redis.Conn
	if connValue == nil {
		panic("connValue == nil")
	}
	conn = connValue.(*redis.Conn)
	return ctx, conn, nil
}

func getObjectStorageFromCtx(ctx context.Context) (context.Context, object.ObjectStorage, error) {
	objectStorageValue := ctx.Value("object_storage")
	var objectStorage object.ObjectStorage
	if objectStorageValue == nil {
		panic("objectStorageValue == nil")
	}
	objectStorage = objectStorageValue.(object.ObjectStorage)
	return ctx, objectStorage, nil
}

func SelectTableBy(tx *sql.Tx, table string, keys []string, values []any, by string, byValue any) error {
	if len(keys) == 0 {
		return errors.New("len(keys)==0")
	}

	keyBuilder := strings.Builder{}
	keyBuilder.WriteString(keys[0])
	for i := 1; i < len(keys); i++ {
		keyBuilder.WriteString(",")
		keyBuilder.WriteString(keys[i])
	}

	query := fmt.Sprintf("SELECT %v from %v WHERE %v=?", keyBuilder.String(), table, by)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows := stmt.QueryRow(byValue)
	err = rows.Scan(values...)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTableBy(tx *sql.Tx, table string, keys []string, values []any, by string, byValue any) (error) {
	sets := strings.Builder{}
	sets.WriteString(keys[0] + "=?")
	for i := 1; i < len(keys); i++ {
		sets.WriteString(",")
		sets.WriteString(keys[i] + "=?")
	}

	sql := fmt.Sprintf("UPDATE %v SET %v WHERE %v=?", table, sets.String(), by)
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	values = append(values, byValue)
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	return nil
}

func InsertTableWith(tx *sql.Tx, table string, keys []string, values []any) (int64,error) {
	// 构造field
	sets := strings.Builder{}
	sets.WriteString(keys[0])
	for i := 1; i < len(keys); i++ {
		sets.WriteString(",")
		sets.WriteString(keys[i])
	}
	// 构造问号
	qms := strings.Builder{}
	qms.WriteString("?")
	for i := 1; i < len(keys); i++ {
		qms.WriteString(",")
		qms.WriteString("?")
	}

	sql := fmt.Sprintf("INSERT INTO %v(%v) VALUES(%v)", table, sets.String(),qms.String())
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return 0,err
	}
	defer stmt.Close()

	res, err := stmt.Exec(values...)
	if err != nil {
		return 0,err
	}

	return res.LastInsertId()
}
