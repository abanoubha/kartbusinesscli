package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	list   bool
	search string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kartbusiness",
		Short: "See all digital business cards published on KartBusiness.com",
		Long:  `See all digital business cards published on KartBusiness.com right from your terminal.`,
		Example: `
kartbusiness               # show business cards one by one.
kartbusiness -l            # show business cards one by one.
kartbusiness -s "Abanoub"  # search all cards for "Abanoub".
`,
	}

	rootCmd.Flags().BoolVarP(&list, "list", "l", false, "show business cards one by one")

	rootCmd.Flags().StringVarP(&search, "search", "s", "", "search all cards")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if list {
			showAllCards()
		} else if search != "" {
			searchCards(search)
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
