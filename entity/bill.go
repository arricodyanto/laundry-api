package entity

import (
	"time"
)

type Bill struct {
	Id          string    `json:"id"`
	BillDate    time.Time `json:"billDate"`
	EntryDate   time.Time `json:"entryDate"`
	FinishDate  time.Time `json:"finishDate"`
	Employee_Id string    `json:"employeeId"`
	Customer_Id string    `json:"customerId"`
}
