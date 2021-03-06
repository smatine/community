/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2016 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package qemu

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
)

type machine struct {
	id     string
	server string
}

type vmInfo struct {
	Ico         string `json:"ico"`
	Name        string `json:"-"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	Locked      bool   `json:"locked"`
	CurrentSize string `json:"current_size"`
	TotalSize   string `json:"total_size"`
}

func (m *machine) Status() (vms.MachineStatus, error) {
	ip, _ := m.IP()
	resp, err := http.Get("http://" + string(ip) + ":8080/api/vms")
	if err != nil {
		log.Error(err)
		return vms.StatusUnknown, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return vms.StatusUnknown, err
	}

	var state struct {
		Data []struct {
			Id         string `json:"id"`
			Type       string `json:"type"`
			Attributes vmInfo
		} `json:"data"`
	}

	err = json.Unmarshal(body, &state)
	if err != nil {
		log.Error(err)
		return vms.StatusUnknown, err
	}
	for _, val := range state.Data {
		switch val.Attributes.Status {
		case "running":
			return vms.StatusUp, nil
		case "booting":
			return vms.StatusBooting, nil
		case "download":
			return vms.StatusCreating, nil
		case "available":
			return vms.StatusDown, nil
		}
	}
	return vms.StatusUnknown, nil
}

func (m *machine) IP() (net.IP, error) {
	return net.ParseIP(m.server), nil
}

func (m *machine) Type() (vms.MachineType, error) {
	return defaultType, nil
}

func (m *machine) Platform() string {
	return "qemu"
}

func (m *machine) Start() error {
	ip, _ := m.IP()
	resp, err := http.Post("http://"+string(ip)+":8080/api/vms/"+m.id+"/start", "", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error(err)
		return err
	}
	return nil
}

func (m *machine) Stop() error {
	ip, _ := m.IP()
	resp, err := http.Post("http://"+string(ip)+":8080/api/vms/"+m.id+"/stop", "", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error(err)
		return err
	}
	return nil
}

func (m *machine) Terminate() error {
	return nil
}

func (m *machine) Id() string {
	return m.id
}

func (m *machine) Name() (string, error) {
	return "Windows Active Directory", nil
}

func (m *machine) Progress() (uint8, error) {
	ip, _ := m.IP()
	resp, err := http.Get("http://" + string(ip) + ":8080/api/vms")
	if err != nil {
		log.Error(err)
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	var state struct {
		Data []struct {
			Id         string `json:"id"`
			Type       string `json:"type"`
			Attributes vmInfo
		} `json:"data"`
	}

	err = json.Unmarshal(body, &state)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	for _, val := range state.Data {
		if val.Attributes.Status == "download" {
			currentSize, err := strconv.ParseUint(val.Attributes.CurrentSize, 10, 64)
			if err != nil {
				return 0, err
			}

			totalSize, err := strconv.ParseUint(val.Attributes.TotalSize, 10, 64)
			if err != nil {
				return 0, err
			}

			return uint8(currentSize * 100 / totalSize), nil
		}
	}
	return 100, nil
}

func (m *machine) Credentials() (string, string, error) {
	return "Administrator", "Nanocloud123+", nil
}
