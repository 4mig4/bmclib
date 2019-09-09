package redfish

import (
	"fmt"
	"strings"

	"github.com/bmc-toolbox/bmclib/internal/httpclient"

	gofish "github.com/stmcginnis/gofish/school"
	// this make possible to setup logging and properties at any stage
	_ "github.com/bmc-toolbox/bmclib/logging"
)

// Ilo holds the status and properties of a connection to an iLO device
type RedFish struct {
	ip        string
	username  string
	password  string
	apiClient *gofish.ApiClient
	service   *gofish.Service
	sessionID *string
}

// New returns a new Ilo ready to be used
func New(ip string, username string, password string) (r *RedFish, err error) {
	client, err := httpclient.Build()
	if err != nil {
		return r, err
	}

	apiClient, err := gofish.APIClient(fmt.Sprintf("https://%s", ip), client)
	if err != nil {
		return r, err
	}

	return &RedFish{ip: ip, username: username, password: password, apiClient: apiClient}, err
}

// CheckCredentials verify whether the credentials are valid or not
func (r *RedFish) CheckCredentials() (err error) {
	return r.httpLogin()
}

// Serial returns the device serial
func (r *RedFish) Serial() (serial string, err error) {
	err = r.httpLogin()
	if err != nil {
		return serial, err
	}

	entries, err := r.service.Chassis()
	if err != nil {
		return serial, err
	}

	for _, e := range entries {
		if e.ChassisType == "Blade" || e.ChassisType == "RackMount" {
			return strings.ToLower(e.SerialNumber), err
		}
	}
	return serial, err
}

// ChassisSerial returns the serial number of the chassis where the blade is attached
func (r *RedFish) ChassisSerial() (serial string, err error) {
	err = r.httpLogin()
	if err != nil {
		return serial, err
	}

	entries, err := r.service.Chassis()
	if err != nil {
		return serial, err
	}

	for _, e := range entries {
		if e.ChassisType == "Enclosure" {
			return strings.ToLower(e.SerialNumber), err
		}
	}
	return serial, err
}

// Model returns the device model
func (r *RedFish) Model() (model string, err error) {
	err = r.httpLogin()
	if err != nil {
		return model, err
	}

	entries, err := r.service.Chassis()
	if err != nil {
		return model, err
	}

	for _, e := range entries {
		if e.ChassisType == "Blade" || e.ChassisType == "RackMount" {
			return e.Model, err
		}
	}
	return model, err
}

// HardwareType returns the type of bmc we are talking to
func (r *RedFish) HardwareType() (model string, err error) {
	err = r.httpLogin()
	if err != nil {
		return model, err
	}

	entries, err := r.service.Managers()
	if err != nil {
		return model, err
	}

	for _, e := range entries {
		if e.ManagerType == "BMC" {
			return e.Model, err
		}
	}
	return model, err
}

// Version returns the version of the bmc we are running
func (r *RedFish) Version() (version string, err error) {
	err = r.httpLogin()
	if err != nil {
		return version, err
	}

	entries, err := r.service.Managers()
	if err != nil {
		return version, err
	}

	for _, e := range entries {
		if e.ManagerType == "BMC" {
			return e.FirmwareVersion, err
		}
	}
	return version, err
}

// Name returns the version of the bmc we are running
func (r *RedFish) Name() (name string, err error) {
	err = r.httpLogin()
	if err != nil {
		return name, err
	}

	entries, err := r.service.Systems()
	if err != nil {
		return name, err
	}

	for _, e := range entries {
		if e.SystemType == "Physical" {
			return e.HostName, err
		}
	}
	return name, err
}

// Status returns health string status from the bmc
func (r *RedFish) Status() (status string, err error) {
	err = r.httpLogin()
	if err != nil {
		return status, err
	}

	entries, err := r.service.Systems()
	if err != nil {
		return status, err
	}

	for _, e := range entries {
		if e.SystemType == "Physical" {
			return string(e.Status.Health), err
		}
	}
	return status, err
}

// Memory returns the total amount of memory of the server
func (r *RedFish) Memory() (memory int, err error) {
	err = r.httpLogin()
	if err != nil {
		return memory, err
	}

	entries, err := r.service.Systems()
	if err != nil {
		return memory, err
	}

	for _, e := range entries {
		if e.SystemType == "Physical" {
			return int(e.MemorySummary.TotalSystemMemoryGiB), err
		}
	}
	return memory, err
}

// CPU returns the cpu, cores and hyperthreads of the server
func (r *RedFish) CPU() (cpu string, cpuCount int, coreCount int, hyperthreadCount int, err error) {
	err = r.httpLogin()
	if err != nil {
		return cpu, cpuCount, coreCount, hyperthreadCount, err
	}

	entries, err := r.service.Systems()
	if err != nil {
		return cpu, cpuCount, coreCount, hyperthreadCount, err
	}

	for _, e := range entries {
		if e.SystemType == "Physical" {
			return e.ProcessorSummary.Model, e.ProcessorSummary.Count, err
		}
	}
	return cpu, cpuCount, coreCount, hyperthreadCount, err
}