package model

type Account struct {
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	LoginStatus bool   `json:"login_status"`
}
