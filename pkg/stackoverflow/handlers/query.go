package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type UserRegistration struct {
	User User `json:"user"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type QueryHandler struct {
	Path                    string
	StackOverflowRepository StackOverflowRepository
}

func (u *StackOverflowRepository) Register(writer http.ResponseWriter, request *http.Request) {
	requestBody, _ := ioutil.ReadAll(request.Body)
	userRegistrationRequest := UserRegistration{}
	_ = json.Unmarshal(requestBody, &userRegistrationRequest)
	requestUser := userRegistrationRequest.User
	_ = u.UserRepository.RegisterUser(&requestUser)

	writer.WriteHeader(201)
	writer.Header().Add("Content-Type", "application/json")
	userRegistrationResponse := UserRegistration{
		User: User{
			Username: requestUser.Username,
			Email:    requestUser.Email,
		}}
	bytes, _ := json.Marshal(&userRegistrationResponse)
	_, _ = writer.Write(bytes)
}
