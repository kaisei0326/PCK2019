package main

import (
	"database/sql"
	"log"

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
	myScore(r, db)
	// rank(r, db)
	defer db.Close()
}

func sqlInit() *sql.DB {
	db, err := sql.Open("mysql", "root:test@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatal("SQL open error.")
	}
	return db
}

func myScore(r *gin.Engine, db *sql.DB) {
	userResult := getSQL(db)
	//name := r.Query("name")

	for rank, u := range userResult {
		r.GET("/"+u.name, func(c *gin.Context) {
			c.JSON(200, gin.H{
				"rank":  rank,
				"score": u.score,
			})
		})
	}
	r.Run(":8080")
}

/*
func mapData {

}
*/
func rank(r *gin.Engine, db *sql.DB) {
	userResult := getSQL(db)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name_1st":  userResult[0].name,
			"score_1st": userResult[0].score,
			"name_2nd":  userResult[1].name,
			"score_2nd": userResult[1].score,
			"name_3rd":  userResult[2].name,
			"score_3rd": userResult[2].score,
		})
	})
	r.Run(":8080")
}

func getSQL(db *sql.DB) []User {
	rows, err := db.Query("select * from test.userdata order by score desc")
	if err != nil {
		log.Fatal("SQL fetch error.")
	}

	var userResult []User
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.name, &user.score); err != nil {
			log.Fatal("rows fetch error.")
		}
		userResult = append(userResult, user)
	}
	return (userResult)
}

/*
for _, u := range userResult {
	// ここでQuery呼び出す必要あるけど、r.GETの中でしかQuery呼び出せない
	//forをgetの中に入れる...？
	if key = u.name {
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"rank":  rank,
				"score": u.score,
			})
		})
		r.Run(":8080")
	}
}
*/
