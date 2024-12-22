package models

type Expence struct {
    Id int	`json:"id"`
    Description string `json:"desc"`
    Amount float64 `json:"amount"`
    CreatedAt int64 `json:"created_at"`
}
