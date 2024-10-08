package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ParkWithEase/parkeasy/backend/internal/app/parkserver"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	debugMode bool
	insecure  bool
	port      uint16
	dbURL     string // New flag to get the Postgres connection URL
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "parkserver",
	Short: "ParkEasy API server",
	Long:  `The API server for ParkEasy app.`,
	Run: func(_ *cobra.Command, _ []string) {
		if debugMode {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		// see init() for insecure definition
		if insecure {
			log.Warn().Msg("running in insecure mode")
		}

		// Shutdown on Ctrl-C
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Establish a database connection
		pool, err := pgxpool.New(context.Background(), dbURL)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to connect to the database")
		}
		defer pool.Close()


		config := parkserver.Config{
			Addr:     fmt.Sprintf(":%v", port),
			Insecure: insecure,
			DBPool:   pool, // Pass pool to config
		}
		log.Info().Uint16("port", port).Msg("server started")

		// Start the server and pass the dbPool connection
		if err := config.ListenAndServe(ctx); err != nil {
			log.Fatal().Err(err).Msg("")
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "show debug logs")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "run in insecure mode for development (ie. allow cookies to be sent over HTTP)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/parkserver.toml)")
	rootCmd.PersistentFlags().Uint16VarP(&port, "port", "p", 8080, "port to serve on")
	rootCmd.PersistentFlags().StringVar(&dbURL, "db-url", "postgres://testuser:testpassword@db:5432/testdb?sslmode=disable", "Database connection URL")

	// Bind flags to viper for configuration
	err := viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	// Panic since errors here can only happen due to programming mistakes
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	// Panic since errors here can only happen due to programming mistakes
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("db.url", rootCmd.PersistentFlags().Lookup("db-url"))
	// Panic since errors here can only happen due to programming mistakes
	if err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".backend" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName("parkserver")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
