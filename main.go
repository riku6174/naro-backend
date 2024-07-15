package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)
type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}
var db *sqlx.DB
func main(){
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	conf := mysql.Config{
		User:      os.Getenv("DB_USERNAME"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_HOSTNAME") + ":" + os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_DATABASE"),
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}
	_db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")
	db = _db
	e := echo.New()
	e.GET("/cities/:cityName", getCityInfoHandler)
	e.POST("/cities", addCityHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")
	log.Println(cityName)
	var  city City
	err := db.Get(&city,"SELECT * FROM city WHERE Name = ?",cityName)
	if errors.Is(err,sql.ErrNoRows){
		return echo.NewHTTPError(http.StatusNotFound,fmt.Sprintf("No such city Name = %s",cityName))
	}
	if err != nil {
		log.Printf("DB Error: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, city)
}
func addCityHandler(c echo.Context) error {
	city := City{}
	err := c.Bind(&city)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,fmt.Sprintf("%+v",err))
	}
	result,err := db.Exec("INSERT INTO city (Name, CountryCode, District, Population) VALUES (?,?,?,?)",city.Name,city.CountryCode,city.District,city.Population)
	if err != nil {
		log.Printf("DB Error:%s",err)
		return echo.NewHTTPError(http.StatusBadRequest,"internal server error")
	}
	id ,_ := result.LastInsertId()
	city.ID = int(id)
	return c.JSON(http.StatusOK,city)
}