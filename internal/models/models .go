package models

import "time"

type TotalReport struct {
	TotalSpend int64
	SpendsSet []*Report
}

type Report struct {
	Category string
	Amount int64
}

type Buttons struct {
	ID                int
	MessageID         int
	Amount            int
	Firstname         string
	Lastname          string
	Username          string
}

type Transaction struct {
	ButtonID int
	Amount   int
	Category string
	Time     time.Time
}
