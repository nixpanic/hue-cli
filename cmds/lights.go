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
	"fmt"

	hue "github.com/collinux/GoHue"
	"github.com/spf13/cobra"
)

type LightOptions struct {
	light     string
	toggle    bool
	colorLoop bool
	blink     int
}

var (
	lightOptions LightOptions
)

func initLights(cmd *cobra.Command) {
	// hue-cli list-lights
	cmd.AddCommand(cmdListLights)
	addBridgeOptions(cmdListLights)
	cmdListLights.SilenceUsage = true

	// hue-cli light
	cmd.AddCommand(cmdLight)
	addBridgeOptions(cmdLight)
	cmdLight.Flags().StringVar(&lightOptions.light, "light", "",
		"act on the given light by name")
	cmdLight.Flags().BoolVar(&lightOptions.toggle, "toggle", false,
		"Toggle light switch")
	cmdLight.Flags().BoolVar(&lightOptions.colorLoop, "colorloop", false,
		"enable/disable color-loop for a light")
	cmdLight.Flags().IntVar(&lightOptions.blink, "blink", -1,
		"blink a light for the given number of seconds")
	cmdLight.SilenceUsage = true

}

var cmdListLights = &cobra.Command{
	Use:   "list-lights",
	Short: "list all lights attached to the bright",
	Long:  "list all lights attached to the bridge",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		lights, err := bridge.GetAllLights()

		fmt.Printf("Found %d lights\n", len(lights))
		for _, light := range lights {
			fmt.Printf("%s\n", lightToString(light))
		}
		return nil
	},
}

var cmdLight = &cobra.Command{
	Use:   "lights",
	Short: "list all lights attached to the bright",
	Long:  "list all lights attached to the bridge",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		light, err := bridge.GetLightByName(lightOptions.light)
		if err != nil {
			return err
		}

		if lightOptions.toggle {
			err = light.Toggle()
			if err != nil {
				return err
			}
			return nil
		}

		// TODO: split colorLoop into its own function
		err = light.ColorLoop(lightOptions.colorLoop)
		if err != nil {
			return err
		}

		var action string
		if lightOptions.colorLoop {
			action = "Activated"
		} else {
			action = "Deactivated"
		}
		fmt.Printf("%s color-loop for '%s'\n",  action, light.Name)

		// TODO: split blink into its own function
		if lightOptions.blink != -1 {
			fmt.Printf("blinking %s for %d seconds\n",  light.Name, lightOptions.blink)

			err = light.Blink(lightOptions.blink)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func lightToString(light hue.Light) string {
	s := fmt.Sprintf("Light: %s\n"+
		"\tIndex: %d\n"+
		"\tType: %s\n"+
		"\tModel: %s\n"+
		"\tUniqueID: %s",
		light.Name, light.Index, light.Type, light.ModelID, light.UniqueID)

	return s
}
