package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro-gis/users-api/domain/users"
	"github.com/micro-gis/users-api/services"
	"github.com/micro-gis/users-api/utils/errors"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(fmt.Sprintf("invalid json : %s", err.Error()))
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		//TODO : handle user creation error
		return
	}
	c.JSON(http.StatusCreated, result)

}

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implement Me !")
}
