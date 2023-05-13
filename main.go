package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type SensorData struct {
	TDS       float64   `json:"tds"`
	Moisture  float64   `json:"moisture"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func main() {
	router := gin.Default()

	opt := option.WithCredentialsFile("serviceAccount.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize Firebase client: %v", client)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Success fetching the API!", "statusCode": 200})
	})

	router.POST("/data", func(c *gin.Context) {
		handleSensorData(c, client)
	})

	router.Run(":8080")
}

func handleSensorData(c *gin.Context, client *firestore.Client) {
	var data SensorData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Received sensor data: ")
	fmt.Println("TDS: ", data.TDS)
	fmt.Println("Moisture: ", data.Moisture)

	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now

	userID := "s8uiODzNahcLoMwc27LNH2EkWTt1"

	// Create a reference to the user's document in the "users" collection
	userRef := client.Collection("users").Doc(userID)

	// Create a reference to the "devices" subcollection within the user's document
	devicesRef := userRef.Collection("devices")

	// Store the received data in the "devices" subcollection
	docRef, _, err := devicesRef.Add(context.Background(), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store data"})
		return
	}

	fmt.Println("Data stored in Firestore with ID:", docRef.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Data received successfully"})
}
