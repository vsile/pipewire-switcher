package main

import (
	"flag"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func getNextID(devices [][2]string) string {
	for i, device := range devices {
		if device[1] == "*" {
			if i == len(devices)-1 {
				return devices[0][0]
			}
			return devices[i+1][0]
		}
	}
	return devices[0][0]
}

func main() {
	skip := flag.String("skip", "", "pipewire-switcher --skip 'sof-hda-dsp HDMI'")
	flag.Parse()

	out, err := exec.Command("/bin/sh", "-c", "wpctl status").Output()
	if err != nil {
		log.Fatal(err)
	}
	// First filter only Sinks devices
	// re := regexp.MustCompile(`(?s)Sinks.*?├─`)
	re := regexp.MustCompile(`(?s)Sinks.*?├─`)
	generalMatch := re.FindSubmatch(out)
	if len(generalMatch) == 0 {
		log.Fatal("No devices found")
	}

	re = regexp.MustCompile(`(\*)?\s+(\d+)\.\s+.*`)
	matches := re.FindAllSubmatch(generalMatch[0], -1)

	// Devices contains the id of device and marker indicating if the device is selected
	// For expamle: [[33 ] [42 ] [43 *] [44 ]]
	devices := [][2]string{}
	for _, match := range matches {
		if *skip != "" && strings.Contains(string(match[0]), *skip) {
			continue
		}
		devices = append(devices, [2]string{string(match[2]), string(match[1])})
	}

	// Add another Audio/Sink devices
	re = regexp.MustCompile(`(\*)?\s+(\d+)\.\s+.*\[Audio/Sink\]`)
	advancedMatch := re.FindSubmatch(out)
	if len(advancedMatch) != 0 {
		devices = append(devices, [2]string{string(advancedMatch[2]), string(advancedMatch[1])})
	}

	err = exec.Command("/bin/sh", "-c", "wpctl set-default "+getNextID(devices)).Run()
	if err != nil {
		log.Fatal(err)
	}
}
