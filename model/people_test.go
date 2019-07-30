package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"../database"
)

func TestValidate(t *testing.T) {
	var people People

	expected := "name cannot be empty"
	if err := people.Validate(); err != nil && err.Error() != expected {
		t.Errorf("Validate() = %v; want %v", err, expected)
	}

	people.Name = "Toto"

	if err := people.Validate(); err != nil {
		t.Errorf("Validate() = %v; want %v", err, expected)
	}

	height := 0
	people.Height = &height

	expected = "height must be greater than 0 if provided"
	if err := people.Validate(); err != nil && err.Error() != expected {
		t.Errorf("Validate() = %v; want %v", err, expected)
	}

	height = 1

	if err := people.Validate(); err != nil {
		t.Errorf("Validate() = %v; want %v", err, expected)
	}

	mass := -0.
	people.Mass = &mass

	expected = "mass must be greater than 0 if provided"
	if err := people.Validate(); err != nil && err.Error() != expected {
		t.Errorf("Validate() = %v; want %v", err, expected)
	}

	mass = 1

	if err := people.Validate(); err != nil {
		t.Errorf("Validate() = %v; want %v", err, expected)
	}
}

func TestFindAndUpdate(t *testing.T) {
	copyDb()

	_, err := FindPeopleByID(100)

	expected := sql.ErrNoRows
	if err != expected {
		t.Errorf("TestFindAndUpdate() = %v; want %v", err, expected)
	}

	people, err := FindPeopleByID(1)

	expected = nil
	if err != expected {
		t.Errorf("TestFindAndUpdate() = %v; want %v", err, expected)
	}

	expectedStr := `{
  "id": 1,
  "name": "Luke Skywalker",
  "height": 172,
  "mass": 77,
  "hair_color": "blond",
  "skin_color": "fair",
  "eye_color": "blue",
  "birth_year": "19BBY",
  "gender": "male",
  "created": "2014-12-09T13:50:51.644000Z",
  "edited": "2014-12-20T21:17:56.891000Z"
}`

	peopleStr, _ := json.MarshalIndent(people, "", "  ")

	if string(peopleStr) != expectedStr {
		t.Errorf("TestFindAndUpdate() = %v; want %v", string(peopleStr), expectedStr)
	}

	people.Name = "Toto"
	*people.Gender = "female"
	*people.Mass = 3

	ok, err := UpdatePeople(people)
	if ok != true {
		t.Errorf("TestFindAndUpdate() = %v; want %v", ok, true)
	}

	people, err = FindPeopleByID(1)

	expected = nil
	if err != expected {
		t.Errorf("TestFindAndUpdate() = %v; want %v", err, expected)
	}

	expectedStr = fmt.Sprintf(`{
  "id": 1,
  "name": "Toto",
  "height": 172,
  "mass": 3,
  "hair_color": "blond",
  "skin_color": "fair",
  "eye_color": "blue",
  "birth_year": "19BBY",
  "gender": "female",
  "created": "2014-12-09T13:50:51.644000Z",
  "edited": "%s"
}`, people.Edited)

	peopleStr, _ = json.MarshalIndent(people, "", "  ")

	if string(peopleStr) != expectedStr {
		t.Errorf("TestFindAndUpdate() = %v; want %v", string(peopleStr), expectedStr)
	}
}

func TestDelete(t *testing.T) {
	copyDb()

	_, err := FindPeopleByID(1)

	if err != nil {
		t.Errorf("TestDelete() = %v; want %v", err, nil)
	}

	DeletePeopleByID(1)

	_, err = FindPeopleByID(1)

	expected := sql.ErrNoRows
	if err != expected {
		t.Errorf("TestDelete() = %v; want %v", err, expected)
	}
}

func copyDb() {
	os.Remove("test.dat")

	in, _ := os.Open("../swapi.dat")
	defer in.Close()

	out, _ := os.Create("test.dat")
	defer out.Close()

	_, _ = io.Copy(out, in)
	out.Close()

	database.Open("test.dat")
}
