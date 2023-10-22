package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Urls struct {
	gorm.Model
	LongUrl string
	ShortUrl string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var shortUrlRunes = []rune("ABCDEFGHKMNPQRSTUVWXYZ23456789")

func createNewShortUrl() string {
	b := make([]rune, 6)

	for i := range b {
		b[i] = shortUrlRunes[rand.Intn(len(shortUrlRunes))]
	}

	return string(b)
}

func main() {
	fmt.Println("⭐ supershort now running...")
	
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
  if err != nil {
    panic("failed to open database file")
  }

  db.AutoMigrate(&Urls{})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.RawQuery
		if len(query) == 0 {
			fmt.Fprintln(w, "[supershort] nothing here...")
			return
		}

		if strings.HasPrefix(query, "http://") || strings.HasPrefix(query, "https://") {
			url := Urls{
				ShortUrl: createNewShortUrl(),
				LongUrl: query,
			}

			db.Create(&url)
			fmt.Fprintf(w, "Done... #%d : %s -> %s", url.ID, url.ShortUrl, url.LongUrl)
			return
		}

		var url Urls
		result := db.First(&url, "short_url = ?", query)

		if result.RowsAffected == 0 {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		http.Redirect(w, r, url.LongUrl, http.StatusFound)
	})

	panic(http.ListenAndServe("0.0.0.0:8080", nil))
}
