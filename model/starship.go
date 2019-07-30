package model

import (
	"../database"
)

// Starship is the representation of the table starships
type Starship struct {
	ID                   int64    `json:"id"`
	Name                 string   `json:"name"`
	Model                string   `json:"model"`
	Manufacturer         string   `json:"manufacturer"`
	CostInCredits        *int64   `json:"cost_in_credits"`
	Length               *float64 `json:"length"`
	MaxAtmospheringSpeed *int     `json:"max_atmosphering_speed"`
	Crew                 *int     `json:"crew"`
	Passengers           *int     `json:"passengers"`
	CargoCapacity        *int64   `json:"cargo_capacity"`
	Consumables          *string  `json:"consumables"`
	HyperdriveRating     *float64 `json:"hyperdrive_rating"`
	MGLT                 *int     `json:"MGLT"`
	StarshipClass        string   `json:"starship_class"`
	Created              string   `json:"created"`
	Edited               string   `json:"edited"`
}

// FindPeopleStarshipList returns all Starship owned by People with ID id
func FindPeopleStarshipList(id int) ([]Starship, error) {
	objs := make([]Starship, 0, 5)

	rows, err := database.DB.Query(`
		SELECT
			s.name,
			s.model,
			s.manufacturer,
			CASE s.cost_in_credits WHEN 'unknown' THEN NULL ELSE s.cost_in_credits END,
			CASE s.length WHEN 'unknown' THEN NULL ELSE s.length END,
			CASE s.max_atmosphering_speed WHEN 'unknown' THEN NULL WHEN 'na' THEN NULL ELSE s.max_atmosphering_speed END,
			CASE s.crew WHEN 'unknown' THEN NULL ELSE s.crew END,
			CASE s.passengers WHEN 'unknown' THEN NULL ELSE s.passengers END,
			CASE s.cargo_capacity WHEN 'unknown' THEN NULL ELSE s.cargo_capacity END,
			CASE s.consumables WHEN 'unknown' THEN NULL ELSE s.consumables END,
			CASE s.hyperdrive_rating WHEN 'unknown' THEN NULL ELSE s.hyperdrive_rating END,
			CASE s.MGLT WHEN 'unknown' THEN NULL ELSE s.MGLT END,
			s.starship_class,
			s.created,
			s.edited,
			s.id
		FROM starships s
		INNER JOIN people_starships ps ON s.id = ps.starships
		WHERE ps.people = ?
	`, id)

	if err != nil {
		return objs, err
	}

	defer rows.Close()
	for rows.Next() {
		var obj Starship
		rows.Scan(
			&obj.Name,
			&obj.Model,
			&obj.Manufacturer,
			&obj.CostInCredits,
			&obj.Length,
			&obj.MaxAtmospheringSpeed,
			&obj.Crew,
			&obj.Passengers,
			&obj.CargoCapacity,
			&obj.Consumables,
			&obj.HyperdriveRating,
			&obj.MGLT,
			&obj.StarshipClass,
			&obj.Created,
			&obj.Edited,
			&obj.ID,
		)

		objs = append(objs, obj)
	}

	return objs, nil
}
