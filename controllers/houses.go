// Package controllers provides HTTP request handlers (controllers)
// for managing houses in the house-cup application.
package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gambinish/house-cup/config"
	"github.com/gambinish/house-cup/models"

	"github.com/gin-gonic/gin"
)

// GetHouses retrieves a list of all houses.
// It queries the database for all houses and returns the results in JSON format.
func GetHouses(c *gin.Context) {
	db := config.ConnectToDB()
	// An albums slice to hold data from returned rows.
	var houses []models.House

	rows, err := db.Query("SELECT * FROM Houses")

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getHouses: %v", err)})
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var house models.House
		if err := rows.Scan(&house.ID, &house.House_Name, &house.House_Points, &house.Tournament_ID); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getHouses: %v", err)})
		}
		houses = append(houses, house)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getHouses: %v", err)})
	}
	c.IndentedJSON(http.StatusOK, houses)
	return
}

// GetHousesByTournamentId retrieves a list of houses associated with a specific tournament.
// It takes the tournament ID as a URL parameter, queries the database for the corresponding houses,
// and returns the results in JSON format.
func GetHousesByTournamentId(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Print(err)
	}
	// An albums slice to hold data from returned rows.
	var houses []models.House

	rows, err := db.Query("SELECT * FROM Houses WHERE Tournament_ID = ?", parsedID)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getHouses: %v", err)})
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var house models.House
		if err := rows.Scan(&house.ID, &house.House_Name, &house.House_Points, &house.Tournament_ID); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getHouses: %v", err)})
		}
		houses = append(houses, house)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getHouses: %v", err)})
	}
	c.IndentedJSON(http.StatusOK, houses)
	return
}

// PostHouseByTournamentId creates a new house associated with a specific tournament.
// It takes the tournament ID as a URL parameter, parses the JSON payload from the request,
// inserts a new house record into the database, and returns the created house in JSON format.
func PostHouseByTournamentId(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	parsedID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		log.Print(err)
	}

	var newHouse models.House

	if err := c.BindJSON(&newHouse); err != nil {
		return
	}

	rows, err := db.Query("INSERT INTO Houses (house_name, house_points, tournament_id) VALUES (?, 0, ?)", newHouse.House_Name, parsedID)

	if err != nil {
		log.Print(err)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, newHouse)
	return
}

// UpdateHouseById updates a specific house by its ID.
// It takes the house ID as a URL parameter, parses the JSON payload from the request,
// updates the corresponding house in the database, and returns the updated house in JSON format.
func UpdateHouseById(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	parsedID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		log.Print(err)
	}

	var newHouse models.House

	if err := c.BindJSON(&newHouse); err != nil {
		log.Print("ERROR: ", err)
		return
	}
	rows, err := db.Query("UPDATE Houses SET house_name = ?, house_points = ?, tournament_id = ? WHERE id = ?", newHouse.House_Name, newHouse.House_Points, newHouse.Tournament_ID, parsedID)

	if err != nil {
		log.Print(err)
	}

	row := db.QueryRow("SELECT * FROM Houses where id = ?", parsedID)
	var updatedHouse models.House
	if err := row.Scan(&updatedHouse.ID, &updatedHouse.House_Name, &updatedHouse.House_Points, &updatedHouse.Tournament_ID); err != nil {
		log.Print(err)
	}
	defer rows.Close()

	log.Print(updatedHouse)
	c.IndentedJSON(http.StatusOK, updatedHouse)
	return
}
