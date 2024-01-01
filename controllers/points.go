// Package controllers provides HTTP request handlers (controllers)
// for managing points in the house-cup application.
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

// GetPoints retrieves a list of all points records.
// It queries the database for all points and returns the results in JSON format.
func GetPoints(c *gin.Context) {
	db := config.ConnectToDB()
	// An albums slice to hold data from returned rows.
	var points []models.Point

	rows, err := db.Query("SELECT * FROM Points")

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getPoints: %v", err)})
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var point models.Point
		if err := rows.Scan(&point.ID, &point.Points, &point.Notes, &point.Student_ID); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getPoints: %v", err)})
		}
		points = append(points, point)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getPoints: %v", err)})
	}
	c.IndentedJSON(http.StatusOK, points)
	return
}

// GetPointsByStudentId retrieves a list of points records associated with a specific student.
// It takes the student ID as a URL parameter, queries the database for the corresponding points,
// and returns the results along with the total points in JSON format.
func GetPointsByStudentId(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	var points []models.Point
	var total int64 = 0

	parsedID, err := strconv.ParseInt(id, 10, 64)

	rows, err := db.Query("SELECT * FROM Points WHERE Student_ID = ?", parsedID)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getPoints: %v", err)})
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var point models.Point

		if err := rows.Scan(&point.ID, &point.Points, &point.Notes, &point.Student_ID); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getPoints: %v", err)})
		}
		total += point.Points
		points = append(points, point)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getPoints: %v", err)})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"points": points, "total": total})
	return
}

// GetPointsByHouseId retrieves a list of points records associated with a specific house.
// It takes the house ID as a URL parameter, queries the database for the corresponding points,
// and returns the results along with the total points in JSON format.
func GetPointsByHouseId(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	var points []models.Point
	var total int64 = 0

	parsedID, err := strconv.ParseInt(id, 10, 64)

	rows, err := db.Query(`SELECT Points.*
							FROM Points
							JOIN Students ON Points.Student_ID = Students.ID
							WHERE Students.House_ID = ?`, parsedID)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getPoints: %v", err)})
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var point models.Point

		if err := rows.Scan(&point.ID, &point.Points, &point.Notes, &point.Student_ID); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getPoints: %v", err)})
		}
		total += point.Points
		points = append(points, point)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("getPoints: %v", err)})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"points": points, "total": total})
	return
}

// PostPoints creates a new points record.
// It parses the JSON payload from the request, inserts a new points record into the database,
// and updates the corresponding student and house points in the database.
func PostPoints(c *gin.Context) {
	db := config.ConnectToDB()

	var newPoints models.Point

	if err := c.BindJSON(&newPoints); err != nil {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		// Handle error
		return
	}

	_, err = tx.Exec(`
		INSERT INTO Points (Points, Notes, Student_ID, House_ID) VALUES (?, ?, ?, ?)
	`, newPoints.Points, newPoints.Notes, newPoints.Student_ID, newPoints.House_ID)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return
	}

	// Update the Points row in the Students table by incrementing it
	_, err = tx.Exec(`UPDATE Students SET Points = Points + ? WHERE ID = ?;`, newPoints.Points, newPoints.Student_ID)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return
	}

	// Update the House_Points in the Houses table by incrementing it
	_, err = tx.Exec(`UPDATE Houses
					    SET House_Points = House_Points + ?
						WHERE ID = ?;`, newPoints.Points, newPoints.House_ID)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"succes": true})
	return
}
