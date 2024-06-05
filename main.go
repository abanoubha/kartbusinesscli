package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"unicode/utf8"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var (
	list   bool
	update bool
	search string
)

func main() {
	var db *sql.DB
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	initializeLocalDatabase(db)

	rootCmd := &cobra.Command{
		Use:   "kartbusiness",
		Short: "See all digital business cards published on KartBusiness.com",
		Long:  `See all digital business cards published on KartBusiness.com right from your terminal.`,
		Example: `
kartbusiness               # show business cards one by one.
kartbusiness -l            # show business cards one by one.
kartbusiness -s "Abanoub"  # search all cards for "Abanoub".
kartbusiness -u            # update/sync local database with new cards.
`,
	}

	rootCmd.Flags().BoolVarP(&list, "list", "l", false, "show business cards one by one")

	rootCmd.Flags().BoolVarP(&update, "update", "u", false, "update local database with new cards")

	rootCmd.Flags().StringVarP(&search, "search", "s", "", "search all cards")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if list {
			showAllCards()
		} else if search != "" {
			searchCards(search)
		} else if update {
			syncData(db)
		} else {
			fmt.Println("KartBusiness : see all digital business cards published on KartBusiness.com")
			showOneCards(db)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func showAllCards() {
	fmt.Println("All cards")
}

func searchCards(search string) {
	fmt.Println("cards containing ", search)
}

func syncData(db *sql.DB) {
	getCards(db)
	fmt.Println("saved top 50 cards into local database")
}

func showOneCards(db *sql.DB) {
	rows, err := db.Query(`SELECT * FROM cards LIMIT 1;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var id, cat_id, country_id, gov_id, city_id int
	var slug, name, slogan, mob, whatsapp, mail, web, blog, created_at, updated_at interface{}

	for rows.Next() {
		if err := rows.Scan(
			&id, &cat_id, &country_id, &gov_id, &city_id,
			&slug, &name, &slogan, &mob, &whatsapp, &mail, &web, &blog,
			&created_at, &updated_at,
		); err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"Card #%v\n+-----------------------------+\n| %v |\n| %v |\n| category: %v country: %v |\n|  governorate: %v city: %v |\n| slug: %v |\n| mobile: %v |\n| whatsapp: %v |\n| mail: %v |\n| website: %v |\n| blog: %v |\n| %v | %v |",
			id, name, slogan, cat_id, country_id, gov_id, city_id,
			slug, mob, whatsapp, mail, web, blog,
			created_at, updated_at)

		fmt.Println("\n██████████████████████████████")

		nameLength, _ := getStringLength(name)
		spaces := (30 - nameLength) / 2

		for i := 0; i < spaces; i++ {
			fmt.Printf("█")
		}
		fmt.Printf("%v", name)
		for i := 0; i < spaces; i++ {
			fmt.Printf("█")
		}

		fmt.Println()

		sloganLength, _ := getStringLength(slogan)
		spaces = (30 - sloganLength) / 2

		for i := 0; i < spaces; i++ {
			fmt.Printf("█")
		}
		fmt.Printf("%v", slogan)
		for i := 0; i < spaces; i++ {
			fmt.Printf("█")
		}

		fmt.Println()

		var mobLength int
		if mob != nil && mob != "" {
			mobLength, _ = getStringLength(mob)
		} else {
			mobLength = 10
			mob = "0000000000"
		}

		var waLength int
		if whatsapp != nil && whatsapp != "" {
			waLength, _ = getStringLength(whatsapp)
		} else {
			waLength = 10
			whatsapp = "0000000000"
		}

		mobWhatsappLength := mobLength + waLength

		if mobWhatsappLength < 30 {
			spaces = (30 - mobWhatsappLength - 1) / 2
			for i := 0; i < spaces; i++ {
				fmt.Printf("█")
			}
			fmt.Printf("%v ", mob)
			fmt.Printf("%v", whatsapp)
			for i := 0; i < spaces; i++ {
				fmt.Printf("█")
			}
		}

		fmt.Println("\n██████████████████████████████")
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func getStringLength(s interface{}) (int, error) {
	rv := reflect.ValueOf(s)
	switch rv.Kind() {
	case reflect.String:
		str := rv.String()
		if !utf8.ValidString(str) {
			return 0, fmt.Errorf("input contains invalid UTF-8 sequence")
		}
		return utf8.RuneCountInString(str), nil
	default:
		return 0, fmt.Errorf("input is not a string")
	}
}
