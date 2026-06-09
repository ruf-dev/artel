package couchdb

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

var ErrUnknownError = rerrors.New("unknown error")

const (
	ErrorCodeIllegalDatabaseName = "illegal_database_name"
	ErrorCodeAlreadyExists       = "file_exists"
	ErrorCodeForbidden           = "forbidden"
)

type forbiddenResp struct {
	Reason string `json:"reason"`
}

func checkAccountLocked(statusCode int, body []byte) error {
	if statusCode != http.StatusForbidden {
		return nil
	}
	var resp forbiddenResp
	err := json.Unmarshal(body, &resp)
	if err != nil || !strings.Contains(resp.Reason, "locked") {
		return nil
	}
	return rerrors.Wrap(user_errors.CouchDbAccountLocked)
}
