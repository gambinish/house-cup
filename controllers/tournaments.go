// Package controllers provides HTTP request handlers (controllers)
// for managing tournaments in the house-cup application.
package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gambinish/house-cup/config"
	"github.com/gambinish/house-cup/models"

	"github.com/gin-gonic/gin"
)

// GetTournaments retrieves a list of all tournaments.
// It queries the database for all tournaments and returns the results in JSON format.
func GetTournaments(c *gin.Context) {
	db := config.ConnectToDB()

	var tournaments []models.Tournament
	rows, err := db.Query("SELECT * FROM Tournaments")
	if err != nil {
		panic(err)
	}
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var tournament models.Tournament
		if err := rows.Scan(&tournament.ID, &tournament.Tournament_Name, &tournament.Created_At, &tournament.Ended_At); err != nil {
			panic(err)
		}
		tournaments = append(tournaments, tournament)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, tournaments)
}

// GetTournamentById retrieves a specific tournament by its ID.
// It takes the tournament ID as a URL parameter, queries the database for the corresponding tournament,
// and returns the result in JSON format.
func GetTournamentById(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(err)
	}

	row := db.QueryRow("SELECT * From Tournaments where id = ?", parsedID)
	log.Print(row, err)
	var tournament models.Tournament
	if err := row.Scan(&tournament.ID, &tournament.Tournament_Name, &tournament.Created_At, &tournament.Ended_At); err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, tournament)
}

// PostTournament creates a new tournament.
// It parses the JSON payload from the request, inserts a new tournament into the database,
// and returns the ID of the newly created tournament in JSON format.
func PostTournament(c *gin.Context) {
	db := config.ConnectToDB()

	var newTournament models.Tournament

	if err := c.BindJSON(&newTournament); err != nil {
		panic(err)
	}

	result, err := db.Exec("INSERT INTO Tournaments (tournament_name, created_at, ended_at) VALUES (?, ?, ?)", newTournament.Tournament_Name, newTournament.Created_At, newTournament.Ended_At)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusCreated, id)
	return
}

// UpdateTournamentById updates a specific tournament by its ID.
// It takes the tournament ID as a URL parameter, parses the JSON payload from the request,
// updates the corresponding tournament in the database, and returns the updated tournament in JSON format.
func UpdateTournamentById(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(err)
	}

	var newTournament models.Tournament

	if err := c.BindJSON(&newTournament); err != nil {
		panic(err)
	}

	result, err := db.Exec("UPDATE Tournaments SET tournament_name = ?, created_at = ?, ended_at = ? WHERE id = ?", newTournament.Tournament_Name, newTournament.Created_At, newTournament.Ended_At, parsedID)

	if err != nil {
		log.Print(result, err)
		panic(err)
	}

	row := db.QueryRow("SELECT * FROM Tournaments where id = ?", parsedID)
	var updatedTournament models.Tournament
	if err := row.Scan(&updatedTournament.ID, &updatedTournament.Tournament_Name, &updatedTournament.Created_At, &updatedTournament.Ended_At); err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, updatedTournament)
	return
}
