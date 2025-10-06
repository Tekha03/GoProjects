package router

import (
	"bank_server/account"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/account/:id", account.GetAccount)
	r.POST("/account", account.CreateAccount)
	r.PUT("/account/:id/name", account.CreateAccount)
	r.PUT("/account/:id/balance", account.UpdateAccountBalance)
	r.DELETE("/account/:id", account.DeleteAccount)

	return r
}
