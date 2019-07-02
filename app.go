package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"log"
)

// ProjectID name of the project
const ProjectID = "bookshelf-2019-3"

var ctx = context.Background()

func main() {
	r := gin.Default()
	r.GET("/decode/:token", func(c *gin.Context) {
		encoded := c.Params.ByName("token")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(200, gin.H{
			"result": string(decoded),
		})
	})
	r.GET("/restaurants", getRestaurants)
	_ = r.Run()
}

func getRestaurants(c *gin.Context) {
	client, err := firestore.NewClient(ctx, ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	var arry []map[string]interface{}
	col := client.Collection("restaurants")
	iter := col.Documents(ctx)
	for {
		dr, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if dr != nil {
			arry = append(arry, dr.Data())
		}
	}
	c.JSON(200, arry)
}
