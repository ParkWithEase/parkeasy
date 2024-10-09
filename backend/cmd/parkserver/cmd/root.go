package cmd

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/ParkWithEase/parkeasy/backend/internal/app/parkserver"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	debugMode bool
	insecure  bool
	port      uint16
	dbURL     string // New flag to get the Postgres connection URL
	dbHost    string
	dbPort    uint16
	dbUser    string
	dbPass    string
	dbName    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "parkserver",
	Short: "ParkEasy API server",
	Long:  `The API server for ParkEasy app.`,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		bindConfig(cmd)
	},
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

		if dbURL == "" {
			host := dbHost
			if dbPort != 0 {
				host = net.JoinHostPort(host, strconv.Itoa(int(dbPort)))
			}
			var user *url.Userinfo
			if dbPass != "" {
				user = url.UserPassword(dbUser, dbPass)
			} else if dbUser != "" {
				user = url.User(dbUser)
			}
			url := url.URL{
				Scheme: "postgres",
				User:   user,
				Host:   host,
				Path:   path.Join("/", dbName),
			}
			dbURL = url.String()
		}

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

		// log.Info().Msg("running migrations")
		// err = config.RunMigrations()
		// if err != nil {
		// 	log.Fatal().Err(err).Msg("")
		// }

		log.Info().Uint16("port", port).Msg("server started")
		// Start the server and pass the dbPool connection
		err = config.ListenAndServe(ctx)
		if err != nil {
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
	rootCmd.PersistentFlags().StringVar(&dbURL, "db-url", "", "Database connection URL, preferred over other DB flags if provided")
	rootCmd.PersistentFlags().StringVar(&dbHost, "db-host", "", "Database connection host")
	rootCmd.PersistentFlags().Uint16Var(&dbPort, "db-port", 0, "Database connection port")
	rootCmd.PersistentFlags().StringVar(&dbUser, "db-user", "", "Database connection user")
	rootCmd.PersistentFlags().StringVar(&dbPass, "db-password", "", "Database connection password")
	rootCmd.PersistentFlags().StringVar(&dbName, "db-name", "", "Database connection database name")

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
	err = viper.BindPFlag("db.host", rootCmd.PersistentFlags().Lookup("db-host"))
	// Panic since errors here can only happen due to programming mistakes
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("db.port", rootCmd.PersistentFlags().Lookup("db-port"))
	// Panic since errors here can only happen due to programming mistakes
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("db.user", rootCmd.PersistentFlags().Lookup("db-user"))
	// Panic since errors here can only happen due to programming mistakes
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("db.password", rootCmd.PersistentFlags().Lookup("db-password"))
	// Panic since errors here can only happen due to programming mistakes
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("db.name", rootCmd.PersistentFlags().Lookup("db-name"))
	// Panic since errors here can only happen due to programming mistakes
	if err != nil {
		panic(err)
	}
}

// setup viper for reading configuration
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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// read configuration from viper
func bindConfig(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", ".")
		if !f.Changed && viper.IsSet(configName) {
			_ = f.Value.Set(viper.GetString(configName))
		}
	})
}
