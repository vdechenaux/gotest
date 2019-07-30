package main

import (
	"net/http"

	"./controller"
	"./database"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := database.Open("swapi.dat")

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/peoples", controller.GetPeopleList).Methods("GET")
	r.HandleFunc("/peoples", controller.CreatePeople).Methods("POST")

	r.HandleFunc("/peoples/{id}", controller.GetPeople).Methods("GET")
	r.HandleFunc("/peoples/{id}", controller.DeletePeople).Methods("DELETE")
	r.HandleFunc("/peoples/{id}", controller.UpdatePeople).Methods("PUT")

	r.HandleFunc("/peoples/{id}/starships", controller.GetPeopleStarshipList).Methods("GET")
	r.HandleFunc("/peoples/{id}/vehicles", controller.GetPeopleVehicleList).Methods("GET")

	r.PathPrefix("/doc").Handler(http.StripPrefix("/doc", http.FileServer(http.Dir("doc"))))

	http.ListenAndServe(":8001", r)
}
