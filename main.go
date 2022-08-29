package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

type transaction struct {
	TransactionId     uuid.UUID `json:"transaction_id"`
	BankId            string    `json:"bank_id"`
	TransactionStatus string    `json:"transaction_status"`
	TransactionAmount float32   `json:"transaction_amount"`
	TransactionFee    float32   `json:"transaction_fee"`
	TransactionTime   time.Time `json:"transaction_time"`
}

type bank struct {
	BankId         string  `json:"bank_id"`
	BankPercentage float32 `json:"bank_percentage"`
}

type account struct {
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	LoginStatus bool   `json:"login_status"`
}

var banks = []bank{
	{BankId: "BCA", BankPercentage: 0.2},
	{BankId: "BRI", BankPercentage: 0.1},
}

var transactions []transaction
var AccountSession account

func main() {

	db, err := sql.Open("mysql", "root:@Admin123@tcp(localhost:3306)/go-training-payment")
	if err != nil {
		fmt.Println("Error validating sql.Open arguments")
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("Connected to database")

	err = db.Ping()
	if err != nil {
		fmt.Println("Error validating db.Ping")
		panic(err.Error())
	}

	router := gin.Default()
	router.GET("/transactions/:id", getTransaction)
	router.GET("/transactions", getTransactions)
	router.GET("/banks/", getBanks)
	router.GET("/banks/:bank_id", getBankPercentage)
	router.GET("/banks/fee/:bank_id/:transactionAmount", getTransactionFeeByBank)
	router.GET("/account", getAccount)
	router.POST("/transactions", addTransactions)
	router.POST("/transactions/confirmTransaction/:id", confirmTransaction)
	router.POST("/transactions/cancelTransaction/:id", cancelTransaction)
	router.POST("/banks", addBanks)
	router.POST("/account", addAccount)
	router.Run("localhost:9090")
}

func isLoggedIn(context *gin.Context) bool {
	if len(AccountSession.AccountId) > 0 {
		return true
	}
	return false
}

func addAccount(context *gin.Context) {
	var newAccount account
	newAccount.LoginStatus = true

	err := context.BindJSON(&newAccount)
	if err != nil {
		return
	}
	AccountSession = newAccount
	context.IndentedJSON(http.StatusOK, AccountSession)
}

func getAccount(context *gin.Context) {
	if isLoggedIn(context) {
		context.IndentedJSON(http.StatusOK, AccountSession)
	} else {
		context.IndentedJSON(http.StatusOK, "Login Required")
	}
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

func confirmTransaction(context *gin.Context) {
	id := context.Param("id")
	trId, err := uuid.Parse(id)
	transaction, err := getTransactionById(trId)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
		return
	}

	if transaction.TransactionStatus == "ready" {
		transaction.TransactionStatus = "paid"
		context.IndentedJSON(http.StatusOK, transaction)

	} else if transaction.TransactionStatus == "cancelled" {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "transaction has been cancelled. Cannot confirm transaction"})
		return
	} else {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "this transaction has been confirmed"})
		return
	}
}

func cancelTransaction(context *gin.Context) {
	id := context.Param("id")
	trId, err := uuid.Parse(id)
	transaction, err := getTransactionById(trId)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
		return
	}

	if transaction.TransactionStatus == "ready" {
		transaction.TransactionStatus = "cancelled"
		context.IndentedJSON(http.StatusOK, transaction)

	} else if transaction.TransactionStatus == "paid" {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "transaction has been paid. Cannot cancel transaction"})
		return
	} else {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "this transaction has been cancelled"})
		return
	}
}

func addTransactions(context *gin.Context) {
	if !isLoggedIn(context) {
		context.IndentedJSON(http.StatusOK, "Login Required")
		return
	}
	var newTransaction transaction
	newTransaction.TransactionId = uuid.New()
	newTransaction.TransactionTime = time.Now()

	err := context.BindJSON(&newTransaction)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, "Please fill in required fields")
		return
	}

	newTransactionBankId := newTransaction.BankId
	newTransactionPercentage, _ := getBankPercentageByBankId(newTransactionBankId)
	if newTransactionPercentage == nil {
		context.IndentedJSON(http.StatusOK, "The bank is not available")
		return
	}
	newTransaction.TransactionFee = newTransaction.TransactionAmount * newTransactionPercentage.BankPercentage
	transactions = append(transactions, newTransaction)
	context.IndentedJSON(http.StatusOK, newTransaction)
}

func addBanks(context *gin.Context) {
	var newBank bank

	err := context.BindJSON(&newBank)
	if err != nil {
		return
	}
	for _, a := range banks {
		if a.BankId == newBank.BankId {
			context.IndentedJSON(http.StatusConflict, gin.H{"message": "bank already exists"})
			return
		}
	}
	banks = append(banks, newBank)
	context.IndentedJSON(http.StatusOK, newBank)
}

func getBanks(context *gin.Context) {
	if len(banks) == 0 {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "You have not assigned any banks"})
	} else {
		context.IndentedJSON(http.StatusOK, banks)
	}
}

func getTransactionFeeByBank(context *gin.Context) {
	bankID := context.Param("bank_id")
	transactionAmount := context.Param("transactionAmount")
	value, _ := strconv.ParseFloat(transactionAmount, 32)
	transactionAmountFloat := float32(value)
	bank, err := getBankPercentageByBankId(bankID)
	if err == nil {
		percentage := bank.BankPercentage
		temp := percentage * transactionAmountFloat
		context.IndentedJSON(http.StatusOK, temp)
		return
	} else {
		context.IndentedJSON(http.StatusOK, "bank not found")
	}
}
