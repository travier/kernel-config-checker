// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var config string
var profile string

type Recommendation struct {
	name    string
	options map[string]string
}

type Config struct {
	name    string
	options map[string]string
}

func parseKernelConfig(filename string) (map[string]string, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	if !stat.Mode().IsRegular() {
		return nil, fmt.Errorf("Not a file: %s", filename)
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	options := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(string(bytes)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "CONFIG_") {
			split := strings.Split(line, "=")
			options[split[0]] = strings.TrimPrefix(line, fmt.Sprintf("%s=", split[0]))
		} else if strings.HasPrefix(line, "# CONFIG_") && strings.HasSuffix(scanner.Text(), "is not set") {
			opt := strings.TrimPrefix(line, "# ")
			opt = strings.TrimSuffix(opt, " is not set")
			options[opt] = "is not set"
		}
	}

	return options, nil
}

func main() {
	flag.StringVar(&config, "config", "", "Kernel config to check")
	flag.StringVar(&profile, "profile", "", "Profile to use for checking")
	flag.Parse()

	if config == "" {
		config = "examples/c9s/kernel-x86_64-rhel.config"
	}
	if profile == "" {
		profile = "profiles/ANSSI-BP-028"
	}

	// Parse kernel config from profile recommendations
	recommendations := make([]Recommendation, 0)
	files, err := os.ReadDir(filepath.Join(profile, "kconfig"))
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		filename := filepath.Join(profile, "kconfig", file.Name())
		opts, err := parseKernelConfig(filename)
		if err != nil {
			log.Fatal(err)
		}
		var reco = Recommendation{name: strings.ToUpper(file.Name()), options: opts}
		recommendations = append(recommendations, reco)
	}

	// Parse kernel configs to check
	kconfigs := make([]Config, 0)
	for _, conf := range strings.Split(config, ",") {
		opts, err := parseKernelConfig(conf)
		if err != nil {
			log.Fatal(err)
		}
		kconfigs = append(kconfigs, Config{name: conf, options: opts})
	}

	// Match recommendations and config
	var sorted sort.StringSlice
	for _, r := range recommendations {
		for k, v := range r.options {
			line := fmt.Sprintf("%s;%s;%s", r.name, k, v)
			for _, kconfig := range kconfigs {
				opt := kconfig.options[k]
				if opt == "" {
					opt = "is not set"
				}
				line = fmt.Sprintf("%s;%s", line, opt)
			}
			sorted = append(sorted, line)
		}
	}
	sorted.Sort()
	fmt.Printf("RecoNum;KCONFIG;ANSSI")
	for _, kconfig := range kconfigs {
		fmt.Printf(";%s", kconfig.name)
	}
	fmt.Printf("\n")
	fmt.Println(strings.Join(sorted, "\n"))
}
