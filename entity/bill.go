package entity

import (
	"time"
)

type Bill struct {
	Id          int       `json:"id"`
	BillDate    time.Time `json:"billDate"`
	EntryDate   time.Time `json:"entryDate"`
	FinishDate  time.Time `json:"finishDate"`
	Employee_Id int       `json:"employeeId"`
	Customer_Id int       `json:"customerId"`
}
