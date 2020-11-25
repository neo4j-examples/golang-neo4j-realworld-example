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

type FakeUserRepository struct {
}

func (FakeUserRepository) RegisterUser(user *users.User) error {
	return nil
}

func (FakeUserRepository) FindByEmailAndPassword(email, password string) (*users.User, error) {
	return &users.User{
		Username: "flo",
		Email:    email,
	}, nil
}

var _ = Describe("Users", func() {

	var userRequest = users.UserRegistration{
		User: users.User{
			Username: "user",
			Email:    "user@example.com",
			Password: "s3cr3t",
		}}

	var expectedUserResponse = users.UserRegistration{
		User: users.User{
			Username: "user",
			Email:    "user@example.com",
		}}

	It("should register", func() {
		handler := users.UserRegistrationHandler{
			Path:           "/users",
			UserRepository: &FakeUserRepository{},
		}
		testResponseWriter := httptest.NewRecorder()

		handler.Register(
			testResponseWriter,
			httptest.NewRequest("POST", "/users", strings.NewReader(marshalRegistration(userRequest))))

		Expect(testResponseWriter.Code).To(Equal(201))
		Expect(unmarshalRegistration(testResponseWriter.Body)).To(Equal(expectedUserResponse))
	})

})

func marshalRegistration(registration users.UserRegistration) string {
	payload, err := json.Marshal(registration)
	Expect(err).To(BeNil(), "JSON marshalling should work")
	return string(payload)
}

func unmarshalRegistration(payload *bytes.Buffer) *users.UserRegistration {
	var result users.UserRegistration
	err := json.Unmarshal(payload.Bytes(), &result)
	Expect(err).To(BeNil(), "JSON unmarshalling should work")
	return &result
}
