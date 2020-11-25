package users_test

import (
	"bytes"
	"encoding/json"
	"github.com/neo4j-examples/golang-neo4j-realworld-example/pkg/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Users", func() {
	userLoginRequest := users.UserLogin{
		User: users.User{
			Email: "florent@example.org",
			Password: "very-secure",
		},
	}

	FIt("should log in", func() {
		handler := users.UserLoginHandler{
			Path:           "/users/login",
			UserRepository: &FakeUserRepository{},
		}
		testResponseWriter := httptest.NewRecorder()

		handler.Login(
			testResponseWriter,
			httptest.NewRequest(
				"POST",
				"/users/login",
				strings.NewReader(marshalLogin(&userLoginRequest))))

		Expect(testResponseWriter.Code).To(Equal(200))
		login := unmarshalLogin(testResponseWriter.Body)
		Expect(login.Email).To(Equal("florent@example.org"))
		Expect(login.Username).To(Equal("flo"))
		Expect(login.Token).NotTo(Equal(""), "token should be set")
	})
})

func unmarshalLogin(payload *bytes.Buffer) *users.LoggedInUser {
	var result users.LoggedInUser
	err := json.Unmarshal(payload.Bytes(), &result)
	Expect(err).To(BeNil(), "JSON unmarshalling should work")
	return &result
}

func marshalLogin(login *users.UserLogin) string {
	payload, err := json.Marshal(login)
	Expect(err).To(BeNil(), "JSON marshalling should work")
	return string(payload)
}