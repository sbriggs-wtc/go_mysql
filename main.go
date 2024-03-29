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
	//"strings"
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

//var db_pool *sql.DB

func initiate_db_pool() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	var err error
	db_pool, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db_pool.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

type Album struct {
	ID     string
	Title  string
	Artist string
	Price  string
}

func db_fetch_albums() ([]Album, error) {
	//next() + scan() + struct
	rows, err := db_pool.Query("select * from album")
	if err != nil {
		fmt.Println(err, "Query error")
	}
	var albums []Album
	for rows.Next() {
		var row Album
		err := rows.Scan(&row.ID, &row.Title, &row.Artist, &row.Price)
		if err != nil {
			fmt.Println(err, "Row Scan error")
		}
		albums = append(albums, row)
	}
	fmt.Println(albums, "albums")
	return albums, nil
}

func fetch_albums(w http.ResponseWriter, r *http.Request) {

	search_params := r.URL.Query()
	fmt.Println("Search params:", search_params)
	fmt.Println("GET params were:", r.URL.Query())
	fmt.Println(r.URL, "r.URL")

	//var s2 = []string
	//fmt.Println(search_params["artists"], "artists")

	var s1 = []string{"a", "b"}

	//var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

	for i, v := range s1 {
		fmt.Println(i, v, "GET param")
	}

	//url_parsed, err := url.Parse(r.URL)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println(url_parsed.Query, "r.URL.Query")

	albums, err := db_fetch_albums()
	if err != nil {
		fmt.Println(err, "err")
	}
	json_bytes, err := json.Marshal(albums)
	if err != nil {
		fmt.Println(err, "JSON Marshall err")
	}
	fmt.Println(albums)
	w.Write(json_bytes)
}

var db_pool *sql.DB

func main() {

	initiate_db_pool()
	print_db_pointer(db_pool)

	//db_fetch_all(db_pool)
	//"/albums?artist=John%20Coltane"
	//"/albums?artist=john%20coltane"

	http.HandleFunc("/albums", fetch_albums)
	http.HandleFunc("/print_pointer/", print_pointer)
	http.HandleFunc("/decoder_decode/", decoder_decode)
	http.HandleFunc("/read_unmarshall/", read_unmarshall)

	println("Starting Server")
	http.ListenAndServe(":8080", nil)
}
