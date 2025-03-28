package models

type Student struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Sex    bool   `json:"sex"`
	CardId int    `json:"card_id"`
}


type CardCredit struct {
	Id         int `json:"id"`
	StudentId  int `json:"student_id"`
	CardNumber int `json:"card_number"`
	Expiration int `json:"expiration"`
	Cvv        int `json:"cvv"`
}