package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

type Coil struct {
    ID          int     `json:"id"`
    OrderPrice  float64 `json:"orderPrice"`
    Sets        int     `json:"sets"`
    CoilName    string  `json:"coilName"`
    Weight      float64 `json:"weight"`
    WireReq     float64 `json:"wireReq"`
    WireGauge   float64 `json:"wireGauge"`
    Delivered   bool    `json:"delivered"`
    DeliveredOn string  `json:"deliveredOn"`
}

type Wire struct {
	ID          int     `json:"id"`
	Sets        int     `json:"sets"`
	Weight      float64 `json:"weight"`
	WireReq     float64 `json:"wireReq"`
	WireGauge   float64 `json:"wireGauge"`
	Delivered   bool    `json:"delivered"`
	DeliveredOn string  `json:"deliveredOn"`
}

type Order struct {
	ID          int     `json:"id"`
	Sets        int     `json:"sets"`
	Weight      float64 `json:"weight"`
	WireReq     float64 `json:"wireReq"`
	WireGauge   float64 `json:"wireGauge"`
	ToDeliverOn string  `json:"toDeliverOn"`
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	conn := os.Getenv("CONN_URL")
	db, err := sql.Open("pgx", conn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
		fmt.Println("Error connecting to database")
	}
	defer db.Close()

	createCoilTable(db)
	createWireTable(db)
	createOrdersTable(db)

	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET /coil")
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "Welcome to the Coil API"}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
	})



	http.ListenAndServe(":3000", nil)
}

func createCoilTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS coil (
        id SERIAL PRIMARY KEY,
        orderPrice FLOAT NOT NULL,
        sets INT NOT NULL,
        coilName VARCHAR NOT NULL,
        weight FLOAT NOT NULL,
        wireReq FLOAT NOT NULL,
        wireGauge FLOAT NOT NULL,
        delivered BOOLEAN NOT NULL DEFAULT FALSE,
        deliveredOn TIMESTAMP DEFAULT NOW()
    )`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func createWireTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS wire (
        id SERIAL PRIMARY KEY,
        weight FLOAT NOT NULL,
        wireReq FLOAT NOT NULL,
        wireGauge FLOAT NOT NULL,
        used BOOLEAN DEFAULT FALSE,
        deliveredOn TIMESTAMP
    )`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func createOrdersTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS orders (
        id SERIAL PRIMARY KEY,
        coilName VARCHAR NOT NULL,
        sets INT NOT NULL,
        weight FLOAT NOT NULL,
        wireReq FLOAT NOT NULL,
        wireGauge FLOAT NOT NULL,
        toDeliverOn TIMESTAMP
    )`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func displayCoils (db *sql.DB) {
	// Display all coils
	Query:= `SELECT * FROM coil`

	rows , err := db.Query(Query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Coils:",rows)
	defer rows.Close()
}
