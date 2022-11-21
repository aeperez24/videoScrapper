/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"aeperez24/videoScrapper/application"
	"aeperez24/videoScrapper/cmd/editor"
	"aeperez24/videoScrapper/service"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editSerieName string
var editSerieNewLink string
var editSerieNewName string
var editSerieNewProvider string

var addSerieNewLink string
var addSerieNewName string
var addSerieNewProvider string

var rootCmd = &cobra.Command{
	Use: "wch",

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(args[0])
	},
}

var listProviders = &cobra.Command{
	Use: "providers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(editor.GetProviders())
	},
}

var listSeries = &cobra.Command{
	Use: "series",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := editor.GetConfigPath()
		if err != nil {
			fmt.Println(err)
			return
		}
		configSeries := application.LoadConfigurationWithPath(configPath).SerieConfigurations

		for _, serie := range configSeries {
			fmt.Printf("%v\t%v\t%v\n", serie.SerieName, serie.SerieLink, serie.Provider)
		}
	},
}

var editSerie = &cobra.Command{
	Use: "editSerie",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := editor.GetConfigPath()
		if err != nil {
			fmt.Println(err)
			return
		}
		config := application.LoadConfigurationWithPath(configPath)
		configSeries := config.SerieConfigurations
		serieName := editSerieName
		selectedConfig := -1

		for i, serie := range configSeries {
			if serie.SerieName == serieName {
				selectedConfig = i
			}
		}
		if selectedConfig == -1 {
			fmt.Printf("no config found for seriename \"%v\"\n", serieName)
			return
		}
		if editSerieNewLink != "" {
			config.SerieConfigurations[selectedConfig].SerieLink = editSerieNewLink

		}
		if editSerieNewName != "" {
			alreadyInUse := nameAlreadyInUse(editSerieNewName, selectedConfig, config.SerieConfigurations)
			if alreadyInUse {
				fmt.Printf("seriename \"%v\" is already in use by another configuration\n", editSerieNewName)
			}
			config.SerieConfigurations[selectedConfig].SerieName = editSerieNewName

		}
		viper.Set("SerieConfigurations", config.SerieConfigurations)
		viper.WriteConfig()

	},
}

var addSerie = &cobra.Command{
	Use: "addSerie",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := editor.GetConfigPath()
		if err != nil {
			fmt.Println(err)
			return
		}
		config := application.LoadConfigurationWithPath(configPath)

		validProvider := isValidProvider(addSerieNewProvider, editor.GetProviders())
		if !validProvider {
			fmt.Printf("provider \"%v\" is not valid\n", addSerieNewProvider)
			return
		}
		nameInUse := nameAlreadyInUse(addSerieNewName, -1, config.SerieConfigurations)
		if nameInUse {
			fmt.Printf("serieName \"%v\" is alredy in use\n", addSerieNewName)
			return

		}
		if addSerieNewName == "" {
			fmt.Println("empty serieName is not valid", addSerieNewName)
			return
		}

		newConfig := service.SerieConfiguration{
			SerieLink: addSerieNewLink,
			SerieName: addSerieNewName,
			Provider:  addSerieNewProvider,
		}
		config.SerieConfigurations = append(config.SerieConfigurations, newConfig)
		fmt.Println(config)
		viper.Set("SerieConfigurations", config.SerieConfigurations)
		viper.WriteConfig()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	editSerie.Flags().StringVar(&editSerieName, "n", "", "Name of configuration to edit")
	editSerie.Flags().StringVar(&editSerieNewLink, "link", "", "Link to set")
	editSerie.Flags().StringVar(&editSerieNewName, "name", "", "Name to set")
	editSerie.Flags().StringVar(&editSerieNewProvider, "provider", "", "provider to set")

	addSerie.Flags().StringVar(&addSerieNewLink, "link", "", "Link to set")
	addSerie.Flags().StringVar(&addSerieNewName, "name", "", "Name to set")
	addSerie.Flags().StringVar(&addSerieNewProvider, "provider", "", "Provider to set")
	addSerie.MarkFlagRequired("provider")
	addSerie.MarkFlagRequired("link")
	addSerie.MarkFlagRequired("name")

	rootCmd.AddCommand(listProviders)
	rootCmd.AddCommand(listSeries)
	rootCmd.AddCommand(editSerie)
	rootCmd.AddCommand(addSerie)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func nameAlreadyInUse(newName string, selectedConfig int, SerieConfigurations []service.SerieConfiguration) bool {
	for index, config := range SerieConfigurations {
		if (config.SerieName) == newName && index != selectedConfig {
			return true
		}
	}
	return false
}

func isValidProvider(provider string, providers []string) bool {
	for _, prov := range providers {
		if prov == provider {
			return true
		}
	}
	return false
}
