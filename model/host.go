package model

type Host struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
