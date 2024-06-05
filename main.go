package main

import (
	"database/sql"
	"fmt"
	"log"

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
			showAllCards()
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
