package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AaryanO2/go_projects/project_5_postgres/pkg/config"
	models "github.com/AaryanO2/go_projects/project_5_postgres/pkg/models"
	"github.com/gorilla/mux"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body: %s", err)
	}

	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "Stock created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert string to int. %v", err)
	}
	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatalf("Unable to get stock. %v", err)
	}
	json.NewEncoder(w).Encode(stock)
}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()
	if err != nil {
		log.Fatalf("Unable to get all stocks. %v", err)
	}
	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert string to int. %v", err)
	}
	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}
	updatedRows := updateStock(int64(id), stock)
	msg := fmt.Sprintf("Stock updated successfully. %v", updatedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert string to int. %v", err)
	}
	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock deleted Successfully.\nTotal rows/records affected: %v", deletedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := config.CreateConnection()
	defer db.Close()
	var id int64

	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stockid`
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record.")
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := config.CreateConnection()
	defer db.Close()
	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned.")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the rows")
	}
	return stock, err
}

func getAllStocks() ([]models.Stock, error) {
	db := config.CreateConnection()
	var stocks []models.Stock
	defer db.Close()

	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute query. %v", err)
	}

	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("Unable to scan row %v", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {
	db := config.CreateConnection()
	defer db.Close()

	sqlStatement := `UPDATE stocks SET name=$1, price=$2, company=$3 WHERE stockid=$4`
	res, err := db.Exec(sqlStatement, &stock.Name, &stock.Price, &stock.Company, &id)
	if err != nil {
		log.Fatalf("Unable to execute query. %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking affected rows. %v", err)
	}
	fmt.Printf("Rows affected: %v", rowsAffected)
	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := config.CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stockid =$1`
	res, err := db.Exec(sqlStatement, &id)
	if err != nil {
		log.Fatalf("Unable to execute query.%v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking rows affected. %v", err)
	}
	return rowsAffected
}
