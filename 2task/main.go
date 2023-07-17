package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"math"
	"net/http"
)

type Spot struct {
	Id          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Website     sql.NullString `json:"website"`
	Coordinates string         `json:"coordinates"`
	Description sql.NullString `json:"description"`
	Rating      float64        `json:"rating"`
}

type Distance struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
	Type      string  `json:"type"`
}

func main() {
	router := gin.Default()

	router.GET("/", SpotsArea)

	err := router.Run(":8080")
	if err != nil {
		return
	}

}

// SpotsArea Handler
func SpotsArea(c *gin.Context) {
	var data *Distance

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"error": "incorrect data"})
		return
	}

	spotsRows, err := getAllSpotsWithParam(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, spotsRows)
}

// getAllSpotsWithParam function sends requests to the DB
func getAllSpotsWithParam(data *Distance) (*[]Spot, error) {
	db, err := sql.Open("postgres", "user=postgres password=123 port=5432 dbname=spots sslmode=disable")
	if err != nil {
		return &[]Spot{}, err
	}
	defer db.Close()

	var query string
	point := degreeMath(data)

	switch data.Type {
	case "square":
		query = fmt.Sprintf(FIND_BY_SQUARE, point.minLong, point.minLat, point.maxLong, point.maxLat, data.Longitude, data.Latitude, data.Longitude, data.Latitude)
	case "circle":
		query = fmt.Sprintf(FIND_BY_CIRCLE, data.Longitude, data.Latitude, data.Radius, data.Longitude, data.Latitude, data.Longitude, data.Latitude)
	}

	rows, err := db.Query(query)
	if err != nil {
		return &[]Spot{}, err
	}
	defer rows.Close()

	var spotsRows []Spot
	var spot Spot
	for rows.Next() {
		if err := rows.Scan(&spot.Id, &spot.Name, &spot.Website, &spot.Coordinates, &spot.Description, &spot.Rating); err != nil {
			println("db Next error", err.Error())
			return &[]Spot{}, err
		}
		spotsRows = append(spotsRows, spot)
	}
	println(len(spotsRows))

	return &spotsRows, nil
}

// CoordinatesWithOffset struct contains points for ST_MakeEnvelope function
type CoordinatesWithOffset struct {
	minLong float64
	minLat  float64
	maxLong float64
	maxLat  float64
}

// degreeMath function returns 4 points for postGis function ST_MakeEnvelope with radius offset
func degreeMath(d *Distance) *CoordinatesWithOffset {
	var R float64 = 6378137 //Earth’s radius, sphere

	dLat := d.Radius / R
	dLon := d.Radius / (R * math.Cos(math.Pi*d.Longitude/180))

	result := CoordinatesWithOffset{
		minLong: d.Longitude - dLon*180/math.Pi,
		minLat:  d.Latitude - dLat*180/math.Pi,
		maxLong: d.Longitude + dLon*180/math.Pi,
		maxLat:  d.Latitude + dLat*180/math.Pi,
	}
	return &result
}
