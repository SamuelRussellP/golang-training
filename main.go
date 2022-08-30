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
	TransactionTime   string    `json:"transaction_time"`
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

	router := gin.Default()
	router.GET("/transactions/:id", getTransaction)
	router.GET("/transactions", fetchTransactions)
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

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@Admin123@tcp(localhost:3306)/go-training-payment?parseTime=true")
	if err != nil {
		fmt.Println("Error validating sql.Open arguments")
		panic(err.Error())
	}
	fmt.Println("Successful connection to database")

	err = db.Ping()
	if err != nil {
		fmt.Println("Error validating db.Ping")
		panic(err.Error())
	}
	return db
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

func fetchTransactions(context *gin.Context) {
	db := connect()

	var tempTransaction transaction
	var tempTransactionList []transaction
	queryText := fmt.Sprintf("SELECT * FROM `go-training-payment`.transaction_details")

	rows, err := db.Query(queryText)
	if err != nil {
		fmt.Println("Error validating sql.Query arguments")
		panic(err.Error())
	}
	for rows.Next() {
		if err := rows.Scan(&tempTransaction.TransactionId, &tempTransaction.BankId, &tempTransaction.TransactionStatus, &tempTransaction.TransactionAmount, &tempTransaction.TransactionFee, &tempTransaction.TransactionTime); err != nil {
			fmt.Println("Error validating sql.Query arguments")
			panic(err.Error())
		} else {
			transactions = append(transactions, tempTransaction)
			tempTransactionList = append(tempTransactionList, tempTransaction)
		}
	}
	if len(tempTransactionList) > 0 {
		context.IndentedJSON(http.StatusOK, tempTransactionList)
	} else {
		context.IndentedJSON(http.StatusOK, "No transactions found")
	}
}

func fetchTransactionById(id uuid.UUID) (transaction, error) {
	var tempTransaction transaction
	db := connect()
	queryText := fmt.Sprintf("SELECT * FROM `go-training-payment`.transaction_details WHERE transaction_id = '%v'", id)

	rows, err := db.Query(queryText)
	if err != nil {
		fmt.Println("Error validating sql.Query arguments")
		panic(err.Error())
	}
	for rows.Next() {
		if err := rows.Scan(&tempTransaction.TransactionId, &tempTransaction.BankId, &tempTransaction.TransactionStatus, &tempTransaction.TransactionAmount, &tempTransaction.TransactionFee, &tempTransaction.TransactionTime); err != nil {
			fmt.Println("Error validating sql.Query arguments")
			panic(err.Error())
		} else {
			return tempTransaction, nil
		}
	}
	return tempTransaction, errors.New("transaction not found")
}

func updateTransaction(transaksi transaction, context *gin.Context) {
	db := connect()
	defer db.Close()

	transaksi.TransactionTime = time.Now().Format("2006-01-02 15:04:05")

	queryText := fmt.Sprintf("UPDATE `go-training-payment`.transaction_details SET transaction_status = '%v', transaction_time = '%v' WHERE transaction_id = '%v'",
		transaksi.TransactionStatus,
		transaksi.TransactionTime,
		transaksi.TransactionId,
	)
	_, err := db.Query(queryText)
	if err != nil {
		fmt.Println("Error validating sql.Query arguments")
		panic(err.Error())
	}
	context.IndentedJSON(http.StatusOK, transaksi)
}

func insertTransaction(transaksi transaction) error {
	db := connect()
	defer db.Close()

	queryText := fmt.Sprintf("INSERT INTO `go-training-payment`.transaction_details  (transaction_id, bank_id, transaction_status, transaction_amount, transaction_fee, transaction_time)"+"VALUES ('%v','%v','%v','%v','%v','%v')",
		transaksi.TransactionId,
		transaksi.BankId,
		transaksi.TransactionStatus,
		transaksi.TransactionAmount,
		transaksi.TransactionFee,
		transaksi.TransactionTime,
	)

	_, err := db.Query(queryText)

	if err != nil {
		fmt.Println("Error validating sql.Query arguments")
		panic(err.Error())
	}
	return nil
}

func getTransaction(context *gin.Context) {
	id := context.Param("id")
	trId, err := uuid.Parse(id)
	transaction, err := fetchTransactionById(trId)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, transaction)
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
	transaction, err := fetchTransactionById(trId)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
		return
	}

	if transaction.TransactionStatus == "ready" {
		transaction.TransactionStatus = "paid"
		updateTransaction(transaction, context)

	} else if transaction.TransactionStatus == "cancelled" {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "transaction has been cancelled. Cannot confirm transaction"})
		return
	} else {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "this transaction has been confirmed."})
		return
	}
}

func cancelTransaction(context *gin.Context) {
	id := context.Param("id")
	trId, err := uuid.Parse(id)
	transaction, err := fetchTransactionById(trId)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
		return
	}

	if transaction.TransactionStatus == "ready" {
		transaction.TransactionStatus = "cancelled"
		updateTransaction(transaction, context)

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
	newTransaction.TransactionTime = time.Now().Format("2006-01-02 15:04:05")

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
	insertTransaction(newTransaction)
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
