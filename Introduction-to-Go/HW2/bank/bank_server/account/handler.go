package account

import (
	"bank_server/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var accounts = make(map[string]models.Account)

func CreateAccount(c *gin.Context) {
	var account models.Account
	if err := c.BindJSON(&account); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	accounts[account.ID] = account
	c.JSON(http.StatusCreated, account)
}

func GetAccount(c *gin.Context) {
	id := c.Param("id")
	account, exists := accounts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	} else { // здесь можно убрать else
		c.JSON(http.StatusOK, account)
	}
}

func UpdateAccountName(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	account, exists := accounts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	account.Name = request.Name
	c.JSON(http.StatusOK, account)
}

func UpdateAccountBalance(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Balance float64 `json:"balance"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	account, exists := accounts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	account.Balance = request.Balance
	c.JSON(http.StatusOK, account)
}

func DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	_, exists := accounts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}
	delete(accounts, id)
	c.JSON(http.StatusNoContent, nil)
}
