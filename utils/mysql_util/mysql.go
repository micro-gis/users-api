package mysql_util

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	errors "github.com/micro-gis/utils/rest_errors"
	"strings"
)

const (
	ErrNoRows = "no rows in result set"
)

func ParseError(err error) errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("no record matching given id"))
		}
		return errors.NewInternalServerError(fmt.Sprintf("Error when trying to save user : %s\n ", err.Error()), err)
	}
	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestError(fmt.Sprintf("email already exists"))
	}
	return errors.NewInternalServerError("error processing request", err)
}
