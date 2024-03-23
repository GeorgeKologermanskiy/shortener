package main;

type UserInfo struct {
	userID int
	email string
	login string
	password string
}

type UserStorage struct {
	users map[string]UserInfo
	nextUserID int
}

func NewUserStorage() (us *UserStorage) {
	us = new(UserStorage)
	us.users = make(map[string]UserInfo)
	us.nextUserID = 0
	return
}

func (us UserStorage) addUser(email string, login string, password string) bool {
	if _, ok := us.users[login]; ok {
		return true
	}

	us.users[login] = UserInfo {
		userID: us.nextUserID,
		email: email,
		login: login,
		password: password,
	}
	us.nextUserID++
	return false
}

func (us UserStorage) getUser(login string) (UserInfo, bool) {
	val, ok := us.users[login]
	return val, ok
}
