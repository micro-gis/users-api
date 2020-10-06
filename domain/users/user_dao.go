package users

import (
	"fmt"
	"github.com/micro-gis/users-api/datasources/mysql/users_db"
	"github.com/micro-gis/users-api/utils/mysql_util"
	"github.com/micro-gis/utils/logger"
	errors "github.com/micro-gis/utils/rest_errors"
	"strings"
)

const (
	queryInsertUser             = "INSERT INTO users (first_name, last_name, email, date_created, status, password) VALUES (?,?,?,?,?,?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?;"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("error when trying to get user by id", mysql_util.ParseError(err))
		return errors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) FindByEmailAndPassword() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare find user by email and password", err)
		return errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		if strings.Contains(err.Error(), mysql_util.ErrNoRows) {
			return errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by id", mysql_util.ParseError(err))
		return errors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare insert user statement", err)
		return errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if err != nil {
		logger.Error("error when trying to insert new user", mysql_util.ParseError(err))
		return errors.NewInternalServerError("database error", err)
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to retrieve last id", mysql_util.ParseError(err))
		return errors.NewInternalServerError("database error", err)
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)

	if err != nil {
		logger.Error("error when trying to update user", mysql_util.ParseError(err))
		return errors.NewInternalServerError("database error", err)
	}
	return nil

}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", mysql_util.ParseError(err))
		return errors.NewInternalServerError("database error", err)
	}
	return nil
}

func (user *User) FindByStatus(status string) (Users, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare Find user statement", err)
		return nil, errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find user", mysql_util.ParseError(err))
		return nil, errors.NewInternalServerError("database error", err)
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to create find user response ", mysql_util.ParseError(err))
			return nil, errors.NewInternalServerError("database error", err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matchin status %s", status))
	}
	return results, nil

}
