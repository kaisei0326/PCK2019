package main

import (
	"database/sql"
	"log"

	//"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

//latは緯度, lngは経度

//User is structure fo storing userdata
type User struct {
	name  string
	score int
}

func main() {
	r := gin.Default()
	db := sqlInit()
	rank(r, db)
	defer db.Close()
}

func sqlInit() *sql.DB {
	db, err := sql.Open("mysql", "root:test@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatal("db error.")
	}
	return db
}

func rank(r *gin.Engine, db *sql.DB) {
	userResult := getSQL(db)
	for i := 1; i <= 3; i++ {
		//name := "name" + strconv.Itoa(i)
		//score := "score" + strconv.Itoa(i)
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"name":  userResult[i-1].name,
				"score": userResult[i-1].score,
			})
		})
		r.Run(":8080")
	}
	/*
		for _, u := range userResult {
			r.GET("/", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"name": u.name,
					"score": u.score,
				})
			})
			r.Run(":8080")
		}
	*/
}

func getSQL(db *sql.DB) []User {
	rows, err := db.Query("select * from test.userdata order by score desc")
	if err != nil {
		log.Fatal("db error.")
	}

	var userResult []User
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.name, &user.score); err != nil {
			log.Fatal("db error.")
		}
		userResult = append(userResult, user)
	}
	return (userResult)
}
