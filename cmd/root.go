/*
Copyright © 2021 Isaac Daly <isaac.daly@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	ConfigDir    string `mapstructure:"config-dir"`
	ClientID     string `mapstructure:"client-id"`
	ClientSecret string `mapstructure:"client-secret"`
	AuthTimeout  int    `mapstructure:"auth-timeout"`
	Port         int    `mapstructure:"port"`
}

var (
	configDir      string
	cliConfig      Config
	portStr        string
	authTimeoutStr string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mstodo",
	Short: "mstodo is a CLI program for Microsoft To Do",
	Long: `mstodo is a CLI program for using Microsoft To Do.
Built by @dalyIsaac at https://github.com/dalyIsaac/mstodo

To see available commands, type mstodo help`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// config dir
	rootCmd.PersistentFlags().StringVar(&configDir, "config-dir", defaultConfigDir(), "config directory")
	viper.BindPFlag("config-dir", rootCmd.PersistentFlags().Lookup("config-dir"))

	// port
	rootCmd.PersistentFlags().StringVar(&portStr, "port", "", "port for mstodo")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	// auth timeout
	rootCmd.PersistentFlags().StringVar(&authTimeoutStr, "auth-timeout", "", "seconds to wait before giving up on authentication and exiting")
	viper.BindPFlag("auth-timeout", rootCmd.PersistentFlags().Lookup("auth-timeout"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// config dir
	viper.AddConfigPath(configDir)

	// config file
	viper.SetConfigName("config")

	// read in environment variables that match
	viper.AutomaticEnv()

	// Read in config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Could not read config:", err)
		os.Exit(1)
	}

	// Validate config
	if err := validateConfig(); err != nil {
		fmt.Println("Could not read config:", err)
		os.Exit(1)
	}
}

// Validates the viper config. This should be called after viper has read the
// config.
func validateConfig() error {
	if err := viper.UnmarshalExact(&cliConfig); err != nil {
		return err
	}

	// client id
	if len(cliConfig.ClientID) == 0 {
		return errors.New("client-id must not be empty")
	}

	// client secret
	if len(cliConfig.ClientSecret) == 0 {
		return errors.New("client-secret must not be empty")
	}

	// validate port
	if cliConfig.Port <= 1023 {
		return errors.New("port must be greater than 1023")
	}

	// auth timeout
	if cliConfig.AuthTimeout <= 0 {
		return errors.New("auth-timeout must be greater than 0")
	}

	return nil
}

func defaultConfigDir() string {
	// Find home directory.
	dir, err := homedir.Dir()
	cobra.CheckErr(err)

	// .mstodo folder
	dir = path.Join(dir, ".mstodo")

	return dir
}
