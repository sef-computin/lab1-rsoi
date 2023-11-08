package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getPersons(c *gin.Context) {

	db, err := connect_to_db()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	defer db.Close()

	var persons []person

	fmt.Println(db)

	rows, err := db.Query(`SELECT * FROM persons.persons`)
	if err != nil {
		// fmt.Println("\n\n[][]", err)
		c.IndentedJSON(http.StatusInternalServerError, persons)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var man person

		err = rows.Scan(&man.ID, &man.Name, &man.Age, &man.Address, &man.Work)
		if err != nil {
			// fmt.Println("\n\n[][] ", err)
			c.IndentedJSON(http.StatusInternalServerError, persons)
			return
		}

		persons = append(persons, man)
	}

	c.IndentedJSON(http.StatusOK, persons)
}

func postPerson(c *gin.Context) {

	db, err := connect_to_db()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	defer db.Close()

	var newPerson person

	err = c.BindJSON(&newPerson)
	if err == nil {

		newID := 0

		err = db.QueryRow("insert into persons.persons (Name, Age, Address, Work) values($1, $2, $3, $4) Returning id", newPerson.Name, newPerson.Age, newPerson.Address, newPerson.Work).Scan(&newID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, nil)
		}

		newPerson.ID = int(newID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		c.Writer.Header().Set("Location", fmt.Sprintf("/api/v1/persons/%d", int(newID)))

		c.IndentedJSON(http.StatusCreated, newPerson)
		return
	}

	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})

}

func getPersonByID(c *gin.Context) {

	db, err := connect_to_db()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	defer db.Close()

	idparam := c.Param("id")

	var res []person

	id, err := strconv.Atoi(idparam)

	if err == nil {

		rows, err := db.Query(fmt.Sprintf("SELECT Id, Name, Age, Address, Work FROM persons.persons WHERE Id = %d", id))
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
			return
		}

		defer rows.Close()
		for rows.Next() {
			var man person
			err = rows.Scan(&man.ID, &man.Name, &man.Age, &man.Address, &man.Work)
			if err != nil {
				fmt.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, nil)
				return
			}
			res = append(res, man)
		}

		if len(res) == 1 {
			c.IndentedJSON(http.StatusOK, res[0])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "person not found"})
}

func deletePerson(c *gin.Context) {

	db, err := connect_to_db()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	defer db.Close()

	idparam := c.Param("id")

	id, err := strconv.Atoi(idparam)
	if err == nil {
		insertDynStmt := `delete from persons.persons where id=$1`
		res, _ := db.Exec(insertDynStmt, id)
		rowsdeleted, err := res.RowsAffected()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, nil)
		}

		if rowsdeleted > 0 {
			c.IndentedJSON(http.StatusNoContent, "Successfully deleted")
		} else {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "person not found"})
		}
	}
}

func editPerson(c *gin.Context) {

	db, err := connect_to_db()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	defer db.Close()

	idparam := c.Param("id")

	id, err := strconv.Atoi(idparam)

	var temp person

	if err == nil {

		if err := c.BindJSON(&temp); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "body format error"})
			return
		}

		temp.ID = id
		updateStmt := "update persons.persons set "
		flagfirst := true

		if temp.Name != "" {
			if flagfirst {
				updateStmt = updateStmt + fmt.Sprintf(`name='%v'`, temp.Name)
				flagfirst = false
			}
		}
		if temp.Age != 0 {
			if flagfirst {
				updateStmt = updateStmt + fmt.Sprintf(`age=%d`, temp.Age)
				flagfirst = false
			} else {
				updateStmt = updateStmt + fmt.Sprintf(`, age=%d`, temp.Age)
			}

		}
		if temp.Address != "" {
			if flagfirst {
				updateStmt = updateStmt + fmt.Sprintf(`address='%v'`, temp.Address)
				flagfirst = false
			} else {
				updateStmt = updateStmt + fmt.Sprintf(`, address='%v'`, temp.Address)
			}

		}
		if temp.Work != "" {
			if flagfirst {
				updateStmt = updateStmt + fmt.Sprintf(`work='%v'`, temp.Work)
				flagfirst = false
			} else {
				updateStmt = updateStmt + fmt.Sprintf(`, work='%v'`, temp.Work)
			}
		}

		updateStmt = updateStmt + " where id=$1 returning Id, Name, Age, Address, Work"

		err := db.QueryRow(updateStmt, id).Scan(&temp.ID, &temp.Name, &temp.Age, &temp.Address, &temp.Work)
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "person not found"})
			return
		}

		c.IndentedJSON(http.StatusOK, temp)

		// for i, a := range persons {
		// 	if a.ID == id {

		// 		c.IndentedJSON(http.StatusOK, persons[i])
		// 		return
		// }
		// }
	}

}
