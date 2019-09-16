package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

// User の情報をjson形式で通信するときに使う
type User struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// Map の情報をjson形式で通信するときに使う
type Map struct {
	PictID int     `json:"pictID"`
	Date   int     `json:"date"`
	Time   int     `json:"time"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Name   string  `json:"name"`
}

func main() {
	r := gin.Default()
	db := sqlInit()
	ranking(r, db)
	//mapCollection(r, db)
	r.Run(":8080")
	defer db.Close()
}

func sqlInit() *sql.DB {
	db, err := sql.Open("mysql", "root:test@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatal("SQL open error.")
	}
	return db
}

func ranking(r *gin.Engine, db *sql.DB) {
	userResult := getUserdata(db)
	// 名前を入力して順位を調べる
	r.GET("/ranking/search", func(c *gin.Context) {
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
					"name":  name,
					"rank":  rank,
					"score": u.Score,
				})
			}
		}
	})
	// 新しくユーザーを作成する
	ins, err := db.Prepare("INSERT INTO userdata(name, score) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	r.POST("/ranking/insert", func(c *gin.Context) {
		var newname User
		c.BindJSON(&newname)
		ins.Exec(newname.Name, 0)
		userResult = getUserdata(db)
	})
	// POSTされたscoreを受け取ってDBを更新する
	upd, err := db.Prepare("UPDATE userdata set score=? where name=?")
	if err != nil {
		log.Fatal(err)
	}
	r.POST("/ranking/update", func(c *gin.Context) {
		var newscore User
		c.BindJSON(&newscore)
		upd.Exec(newscore.Score, newscore.Name)
		userResult = getUserdata(db)
	})

	// 上位3人を表示
	r.GET("/ranking/top", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name_1st":  userResult[0].Name,
			"name_2nd":  userResult[1].Name,
			"name_3rd":  userResult[2].Name,
			"score_1st": userResult[0].Score,
			"score_2nd": userResult[1].Score,
			"score_3rd": userResult[2].Score,
		})
	})
	r.GET("/ranking/top/1st", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":  userResult[0].Name,
			"rank":  1,
			"score": userResult[0].Score,
		})
	})
	r.GET("/ranking/top/2nd", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":  userResult[1].Name,
			"rank":  2,
			"score": userResult[1].Score,
		})
	})
	r.GET("/ranking/top/3rd", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":  userResult[2].Name,
			"rank":  3,
			"score": userResult[2].Score,
		})
	})
}

func getUserdata(db *sql.DB) []User {
	// 降順にすべてのデータを格納する
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
	return userResult
}

/*
func mapCollection(r *gin.Engine, db *sql.DB) {
	mapResult := getMapdata(db)
	// 位置情報などをjson形式で受け取りDBに保存する
	// 今の位置から1km以内のデータを読み出す
	for i := 0; i < 10; i++ {
		r.GET("mapcollection/near/"+string(i), func(c *gin.Context) {
			c.JSON(200, gin.H{
				"pickID": mapResult[i].PictID,
				"date":   mapResult[i].Date,
				"time":   mapResult[i].Time,
				"lat":    mapResult[i].Lat,
				"lng":    mapResult[i].Lng,
				"name":   mapResult[i].Name,
			})
		})
	}
}

func getMapdata(db *sql.DB) []Map {
	// 降順にすべてのデータを格納する
	rows, err := db.Query("select * from test.mapdata")
	if err != nil {
		log.Fatal("SQL fetch error.")
	}
	var mapResult []Map
	for rows.Next() {
		maps := Map{}
		if err := rows.Scan(&maps.PictID, &maps.Date, &maps.Time, &maps.Lat, &maps.Lng, &maps.Name); err != nil {
			log.Fatal("maps fetch error.")
		}
		mapResult = append(mapResult, maps)
	}
	return mapResult
}
*/
