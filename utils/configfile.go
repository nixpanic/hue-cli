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

package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	// for yaml conversion of the ConfigFile
	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	Bridges []BridgeConfig
}

type BridgeConfig struct {
	IPAddress string `yaml:"ipaddress"`
	User      string `yaml:"user"`
}

func (config *ConfigFile) String() ([]byte, error) {
	s, err := yaml.Marshal(&config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to convert to yaml (%s)", err))
	}

	return s, nil
}

func NewConfigFile(filename string) (*ConfigFile, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	conf, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	config := ConfigFile{}
	err = yaml.Unmarshal(conf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
