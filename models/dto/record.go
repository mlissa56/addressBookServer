package dto

type Record struct {
	ID         int64  `json:"-" sql.field:"id"`
	Name       string `json:"name" sql.field:"name"`
	LastName   string `json:"last_name" sql.field:"last_name"`
	MiddleName string `json:"middle_name" sql.field:"middle_name"`
	Address    string `json:"address" sql.field:"address"`
	Phone      string `json:"phone" sql.field:"phone"`
}
