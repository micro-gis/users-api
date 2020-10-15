package authenticate_utils

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-gis/oauth-go/oauth"
	errors "github.com/micro-gis/utils/rest_errors"
	"net/http"
)

func AuthenticateRequest(c *gin.Context, forceAuth bool, forceSameUserId int64) errors.RestErr {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		return err
	}
	// For forcing authentication
	if forceAuth {
		if callerId := oauth.GetCallerId(c.Request); callerId == 0 {
			err := errors.NewRestError("Authentication required", http.StatusUnauthorized, "unauthorized", nil)
			return err
		}
	}

	if forceSameUserId != 0 {
		if callerId := oauth.GetCallerId(c.Request); callerId != forceSameUserId {
			err := errors.NewRestError("Authentication required", http.StatusUnauthorized, "unauthorized", nil)
			return err
		}
	}
	return nil
}

