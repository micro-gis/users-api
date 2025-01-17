package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro-gis/oauth-go/oauth"
	"github.com/micro-gis/users-api/domain/users"
	"github.com/micro-gis/users-api/services"
	"github.com/micro-gis/users-api/utils/authenticate_utils"
	errors "github.com/micro-gis/utils/rest_errors"
	"net/http"
	"strconv"
)

func getUserId(userIdParam string) (int64, errors.RestErr) {
	// Check for userid passed as number
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func Create(c *gin.Context) {
	if err := authenticate_utils.AuthenticateRequest(c, true, 0); err != nil {
		c.JSON(err.Status(), err)
		return
	}
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(fmt.Sprintf("invalid json : %s", err.Error()))
		c.JSON(restErr.Status(), restErr)
		return
	}
	result, saveErr := services.UserService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status(), saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))

}

func Get(c *gin.Context) {
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	if err := authenticate_utils.AuthenticateRequest(c, false, 0); err != nil {
		c.JSON(err.Status(), err)
		return
	}



	user, getErr := services.UserService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Update(c *gin.Context) {
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	if err := authenticate_utils.AuthenticateRequest(c, true, userId); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("Invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	user.Id = userId

	isPartial := c.Request.Method == http.MethodPatch
	result, resterr := services.UserService.UpdateUser(isPartial, user)
	if resterr != nil {
		c.JSON(resterr.Status(), resterr)
		return
	}


	c.JSON(http.StatusOK, result.Marshall(oauth.IsPublic(c.Request)))
}

func Delete(c *gin.Context) {
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	if err := authenticate_utils.AuthenticateRequest(c, true, userId); err != nil {
		c.JSON(err.Status(), err)
		return
	}
	if err := services.UserService.DeleteUser(userId); err != nil {
		c.JSON(err.Status(), err)
		return
	}
	err = oauth.DeleteAllAccessToken(c.Request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	if err := authenticate_utils.AuthenticateRequest(c, false, 0); err != nil {
		c.JSON(err.Status(), err)
		return
	}
	status := c.Query("status")
	users, err := services.UserService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, users.Marshall(oauth.IsPublic(c.Request)))
}

func Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}
	user, err := services.UserService.LoginUser(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}
