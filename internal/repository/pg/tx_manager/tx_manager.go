package tx_manager

import (
	"database/sql"
	stdErrors "errors"

	errors "go.redsock.ru/rerrors"
)

type TxManager interface {
	Execute(do func(tx *sql.Tx) error) error
}

type txManager struct {
	conn *sql.DB
}

func New(conn *sql.DB) TxManager {
	return &txManager{
		conn: conn,
	}
}

func (m *txManager) Execute(do func(tx *sql.Tx) error) error {
	tx, err := m.conn.Begin()
	if err != nil {
		return errors.Wrap(err)
	}

	err = do(tx)

	var commitErr error
	if err != nil {
		commitErr = tx.Rollback()
	} else {
		commitErr = tx.Commit()
	}

	if commitErr != nil || err != nil {
		return stdErrors.Join(err, commitErr)
	}

	return nil
}
