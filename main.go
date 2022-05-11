package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.Run()
}

type Locality struct {
	Code            int    `json:"code"`
	StatisticalCode int    `json:"statisticalCode"`
	Name            string `json:"name"`
	Status          int    `json:"status"`
	ParentCode      int    `json:"parentCode"`
}
