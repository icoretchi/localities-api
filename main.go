// Localities API
//
// This is a sample localities API.
//
//	Schemes: http
//  Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Iulian Coretchi <iulian.coretchi@gmail.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

// swagger:parameters localities newLocality
type Locality struct {
	//swagger:ignore
	Code            int    `json:"code"`
	StatisticalCode int    `json:"statisticalCode"`
	Name            string `json:"name"`
	Status          int    `json:"status"`
	ParentCode      int    `json:"parentCode"`
}

var localities []Locality

func init() {
	localities = make([]Locality, 0)
	file, _ := ioutil.ReadFile("localities.json")
	_ = json.Unmarshal([]byte(file), &localities)
}

// swagger:operation POST /localities  newLocality
// Create a new locality
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
func NewLocalityHandler(c *gin.Context) {
	var locality Locality
	if err := c.ShouldBindJSON(&locality); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	localities = append(localities, locality)
	c.JSON(http.StatusOK, locality)
}

// swagger:operation GET /localities  listLocalities
// Returns list of localities
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
func ListLocalitiesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, localities)
}

// swagger:operation PUT /localities/{code} updateLocality
// Update an existing locality
// ---
// parameters:
// - name: code
//   in: path
//   description: Code of the locality
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
//     '404':
//         description: Invalid locality Code
func UpdateLocalityHandler(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var locality Locality
	if err := c.ShouldBindJSON(&locality); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	index := -1
	for i := 0; i < len(localities); i++ {
		if localities[i].Code == code {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Locality not found"})
		return
	}

	localities[index] = locality

	c.JSON(http.StatusOK, locality)
}

// swagger:operation DELETE /localities/{code} deleteLocality
// Delete an existing locality
// ---
// produces:
// - application/json
// parameters:
//   - name: code
//     in: path
//     description: Code of the locality
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid recipe ID
func DeleteLocalityHandler(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	index := -1
	for i := 0; i < len(localities); i++ {
		if localities[i].Code == code {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Locality not found"})
		return
	}

	localities = append(localities[:index], localities[index+1:]...)

	c.JSON(http.StatusOK, gin.H{"message": "Locality has been deleted"})
}

// swagger:operation GET /localities/{code} getLocality
// Get one locality
// ---
// produces:
// - application/json
// parameters:
//   - name: code
//     in: path
//     description: Code of the locality
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid recipe ID
func GetLocalityHandler(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(localities); i++ {
		if localities[i].Code == code {
			c.JSON(http.StatusOK, localities[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Locality not found"})
}

func main() {
	router := gin.Default()
	router.POST("/localities", NewLocalityHandler)
	router.GET("/localities", ListLocalitiesHandler)
	router.PUT("/localities/:code", UpdateLocalityHandler)
	router.DELETE("/localities/:code", DeleteLocalityHandler)
	router.GET("/localities/:code", GetLocalityHandler)
	router.Run()
}
