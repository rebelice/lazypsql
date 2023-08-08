/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rebelice/lazypsql/app"
	"github.com/rebelice/lazypsql/postgres"
	"github.com/spf13/cobra"
)

var (
	flags struct {
		host         string
		port         string
		username     string
		password     string
		noPassword   bool
		databaseName string
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lazypsql",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		database := &postgres.Database{
			DataSource: &postgres.DataSource{
				Host:         flags.host,
				Port:         flags.port,
				Username:     flags.username,
				Password:     flags.password,
				NoPassword:   flags.noPassword,
				DatabaseName: flags.databaseName,
			},
		}

		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
		model := app.NewModel(database, f)
		p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lazypsql.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	usage := "help for "
	if rootCmd.Name() == "" {
		usage += "this command"
	} else {
		usage += rootCmd.Name()
	}
	rootCmd.Flags().BoolP("help", "", false, usage)
	_ = rootCmd.Flags().SetAnnotation("help", cobra.FlagSetByCobraAnnotation, []string{"true"})

	rootCmd.Flags().StringVarP(&flags.databaseName, "dbname", "d", "", `database name to connect to`)
	rootCmd.Flags().StringVarP(&flags.host, "host", "h", "localhost", `database server host or socket directory (default: "local socket")`)
	rootCmd.Flags().StringVarP(&flags.port, "port", "p", "5432", `database server port (default: "5432")`)
	rootCmd.Flags().StringVarP(&flags.username, "username", "U", "", `database user name`)
	rootCmd.Flags().BoolVarP(&flags.noPassword, "no-password", "w", false, `never prompt for password`)
	rootCmd.Flags().StringVarP(&flags.password, "password", "W", "", `force password prompt (should happen automatically)`)
}
