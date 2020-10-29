package users_test

import (
	"github.com/neo4j-examples/golang-neo4j-realworld-example/pkg/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http/httptest"
	"strings"
)

type FakeUserRepository struct {

}

func (FakeUserRepository) RegisterUser(user users.User) error {
	return nil
}

var _ = Describe("Users", func() {

	It("should register", func() {
		handler := users.UserHandler{
			Path:           "/users",
			UserRepository: &FakeUserRepository{},
		}
		testResponseWriter := httptest.NewRecorder()
		requestBody := strings.NewReader(
			"{\"user\":{\"email\":\"user@example.com\", \"password\":\"s3cr3t\", \"username\":\"user\"}}")

		handler.Register(testResponseWriter,
			httptest.NewRequest("POST", "/users", requestBody))

		Expect(testResponseWriter.Code).To(Equal(201))
		responseBody, _ := ioutil.ReadAll(testResponseWriter.Body)
		Expect(string(responseBody)).To(Equal("{\"user\":{\"username\":\"user\",\"email\":\"user@example.com\"}}"))
	})

})
