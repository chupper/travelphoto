package user

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"

	"log"

	"github.com/chupper/travelphoto/shared/database"
)

const (
	table = "admin_user"
)

// User Model
type User struct {
	ID           int
	UserName     string
	PasswordHash []byte
	Salt         []byte
}

// Validate fetch the user and validate the password
func Validate(db database.Connection, userName string, password string) (bool, error) {

	// retrieve the user
	res, err := db.Query(fmt.Sprintf(`
		SELECT
			ID,
			PASSWORDHASH,
			SALT
		FROM %v
		WHERE USERNAME = $1
	`, table), userName)
	defer res.Close()

	if err != nil {
		log.Fatal("Error retrieving item", err)
		return false, err
	}

	res.Next()
	var id int
	var passwordhash, salt []byte
	if err = res.Scan(&id, &passwordhash, &salt); err != nil {
		return false, err
	}

	hash, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 64)

	log.Println("Hash", string(hash))
	log.Println("Passowrd Hash", string(passwordhash))

	return string(hash) == string(passwordhash), nil
}

// CreateUser Creates the user
func CreateUser(db database.Connection, userName string, password string) error {

	// generate the salt and the password hash
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	hash, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 64)

	if err != nil {
		log.Fatal("Error hashing password")
	}

	// create the password
	_, err = db.Exec(fmt.Sprintf(`
		INSERT INTO %v (USERNAME, PASSWORDHASH, SALT)
		VALUES($1, $2, $3)
	`, table), userName, hash, salt)

	if err != nil {
		return err
	}

	return nil
}
