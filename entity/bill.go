package entity

import (
	"time"
)

type Bill struct {
	Id          int
	Customer_Id int
	Employee_Id int
	BillDate    time.Time
	EntryDate   time.Time
	FinishDate  time.Time
}
