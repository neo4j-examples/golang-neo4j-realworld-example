package stackoverflow

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type StackOverflowRepository interface {
}

type StackOverflowNeo4jRepository struct {
	Driver neo4j.Driver
}
