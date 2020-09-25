package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro-gis/users-api/domain/users"
	"github.com/micro-gis/users-api/services"
	"github.com/micro-gis/users-api/utils/errors"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(fmt.Sprintf("invalid json : %s", err.Error()))
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)

}

func Get(c *gin.Context) {

	// Check for userid passed as number
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status, err)
		return
	}

	user, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user)
}

func Update(c *gin.Context) {
	// Check for userid passed as number
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status, err)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("Invalid json body")
		c.JSON(restErr.Status, restErr)
	}

	user.Id = userId

	isPartial := c.Request.Method == http.MethodPatch
	result, resterr := services.UpdateUser(isPartial, user)
	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)
}
