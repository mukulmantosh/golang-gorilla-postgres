package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go-stocks/models"
	"log"
	"net/http"
	"os"
	"strconv"
)

type response struct {
	ID      int64  `json:"id"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database.")
	return db
}

func GetStock(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.ParseInt(params["id"], 0, 64)
	if err != nil {
		log.Fatalf("Unable to convert the string into int %v", err)
	}
	stock, err := getStock(id)

	if err != nil {
		log.Fatalf("Unable to get stock %v", err)
	}
	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(stock)

}

func GetAllStock(writer http.ResponseWriter, request *http.Request) {
	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("Unable to get all the stocks %v", err)
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(stocks)
}

func CreateStock(writer http.ResponseWriter, request *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(request.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the request body %v", err)
	}
	defer request.Body.Close()

	insertID := insertStock(stock)
	res := response{
		ID:      insertID,
		Message: "Stock created successfully",
	}
	json.NewEncoder(writer).Encode(res)
}

func UpdateStock(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.ParseInt(params["id"], 0, 64)
	if err != nil {
		log.Fatalf("Unable to convert the string into int %v", err)
	}
	var stock models.Stock
	json.NewDecoder(request.Body).Decode(&stock)
	defer request.Body.Close()

	if err != nil {
		log.Fatalf("Unable to decode the request body %v", err)
	}
	updateRows := updateStock(id, stock)

	msg := fmt.Sprintf("Stock updated successfully. Total rows affected : %v", updateRows)
	res := response{
		ID:      id,
		Message: msg,
	}
	json.NewEncoder(writer).Encode(res)

}

func DeleteStock(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, err := strconv.ParseInt(params["id"], 0, 64)
	if err != nil {
		log.Fatalf("Unable to convert the string into int %v", err)
	}

	deletedRows := deleteStock(id)
	msg := fmt.Sprintf("Stock deleted successfully. Total rows affected %v", deletedRows)
	res := response{
		ID:      id,
		Message: msg,
	}
	json.NewEncoder(writer).Encode(res)
}
