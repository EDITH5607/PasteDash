package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-sql-driver/mysql"
)


type User struct {
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
}

type UserModel struct {
	DB *sql.DB
}


// insert user data to the db 
func (m *UserModel) Insert(name, email, password  string ) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password),12)
	if err != nil {
		return  err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt,name,email,string(hashPassword))
	if err!= nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			if mysqlError.Number == 1062 && strings.Contains(mysqlError.Message, "users_uc_email"){
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return  nil
}


//checking the password is correct or not , authenticating the user 
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	// execute the query and retrive the data to the id and hashpassword 
	err := m.DB.QueryRow(stmt,email).Scan(&id, &hashPassword)
	if err!=nil {
		if errors.Is(err,sql.ErrNoRows) {
			return 0, ErrInvalidCredential
		} else {
			return 0,err
		}
	}
	// convert the given password to the hash and compare the db user password and current hashed password 
	err = bcrypt.CompareHashAndPassword(hashPassword,[]byte(password))
	if err !=nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0,ErrInvalidCredential
		} else {
			return 0, err
		}
	}
	return id, nil
}

// check the user id is valid or not using db data
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	// EXIST is the best method for checking the valid user and is the best optimized query
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	// execute the query and scan to exist bool value
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}