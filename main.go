package main

import (
	"encoding/json"
	"fmt"

	"database/sql"
	"github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"os"
)

// can't use println() for structs
type test_struct struct {
	Test string
}

func print_db_pointer(ptr *sql.DB) {
	fmt.Printf("%+v\n\n", ptr) //+v used for structs
	fmt.Printf("1 r %T\n\n", ptr)
	fmt.Printf("2 r &r=%r\n\n", &ptr)
	fmt.Printf("3 r %r=&i=%r\n\n", ptr)
	fmt.Printf("4 r *r=i=%v\n\n", *ptr)
	println(&ptr)
	println("Pointer")
}

func print_pointer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n\n", r) //+v used for structs
	fmt.Printf("1 r %T\n\n", r)
	fmt.Printf("2 r &r=%r\n\n", &r)
	fmt.Printf("3 r %r=&i=%r\n\n", r)
	fmt.Printf("4 r *r=i=%v\n\n", *r)
	println(&r)
	w.Write([]byte("pointer logged"))
}

func decoder_decode(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	fmt.Printf("%T decoder type\n\n", decoder)
	fmt.Printf("%+v decoder value\n\n", decoder)
	var t test_struct
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)
}

func read_unmarshall(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body), "io.ReadAll(r.Body)")
	var t test_struct
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	log.Println(t)
}

func db_connect(cfg mysql.Config) *sql.DB {
	var pool *sql.DB
	var err error
	pool, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := pool.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return pool
}

type Album struct {
	ID     string
	Title  string
	Artist string
	Price  string
}

func db_fetch_all(db_pool *sql.DB) (string, error) {
	//next() + scan() + struct
	rows, err := db_pool.Query("select * from album")
	if err != nil {
		fmt.Println(err, "Query error")
	}
	fmt.Println(rows, "rows")
	
	var albums []Album
	for rows.Next() {
		var row Album
		err := rows.Scan(&row.ID, &row.Title, &row.Artist, &row.Price)
		if err != nil {
			fmt.Println(err, "err")
		}
		albums = append(albums, row);
		fmt.Println(albums, "albums")
	}
	return "hello world", nil
}

func fetch_all(db_pool *sql.DB) func(http.ResponseWriter, *http.Request) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("fetch all")

		//busy here
		test, err := db_fetch_all(db_pool)
		if err != nil {
			fmt.Println(err, "err")
		}
		println(test, "test")

		//busy here
		w.Write([]byte("fetch all"))
	}
	return handler
}

func main() {

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	pool := db_connect(cfg)
	print_db_pointer(pool)
	db_fetch_all(pool)

	http.HandleFunc("/fetch_all/", fetch_all(pool))
	http.HandleFunc("/print_pointer/", print_pointer)
	http.HandleFunc("/decoder_decode/", decoder_decode)
	http.HandleFunc("/read_unmarshall/", read_unmarshall)

	println("Starting Server")
	http.ListenAndServe(":8080", nil)
}
