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
	"strconv"
	"strings"

	hue "github.com/collinux/GoHue"
	"github.com/spf13/cobra"
)

type GroupOptions struct {
	name   string
	class  string
	lights string
}

var (
	groupOptions GroupOptions
)

func initGroup(cmd *cobra.Command) {
	// hue-cli list-groups
	cmd.AddCommand(cmdListGroups)
	addBridgeOptions(cmdListGroups)
	cmdListGroups.SilenceUsage = true

	// hue-cli new-group
	cmd.AddCommand(cmdNewGroup)
	addBridgeOptions(cmdNewGroup)
	// hue-cli new-group --name=newgroupname
	cmdNewGroup.Flags().StringVar(&groupOptions.name, "name", "",
		"name of the new group")
	// hue-cli new-group --class=Bedroom
	cmdNewGroup.Flags().StringVar(&groupOptions.class, "class", "Other",
		"type of the room (Bedroom, Kitchen, ...)")
	// hue-cli new-group --lights=2,3,4
	cmdNewGroup.Flags().StringVar(&groupOptions.lights, "lights", "",
		"list of indexes with the lights that should get added to the group")
	cmdNewGroup.SilenceUsage = true

	// hue-cli delete-group
	cmd.AddCommand(cmdDeleteGroup)
	addBridgeOptions(cmdDeleteGroup)
	// hue-cli delete-group --bridge=<ip-address>
	// hue-cli delete-group --name=newgroupname
	cmdDeleteGroup.Flags().StringVar(&groupOptions.name, "name", "",
		"name of the group to delete")

	// hue-cli toggle-group
	cmd.AddCommand(cmdToggleGroup)
	addBridgeOptions(cmdToggleGroup)
	// hue-cli toggle-group --name=groupname
	cmdToggleGroup.Flags().StringVar(&groupOptions.name, "name", "",
		"name of the new group")
	cmdToggleGroup.SilenceUsage = true
}

var cmdListGroups = &cobra.Command{
	Use:   "list-groups",
	Short: "list all lights attached to the bright",
	Long:  "list all lights attached to the bridge",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		groups, err := bridge.GetAllGroups()
		if err != nil {
			return err
		}

		fmt.Printf("Found %d groups\n", len(groups))
		for _, group := range groups {
			fmt.Printf("%s\n", groupToString(group))
		}
		return nil
	},
}

var cmdNewGroup = &cobra.Command{
	Use:   "new-group",
	Short: "create a new group",
	Long:  "create a new group with selected lights",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		if groupOptions.name == "" {
			return errors.New("can create group, no --name=newgroupname passed")
		}

		if groupOptions.lights == "" {
			return errors.New("can create group, no --lights=2,3,4 passed")
		}

		lights := []hue.Light{}
		lightsStr := strings.Split(groupOptions.lights, ",")
		for _, indexStr := range(lightsStr) {
			indexInt, err := strconv.Atoi(indexStr)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to convert light-index %s to integer", indexStr))
			}

			light, err := bridge.GetLightByIndex(indexInt)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to get light for index %d: %s", indexInt, err))
			}

			lights = append(lights, light)
		}

		_, err = bridge.NewGroup(groupOptions.name, groupOptions.class, lights)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create group: %s", err))
		}

		return nil
	},
}

var cmdDeleteGroup = &cobra.Command{
	Use:   "delete-group",
	Short: "delete a group",
	Long:  "delete a group",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		if groupOptions.name == "" {
			return errors.New("can create group, no --name=newgroupname passed")
		}

		group, err := bridge.GetGroupByName(groupOptions.name)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to get group: %s", err))
		}

		err = group.Delete()
		if err != nil {
			return errors.New(fmt.Sprintf("failed to delete group: %s", err))
		}

		return nil
	},
}

var cmdToggleGroup = &cobra.Command{
	Use:   "toggle-group",
	Short: "toggle the light-switch for a group",
	Long:  "toggle the light-switch for a group",

	RunE: func(cmd *cobra.Command, args []string) error {
		bridge, err := getBridge()
		if err != nil {
			return err
		}

		if groupOptions.name == "" {
			return errors.New("can create group, no --name=newgroupname passed")
		}

		group, err := bridge.GetGroupByName(groupOptions.name)
		if err != nil {
			return errors.New(fmt.Sprintf("could not find group %s: %s", groupOptions.name, err))
		}

		// depending on the state, need to turn on or off
		if group.State.AnyOn {
			err = group.Off()
		} else {
			err = group.On()
		}

		if err != nil {
			return errors.New(fmt.Sprintf("failed to toggle light-switch for group %s: %s", groupOptions.name, err))
		}

		return nil
	},
}

func groupToString(group hue.Group) string {
	status := "lights are off"
	if group.State.AllOn {
		status = "lights are on"
	} else if group.State.AnyOn {
		status = "some lights are on"
	}

	s := fmt.Sprintf("Group: %s\n"+
		"Status: %s\n"+
		"Type: %s",
		group.Name, status, group.Type)

	if len(group.Lights) > 0 {
		s += "\nLights:"
	}
	for _, light := range group.Lights {
		s += fmt.Sprintf("\n\t- %s", light.Name)
	}

	return s
}
