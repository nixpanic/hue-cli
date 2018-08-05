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
	"github.com/spf13/cobra"
)

var HueCli = &cobra.Command{
	Use:   "hue-cli",
	Short: "Commandline application to show the capabilities of GoHue",
	Long:  "Commandline application to show the capabilities of GoHue",
}

func init() {
	initBridge(HueCli)
	initDiscover(HueCli)
	initGroup(HueCli)
	initLights(HueCli)
	initSensors(HueCli)
	initUser(HueCli)
}
