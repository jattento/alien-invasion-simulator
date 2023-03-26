package cmd

import (
	"log"

	"github.com/jattento/alien-invasion-simulator/cmd/client"

	"github.com/spf13/cobra"
)

var (
	_aliens     *int
	_days       *int
	_cityConfig *string
	_matrix     *int
	_cities     *int

	rootCmd = &cobra.Command{
		Use:   "alien-sim",
		Short: "An alien invasion simulator",
		Long:  "An alien invasion simulator with 99% accuracy.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := client.Run(*_days, *_aliens, *_cities, *_matrix, *_cityConfig); err != nil {
				log.Fatal(err.Error())
			}
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	_aliens = rootCmd.Flags().IntP("aliens", "a", 15, "Amount of aliens to spawn")
	_days = rootCmd.Flags().IntP("days", "d", 10000, "Days until simulation ends.")

	_cityConfig = rootCmd.Flags().String("city-config", "", "path where to find the city config file.")
	_matrix = rootCmd.Flags().IntP("matrix", "m", 5, "Matrix size where the value is N when N*N=total matrix size.")
	_cities = rootCmd.Flags().IntP("cities", "c", 20, "Amount of cities deployed in the matrix.")

	rootCmd.MarkFlagsMutuallyExclusive("city-config", "matrix")
	rootCmd.MarkFlagsMutuallyExclusive("city-config", "cities")
}
