package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type transaction struct {
	transactionId     string  `json:"transaction_id"`
	loginStatus       bool    `json:"login_status"`
	transactionStatus bool    `json:"transaction_status"`
	transactionAmount float32 `json:"transaction_amount"`
	balance           float32 `json:"balance"`
}

var transactions = []transaction{
	{"1", true, true, 50.0, 100.0}}

func main() {
	router := gin.Default()
	router.GET("/transactions", getTransactions)
	router.POST("/transactions", addTransactions)
	router.Run("localhost:9090")
}

func getTransactions(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, transactions)
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
