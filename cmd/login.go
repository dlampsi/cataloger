package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

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
			// Create config in not exists
			if viper.ConfigFileUsed() == "" {
				askUser()

				file := cfgFilename + "." + cfgExtention
				folder := "./" + cfgFolder
				// osuser, err := user.Current()
				// if err == nil {
				// 	folder = osuser.HomeDir + "/" + cfgFolder
				// }
				path := folder + "/" + file

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
				viper.SetConfigType(cfgExtention)
				if err := viper.WriteConfig(); err != nil {
					log.Fatal(err)
				}
			}
			log.Info("Login successfull")
		},
	}
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

func askUser() {
	if viper.GetString("server.host") == "" {
		promptString("Host", "server.host")
	}

	if viper.GetInt("server.port") == 0 {
		fmt.Print("Port: ")
		port := 0
		fmt.Scanln(&port)
		if port == 0 {
			log.Fatal("Catalog host port can't be 0")
		}
		viper.Set("server.port", port)
	}

	if viper.GetString("auth.bind_dn") == "" {
		promptString("BindDN", "auth.bind_dn")
	}

	if viper.GetString("auth.bind_pass") == "" {
		fmt.Printf("BindDN password: ")
		pass, _ := gopass.GetPasswdMasked()
		viper.Set("auth.bind_pass", base64Encode(string(pass)))
	}

	if viper.GetString("params.search_base") == "" {
		promptString("Search base:", "params.search_base")
	}
}

func promptString(label string, param string) {
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
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("prompt failed %s", err)
	}
	viper.Set(param, result)
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
