package entity

type Account struct {
	Username string
	Password string
}

func NewAccount(username, password string) *Account {
	return &Account{
		Username: username,
		Password: password,
	}
}
