// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package nvidia

import (
	"os"
	"os/exec"

	"cuda"
	"nvml"
)

func Init() error {
	return nvml.Init()
}

func Shutdown() error {
	return nvml.Shutdown()
}

func LoadUVM() error {
	if _, err := os.Stat("/dev/nvidia-uvm"); err == nil {
		return nil
	}
	return exec.Command("nvidia-modprobe", "-u", "-c=0").Run()
}

func GetDriverVersion() (string, error) {
	return nvml.GetDriverVersion()
}

func GetCUDAVersion() (string, error) {
	return cuda.GetDriverVersion()
}
