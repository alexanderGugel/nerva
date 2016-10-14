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
	"github.com/alexanderGugel/nerva/registry"
	"github.com/alexanderGugel/nerva/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// registryCmd represents the registry command
var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Start a new registry server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info(util.Logo)

		storageDir := viper.GetString("backend.storageDir")
		upstreamURL := viper.GetString("backend.upstreamURL")
		shaCacheSize := viper.GetInt("cache.shaCacheSize")
		addr := viper.GetString("listener.addr")
		frontAddr := viper.GetString("listener.frontAddr")
		certFile := viper.GetString("listener.certFile")
		keyFile := viper.GetString("listener.keyFile")

		contextLog := log.WithFields(log.Fields{
			"storageDir":   storageDir,
			"upstreamURL":  upstreamURL,
			"addr":         addr,
			"frontAddr":    frontAddr,
			"certFile":     certFile,
			"keyFile":      keyFile,
			"shaCacheSize": shaCacheSize,
		})

		registryConfig := &registry.Config{
			StorageDir:   storageDir,
			UpstreamURL:  upstreamURL,
			ShaCacheSize: shaCacheSize,
			Addr:         addr,
			CertFile:     certFile,
			KeyFile:      keyFile,
			FrontAddr:    frontAddr,
			Logger:       log.StandardLogger(),
		}
		registry, err := registry.New(registryConfig)
		if err != nil {
			util.LogFatal(contextLog, err, "failed to create registry")
		}

		if err := registry.Start(); err != nil {
			util.LogFatal(contextLog, err, "failed to start registry")
		}
	},
}

func init() {
	RootCmd.AddCommand(registryCmd)

	registryCmd.Flags().String("addr", ":8200", "address to bind to for listening")
	registryCmd.Flags().String("frontAddr", "http://127.0.0.1:8200", "full url of front-facing host that nerva will run on")
	registryCmd.Flags().String("certFile", "", "path to TLS certificate file")
	registryCmd.Flags().String("keyFile", "", "path to TLS key file")

	registryCmd.Flags().String("storageDir", "./packages", "storage directory to use for Git repositories")
	registryCmd.Flags().String("upstreamURL", "http://registry.npmjs.com", "upstream Common JS registry")
	registryCmd.Flags().Int("shaCacheSize", 500, "size of SHA1-cache")

	viper.BindPFlag("listener.addr", registryCmd.Flags().Lookup("addr"))
	viper.BindPFlag("listener.frontAddr", registryCmd.Flags().Lookup("frontAddr"))
	viper.BindPFlag("listener.certFile", registryCmd.Flags().Lookup("certFile"))
	viper.BindPFlag("listener.keyFile", registryCmd.Flags().Lookup("keyFile"))

	viper.BindPFlag("backend.storageDir", registryCmd.Flags().Lookup("storageDir"))
	viper.BindPFlag("backend.upstreamURL", registryCmd.Flags().Lookup("upstreamURL"))

	viper.BindPFlag("cache.shaCacheSize", registryCmd.Flags().Lookup("shaCacheSize"))
}
