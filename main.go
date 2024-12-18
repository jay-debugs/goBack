package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	//apis

	app := fiber.New()

	app.Use(cors.New(cors.Config{
	AllowOrigins: "*",
}))

app.Post("/coils", func(c *fiber.Ctx) error {
    coil := new(Coil)
    err := c.BodyParser(coil)
    if err != nil {
        return c.Status(400).SendString("Bad Request")
    }
guery := `INSERT INTO coil (orderPrice, sets, coilName, weight, wireReq, wireGauge) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = db.Exec(guery, coil.OrderPrice, coil.Sets, coil.CoilName, coil.Weight, coil.WireReq, coil.WireGauge)
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}

	// Additional logic to handle the parsed coil data can be added here
	return c.JSON(coil) // Placeholder for successful handling

})

	log.Fatal(app.Listen(":3000"))

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






