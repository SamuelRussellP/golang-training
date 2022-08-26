package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

type transaction struct {
	TransactionId     uuid.UUID `json:"transaction_id"`
	BankId            string    `json:"bank_id"`
	TransactionStatus bool      `json:"transaction_status"`
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

var transactions []transaction

var banks = []bank{
	{BankId: "BCA", BankPercentage: 0.2},
	{BankId: "BRI", BankPercentage: 0.1},
}

var AccountSession account

func main() {
	router := gin.Default()
	router.GET("/transactions/:id", getTransaction)
	router.GET("/transactions", getTransactions)
	router.GET("/banks/", getBanks)
	router.GET("/banks/:bank_id", getBankPercentage)
	router.GET("/banks/fee/:bank_id/:transactionAmount", getTransactionFeeByBank)
	router.GET("/account", getAccount)
	router.POST("/transactions", addTransactions)
	router.POST("/banks", addBanks)
	router.POST("/account", addAccount)
	//router.PATCH("/account", toggleLogInStatus)
	router.Run("localhost:9090")
}

func isLoggedIn(context *gin.Context) bool {
	if len(AccountSession.AccountId) > 0 {
		return true
	}
	return false
}

//func toggleLogInStatus(context *gin.Context) {
//	err := context.BindJSON(&AccountSession)
//	if err != nil {
//		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "account not found"})
//		return
//	}
//	AccountSession.LoginStatus = !AccountSession.LoginStatus
//}

func addAccount(context *gin.Context) {
	var newAccount account
	//newAccount.AccountId = "1"
	//newAccount.AccountName = "Samuel"
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
