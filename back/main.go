package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	// "github.com/gin-gonic/gin/render"
)

// var db *sql.DB

func main() {

	router := gin.Default()

	router.GET("api/v1/persons", getPersons)
	router.GET("api/v1/persons/:id", getPersonByID)
	router.POST("api/v1/persons", postPerson)
	router.PATCH("api/v1/persons/:id", editPerson)
	router.DELETE("api/v1/persons/:id", deletePerson)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("Error %s", err)
	}

}

func connect_to_db() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		// panic("Error connecting to DB")
		return nil, fmt.Errorf("error connecting to DB")
	}

	err = db.Ping()
	if err != nil {
		// panic("Error sending initial Ping")
		return nil, fmt.Errorf("error getting initial ping")
	}
	fmt.Println("Connected!")
	return db, nil
}
