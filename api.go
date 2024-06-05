package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

type ApiResponse struct {
	Data []map[string]interface{} `json:"data"`
}

func getCards() {
	resp, err := http.Get("https://kartbusiness.com/api/v1/index?page=1")

	if err != nil {
		log.Fatalf("Error: can not reach API endpoint: %v", err.Error())
		return
	}

	defer resp.Body.Close()

	// fmt.Println(resp.Status, resp.Header, resp.Body)

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err.Error())
	}

	// fmt.Printf("Body type is %T", responseData)

	var apiResponse ApiResponse
	err = json.Unmarshal(responseData, &apiResponse)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// fmt.Printf("type: %T", apiResponse.Data)

	card := make(map[string]interface{})

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO cards (id, cat_id, country_id, gov_id, city_id, slug, name, slogan, mob, whatsapp, mail, web, blog, created_at,updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, v := range apiResponse.Data {

		// id
		value, ok := v["id"].(string)
		if !ok {
			log.Fatalf("Expected id to be a string, got %T", v["id"])
		}
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Failed to convert id to integer: %v", err)
		}
		// fmt.Printf("id: %T\n", id)

		// cat_id
		value, ok = v["cat_id"].(string)
		if !ok {
			log.Fatalf("Expected cat_id to be a string, got %T", v["cat_id"])
		}
		cat_id, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Failed to convert cat_id to integer: %v", err)
		}
		// fmt.Printf("cat_id: %T\n", cat_id)

		// country_id
		value, ok = v["country_id"].(string)
		if !ok {
			log.Fatalf("Expected country_id to be a string, got %T", v["country_id"])
		}
		country_id, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Failed to convert country_id to integer: %v", err)
		}
		// fmt.Printf("country_id: %T\n", country_id)

		// gov_id
		value, ok = v["gov_id"].(string)
		if !ok {
			log.Fatalf("Expected gov_id to be a string, got %T", v["gov_id"])
		}
		gov_id, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Failed to convert gov_id to integer: %v", err)
		}
		// fmt.Printf("gov_id: %T\n", gov_id)

		// city_id
		var city_id int
		if v["city_id"] != nil {
			value, ok = v["city_id"].(string)
			if !ok {
				log.Fatalf("Expected city_id to be a string, got %T", v["city_id"])
			}
			city_id, err = strconv.Atoi(value)
			if err != nil {
				log.Fatalf("Failed to convert city_id to integer: %v", err)
			}
		} else {
			city_id = 1
		}
		// fmt.Printf("city_id: %T\n", city_id)

		// fmt.Printf("slug: %T\n", v["slug"])
		// fmt.Printf("name: %T\n", v["name"])
		// fmt.Printf("slogan: %T\n", v["slogan"])
		// fmt.Printf("mob: %T\n", v["mob"])
		// fmt.Printf("created_at: %T\n", v["created_at"])
		// fmt.Printf("updated_at: %T\n", v["updated_at"])

		_, err = stmt.Exec(
			id, cat_id, gov_id, city_id, country_id,
			v["slug"], v["name"], v["slogan"],
			v["mob"], v["whatsapp"],
			v["mail"], v["web"], v["blog"],
			v["created_at"], v["updated_at"])

		if err != nil {
			log.Fatal(err)
		}
	}

	// commit all cards (50) to local SQLite db
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(card)
}
