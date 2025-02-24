package router

import (
	controller "github.com/AaryanO2/go_projects/project_5_postgres/pkg/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/stock/{id}", controller.GetStock).Methods("GET")
	r.HandleFunc("/api/stock", controller.GetAllStocks).Methods("GET")
	r.HandleFunc("/api/newstock", controller.CreateStock).Methods("POST")
	r.HandleFunc("/api/stock/{id}", controller.UpdateStock).Methods("PUT")
	r.HandleFunc("/api/deletestock/{id}", controller.DeleteStock).Methods("DELETE")
	return r
}
