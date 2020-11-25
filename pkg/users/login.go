package users

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type UserLogin struct {
	User User `json:"user"`
}

type LoggedInUser struct {
	Username string
	Email string
	Token string
}

type UserLoginHandler struct {
	Path           string
	UserRepository UserRepository
}

func (u *UserLoginHandler) Login(writer http.ResponseWriter, request *http.Request) {
	requestBody, _ := ioutil.ReadAll(request.Body)
	userLoginRequest := UserLogin{}
	_ = json.Unmarshal(requestBody, &userLoginRequest)
	requestUser := userLoginRequest.User
	user, _ := u.UserRepository.FindByEmailAndPassword(
		requestUser.Email,
		requestUser.Password)

	var token string // ???
	loggedInUser := &LoggedInUser{
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}
	writer.WriteHeader(200)
	bytes, _ := json.Marshal(&loggedInUser)
	_, _ = writer.Write(bytes)
}
