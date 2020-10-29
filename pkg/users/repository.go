package users

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	RegisterUser(user *User) error
	FindByEmailAndPassword(email string, password string) (*User, error)
}

type UserNeo4jRepository struct {
	Driver neo4j.Driver
}

func (u *UserNeo4jRepository) RegisterUser(user *User) (err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close()
	}()
	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.persistUser(tx, user)
		}); err != nil {
		return err
	}
	return nil
}

func (u *UserNeo4jRepository) FindByEmailAndPassword(email string, password string) (user *User, err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = session.Close()
	}()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return u.findUser(tx, email, password)
	})
	if result == nil {
		return nil, err
	}
	user = result.(*User)
	return user, err
}

func (u *UserNeo4jRepository) persistUser(tx neo4j.Transaction, user *User) (interface{}, error) {
	query := "CREATE (:User {email: $email, username: $username, password: $password})"
	hashedPassword, err := hash(user.Password)
	if err != nil {
		return nil, err
	}
	parameters := map[string]interface{}{
		"email":    user.Email,
		"username": user.Username,
		"password": hashedPassword,
	}
	_, err = tx.Run(query, parameters)
	return nil, err
}

func (u *UserNeo4jRepository) findUser(tx neo4j.Transaction, email string, password string) (*User, error) {
	result, err := tx.Run(
		"MATCH (u:User {email: $email}) RETURN u.username AS username, u.password AS password",
		map[string]interface{}{
			"email": email,
		},
	)
	if err != nil {
		return nil, err
	}
	record, err := result.Single()
	if err != nil {
		return nil, err
	}
	hashedPassword, _ := record.Get("password")
	if !passwordsMatch(hashedPassword.(string), password) {
		return nil, nil
	}
	username, _ := record.Get("username")
	return &User{
		Username: username.(string),
		Email:    email,
	}, nil
}

func passwordsMatch(hashedPassword string, clearTextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(clearTextPassword))
	return err == nil
}

func hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
