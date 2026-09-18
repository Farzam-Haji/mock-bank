package models

import "time"

type Account struct {
	ID 			int64 	`json:"id"`
	Name 		string 	`json:"name"`
	CardNumber 	string 	`json:"card_number"`
	Password 	string 	`json:"password"`
	Balance 	int64 	`json:"balance"`  //minor units
}

type Payment struct {
	ID  				int64  		`json:"id"`
	MerchantOrderID  	int64  		`json:"merchant_order_id"`
	AccountID  			*int64  	`json:"account_id"`
	Status  			string  	`json:"status"`
	Amount				int64  		`json:"amount"`   //minor units
	CreatedAt  			time.Time  	`json:"created_at"`
	UpdatedAt  			time.Time  	`json:"updated_at"`
}