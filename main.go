package main

import (
	"fmt"

	"database/sql"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

// album represents data about a record album.
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

var db *sql.DB

func main() {
    router := gin.Default()
    router.GET("/albums", getAlbums)
    router.POST("/albums", postAlbum)
    var err error
    db, err = sql.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/goprac")

    if err != nil {
        panic(err.Error())
    }
    router.Run("localhost:8080")

    defer db.Close()
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
    var response []album = make([]album, 0, 5)
    results, err := db.Query("SELECT ID, ARTIST, TITLE, PRICE FROM ALBUM")
    if err != nil {
        panic(err.Error())
    }
    for results.Next() {
        var albumResponse album
        err = results.Scan(&albumResponse.ID, &albumResponse.Artist, &albumResponse.Title, &albumResponse.Price)
        if err != nil {
            panic(err.Error())
        }
        response = append(response, albumResponse)
    }
    c.JSON(201, response)
}

// Add an album to the database
func postAlbum(c *gin.Context) {
    var newAlbum *album

    if err := c.BindJSON(&newAlbum); err != nil {
        panic(err.Error())
    }

    if insert, _ := db.Exec(
        `INSERT INTO ALBUM (TITLE, ARTIST, PRICE) VALUES (?, ?, ?)`, 
        newAlbum.Title, 
        newAlbum.Artist, 
        newAlbum.Price,
    ); insert != nil {
        if _, err := insert.LastInsertId(); err != nil {
            c.JSON(400, "insert failed")
        } else {
            c.JSON(201, newAlbum)
        }
    } else {
        c.JSON(400, "insert failed")
    }
    fmt.Println(newAlbum)
}