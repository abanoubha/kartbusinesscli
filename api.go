package main

import (
	"fmt"
	"io"
	"net/http"
)

func getCards() {
	resp, err := http.Get("https://kartbusiness.com/api/v1/index?page=1")

	if err != nil {
		fmt.Println("Error: can not reach API endpoint", err.Error())
		return
	}

	defer resp.Body.Close()

	fmt.Println(resp.Status, resp.Header)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body", err.Error())
		return
	}

	fmt.Println(body)
}
