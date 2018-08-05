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
)

type DiscoverOptions struct {
	newSensors bool
}

var (
	discoverOptions DiscoverOptions
)

func initDiscover(cmd *cobra.Command) {
	// hue-cli discover-bridges
	cmd.AddCommand(cmdDiscover)
	addBridgeOptions(cmdDiscover)
	cmdDiscover.SilenceUsage = true

	// hue-cli discover-lights
	cmd.AddCommand(cmdDiscoverLights)
	addBridgeOptions(cmdDiscoverLights)
	cmdDiscoverLights.SilenceUsage = true

	// hue-cli discover-sensors
	cmd.AddCommand(cmdDiscoverSensors)
	addBridgeOptions(cmdDiscoverSensors)
	cmdDiscoverSensors.Flags().BoolVar(&discoverOptions.newSensors, "new", false,
		"list the newly detected sensors only")
	cmdDiscoverSensors.SilenceUsage = true
}

var cmdDiscover = &cobra.Command{
	Use:   "discover-bridges",
	Short: "discover bridges",
	Long:  "request the known bridges in this network from https://discovery.meethue.com/",

	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var bridges []hue.Bridge

		if bridgeOptions.ipaddress != "" {
			// if we know the IP-addres, we dont do any discovery
			bridge, err := hue.NewBridge(bridgeOptions.ipaddress)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to find bridge %s: %s", bridgeOptions.ipaddress, err))
			}

			if bridgeOptions.username != "" {
				err = bridge.Login(bridgeOptions.username)
				if err != nil {
					return errors.New(fmt.Sprintf("failed to login on bridge %s: %s\n", bridge.Info.Device.FriendlyName, err))
				}
			}

			bridges = append(bridges, *bridge)
		} else {
			bridges, err = hue.FindBridges()
			if err != nil {
				return err
			}
		}

		fmt.Printf("Found %d bridges\n", len(bridges))
		for _, bridge := range bridges {
			err := bridge.GetInfo()
			if err != nil {
				fmt.Printf("ERROR: failed to get info for bridge at %s (%s)\n", bridge.IPAddress, err)
				// fall-through, just print few details
			}
			fmt.Printf("%s\n", bridgeToString(bridge))
		}
		return nil
	},
}

var cmdDiscoverLights = &cobra.Command{
	Use:   "discover-lights",
	Short: "discover new lights",
	Long:  "request the bridge to probe for new lights",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		err = bridge.FindNewLights()
		if err != nil {
			return errors.New(fmt.Sprintf("failed to start detecting new lights on %s\n", bridge.Info.Device.FriendlyName))
		}

		fmt.Printf("discovery for new lights on bridge %s started, check for new lights in 1 minute\n", bridge.Info.Device.FriendlyName)

		return nil
	},
}

var cmdDiscoverSensors = &cobra.Command{
	Use:   "discover-sensors",
	Short: "discover new sensors",
	Long:  "request the bridge to probe for new sensors",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		if !discoverOptions.newSensors {
			err = bridge.FindNewSensors()
			if err != nil {
				return errors.New(fmt.Sprintf("failed to start detecting new sensors on %s\n", bridge.Info.Device.FriendlyName))
			}

			fmt.Printf("discovery for new sensors on bridge %s started, check for new sensors in 1 minute\n", bridge.Info.Device.FriendlyName)
		} else {
			sensors, err := bridge.GetNewSensors()
			if err != nil {
				return errors.New(fmt.Sprintf("failed get new sensors from %s\n", bridge.Info.Device.FriendlyName))
			}

			for _, sensor := range(sensors) {
				fmt.Println(sensorToString(sensor))
			}
		}

		return nil
	},
}

func bridgeToString(bridge hue.Bridge) string {
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
