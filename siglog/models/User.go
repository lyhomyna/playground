package models

type User struct {
    Username  string `json:"username"`
    Password  string `json:"password"`
    Firstname string `json:"firstname"`
    Lastname  string `json:"lastname"`
    Role      string `json:"role"`
}
