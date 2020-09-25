package services

import (
	"github.com/micro-gis/users-api/domain/users"
	"github.com/micro-gis/users-api/utils/errors"
	"github.com/micro-gis/users-api/utils/string_utils"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
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
		if string_utils.IsEmptyString(user.FirstName) {
			current.FirstName = user.FirstName
		}
		if string_utils.IsEmptyString(user.LastName) {
			current.LastName = user.LastName
		}
		if string_utils.IsEmptyString(user.Email) {
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
