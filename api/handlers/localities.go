package handlers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"localities-api/api/models"
	"log"
	"net/http"
	"strconv"
)

type LocalitiesHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewLocalitiesHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *LocalitiesHandler {
	return &LocalitiesHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

// swagger:operation GET /localities localities listLocalities
// Returns list of localities
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
func (handler *LocalitiesHandler) ListLocalitiesHandler(c *gin.Context) {
	val, err := handler.redisClient.Get("localities").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(handler.ctx)

		localities := make([]models.Locality, 0)
		for cur.Next(handler.ctx) {
			var locality models.Locality
			cur.Decode(&locality)
			localities = append(localities, locality)
		}

		data, _ := json.Marshal(localities)
		handler.redisClient.Set("localities", string(data), 0)
		c.JSON(http.StatusOK, localities)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		localities := make([]models.Locality, 0)
		json.Unmarshal([]byte(val), &localities)
		c.JSON(http.StatusOK, localities)
	}
}

// swagger:operation POST /localities newLocality
// Create a new locality
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
func (handler *LocalitiesHandler) NewLocalityHandler(c *gin.Context) {
	var locality models.Locality
	if err := c.ShouldBindJSON(&locality); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := handler.collection.InsertOne(handler.ctx, locality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new locality"})
		return
	}

	log.Println("Remove data from Redis")
	handler.redisClient.Del("localities")

	c.JSON(http.StatusOK, locality)
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
func (handler *LocalitiesHandler) UpdateLocalityHandler(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var locality models.Locality
	if err := c.ShouldBindJSON(&locality); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = handler.collection.UpdateOne(handler.ctx, bson.M{
		"code": code,
	}, bson.D{{"$set", bson.D{
		{"code", locality.Code},
		{"statisticalcode", locality.StatisticalCode},
		{"name", locality.Name},
		{"status", locality.Status},
		{"parentcode", locality.ParentCode},
	}}})

	log.Println("Remove data from Redis")
	handler.redisClient.Del("localities")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Locality has been updated"})
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
//         description: Invalid locality Code
func (handler *LocalitiesHandler) DeleteLocalityHandler(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = handler.collection.DeleteOne(handler.ctx, bson.M{
		"code": code,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Locality has been deleted"})
}

// swagger:operation GET /localities/{code} localities
// Get one locality
// ---
// produces:
// - application/json
// parameters:
//   - name: code
//     in: path
//     description: locality Code
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
func (handler *LocalitiesHandler) GetOneLocalityHandler(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"code": code,
	})
	var locality models.Locality
	err = cur.Decode(&locality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, locality)
}
