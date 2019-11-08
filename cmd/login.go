package cmd

import (
	"cataloger/catalogs/ad"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

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
		Run: func(cmd *cobra.Command, args []string) {

			askUser()

			file := cfgFilename + "." + cfgExtention
			folder := "./" + cfgFolder
			osuser, err := user.Current()
			if err == nil {
				folder = osuser.HomeDir + "/" + cfgFolder
			}
			path := folder + "/" + file

			// Create config in not exists
			if viper.ConfigFileUsed() == "" {
				// Create config folder
				if _, err := os.Stat(folder); os.IsNotExist(err) {
					if err := os.Mkdir(folder, 0700); err != nil {
						log.Fatalf("Can't create config direcotry: %s", err.Error())
					}
				}
				log.Debug("Creating config file: " + path)
				f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
				if err != nil {
					log.Fatal(err)
				}
				if err := f.Close(); err != nil {
					log.Fatal(err)
				}
			}
			// Update config file
			viper.SetConfigType(cfgExtention)
			if err := viper.WriteConfig(); err != nil {
				log.Fatal(err)
			}

			log.Debugf("Trying connect to catalog. Source: %s", source)
			// Decode password
			d, err := base64Decode(viper.GetString("auth.bind_pass"))
			if err != nil {
				log.Fatal(err)
			}
			viper.Set("auth.bind_pass", d)
			switch source {
			case "ad":
				if err := ad.CheckConnection(createConfig()); err != nil {
					log.Fatal(err)
				}
				log.Debugf("Successfully connected to %s:%s", viper.GetString("server.host"), viper.GetString("server.port"))
			default:
				log.Fatalf("Unknown source type: %s", source)
			}

			log.Info("Login successfull")
		},
	}
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

func askUser() {
	if viper.ConfigFileUsed() != "" {
		log.Warnf("Already using config file - '%s'. Further answers will overwrite it.", viper.ConfigFileUsed())
	}

	var resp string

	resp = promptString("Host", viper.GetString("server.host"))
	viper.Set("server.host", resp)

	resp = promptString("Port", viper.GetString("server.port"))
	viper.Set("server.port", resp)

	resp = promptString("Use SSL", viper.GetString("server.ssl"))
	switch strings.ToLower(resp) {
	case "true":
		viper.Set("server.ssl", true)
	case "false":
		viper.Set("server.ssl", false)
	default:
		log.Fatal("Unexpected value. Expecting 'true' or 'false'")
	}

	resp = promptString("Insecure SSL", viper.GetString("server.insecure"))
	switch strings.ToLower(resp) {
	case "true":
		viper.Set("server.insecure", true)
	case "false":
		viper.Set("server.insecure", false)
	default:
		log.Fatal("Unexpected value. Expecting 'true' or 'false'")
	}

	resp = promptString("BindDN", viper.GetString("auth.bind_dn"))
	viper.Set("auth.bind_dn", resp)

	fmt.Printf("BindDN password: ")
	pass, _ := gopass.GetPasswdMasked()
	viper.Set("auth.bind_pass", base64Encode(string(pass)))

	resp = promptString("Search base", viper.GetString("params.search_base"))
	viper.Set("params.search_base", resp)
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

func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
