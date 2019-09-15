package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

//User is structure for storing userdata
type User struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

//Map is structure for storing mapdata
type Map struct {
	//ピクトグラムのID
	pictID int
	//日付
	date int
	//時間
	time int
	//X座標
	point_x float64
	//Y座標
	point_y float64
}

func main() {
	r := gin.Default()
	db := sqlInit()
	postRank(r, db)
	updateScore(r, db)
	rank(r, db)
	//mapData
	r.Run(":8080")
	defer db.Close()
}

func postRank(r *gin.Engine, db *sql.DB) {
	r.GET("/myrank", func(c *gin.Context) {
		userResult := getUserdata(db)
		name := c.Query("name")
		rank := 0
		stock := 0
		forwardscore := -1
		for _, u := range userResult {
			if forwardscore != u.Score {
				rank++
				rank += stock
				stock = 0
				forwardscore = u.Score
			} else {
				stock++
			}
			if name == u.Name {
				c.JSON(200, gin.H{
					"rank":  rank,
					"score": u.Score,
				})
			}
		}
	})
}

func updateScore(r *gin.Engine, db *sql.DB) {
	upd, err := db.Prepare("UPDATE userdata set score=? where name=?")
	if err != nil {
		log.Fatal(err)
	}
	r.POST("/update", func(c *gin.Context) {
		var newscore User
		c.BindJSON(&newscore)
		upd.Exec(newscore.Score, newscore.Name)
	})
}

func rank(r *gin.Engine, db *sql.DB) {
	r.GET("/rank", func(c *gin.Context) {
		userResult := getUserdata(db)
		c.JSON(200, gin.H{
			"name_1st":  userResult[0].Name,
			"name_2nd":  userResult[1].Name,
			"name_3rd":  userResult[2].Name,
			"score_1st": userResult[0].Score,
			"score_2nd": userResult[1].Score,
			"score_3rd": userResult[2].Score,
		})
	})
	r.GET("/rank/1st", func(c *gin.Context) {
		userResult := getUserdata(db)
		c.JSON(200, gin.H{
			"name":  userResult[0].Name,
			"score": userResult[0].Score,
		})
	})
	r.GET("/rank/2nd", func(c *gin.Context) {
		userResult := getUserdata(db)
		c.JSON(200, gin.H{
			"name":  userResult[1].Name,
			"score": userResult[1].Score,
		})
	})
	r.GET("/rank/3rd", func(c *gin.Context) {
		userResult := getUserdata(db)
		c.JSON(200, gin.H{
			"name":  userResult[2].Name,
			"score": userResult[2].Score,
		})
	})
}

func postMapData(r *gin.Engine, db *sql.DB) {

}

func sqlInit() *sql.DB {
	db, err := sql.Open("mysql", "root:test@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatal("SQL open error.")
	}
	return db
}

func getUserdata(db *sql.DB) []User {
	rows, err := db.Query("select * from test.userdata order by score desc")
	if err != nil {
		log.Fatal("SQL fetch error.")
	}

	var userResult []User
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.Name, &user.Score); err != nil {
			log.Fatal("rows fetch error.")
		}
		userResult = append(userResult, user)
	}
	return (userResult)
}

func changeMapdata(db *sql.DB) *sql.Rows {
	rows, err := db.Query("select * from map.data order by score desc")
	if err != nil {
		log.Fatal("SQL fetch error.")
	}
	return (rows)
}

func getMapdata(db *sql.DB, name string) []User {
	rows, err := db.Query("select * from map." + name)
	if err != nil {
		log.Fatal("SQL fetch error.")
	}

	var userResult []User
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.Name, &user.Score); err != nil {
			log.Fatal("rows fetch error.")
		}
		userResult = append(userResult, user)
	}
	return (userResult)
}
