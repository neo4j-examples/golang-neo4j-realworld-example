package users

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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

	token, _ := CreateToken(user)
	loggedInUser := &LoggedInUser{
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}
	writer.WriteHeader(200)
	bytes, _ := json.Marshal(&loggedInUser)
	_, _ = writer.Write(bytes)
}

func CreateToken(user *User) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_ACCESS")))
}
