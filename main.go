package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type Data struct {
	Moisture  float64 `json:"moisture"`
	TDS       float64 `json:"tds"`
	PH        float64 `json:"pH"`
	UpdatedAt string  `json:"updatedAt"`
	UserID    string  `json:"user_id"`
	DeviceID  int     `json:"device_id"`
}

func main() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "demoandroidmarketplace", option.WithCredentialsFile("serviceAccount.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	r.POST("/data", func(c *gin.Context) {
		var data Data
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		data.UpdatedAt = time.Now().UTC().String()

		filteredData := Data{
			Moisture:  data.Moisture,
			TDS:       data.TDS,
			PH:        data.PH,
			UpdatedAt: data.UpdatedAt,
		}

		_, err := client.Collection("users").Doc(data.UserID).Collection("devices").Doc(strconv.Itoa(data.DeviceID)).Set(ctx, filteredData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Data received and stored successfully"})
	})

	r.Run()
}
