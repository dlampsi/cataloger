package cmd

import (
	"os"
	"os/user"

	"cataloger/catalog/ad"
	"cataloger/client"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cataloger",
		Short: "Util for interact with ldap catalogs",
	}
	cfgFile string
	verbose bool
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path to config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose console output")

	bindPersistentFlag(rootCmd, "String", &flagAttributes{
		Id:           "host",
		Description:  "LDAP host name or IP",
		DefaultValue: "",
	})
	bindPersistentFlag(rootCmd, "Int", &flagAttributes{
		Id:           "port",
		Description:  "LDAP host port",
		DefaultValue: 0,
	})
	bindPersistentFlag(rootCmd, "Bool", &flagAttributes{
		Id:          "ssl",
		Description: "Use SSL for connect to LDAP server",
	})
	bindPersistentFlag(rootCmd, "Bool", &flagAttributes{
		Id:          "insecure",
		Description: "Insecure ssl connection to LDAP catalog (for self-signet certs)",
	})
	bindPersistentFlag(rootCmd, "StringP", &flagAttributes{
		Id:           "bind",
		Short:        "b",
		Description:  "BindDN user for auth in LDAP catalog",
		DefaultValue: "",
	})
	bindPersistentFlag(rootCmd, "StringP", &flagAttributes{
		Id:           "password",
		Short:        "p",
		Description:  "BindDN user password",
		DefaultValue: "",
	})
	bindPersistentFlag(rootCmd, "String", &flagAttributes{
		Id:           "search-base",
		Description:  "LDAP search base",
		DefaultValue: "",
	})
}

func initConfig() {
	setLogging()
	viper.SetConfigType("json")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".cataloger")
		viper.AddConfigPath("./")
		viper.AddConfigPath("$HOME")
		if osUser, err := user.Current(); err == nil {
			viper.AddConfigPath(osUser.HomeDir)
		}
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debugf("Config file not found: %s", err.Error())
		} else {
			log.Error(err)
		}
	}
	log.Debugf("Using config file: %s", viper.ConfigFileUsed())
}

func setLogging() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
	log.SetOutput(os.Stdout)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

// Returns ldap client config from parsed cmd flags.
func loadClientConfig() *client.Config {
	cfg := &client.Config{
		Host:         viper.GetString("host"),
		Port:         viper.GetInt("port"),
		SSL:          viper.GetBool("insecure"),
		Insecure:     viper.GetBool("insecure"),
		BindDN:       viper.GetString("bind"),
		BindPassword: viper.GetString("password"),
	}
	log.WithFields(log.Fields{
		"host": cfg.Host,
		"port": cfg.Port,
		"bind": cfg.BindDN,
	}).Debug("LDAP client config")
	return cfg
}

func initAdCatalog() (*ad.Catalog, error) {
	config := loadClientConfig()
	searchAttribute := viper.GetString("attribute")
	if searchAttribute == "" {
		searchAttribute = "sAMAccountName"
	}
	searchBase := viper.GetString("search-base")
	return ad.NewCatalog(config, &ad.Attributes{
		SearchBase:      searchBase,
		SearchAttribute: searchAttribute,
	})
}
