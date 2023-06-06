package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)


type User struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Phone  sql.NullString `json:"phone,omitempty"`
}

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("NAME")

	connectionString :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS public.user_table`)
		if err != nil {
			log.Fatal(err)
		}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS public.user_table
		(
			user_id numeric(10,0) NOT NULL,
			name character varying(50) COLLATE pg_catalog."default" NOT NULL,
			age numeric(3,0) NOT NULL,
			phone character varying(20) COLLATE pg_catalog."default",
			CONSTRAINT user_table_pkey PRIMARY KEY (user_id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO public.user_table (user_id, name, age, phone) VALUES (3, 'Jenny', 34, NULL)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO public.user_table (user_id, name, age, phone) VALUES (2, 'Tom', 29, '1-800-123-1234')
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO public.user_table (user_id, name, age, phone) VALUES (1, 'John', 28, NULL)
	`)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/users", func(c *gin.Context) {
		rows, err := db.Query("SELECT user_id, name, age, phone FROM public.user_table")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var user User
			err := rows.Scan(&user.UserID, &user.Name, &user.Age, &user.Phone)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}

		jsonData, err := json.Marshal(users)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusOK,
			"data":        string(jsonData),
		})
	})

	log.Fatal(router.Run(":8080"))
}
