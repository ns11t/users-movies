package datastore

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/ns11t/users-movies/shared/model"
)

func InsertPerson(person model.Person) (sql.Result, error) {
	return db.Exec("INSERT INTO person VALUES (default, $1, $2, $3, $4, $5)",
		person.Login, person.Password, person.Name, person.Age, person.PhoneNumber)
}

func GetPersonPasswordByLogin(login string) (*model.Person, error) {
	var person model.Person
	row := db.QueryRow("SELECT password FROM person WHERE login=$1", login)
	return &person, row.Scan(&person.Password)
}

func GetPersonByLogin(login string) (*model.Person, error) {
	var person model.Person
	row := db.QueryRow("SELECT * FROM person WHERE login=$1", login)
	return &person, row.Scan(&person.Id, &person.Login, &person.Password, &person.Name, &person.Age, &person.PhoneNumber)
}

// Checks if person with given id exists
func CheckPersonExists(login string) (bool, error) {
	row := db.QueryRow("SELECT 1 FROM person WHERE login=$1", login)
	var exists int
	err := row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists > 0, nil
}
