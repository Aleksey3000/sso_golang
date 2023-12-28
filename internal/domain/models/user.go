package models

type User struct {
	Id           int64
	AppId        int32
	Login        string
	PasswordHash []byte
}
