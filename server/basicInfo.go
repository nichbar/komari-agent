package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/komari-monitor/komari-agent/cmd/flags"
	monitoring "github.com/komari-monitor/komari-agent/monitoring/unit"
	"github.com/komari-monitor/komari-agent/patch"
	"github.com/komari-monitor/komari-agent/update"
)

func DoUploadBasicInfoWorks() {
	ticker := time.NewTicker(time.Duration(flags.InfoReportInterval) * time.Minute)
	for range ticker.C {
		err := uploadBasicInfo()
		if err != nil {
			log.Println("Error uploading basic info:", err)
		}
	}
}
func UpdateBasicInfo() {
	err := uploadBasicInfo()
	if err != nil {
		log.Println("Error uploading basic info:", err)
	} else {
		log.Println("Basic info uploaded successfully")
	}
}
func uploadBasicInfo() error {
	var kernelVersion string
	var gpuName string

	cpu := monitoring.Cpu()

	osname := monitoring.OSName()

	// access kernel and gpu info will trigger SIGSYS in Android
	if flags.IsAndroid {
		kernelVersion = "Unknown"
		gpuName = "None"
	} else {
		kernelVersion = monitoring.KernelVersion()
		gpuName = monitoring.GpuName()
	}

	ipv4, ipv6, _ := monitoring.GetIPAddress()

	data := map[string]interface{}{
		"cpu_name":       cpu.CPUName,
		"cpu_cores":      cpu.CPUCores,
		"arch":           cpu.CPUArchitecture,
		"os":             osname,
		"kernel_version": kernelVersion,
		"ipv4":           ipv4,
		"ipv6":           ipv6,
		"mem_total":      monitoring.Ram().Total,
		"swap_total":     monitoring.Swap().Total,
		"disk_total":     monitoring.Disk().Total,
		"gpu_name":       gpuName,
		"virtualization": monitoring.Virtualized(),
		"version":        update.CurrentVersion,
	}

	// 尝试上传完整数据
	err := tryUploadData(data)
	if err != nil {
		// 兼容 <= 1.0.2
		delete(data, "kernel_version")
		err = tryUploadData(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func tryUploadData(data map[string]interface{}) error {
	endpoint := strings.TrimSuffix(flags.Endpoint, "/") + "/api/clients/uploadBasicInfo?token=" + flags.Token
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := patch.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	message := string(body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d,%s", resp.StatusCode, message)
	}

	return nil
}
