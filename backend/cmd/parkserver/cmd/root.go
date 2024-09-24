package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ParkWithEase/parkeasy/backend/internal/app/parkserver"
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
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "parkserver",
	Short: "ParkEasy API server",
	Long:  `The API server for ParkEasy app.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debugMode {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		// Shutdown on Ctrl-C
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		config := parkserver.Config{
			Addr:     fmt.Sprintf(":%v", port),
			Insecure: insecure,
		}
		log.Info().Uint16("port", port).Msg("server started")

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
