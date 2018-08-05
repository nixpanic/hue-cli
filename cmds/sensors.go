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

type SensorOptions struct {
	index int
	name  string
}

var (
	sensorOptions SensorOptions
)

func initSensors(cmd *cobra.Command) {
	// hue-cli list-sensors
	cmd.AddCommand(cmdListSensors)
	addBridgeOptions(cmdListSensors)
	cmdListSensors.SilenceUsage = true

	// hue-cli sensor-set
	cmd.AddCommand(cmdSensorSet)
	addBridgeOptions(cmdSensorSet)
	cmdSensorSet.Flags().IntVar(&sensorOptions.index, "index", -1,
		"index of the sensor to modify")
	cmdSensorSet.Flags().StringVar(&sensorOptions.name, "name", "",
		"name to set for the sensor")
	cmdSensorSet.SilenceUsage = true
}

var cmdListSensors = &cobra.Command{
	Use:   "list-sensors",
	Short: "list all sensors",
	Long:  "list all sensors attached to the bridge",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		sensors, err := bridge.GetAllSensors()

		fmt.Printf("Found %d sensors\n", len(sensors))
		for _, sensor := range sensors {
			fmt.Printf("%s\n", sensorToString(sensor))
		}
		return nil
	},
}

var cmdSensorSet = &cobra.Command{
	Use:   "sensor-set",
	Short: "set attributes of a sensor",
	Long:  "set attributes of a sensor",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		if sensorOptions.index == -1 {
			return errors.New("--index=... is required")
		}

		if sensorOptions.name == "" {
			return errors.New("--name=... is required")
		}

		sensor, err := bridge.GetSensorByIndex(sensorOptions.index)
		if err != nil {
			return err
		}

		err = sensor.SetName(sensorOptions.name)
		if err != nil {
			return err
		}

		return nil
	},
}

func sensorToString(sensor hue.Sensor) string {
	s := fmt.Sprintf("Sensor: %s\n"+
		"\tIndex: %d\n"+
		"\tType: %s\n"+
		"\tProductName: %s\n"+
		"\tUniqueID: %s",
		sensor.Name, sensor.Index, sensor.Type, sensor.ModelID, sensor.UniqueID)

	return s
}
