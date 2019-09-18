package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"

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
	PictID   int       `json:"pictID"`
	DateTime time.Time `json:"datetime"`
	Lat      float64   `json:"lat"`
	Lng      float64   `json:"lng"`
	Name     string    `json:"name"`
}

func main() {
	r := gin.Default()
	db := sqlConnect()
	ranking(r, db)
	mapCollection(r, db)
	r.Run(":8080")
	defer db.Close()
}

func sqlConnect() *sql.DB {
	db, err := sql.Open("mysql", "root:test@tcp(localhost:3306)/test?parseTime=true")
	if err != nil {
		log.Fatal("SQL open error.")
	}
	db.SetConnMaxLifetime(time.Second * 5)
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
		for i := 0; i < 3; i++ {
			c.JSON(200, gin.H{
				"name":  userResult[i].Name,
				"rank":  i + 1,
				"score": userResult[i].Score,
			})
		}
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
			log.Fatal("userdata fetch error.")
		}
		userResult = append(userResult, user)
	}
	return userResult
}

func mapCollection(r *gin.Engine, db *sql.DB) {
	// 位置情報などをjson形式で受け取りDBに保存する
	ins, err := db.Prepare("insert into mapdata (date_time, lat_lng, name, pictID) values (?, ST_GeomFromText(?), ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	r.POST("/mapcollection/insert", func(c *gin.Context) {
		var newpict Map
		c.BindJSON(&newpict)
		tmp := "POINT(" + strconv.FormatFloat(newpict.Lat, 'f', 6, 64) + " " + strconv.FormatFloat(newpict.Lng, 'f', 6, 64) + ")"
		ins.Exec(newpict.DateTime, tmp, newpict.Name, newpict.PictID)
	})
	// 今の位置から1km以内のデータを読み出す
	r.GET("/mapcollection/near", func(c *gin.Context) {
		nameFlag := false
		locFlag := false
		if c.Query("name") != "" {
			nameFlag = true
		}
		if c.Query("lat") != "" && c.Query("lng") != "" {
			locFlag = true
		}
		con := "where "
		if nameFlag && locFlag {
			con += "name = \"" + c.Query("name") + "\" and " + "ST_Within(lat_lng, ST_Buffer(POINT(" + c.Query("lat") + ", " + c.Query("lng") + "), 0.009))"
		} else if nameFlag {
			con += "name = \"" + c.Query("name") + "\""
		} else if locFlag {
			con += "ST_Within(lat_lng, ST_Buffer(POINT(" + c.Query("lat") + ", " + c.Query("lng") + "), 0.009))"
		} else {
			con = ""
		}
		mapResult, cnt := getMapdata(db, con)
		for i := 0; i < cnt; i++ {
			c.JSON(200, gin.H{
				"datetime": mapResult[i].DateTime,
				"lat":      mapResult[i].Lat,
				"lng":      mapResult[i].Lng,
				"name":     mapResult[i].Name,
				"pictID":   mapResult[i].PictID,
			})
		}
	})
}

// selectに1km以内の条件を追加
func getMapdata(db *sql.DB, con string) ([]Map, int) {
	// 降順にすべてのデータを格納する
	log.Print("select pictID, date_time, ST_X(lat_lng), ST_Y(lat_lng), name from test.mapdata " + con)
	rows, err := db.Query("select pictID, date_time, ST_X(lat_lng), ST_Y(lat_lng), name from test.mapdata " + con)
	if err != nil {
		log.Fatal(err)
	}
	var mapResult []Map
	cnt := 0
	for rows.Next() {
		maps := Map{}
		if err := rows.Scan(&maps.PictID, &maps.DateTime, &maps.Lat, &maps.Lng, &maps.Name); err != nil {
			log.Fatal("mapdata fetch error.")
		}
		mapResult = append(mapResult, maps)
		cnt++
	}
	return mapResult, cnt
}
