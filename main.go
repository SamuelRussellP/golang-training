package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type transaction struct {
	TransactionId     uuid.UUID `json:"transaction_id"`
	LoginStatus       bool      `json:"login_status"`
	BankId            string    `json:"bank_id"`
	TransactionStatus bool      `json:"transaction_status"`
	TransactionAmount float32   `json:"transaction_amount"`
	Balance           float32   `json:"balance"`
}

type bank struct {
	BankId         string  `json:"bank_id"`
	BankPercentage float32 `json:"bank_percentage"`
}

var transactions []transaction

var banks = []bank{
	{BankId: "BCA", BankPercentage: 0.2},
	{BankId: "BRI", BankPercentage: 0.1},
}

func main() {
	router := gin.Default()
	router.GET("/transactions/:id", getTransaction)
	router.GET("/transactions", getTransactions)
	router.GET("/transactions/bank-percentage/:bank_id", getBankPercentage)
	router.POST("/transactions", addTransactions)
	router.Run("localhost:9090")
}

func getTransactions(context *gin.Context) {
	if len(transactions) == 0 {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "You have no transactions"})
	} else {
		context.IndentedJSON(http.StatusOK, transactions)
	}
}

func getTransaction(context *gin.Context) {
	id := context.Param("id")
	trId, err := uuid.Parse(id)
	transaction, err := getTransactionById(trId)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, transaction)
}

func getTransactionById(trId uuid.UUID) (*transaction, error) {
	for i, a := range transactions {
		if a.TransactionId == trId {
			return &transactions[i], nil
		}
	}
	return nil, errors.New("transaction not found")
}

func getBankPercentage(context *gin.Context) {
	bankId := context.Param("bank_id")
	bank, err := getBankPercentageByBankId(bankId)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "bank not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, bank)
}

func getBankPercentageByBankId(bankId string) (*bank, error) {
	for i, a := range banks {
		if a.BankId == bankId {
			return &banks[i], nil
		}
	}
	return nil, errors.New("bank not found")
}

func addTransactions(context *gin.Context) {
	var newTransaction transaction
	newTransaction.TransactionId = uuid.New()

	err := context.BindJSON(&newTransaction)
	if err != nil {
		return
	}
	transactions = append(transactions, newTransaction)
	context.IndentedJSON(http.StatusOK, newTransaction)
}
