package cmd

import (
	"cataloger/catalogs"
	"cataloger/info"
	"fmt"
	"os"
	"os/user"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgFilename  = "config"
	cfgExtention = "json"
	cfgFolder    = ".cataloger"
)

var (
	fullVersion bool
	configFile  string
	source      string
	verbose     bool
	rootCmd     = &cobra.Command{
		Use:   "cataloger",
		Short: "Util for interact with ldap and active directory catalogs.",
		Run: func(cmd *cobra.Command, args []string) {
			if fullVersion {
				fmt.Println(info.ForPrintFull())
			}
		},
	}
)

func init() {
	cobra.OnInitialize(onInit)
	// Flags
	rootCmd.Flags().BoolVar(&fullVersion, "version", false, "Prints full cataloger version")
	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "ad", "Source catalog type. Avalible: ad, ldap.")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "Path to cataloger config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Use ssl for connection")
	rootCmd.PersistentFlags().String("host", "", "Catalog host name or ip")
	rootCmd.PersistentFlags().Int("port", 0, "Catalog ldap port")
	rootCmd.PersistentFlags().Bool("ssl", false, "Use ssl for connection")
	rootCmd.PersistentFlags().Bool("insecure", false, "Insecure ssl connection (for self-signet certs)")
	rootCmd.PersistentFlags().String("bind-dn", "", "BindDN user for auth in catalog")
	rootCmd.PersistentFlags().String("bind-pass", "", "BindDN user password")
	rootCmd.PersistentFlags().String("base", "", "Catalog search base")
	rootCmd.PersistentFlags().String("user-base", "", "Catalog search base for users entries")
	rootCmd.PersistentFlags().String("group-base", "", "Catalog search base for groups entries")
	// Bind flags to viper
	bindPFlag("server.host", "host")
	bindPFlag("server.port", "port")
	bindPFlag("server.ssl", "ssl")
	bindPFlag("server.insecure", "insecure")
	bindPFlag("auth.bind_dn", "bind-dn")
	bindPFlag("auth.bind_pass", "bind-pass")
	bindPFlag("params.search_base", "base")
	bindPFlag("params.user_search_base", "user-base")
	bindPFlag("params.group_search_base", "group-base")
}

func bindPFlag(keyName string, flagName string) {
	if err := viper.BindPFlag(keyName, rootCmd.PersistentFlags().Lookup(flagName)); err != nil {
		log.Fatal(err)
	}
}

func onInit() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
	log.SetOutput(os.Stdout)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	path := fmt.Sprintf("./%s", cfgFolder)

	if configFile != "" {
		viper.SetConfigType(cfgExtention)
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(cfgFilename)
		viper.AddConfigPath(path)
		viper.AddConfigPath(fmt.Sprintf("$HOME/%s", cfgFolder))
		if osUser, err := user.Current(); err == nil {
			viper.AddConfigPath(fmt.Sprintf("%s/%s", osUser.HomeDir, cfgFolder))
		}

		viper.SetConfigType(cfgExtention)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warning("Config file not found")
		} else {
			log.Error(err)
		}
	}

	if viper.ConfigFileUsed() != "" {
		log.Debugf("Using config file: %v", viper.ConfigFileUsed())

		decoded, err := base64Decode(viper.GetString("auth.bind_pass"))
		if err != nil {
			log.Fatalf("Can't decode password from base64: %s", err.Error())
		}
		if decoded != "" {
			viper.Set("auth.bind_pass", decoded)
		}
	}
}

// Execute adds all child commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createConfig() *catalogs.Config {
	return &catalogs.Config{
		Host:            viper.GetString("server.host"),
		Port:            viper.GetInt("server.port"),
		SSL:             viper.GetBool("server.ssl"),
		Insecure:        viper.GetBool("server.insecure"),
		BindDn:          viper.GetString("auth.bind_dn"),
		BindPass:        viper.GetString("auth.bind_pass"),
		SearchBase:      viper.GetString("params.search_base"),
		UserSearchBase:  viper.GetString("params.user_search_base"),
		GroupSearchBase: viper.GetString("params.group_search_base"),
	}
}
