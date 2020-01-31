package cmd

import (
	"os"
	"os/user"

	"github.com/dlampsi/ldapconn"
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./.cataloger.json", "Path to config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose console output")
	rootCmd.PersistentFlags().String("host", "", "LDAP catalog connection host name or ip")
	bindFlag("host")
	rootCmd.PersistentFlags().Int("port", 0, "LDAP catalog connection port")
	bindFlag("port")
	rootCmd.PersistentFlags().Bool("ssl", false, "Use ssl for LDAP catalog connection")
	bindFlag("ssl")
	rootCmd.PersistentFlags().Bool("insecure", false, "Insecure ssl connection to LDAP catalog (for self-signet certs)")
	bindFlag("insecure")
	rootCmd.PersistentFlags().StringP("bind", "b", "", "BindDN user for auth in LDAP catalog")
	bindFlag("bind")
	rootCmd.PersistentFlags().StringP("password", "p", "", "BindDN user password")
	bindFlag("password")
	rootCmd.PersistentFlags().String("search-base", "", "LDAP search base")
	bindFlag("search-base")
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

func bindFlag(flagId string) {
	if err := viper.BindPFlag(flagId, rootCmd.PersistentFlags().Lookup(flagId)); err != nil {
		log.Fatal(err)
	}
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

func createLdapClient() *ldapconn.Config {
	cfg := &ldapconn.Config{
		Host:     viper.GetString("host"),
		Port:     viper.GetInt("port"),
		SSL:      viper.GetBool("ssl"),
		Insecure: viper.GetBool("insecure"),
		BindDN:   viper.GetString("bind"),
		BindPass: viper.GetString("password"),
	}
	log.Debugf("Creating ldap client to %s:%d, bind: %s", cfg.Host, cfg.Port, cfg.BindDN)
	return cfg
}
