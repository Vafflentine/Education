package models

type Post struct {
	Id     int    `json:"id" gorm:"primaryKey" gorm:"autoIncrement:true"`
	UserId int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
