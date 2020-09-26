package services

import (
	"github.com/micro-gis/users-api/domain/users"
	"github.com/micro-gis/users-api/utils/crypto_util"
	"github.com/micro-gis/users-api/utils/date_util"
	"github.com/micro-gis/users-api/utils/errors_util"
	"github.com/micro-gis/users-api/utils/string_util"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.DateCreated = date_util.GetNowDBFormat()
	user.Status = users.StatusActive
	user.Password = crypto_util.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUser(userId int64) (*users.User, *errors.RestErr) {
	result := &users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	current, err := GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	if isPartial {
		if string_util.IsEmptyString(user.FirstName) {
			current.FirstName = user.FirstName
		}
		if string_util.IsEmptyString(user.LastName) {
			current.LastName = user.LastName
		}
		if string_util.IsEmptyString(user.Email) {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func DeleteUser(userId int64) *errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func Search(status string) (users.Users, *errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}
