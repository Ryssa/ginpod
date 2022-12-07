package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"io/ioutil"
	"github.com/gin-gonic/gin"
)

const (
	API_URL = "https://online-movie-database.p.rapidapi.com"
	API_HOST = "online-movie-database.p.rapidapi.com"
	API_KEY = "7d3d096792msh5c3e9c52399075ep1c90e2jsn40a17c6879dd"
)

// HelloHTTP is an HTTP Cloud Function with a request parameter.
// Copied from https://cloud.google.com/functions/docs/create-deploy-http-go
func HelloHTTP(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Hello, World!")
		return
	}

	if d.Name == "" {
		fmt.Fprint(w, "Hello, World!")
		return
	}

	fmt.Fprintf(w, "Hello, %s!", html.EscapeString(d.Name))
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1", "localhost"})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	router.GET("/movies/search", func(c *gin.Context) {
		title := c.Query("title")

		url := API_URL + "/auto-complete?q=" + title
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("X-RapidAPI-Key", API_KEY)
		req.Header.Add("X-RapidAPI-Host", API_HOST)

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		var data interface{}
		_ = json.Unmarshal(body, &data)

		c.JSON(200, gin.H{
			"result": data,
		})
	})

	router.Run(":8080")
}
