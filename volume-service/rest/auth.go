package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (adapter *RESTAdapter) isAdminUser(user string) bool {
	return adapter.config.RestAdminUsername == user
}

func (adapter *RESTAdapter) basicAuthAdmin(ctx *gin.Context) {
	user, password, hasAuth := ctx.Request.BasicAuth()
	if hasAuth && user == adapter.config.RestAdminUsername && password == adapter.config.RestAdminPassword {
		log.WithFields(log.Fields{
			"user": user,
		}).Info("User authenticated")

		ctx.Set(gin.AuthUserKey, user)
		return
	}

	ctx.AbortWithStatus(http.StatusUnauthorized)
	ctx.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
}

func (adapter *RESTAdapter) basicAuthDeviceOrAdmin(ctx *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "basicAuthDeviceOrAdmin",
	})

	logger.Info("received basicAuthDeviceOrAdmin()")

	user, password, hasAuth := ctx.Request.BasicAuth()
	if hasAuth {
		if user == adapter.config.RestAdminUsername && password == adapter.config.RestAdminPassword {
			// admin
			log.WithFields(log.Fields{
				"user": user,
			}).Info("User authenticated")

			ctx.Set(gin.AuthUserKey, user)
			return
		} else {
			device, err := adapter.logic.GetDevice(user)
			if err == nil {
				if device.Password == password {
					// admin
					log.WithFields(log.Fields{
						"user": user,
					}).Info("User authenticated")

					ctx.Set(gin.AuthUserKey, user)
					return
				}
			}
		}
	}

	ctx.AbortWithStatus(http.StatusUnauthorized)
	ctx.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
}
