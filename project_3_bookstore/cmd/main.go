package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AaryanO2/go_projects/project_3_bookstore/pkg/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	http.Handle("/", r)
	fmt.Println("Serving now on port:9010")
	log.Fatal(http.ListenAndServe("0.0.0.0:9010", r))
}
