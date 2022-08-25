package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type transaction struct {
	TransactionId     string  `json:"transaction_id"`
	LoginStatus       bool    `json:"login_status"`
	TransactionStatus bool    `json:"transaction_status"`
	TransactionAmount float32 `json:"transaction_amount"`
	Balance           float32 `json:"balance"`
}

type bank struct {
	BankName       string  `json:"bank_name"`
	BankPercentage float32 `json:"bank_percentage"`
}

var transactions []transaction

func main() {
	router := gin.Default()
	router.GET("/transactions/:id", getTransactionById)
	router.GET("/transactions", getTransactions)
	router.POST("/transactions", addTransactions)
	router.Run("localhost:9090")
}

func getTransactions(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, transactions)
}

func getTransactionById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range transactions {
		if a.TransactionId == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Transcation not found"})
}

func addTransactions(context *gin.Context) {
	var newTransaction transaction
	err := context.BindJSON(&newTransaction)
	if err != nil {
		return
	}
	transactions = append(transactions, newTransaction)
	context.IndentedJSON(http.StatusOK, newTransaction)
}
