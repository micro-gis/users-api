package mysql_util

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/micro-gis/users-api/utils/errors_util"
	"strings"
)

const (
	errNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("no record matching giver id"))
		}
		return errors.NewInternalServerError(fmt.Sprintf("Error when trying to save user : %s\n ", err.Error()))
	}
	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestError(fmt.Sprintf("email already exists"))
	}
	return errors.NewInternalServerError("error processing request")
}
