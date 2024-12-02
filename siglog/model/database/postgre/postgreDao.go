package postgre

import "qqweq/siglog/model/models"

type PostgreDao struct {}

func (*PostgreDao) CreateUser(user *models.User) (string, error) {
    // STUB
    return "STUB", nil
}

func (*PostgreDao) ReadUserByUsername(username string) (*models.User, error) {
    // STUB
    return nil, nil
}
func (*PostgreDao) DeleteUser(user *models.User) error {
    // STUB
    return nil 
} 

func (*PostgreDao) CreateSession(username string) (string, error) {
    // STUB
    return "STUB", nil 
}
func (*PostgreDao) DeleteSession(sessionId string) error {
    // STUB
    return nil 
}
func (*PostgreDao) UsernameFromSessionId(sessionId string) (string, error) {
    // STUB
    return "STUB", nil
}

func connectToDb() { 

}
