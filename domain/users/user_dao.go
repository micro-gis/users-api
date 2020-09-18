package users

import (
	"fmt"
	"github.com/micro-gis/users-api/datasources/mysql/users_db"
	"github.com/micro-gis/users-api/utils/date"
	"github.com/micro-gis/users-api/utils/errors"
)

const (
	queryInsertUser = (
		"INSERT INTO users(first_name, last_name, email, date_created VALUES(?, ?, ?, ?);"
)
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {

	if connErr := users_db.Client.Ping(); connErr != nil {
		panic(connErr)
	}

	result := userDB[user.Id]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
	}

	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	//TODO : execute stmt
	
	current := userDB[user.Id]
	if current != nil {
		if current.Email == user.Email {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already registred", current.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user %d already exist", current.Id))
	}
	user.DateCreated = date.GetNowString()
	userDB[user.Id] = user
	return nil
}
