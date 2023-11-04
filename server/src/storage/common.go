/*
 * @Date: 2023-10-28 11:45:15
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-28 11:48:32
 * @FilePath: /1024/server/src/storage/common.go
 */
package storage

import (
	"database/sql"

	"github.com/pkg/errors"
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
