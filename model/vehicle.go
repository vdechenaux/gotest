package model

import (
	"../database"
)

// Vehicle is the representation of the table vehicles
type Vehicle struct {
	ID                   int64    `json:"id"`
	Name                 string   `json:"name"`
	Model                string   `json:"model"`
	Manufacturer         string   `json:"manufacturer"`
	CostInCredits        *int     `json:"cost_in_credits"`
	Length               *float64 `json:"length"`
	MaxAtmospheringSpeed *int     `json:"max_atmosphering_speed"`
	Crew                 *int     `json:"crew"`
	Passengers           *int     `json:"passengers"`
	CargoCapacity        *int     `json:"cargo_capacity"`
	Consumables          *string  `json:"consumables"`
	VehicleClass         string   `json:"vehicle_class"`
	Created              string   `json:"created"`
	Edited               string   `json:"edited"`
}

// FindPeopleVehicleList returns all Vehicle owned by People with ID id
func FindPeopleVehicleList(id int) ([]Vehicle, error) {
	objs := make([]Vehicle, 0, 5)

	rows, err := database.DB.Query(`
		SELECT
			v.name,
			v.model,
			v.manufacturer,
			CASE v.cost_in_credits WHEN 'unknown' THEN NULL ELSE v.cost_in_credits END,
			CASE v.length WHEN 'unknown' THEN NULL ELSE v.length END,
			CASE v.max_atmosphering_speed WHEN 'unknown' THEN NULL WHEN 'na' THEN NULL ELSE v.max_atmosphering_speed END,
			CASE v.crew WHEN 'unknown' THEN NULL ELSE v.crew END,
			CASE v.passengers WHEN 'unknown' THEN NULL ELSE v.passengers END,
			CASE v.cargo_capacity WHEN 'unknown' THEN NULL ELSE v.cargo_capacity END,
			CASE v.consumables WHEN 'unknown' THEN NULL ELSE v.consumables END,
			v.vehicle_class,
			v.created,
			v.edited,
			v.id
		FROM vehicles v
		INNER JOIN people_vehicles pv ON v.id = pv.vehicles
		WHERE pv.people = ?
	`, id)

	if err != nil {
		return objs, err
	}

	defer rows.Close()
	for rows.Next() {
		var obj Vehicle
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
			&obj.VehicleClass,
			&obj.Created,
			&obj.Edited,
			&obj.ID,
		)

		objs = append(objs, obj)
	}

	return objs, nil
}
