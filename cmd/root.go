// Copyright © 2016 Alexander Gugel <alexander.gugel@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alexanderGugel/nerva/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "nerva",
	Short: "Common JS registry server",
	Long:  ``,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal("failed to run root command")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogFormatter)
	cobra.OnInitialize(initLogLevel)

	configExts := strings.Join(viper.SupportedExts, "|")

	RootCmd.PersistentFlags().StringP("config", "c", "", "config file (default $HOME/.nerva."+configExts+")")
	RootCmd.PersistentFlags().StringP("logLevel", "l", "info", "log level")
	RootCmd.PersistentFlags().StringP("logFormatter", "f", "text", "log formatter")

	viper.BindPFlag("logging.logLevel", RootCmd.PersistentFlags().Lookup("logLevel"))
	viper.BindPFlag("logging.logFormatter", RootCmd.PersistentFlags().Lookup("logFormatter"))
}

func initLogLevel() {
	logLevel := viper.GetString("logging.logLevel")
	level, err := log.ParseLevel(logLevel)
	contextLog := log.WithFields(log.Fields{
		"logLevel": logLevel,
	})
	if err != nil {
		util.LogWarn(contextLog, err, "failed to parse log level")
	} else {
		log.SetLevel(level)
	}
}

func initLogFormatter() {
	logFormatter := viper.GetString("logging.logFormatter")
	contextLog := log.WithFields(log.Fields{
		"logFormatter": logFormatter,
	})
	switch logFormatter {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	default:
		contextLog.Warn("failed to parse log formatter")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if config := viper.GetString("config"); config != "" {
		viper.SetConfigFile(config)
	}

	viper.SetConfigName(".nerva") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{
			"config": viper.ConfigFileUsed(),
		}).Info("using config file")
	}
}
