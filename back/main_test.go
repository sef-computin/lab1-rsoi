package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"

	_ "github.com/lib/pq"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestPostPerson(t *testing.T) {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic("Error connecting to DB")
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic("Error sending initial Ping")
	}
	fmt.Println("Connected!")

	r := SetUpRouter()
	r.POST("/persons", postPerson)

	newPerson := person{
		Name:    "Demo Name",
		Address: "Demo Address",
		Work:    "DEMO WORK",
		Age:     33,
	}
	jsonValue, _ := json.Marshal(newPerson)
	req, _ := http.NewRequest("POST", "/persons", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.MatchRegex(t, w.Header().Get("Location"), "/api/v1/persons/[0-9]+")
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetPersons(t *testing.T) {

	r := SetUpRouter()

	r.GET("/persons", getPersons)

	req, err := http.NewRequest("GET", "/persons", nil)
	if err != nil {
		panic("error creating request")
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var persons []person

	json.Unmarshal(w.Body.Bytes(), &persons)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, persons, []person{})
}

func TestGetPerson(t *testing.T) {
	r := SetUpRouter()
	id_found := 2
	id_not_found := 35

	r.GET("/persons/:id", getPersonByID)

	req_found, _ := http.NewRequest("GET", "/persons/"+fmt.Sprint(id_found), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req_found)

	var person_found person

	json.Unmarshal(w.Body.Bytes(), &person_found)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, person_found, nil)

	req_not_found, _ := http.NewRequest("GET", "/persons/"+fmt.Sprint(id_not_found), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req_not_found)

	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestUpdatePerson(t *testing.T) {
	r := SetUpRouter()
	r.PATCH("/persons/:id", editPerson)

	newPerson := person{
		ID:      2,
		Name:    "Edited Name",
		Address: "Demo Address",
		Work:    "DEMO WORK",
		Age:     33,
	}

	jsonValue, _ := json.Marshal(newPerson)
	reqFound, _ := http.NewRequest("PATCH", "/persons/"+fmt.Sprint(newPerson.ID), bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)
	assert.Equal(t, http.StatusOK, w.Code)

	newPerson.ID = 35 // not found

	reqNotFound, _ := http.NewRequest("PATCH", "/persons/"+fmt.Sprint(newPerson.ID), bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqNotFound)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
