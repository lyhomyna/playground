package models

type User struct {
    Username  string `json:"username" bson:"username"`
    Password  string `json:"password" bson:"password"`
    Firstname string `json:"firstname" bson:"firstname"`
    Lastname  string `json:"lastname" bson:"lastname"`
    Role      string `json:"role" bson:"role"`
}
