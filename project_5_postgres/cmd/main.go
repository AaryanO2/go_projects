package main

import (
	"fmt"
	"log"
	"net/http"

	router "github.com/AaryanO2/go_projects/project_5_postgres/pkg/routers"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on the port 9010...")
	log.Fatal(http.ListenAndServe("0.0.0.0:9010", r))

}
