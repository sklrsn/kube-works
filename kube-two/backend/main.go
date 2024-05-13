package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func init() {
}

var (
	db         *sqlx.DB
	redisCache *cache.Cache
	err        error
)

func CreateTable() (sql.Result, error) {
	q := `
		CREATE TABLE IF NOT EXISTS messages (
			id VARCHAR,
			message VARCHAR
		)
	`
	return db.Exec(q)
}

func ConnectDb() (*sqlx.DB, error) {
	return sqlx.Connect("postgres",
		fmt.Sprintf("host=%v user=postgres password=%v dbname=postgres sslmode=disable port=5432",
			os.Getenv("DB_HOST"), os.Getenv("DB_PASSWORD")))
}

func ConnectRedis() *cache.Cache {
	ring := redis.NewRing(&redis.RingOptions{

		Addrs: map[string]string{
			os.Getenv("REDIS_HOST"): os.Getenv("REDIS_PORT"),
		},
		Password: "",
	})
	return cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Hour),
	})
}

func main() {
	sync.OnceFunc(func() {
		db, err = ConnectDb()
		if err != nil {
			panic(err)
		}
		if _, err := CreateTable(); err != nil {
			panic(err)
		}
		redisCache = ConnectRedis()
	})()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/v1/message", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received a request from %v", r.RemoteAddr)
		var m interface{}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			log.Printf("%v", err)
			http.Error(w, "failed to decode message", http.StatusBadRequest)
			return
		}
		bs, err := json.Marshal(m)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "failed to marshal message", http.StatusBadRequest)
			return
		}
		_data := make(map[string]interface{})
		_data["data"] = base64.StdEncoding.EncodeToString(bs)
		_data["created"] = time.Now().UTC()

		data, err := json.Marshal(_data)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "failed to marshal message", http.StatusBadRequest)
			return
		}

		type Data struct {
			ID      string `db:"id"`
			Message string `db:"message"`
		}

		id := uuid.NewString()
		msg := Data{
			ID:      id,
			Message: string(data),
		}
		msgInBytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "failed to marshal message", http.StatusBadRequest)
			return
		}

		//Cache
		redisCache.Set(&cache.Item{
			Ctx:   context.Background(),
			Key:   id,
			Value: string(msgInBytes),
		})

		//DB
		if _, err := db.NamedExecContext(
			context.Background(),
			"INSERT INTO messages(id, message) VALUES(:id,:message)", msg); err != nil {
			log.Printf("%v", err)
			http.Error(w, "failed to write message to db", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	log.Println("backend is ready to serve traffic")
	log.Fatalf("%v", http.ListenAndServe(":9090", mux))
}
