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

	hue "github.com/collinux/GoHue"
	"github.com/spf13/cobra"

	"github.com/nixpanic/hue-cli/utils"
)

type BridgeOptions struct {
	ipaddress string
	username  string
}

var (
	bridgeOptions BridgeOptions
)

func addBridgeOptions(cmd *cobra.Command) {
	// hue-cli --bridge=<ip-address>
	cmd.Flags().StringVar(&bridgeOptions.ipaddress, "bridge", bridgeOptions.ipaddress,
		"IP-address of the bridge (optional)")
	// hue-cli --username=<username>
	cmd.Flags().StringVar(&bridgeOptions.username, "username", bridgeOptions.username,
		"username for authentication to the bridge (optional)")
}

func initBridge(cmd *cobra.Command) {
	// hue-cli bridge-config
	cmd.AddCommand(cmdBridgeConfig)
	addBridgeOptions(cmdBridgeConfig)
	cmdBridgeConfig.SilenceUsage = true

	// read the hue-cli.yaml configfile
	config, err := utils.NewConfigFile("hue-cli.yaml")
	if err != nil {
		// could not find the hue-cli.yaml
		// TODO: output for debugging only
		fmt.Printf("failed to load hue-cli.yaml: %s\n", err)
	} else {
		if len(config.Bridges) > 0 {
			bridgeOptions.ipaddress = config.Bridges[0].IPAddress
			bridgeOptions.username = config.Bridges[0].User
		}
	}
}

func getBridge() (*hue.Bridge, error) {
	// TODO: check for (--bridge && --username) || --config
	if bridgeOptions.ipaddress == "" {
		return nil, errors.New("--bridge=<ip-address> is required (for now)")
	} else if bridgeOptions.username == "" {
		return nil, errors.New("--username=<username> is required (for now)")
	}

	// if we know the IP-addres, we dont do any discovery
	bridge, err := hue.NewBridge(bridgeOptions.ipaddress)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to connect to bridge %s (%s)", bridgeOptions.ipaddress, err))
	}

	err = bridge.Login(bridgeOptions.username)
	if err != nil {
		return nil, err
	}

	return bridge, nil
}

var cmdBridgeConfig = &cobra.Command{
	Use:   "bridge-config",
	Short: "detailed bridge configuration",
	Long:  "detailed bridge configuration",


	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		err = bridge.GetConfig()
		if err != nil {
			return err
		}

		fmt.Println(bridgeConfigToString(bridge))

		return nil
	},
}

func bridgeConfigToString(bridge *hue.Bridge) string {
	s := fmt.Sprintf("Bridge:\n"+
		"\tIP-address: %s",
		bridge.IPAddress)

	if bridge.Info.Device.DeviceType != "" {
		s += fmt.Sprintf("\n\tDevice Information:\n"+
			"\t\tDeviceType: %s\n"+
			"\t\tFriendlyName: %s\n"+
			"\t\tManufacturer: %s\n"+
			"\t\tManufacturerURL: %s\n"+
			"\t\tModelDescription: %s\n"+
			"\t\tModelName: %s\n"+
			"\t\tModelNumber: %s\n"+
			"\t\tModelURL: %s\n"+
			"\t\tSerialNumber: %s\n"+
			"\t\tUDN: %s",
			bridge.Info.Device.DeviceType,
			bridge.Info.Device.FriendlyName,
			bridge.Info.Device.Manufacturer,
			bridge.Info.Device.ManufacturerURL,
			bridge.Info.Device.ModelDescription,
			bridge.Info.Device.ModelName,
			bridge.Info.Device.ModelNumber,
			bridge.Info.Device.ModelURL,
			bridge.Info.Device.SerialNumber,
			bridge.Info.Device.UDN)
	}

	return s
}
