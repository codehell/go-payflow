package main

import (
	"context"
	"encoding/base64"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"google.golang.org/api/iterator"
)

// ProjectID name of the project
const ProjectID = "go-payflow"

var ctx = context.Background()

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	options := sessions.Options{HttpOnly: true}
	store.Options(options)
	r.Use(sessions.Sessions("mysession", store))

	r.Use(csrf.Middleware(csrf.Options{
		Secret: "secret123",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	r.GET("/decode/:token", decodeToken)
	r.GET("/restaurants", getRestaurants)
	r.GET("/tests", tests)
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

func tests(c *gin.Context) {
	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v == nil {
		count = 0
	} else {
		count = v.(int)
		count++
	}
	session.Set("count", count)
	session.Save()
	c.JSON(200, gin.H{"count": count})
}

func decodeToken(c *gin.Context) {
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
}
