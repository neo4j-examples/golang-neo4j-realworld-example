package main

import (
	"net/http"
	"os"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/zhouanqiNB/neo4j-backend/blob/dev/stackoverflow"
)

func main() {
	neo4jUri, found := os.LookupEnv("NEO4J_URI")
	if !found {
		panic("NEO4J_URI not set")
	}
	neo4jUsername, found := os.LookupEnv("NEO4J_USERNAME")
	if !found {
		panic("NEO4J_USERNAME not set")
	}
	neo4jPassword, found := os.LookupEnv("NEO4J_PASSWORD")
	if !found {
		panic("NEO4J_PASSWORD not set")
	}
	// usersRepository := users.UserNeo4jRepository{
	// 	Driver: driver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, "")),
	// }
	stackOverflowRepository := stackoverflow.StackOverflowNeo4jRepository{
		Driver: driver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, "")),
	}

	// registrationHandler := &users.UserRegistrationHandler{
	// 	Path:           "/users",
	// 	UserRepository: &usersRepository,
	// }
	// loginHandler := &users.UserLoginHandler{
	// 	Path:           "/users/login",
	// 	UserRepository: &usersRepository,
	// }

	queryHandler := &stackoverflow.QueryHandler{
		Path:                    "/users/login",
		StackOverflowRepository: &stackOverflowRepository,
	}

	server := http.NewServeMux()
	// server.HandleFunc(registrationHandler.Path, registrationHandler.Register)
	// server.HandleFunc(loginHandler.Path, loginHandler.Login)
	server.HandleFunc(queryHandler.Path, queryHandler.Login)

	if err := http.ListenAndServe(":3000", server); err != nil {
		panic(err)
	}
}

func driver(target string, token neo4j.AuthToken) neo4j.Driver {
	result, err := neo4j.NewDriver(target, token)
	if err != nil {
		panic(err)
	}
	return result
}
