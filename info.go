package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"

	"golang.org/x/sys/windows"
)

type Device struct {
	HostName        string
	IpAddr          string
	PublicIpAddr    string
	UserName        string
	OperatingSystem string
	ProcessID       int
}

func (d *Device) setPublicIpAddr() {
	api := "https://api.ipify.org"
	resp, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	d.PublicIpAddr = string(body)
}

func (d *Device) setOperatingSystem() {
	version := windows.RtlGetVersion()
	switch {
	case version.MajorVersion == 10 && version.BuildNumber >= 22000:
		d.OperatingSystem = "Windows 11"
	case version.MajorVersion == 10:
		d.OperatingSystem = "Windows 10"
	default:
		d.OperatingSystem = "Other"
	}
}
func (d *Device) setPID() {
	pid := os.Getpid()
	d.ProcessID = pid
}

func (d *Device) SetUsername() {
	current, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	d.UserName = current.Username
}
func (d *Device) SetIpAddr() {
	ipAddr, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range ipAddr {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				d.IpAddr = ipNet.IP.String()
			}
		}
	}
}
func (d *Device) SetHostName() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	d.HostName = hostname
}
func (d *Device) toJSON() string {
	b, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
func (d *Device) run() string {
	d.SetHostName()
	d.SetIpAddr()
	d.SetUsername()
	d.setOperatingSystem()
	d.setPID()
	d.setPublicIpAddr()
	j := d.toJSON()
	return j
}
