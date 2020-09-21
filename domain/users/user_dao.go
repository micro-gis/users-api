package users

import (
	"fmt"
	"github.com/micro-gis/users-api/datasources/mysql/users_db"
	"github.com/micro-gis/users-api/utils/date"
	"github.com/micro-gis/users-api/utils/errors"
	"strings"
)

const (
	indexUniqueEmail = "users_email_uindex"
	errNoRows        = "no rows in result set"
	queryInsertUser  = "INSERT INTO users (first_name, last_name, email, date_created) VALUES (?,?,?,?);"
	queryGetUser     = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id = ?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		if strings.Contains(err.Error(), errNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to get user : %d : %s", user.Id, err.Error()))
	}

	if connErr := users_db.Client.Ping(); connErr != nil {
		panic(connErr)
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date.GetNowString()
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("Error when trying to save user : %s\n ", err.Error()))
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error trying to retrieve last user id : %s\n", err.Error()))
	}
	user.Id = userId
	return nil
}
