package users_test

import (
	"context"
	"fmt"
	"github.com/neo4j-examples/golang-neo4j-realworld-example/pkg/users"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
	"io"
)

var _ = Describe("User repository", func() {

	const username = "neo4j"
	const password = "s3cr3t"

	var ctx context.Context
	var neo4jContainer testcontainers.Container

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		neo4jContainer, err = startContainer(ctx, username, password)
		Expect(err).To(BeNil(), "Container should start")
	})

	AfterEach(func() {
		Expect(neo4jContainer.Terminate(ctx)).To(BeNil(), "Container should stop")
	})

	It("registers users", func() {
		port, err := neo4jContainer.MappedPort(ctx, "7687")
		Expect(err).To(BeNil(), "Port should be resolved")
		address := fmt.Sprintf("bolt://localhost:%d", port.Int())
		driver, err := neo4j.NewDriver(address, neo4j.BasicAuth(username, password, ""))
		Expect(err).To(BeNil(), "Driver should be created")
		defer Close(driver, "Driver")
		repository := &users.UserNeo4jRepository{
			Driver: driver,
		}

		username := "some-user"
		email := "some-user@example.com"
		initialPassword := "some-password"
		err = repository.RegisterUser(&users.User{
			Username: username,
			Email:    email,
			Password: initialPassword,
		})
		Expect(err).To(BeNil(), "User should be registered")

		session, _ := driver.NewSession(neo4j.SessionConfig{})
		defer Close(session, "Session")
		result, err := session.
			ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
				res, err := tx.Run("MATCH (u:User {username: $username, email: $email}) "+
					"RETURN u.username AS username, u.email AS email, u.password AS password",
					map[string]interface{}{
						"username": username,
						"email":    email,
					})
				if err != nil {
					return nil, err
				}
				singleRecord, err := res.Single()
				if err != nil {
					return nil, err
				}
				return &users.User{
					Username: singleRecord.Values[0].(string),
					Email:    singleRecord.Values[1].(string),
					Password: singleRecord.Values[2].(string),
				}, nil
			})
		Expect(err).To(BeNil(), "Transaction should successfully run")
		persistedUser := result.(*users.User)
		Expect(persistedUser.Username).To(Equal(username))
		Expect(persistedUser.Email).To(Equal(email))
		Expect(passwordsMatch(initialPassword, persistedUser.Password)).To(BeTrue())
	})
})

func passwordsMatch(initialPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(initialPassword))
	return err == nil
}

func Close(closer io.Closer, resourceName string) {
	Expect(closer.Close()).
		To(BeNil(), "%s should close", resourceName)
}

func startContainer(ctx context.Context, username, password string) (testcontainers.Container, error) {
	request := testcontainers.ContainerRequest{
		Image:        "neo4j",
		ExposedPorts: []string{"7687/tcp"},
		Env:          map[string]string{"NEO4J_AUTH": fmt.Sprintf("%s/%s", username, password)},
		WaitingFor:   wait.ForLog("Bolt enabled"),
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
}
