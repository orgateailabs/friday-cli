package friday

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "Friday",
	Short: "Get SQL query in everyday language",
	Long: "This is Long Description",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Printf("Running a execute function")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Something is wrong!!", err)
		os.Exit(1)
	}
}