package models

type Message struct {
	ID        int
	Text      string
	Firstname string
	Lastname  string
	Username  string
}

type Buttons struct {
	ID                int
	MessageRelationID int
}