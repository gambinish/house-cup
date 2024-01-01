// Package controllers provides HTTP request handlers (controllers)
// for managing students in the house-cup application.
package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gambinish/house-cup/config"
	"github.com/gambinish/house-cup/models"

	"github.com/gin-gonic/gin"
)

// GetStudents retrieves a list of all students.
// It queries the database for all students and returns the results in JSON format.
func GetStudents(c *gin.Context) {
	db := config.ConnectToDB()
	// An albums slice to hold data from returned rows.
	var students []models.Student

	rows, err := db.Query("SELECT * FROM Students")

	if err != nil {
		panic(err)
	}
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Student_Name, &student.Points, &student.House_ID); err != nil {
			panic(err)
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, students)
}

// GetStudentById retrieves a specific student by their ID.
// It takes the student ID as a URL parameter, queries the database for the corresponding student,
// and returns the result in JSON format.
func GetStudentById(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(err)
	}

	row := db.QueryRow("SELECT * FROM Students WHERE ID = ?", parsedID)

	var student models.Student
	if err := row.Scan(&student.ID, &student.Student_Name, &student.Points, &student.House_ID); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, student)
}

// UpdateStudentById updates a specific student by their ID.
// It takes the student ID as a URL parameter, parses the JSON payload from the request,
// updates the corresponding student in the database, and returns the updated student in JSON format.
func UpdateStudentById(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	parsedID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		panic(err)
	}

	var newStudent models.Student

	if err := c.BindJSON(&newStudent); err != nil {
		panic(err)
	}
	rows, err := db.Query("UPDATE Students SET student_name = ?, points = ?, house_id = ? WHERE id = ?", newStudent.Student_Name, newStudent.Points, newStudent.House_ID, parsedID)

	if err != nil {
		panic(err)
	}

	row := db.QueryRow("SELECT * FROM Students where id = ?", parsedID)
	var updatedStudent models.Student
	if err := row.Scan(&updatedStudent.ID, &updatedStudent.Student_Name, &updatedStudent.Points, &updatedStudent.House_ID); err != nil {
		panic(err)
	}
	defer rows.Close()

	log.Print(updatedStudent)
	c.IndentedJSON(http.StatusOK, updatedStudent)
}

// GetStudentsByHouseId retrieves a list of students belonging to a specific house.
// It takes the house ID as a URL parameter, queries the database for the corresponding students,
// and returns the results in JSON format.
func GetStudentsByHouseId(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(err)
	}
	// An albums slice to hold data from returned rows.
	var students []models.Student

	rows, err := db.Query("SELECT * FROM Students WHERE House_ID = ?", parsedID)

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Student_Name, &student.Points, &student.House_ID); err != nil {
			panic(err)
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, students)
}

// GetStudentsByTournamentId retrieves a list of students associated with a specific tournament.
// It takes the tournament ID as a URL parameter, queries the database for students in houses
// associated with the tournament, and returns the results in JSON format.
func GetStudentsByTournamentId(c *gin.Context) {
	db := config.ConnectToDB()
	id := c.Param("id")

	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(err)
	}
	// An albums slice to hold data from returned rows.
	var students []models.Student

	rows, err := db.Query(`SELECT Students.*
							FROM Students
							JOIN Houses ON Students.House_ID = Houses.ID
							WHERE Houses.Tournament_ID = ?`, parsedID)

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Student_Name, &student.Points, &student.House_ID); err != nil {
			panic(err)
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, students)
}

// PostStudent creates a new student record.
// It parses the JSON payload from the request, inserts a new student record into the database,
// and returns the created student in JSON format.
func PostStudent(c *gin.Context) {
	db := config.ConnectToDB()

	var newStudent models.Student

	if err := c.BindJSON(&newStudent); err != nil {
		panic(err)
	}

	rows, err := db.Query("INSERT INTO Students (student_name, points, house_id) VALUES (?, 0, ?)", newStudent.Student_Name, newStudent.House_ID)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	c.IndentedJSON(http.StatusOK, newStudent)
}

// PostStudents creates multiple new student records in a single transaction.
// It parses the JSON payload from the request, starts a transaction, iterates over the array of students,
// executes the INSERT statement for each one, and commits the transaction if all INSERTs are successful.
// Returns the created students in JSON format.
func PostStudents(c *gin.Context) {
	db := config.ConnectToDB()

	var newStudents []models.Student

	if err := c.BindJSON(&newStudents); err != nil {
		panic(err)
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	// Iterate over the array of students and execute the INSERT statement for each one
	for _, student := range newStudents {
		_, err := tx.Exec("INSERT INTO Students (student_name, points, house_id) VALUES (?, 0, ?)", student.Student_Name, student.House_ID)

		if err != nil {
			// Rollback the transaction in case of an error
			tx.Rollback()
			panic(err)
		}
	}

	// Commit the transaction if all INSERTs are successful
	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, newStudents)
}

// DeleteStudentById deletes a student and associated points by their ID.
// It retrieves the student ID from the request parameter, starts a database transaction,
// deletes the points associated with the student, queries for the house ID and points associated with the student,
// deletes the student record, updates House_Points by subtracting the student's points, and commits the transaction
// if all operations are successful. It responds with a JSON message indicating the success of the deletion.
func DeleteStudentById(c *gin.Context) {
	// Connect to the database
	db := config.ConnectToDB()

	// Get student ID from the request parameter
	studentID := c.Param("id")

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	// Delete points associated with the student
	_, err = tx.Exec("DELETE FROM Points WHERE Student_ID = ?", studentID)
	if err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		panic(err)
	}

	// Query for the house ID and points associated with the student
	var houseID, studentPoints int
	err = tx.QueryRow("SELECT House_ID, Points FROM Students WHERE ID = ?", studentID).Scan(&houseID, &studentPoints)
	if err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		panic(err)
	}

	// Delete the student
	_, err = tx.Exec("DELETE FROM Students WHERE ID = ?", studentID)
	if err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		panic(err)
	}

	// Update House_Points by subtracting the studentPoints
	_, err = tx.Exec("UPDATE Houses SET House_Points = House_Points - ? WHERE ID = ?", studentPoints, houseID)
	if err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		panic(err)
	}

	// Commit the transaction if all operations are successful
	err = tx.Commit()
	if err != nil {
		log.Print(err)
		panic(err)
	}

	// Respond with a JSON message indicating success
	c.JSON(http.StatusOK, gin.H{"message": "Student and associated points deleted successfully"})
	return
}
