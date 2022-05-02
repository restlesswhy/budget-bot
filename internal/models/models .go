package models

import "time"

// type Message struct {
// 	ID        int
// 	Text      string
// 	Firstname string
// 	Lastname  string
// 	Username  string
// }

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
