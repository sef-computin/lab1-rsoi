package main

import "sync"

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "persons"
)

type autoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id         int
}

var autoid autoInc

func (a *autoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()

	id = a.id
	a.id++
	return
}

type person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Work    string `json:"work"`
	Age     int    `json:"age"`
}
