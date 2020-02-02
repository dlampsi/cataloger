package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

	"cataloger/catalog/ad"

	"github.com/howeyc/gopass"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to catalog",
		Run:   loginRun,
	}
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

type catalogerConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	SSL        bool   `json:"ssl"`
	Insecure   bool   `json:"insecure"`
	Bind       string `json:"bind"`
	Password   string `json:"password"`
	SearchBase string `json:"search-base"`
}

func loginRun(cmd *cobra.Command, args []string) {
	askUser()
	file := ".cataloger.json"
	folder := "./"
	osuser, err := user.Current()
	if err == nil {
		folder = osuser.HomeDir
	}
	path := folder + "/" + file
	// Create config in not exists
	if viper.ConfigFileUsed() == "" {
		log.Debug("Creating config file: " + path)
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
	viper.SetConfigFile(path)
	config := catalogerConfig{
		Host:       viper.GetString("host"),
		Port:       viper.GetInt("port"),
		SSL:        viper.GetBool("ssl"),
		Insecure:   viper.GetBool("insecure"),
		Bind:       viper.GetString("bind"),
		Password:   viper.GetString("password"),
		SearchBase: viper.GetString("search-base"),
	}
	b, err := json.Marshal(&config)
	if err != nil {
		log.Fatal(err)
	}
	defaultConfig := bytes.NewReader(b)
	v := viper.New()
	v.SetConfigFile(viper.ConfigFileUsed())
	v.SetConfigType("json")
	if err := v.MergeConfig(defaultConfig); err != nil {
		log.Fatal(err)
	}
	if err := v.WriteConfig(); err != nil {
		log.Fatal(err)
	}
	log.Debugf("Trying to connect to catalog...")
	if err := ad.CheckConnection(loadClientConfig()); err != nil {
		log.Fatal(err)
	}
	log.Debugf("Successfully connected to %s:%s", viper.GetString("host"), viper.GetString("port"))

	log.Info("Login successfull")
}

func askUser() {
	if viper.ConfigFileUsed() != "" {
		log.Warnf("Already using config file - '%s'. Further answers will overwrite it.", viper.ConfigFileUsed())
	}

	var resp string

	resp = promptString("Host", viper.GetString("host"))
	viper.Set("host", resp)

	resp = promptString("Port", viper.GetString("port"))
	viper.Set("port", resp)

	resp = promptString("Use SSL", viper.GetString("ssl"))
	switch strings.ToLower(resp) {
	case "true":
		viper.Set("ssl", true)
	case "false":
		viper.Set("ssl", false)
	default:
		log.Fatal("Unexpected value. Expecting 'true' or 'false'")
	}

	resp = promptString("Insecure SSL", viper.GetString("insecure"))
	switch strings.ToLower(resp) {
	case "true":
		viper.Set("insecure", true)
	case "false":
		viper.Set("insecure", false)
	default:
		log.Fatal("Unexpected value. Expecting 'true' or 'false'")
	}

	resp = promptString("BindDN", viper.GetString("bind"))
	viper.Set("bind", resp)

	fmt.Printf("BindDN password: ")
	pass, _ := gopass.GetPasswdMasked()
	viper.Set("password", string(pass))

	resp = promptString("Search base", viper.GetString("search-base"))
	viper.Set("search-base", resp)
}

func promptString(label string, defVal string) string {
	validate := func(input string) error {
		if input == "" {
			return errors.New("Empty param")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}: ",
		Valid:   "{{ . }}: ",
		Invalid: "{{ . }}: ",
		Success: "{{ . }}: ",
	}

	prompt := promptui.Prompt{
		Label:     label,
		Validate:  validate,
		Templates: templates,
	}

	if defVal != "" {
		prompt.Default = defVal
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("prompt failed %s", err)
	}

	return result
}
