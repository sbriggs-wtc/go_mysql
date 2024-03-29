package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

type FetchAllRes struct {
}

func handle_input(input string) {
	//	fmt.Println(">>>", input);
	if input == "q" {
		os.Exit(0)
	}
	if input == "fetch_all" {
		url := "http://localhost:8080/fetch_all/"
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err, "http.Get error")
		}
		fmt.Println(res, "res")

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err, "res.Body error")
		}
		fmt.Println(body, "res.Body")

		//unmarshalled, err := json.Unmarshall(body, &foobar);
	}

	//res, err := http.Get(input)
	//if err != nil {
	//	fmt.Println("oh no", err)
	//}
	//fmt.Println(res, "res")

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Press \"q\" to exit\n>>> ")
	for scanner.Scan() {
		handle_input(scanner.Text())
		fmt.Print(">>> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
