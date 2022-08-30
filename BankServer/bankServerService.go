package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BankTransactionStatus struct {
	Id      int    `json:"BankTransactionStatus_id"`
	Message string `json:"message"`
}

var responseFromBank BankTransactionStatus

func main() {
	router := gin.Default()
	router.POST("/bankservicestatus/:response_val", addBankTransaction)
	router.GET("/bankservicestatus", getLastTransactionStatus)
	router.Run("localhost:9091")
}

func generateBankMessage(response int) {
	switch response {
	case 200:
		responseFromBank.Message = "Transaction Success!"
	case 503:
		responseFromBank.Message = "Service Unavailable"
	default:
		responseFromBank.Message = "Response is unknown"
	}

}

func addBankTransaction(context *gin.Context) {
	response := context.Param("response_val")
	responseFromBankIdInt, err := strconv.Atoi(response)
	if err != nil {
		panic(err)
	} else {
		responseFromBank.Id = responseFromBankIdInt
		generateBankMessage(responseFromBank.Id)
		context.IndentedJSON(http.StatusOK, responseFromBank)
	}
}

func getLastTransactionStatus(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, responseFromBank)
}
