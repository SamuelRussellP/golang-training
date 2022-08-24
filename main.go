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

var transactions = []transaction{
	{TransactionId: "1", LoginStatus: true, TransactionStatus: true, TransactionAmount: 50.0, Balance: 100.0}}

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
