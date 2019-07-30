package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"../model"

	"database/sql"

	"github.com/gorilla/mux"
)

type ApiError struct {
	Message string `json:"message"`
}

var logger = log.New(os.Stderr, "[People Controller] ", log.LstdFlags|log.Llongfile|log.LUTC)

func GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	people, err := model.FindPeopleByID(id)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ApiError{fmt.Sprintf("People %d not found", id)})
		return
	} else if err != nil {
		w.WriteHeader(500)
		logger.Println(err)
		return
	}

	b, _ := json.Marshal(people)

	w.Write(b)
}

func GetPeopleList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	vars := r.URL.Query()

	pageStr, ok := vars["page"]
	var page int
	if ok {
		page, _ = strconv.Atoi(pageStr[0])
	} else {
		page = 1
	}

	perPageStr, ok := vars["per_page"]
	var perPage int
	if ok {
		perPage, _ = strconv.Atoi(perPageStr[0])
	} else {
		perPage = 30
	}

	if page < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiError{"page must be greater than 0"})
		return
	}

	if perPage < 1 || perPage > 30 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiError{"per_page must be between 1 and 30"})
		return
	}

	peoples, err := model.FindAllPeople(page, perPage)
	if err != nil {
		w.WriteHeader(500)
		logger.Println(err)
		return
	}

	count, err := model.CountPeople()
	if err != nil {
		w.WriteHeader(500)
		logger.Println(err)
		return
	}

	json.NewEncoder(w).Encode(model.PeopleList{
		Metadata: model.PeopleMetadata{
			Page:    page,
			PerPage: perPage,
			Total:   count,
		},
		Data: peoples,
	})
}

func DeletePeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := model.DeletePeopleByID(id)

	if err != nil {
		w.WriteHeader(500)
		logger.Println(err)
		return
	}

	w.WriteHeader(204)
}

func CreatePeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var people model.People

	err := json.NewDecoder(r.Body).Decode(&people)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiError{"Body is not valid JSON"})
		return
	}

	err = people.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiError{err.Error()})
		return
	}

	people, err = model.CreatePeople(people)

	if err != nil {
		w.WriteHeader(500)
		logger.Println(err)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(people)
}

func UpdatePeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var people model.People

	err := json.NewDecoder(r.Body).Decode(&people)
	people.ID, _ = strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiError{"Body is not valid JSON"})
		return
	}

	err = people.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiError{err.Error()})
		return
	}

	ok, err := model.UpdatePeople(people)

	if err != nil {
		w.WriteHeader(500)
		logger.Println(err)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ApiError{fmt.Sprintf("People %d not found", people.ID)})
		return
	}

	people, err = model.FindPeopleByID(people.ID)

	if err != nil {
		w.WriteHeader(500)
		logger.Println(err)
		return
	}

	json.NewEncoder(w).Encode(people)
}

func GetPeopleStarshipList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	objs, err := model.FindPeopleStarshipList(id)

	if err != nil {
		w.WriteHeader(500)
		logger.Println(err)
		return
	}

	json.NewEncoder(w).Encode(objs)
}

func GetPeopleVehicleList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	objs, err := model.FindPeopleVehicleList(id)

	if err != nil {
		w.WriteHeader(500)
		logger.Println(err)
		return
	}

	json.NewEncoder(w).Encode(objs)
}
