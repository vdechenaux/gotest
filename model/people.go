package model

import (
	"errors"
	"time"

	"../database"
)

// PeopleList represents the output format of the People list API response
type PeopleList struct {
	Metadata PeopleMetadata `json:"metadata"`
	Data     []People       `json:"data"`
}

// PeopleMetadata stores pagination informations of PeopleList
type PeopleMetadata struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

// People is the representation of table people
type People struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Height    *int     `json:"height"`
	Mass      *float64 `json:"mass"`
	HairColor *string  `json:"hair_color"`
	SkinColor *string  `json:"skin_color"`
	EyeColor  *string  `json:"eye_color"`
	BirthYear *string  `json:"birth_year"`
	Gender    *string  `json:"gender"`
	Created   string   `json:"created"`
	Edited    string   `json:"edited"`
}

// Validate checks fields requirements
func (people People) Validate() error {
	if len(people.Name) == 0 {
		return errors.New("name cannot be empty")
	}

	if people.Height != nil && *people.Height <= 0 {
		return errors.New("height must be greater than 0 if provided")
	}

	if people.Mass != nil && *people.Mass <= 0 {
		return errors.New("mass must be greater than 0 if provided")
	}

	return nil
}

// FindPeopleByID returns a People by ID
func FindPeopleByID(id int) (People, error) {
	var people People

	row := database.DB.QueryRow(`
		SELECT
			id,
			name,
			CASE height WHEN 'unknown' THEN NULL ELSE height END,
			CASE mass WHEN 'unknown' THEN NULL ELSE REPLACE(mass, ',', '') END,
			CASE hair_color WHEN 'na' THEN NULL WHEN 'none' THEN NULL WHEN 'unknown' THEN NULL ELSE hair_color END,
			CASE skin_color WHEN 'na' THEN NULL WHEN 'none' THEN NULL WHEN 'unknown' THEN NULL ELSE skin_color END,
			CASE eye_color WHEN 'na' THEN NULL WHEN 'none' THEN NULL WHEN 'unknown' THEN NULL ELSE eye_color END,
			CASE birth_year WHEN 'unknown' THEN NULL ELSE birth_year END,
			CASE gender WHEN 'na' THEN NULL WHEN 'none' THEN NULL ELSE gender END,
			created,
			edited
		FROM people
		WHERE id = ?
	`, id)

	err := row.Scan(
		&people.ID,
		&people.Name,
		&people.Height,
		&people.Mass,
		&people.HairColor,
		&people.SkinColor,
		&people.EyeColor,
		&people.BirthYear,
		&people.Gender,
		&people.Created,
		&people.Edited,
	)

	return people, err
}

// FindAllPeople returns paginated list of People
func FindAllPeople(page, limit int) ([]People, error) {
	peoples := make([]People, 0, limit)

	rows, err := database.DB.Query(`
		SELECT
			id,
			name,
			CASE height WHEN 'unknown' THEN NULL ELSE height END,
			CASE mass WHEN 'unknown' THEN NULL ELSE REPLACE(mass, ',', '') END,
			CASE hair_color WHEN 'na' THEN NULL WHEN 'none' THEN NULL WHEN 'unknown' THEN NULL ELSE hair_color END,
			CASE skin_color WHEN 'na' THEN NULL WHEN 'none' THEN NULL WHEN 'unknown' THEN NULL ELSE skin_color END,
			CASE eye_color WHEN 'na' THEN NULL WHEN 'none' THEN NULL WHEN 'unknown' THEN NULL ELSE eye_color END,
			CASE birth_year WHEN 'unknown' THEN NULL ELSE birth_year END,
			CASE gender WHEN 'na' THEN NULL WHEN 'none' THEN NULL ELSE gender END,
			created,
			edited
		FROM people
		ORDER BY CAST(id AS INT) ASC
		LIMIT ?, ?
	`, (page-1)*limit, limit)

	if err != nil {
		return peoples, err
	}

	defer rows.Close()
	for rows.Next() {
		var people People
		rows.Scan(
			&people.ID,
			&people.Name,
			&people.Height,
			&people.Mass,
			&people.HairColor,
			&people.SkinColor,
			&people.EyeColor,
			&people.BirthYear,
			&people.Gender,
			&people.Created,
			&people.Edited,
		)

		peoples = append(peoples, people)
	}

	return peoples, nil
}

// DeletePeopleByID takes an ID of People and deletes matching People from DB
func DeletePeopleByID(id int) error {
	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	tx, err := database.DB.Begin()

	if err != nil {
		return err
	}

	_, err = database.DB.Exec("DELETE FROM people WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = database.DB.Exec("DELETE FROM people_vehicles WHERE people = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = database.DB.Exec("DELETE FROM people_starships WHERE people = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = database.DB.Exec("DELETE FROM people_species WHERE people = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// CreatePeople takes a People and save it in DB
func CreatePeople(people People) (People, error) {
	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	people.Created = time.Now().Format("2006-01-02T08:04:05.999999Z")
	people.Edited = people.Created

	database.DB.QueryRow("SELECT MAX(CAST(id as int))+1 from people").Scan(&people.ID)

	_, err := database.DB.Exec(`
			INSERT INTO people
			(id, url, name, height, mass, hair_color, skin_color, eye_color, birth_year, gender, created, edited)
			VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		people.ID,
		people.ID,
		people.Name,
		people.Height,
		people.Mass,
		people.HairColor,
		people.SkinColor,
		people.EyeColor,
		people.BirthYear,
		people.Gender,
		people.Created,
		people.Edited,
	)

	if err != nil {
		return people, err
	}

	return people, nil
}

// UpdatePeople takes a People and save it in DB
// it returns false if the People is not in DB
func UpdatePeople(people People) (bool, error) {
	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	people.Edited = time.Now().Format("2006-01-02T08:04:05.999999Z")

	result, err := database.DB.Exec(`
			UPDATE people
			SET
				name = ?,
				height = ?,
				mass = ?,
				hair_color = ?,
				skin_color = ?,
				eye_color = ?,
				birth_year = ?,
				gender = ?,
				edited = ?
			WHERE id = ?
		`,
		people.Name,
		people.Height,
		people.Mass,
		people.HairColor,
		people.SkinColor,
		people.EyeColor,
		people.BirthYear,
		people.Gender,
		people.Edited,
		people.ID,
	)

	var rows int64
	if result != nil {
		rows, _ = result.RowsAffected()
	}

	return rows == 1, err
}

// CountPeople returns the number of peoples
func CountPeople() (int, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(1) FROM people").Scan(&count)

	return count, err
}
