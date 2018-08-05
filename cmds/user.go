//
// This file is part of hue-cli
// A program written in the Go Programming Language for the Philips Hue API.
// Copyright (C) 2018 Niels de Vos
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
//

package cmds

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nixpanic/hue-cli/utils"
)

type UserOptions struct {
	deviceName string
}

var (
	userOptions UserOptions
)

func initUser(cmd *cobra.Command) {
	// hue-cli create-user
	cmd.AddCommand(cmdCreateUser)
	addBridgeOptions(cmdCreateUser)

	// hue-cli create-user --device=<name>
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	cmdCreateUser.Flags().StringVar(&userOptions.deviceName, "device", hostname,
		"name of the device hue-cli is running on (optional)")
	cmdCreateUser.SilenceUsage = true
}

var cmdCreateUser = &cobra.Command{
	Use:   "create-user",
	Short: "create a new user on the bridge",
	Long:  "create a new user on the bridge, should have pressed the 'link button' in advance",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		// we got a bridge, create a new user
		user, err := bridge.CreateUser("hue-cli#"+userOptions.deviceName)
		if err != nil {
			return err
		}

		// generate a new config file
		config := &utils.ConfigFile{
			Bridges: []utils.BridgeConfig{{
				IPAddress: bridge.IPAddress,
				User: user,
			}},
		}

		configOut, err := config.String()
		if err != nil {
			return errors.New(fmt.Sprintf("failed to conver config to string (%s)", err))
		}

		fmt.Printf("new configuration: %s\n", configOut)

		return nil
	},
}
