package models

type Student struct {
	Name   string `json:"name"`
	Id     int    `json:"id"`
	Age    int    `json:"age"`
	CardId int    `json:"card_id"`
	Sex    bool   `json:"sex"`
}

type CardCredit struct {
	Id         int `json:"id"`
	StudentId  int `json:"student_id"`
	CardNumber int `json:"card_number"`
	Expiration int `json:"expiration"`
	Cvv        int `json:"cvv"`
}
