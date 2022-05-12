package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Locality struct {
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

func ListLocalitiesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, localities)
}

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

func main() {
	router := gin.Default()
	router.POST("/localities", NewLocalityHandler)
	router.GET("/localities", ListLocalitiesHandler)
	router.PUT("/localities/:code", UpdateLocalityHandler)
	router.Run()
}
