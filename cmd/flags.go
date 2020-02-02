package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type flagAttributes struct {
	Id           string
	Short        string
	Description  string
	DefaultValue interface{}
}

// Creates persistent cobra flag and bind it to viper.
func bindPersistentFlag(c *cobra.Command, flagType string, a *flagAttributes) {
	switch flagType {
	case "Bool":
		c.PersistentFlags().Bool(a.Id, false, a.Description)
	case "BoolP":
		c.PersistentFlags().BoolP(a.Id, a.Short, false, a.Description)
	case "Int":
		c.PersistentFlags().Int(a.Id, a.DefaultValue.(int), a.Description)
	case "IntP":
		c.PersistentFlags().IntP(a.Id, a.Short, a.DefaultValue.(int), a.Description)
	case "String":
		c.PersistentFlags().String(a.Id, a.DefaultValue.(string), a.Description)
	case "StringP":
		c.PersistentFlags().StringP(a.Id, a.Short, a.DefaultValue.(string), a.Description)
	case "StringSlice":
		c.PersistentFlags().StringSlice(a.Id, a.DefaultValue.([]string), a.Description)
	case "StringSliceP":
		c.PersistentFlags().StringSliceP(a.Id, a.Short, a.DefaultValue.([]string), a.Description)
	default:
		log.Fatalf("Unsupported flag type: %s", flagType)
	}
	// Bind
	if err := viper.BindPFlag(a.Id, c.PersistentFlags().Lookup(a.Id)); err != nil {
		log.Fatal(err)
	}
}

// Creates cobra flag and bind it to viper.
func bindFlag(c *cobra.Command, flagType string, a *flagAttributes) {
	switch flagType {
	case "Bool":
		c.Flags().Bool(a.Id, false, a.Description)
	case "BoolP":
		c.Flags().BoolP(a.Id, a.Short, false, a.Description)
	case "Int":
		c.Flags().Int(a.Id, a.DefaultValue.(int), a.Description)
	case "IntP":
		c.Flags().IntP(a.Id, a.Short, a.DefaultValue.(int), a.Description)
	case "String":
		c.Flags().String(a.Id, a.DefaultValue.(string), a.Description)
	case "StringP":
		c.Flags().StringP(a.Id, a.Short, a.DefaultValue.(string), a.Description)
	case "StringSlice":
		c.Flags().StringSlice(a.Id, a.DefaultValue.([]string), a.Description)
	case "StringSliceP":
		c.Flags().StringSliceP(a.Id, a.Short, a.DefaultValue.([]string), a.Description)
	default:
		log.Fatalf("Unsupported flag type: %s", flagType)
	}
	// Bind
	if err := viper.BindPFlag(a.Id, c.Flags().Lookup(a.Id)); err != nil {
		log.Fatal(err)
	}
}
