package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gambinish/house-cup/controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func sanityCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"hello": "world"})
}

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")

	if dbHost == "" {
		log.Print("DB_HOST environment variable is not set.")
		// Handle the error appropriately, e.g., log and exit the application.
		return
	}

	// Use dbHost to connect to the database.
	log.Print("Connecting to database at: ", dbHost)

	// config.ConnectToDB()

	router := gin.Default()
	router.GET("/", sanityCheck)

	log.Print(os.Getenv("API_HOST"))
	// tournament routes
	router.GET("/tournaments", controllers.GetTournaments)
	router.GET("/tournaments/:id", controllers.GetTournamentById)
	router.POST("/tournaments", controllers.PostTournament)
	router.PUT("/tournaments/:id", controllers.UpdateTournamentById)

	// houses routes
	router.GET("/houses", controllers.GetHouses)
	router.GET("/houses/:id", controllers.GetHousesByTournamentId)
	router.POST("/houses/:id", controllers.PostHouseByTournamentId)
	router.PUT("/houses/:id", controllers.UpdateHouseById)

	// student routes
	router.GET("/students", controllers.GetStudents)
	router.GET("/students/:id", controllers.GetStudentById)
	router.GET("/students/tournament/:id", controllers.GetStudentsByTournamentId)
	router.GET("/students/house/:id", controllers.GetStudentsByHouseId)
	router.PUT("/students/:id", controllers.UpdateStudentById)
	router.POST("/student", controllers.PostStudent)
	router.POST("/students", controllers.PostStudents)
	router.DELETE("/students/:id", controllers.DeleteStudentById)

	// points routes
	router.GET("/points", controllers.GetPoints)
	router.GET("/points/:id", controllers.GetPointsByStudentId)
	router.GET("/points/house/:id", controllers.GetPointsByHouseId)
	router.POST("/points", controllers.PostPoints)

	apiHost := os.Getenv("API_HOST")
	apiPort := os.Getenv("API_PORT")
	apiAddr := apiHost + ":" + apiPort

	log.Print("API ADDRESS: ", apiAddr)
	router.Run(apiAddr)
}
