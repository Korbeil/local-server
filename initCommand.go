package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/mkideal/cli"
	"gopkg.in/yaml.v2"
)

type initT struct {
	cli.Helper
	Config string `cli:"config" usage:"Your configuration file" dft:"pomdok.yaml"`
}

var initCommand = &cli.Command{
	Name: "init",
	Desc: "init your local symfony binary environment to work with a given project",
	Argv: func() interface{} { return new(initT) },
	Fn: func(ctx *cli.Context) error {
		printHeader()

		argv := ctx.Argv().(*initT)
		if _, err := os.Stat(argv.Config); os.IsNotExist(err) {
			fmt.Printf("%s configuration file does not exists 🙊. Maybe you should create or rename your configuration file ? 🧐\n", bold(argv.Config))
			return nil
		}

		data, _ := ioutil.ReadFile(argv.Config)
		config := PomdokYamlConfig{}
		yaml.Unmarshal([]byte(data), &config)
		if config.Pomdok.Tld == "" {
			fmt.Printf("Configuration file error 🙊. Maybe you should give a %s to your domains 🧐\n", yellow("tld"))
			return nil
		}
		if config.Pomdok.Projects == nil {
			fmt.Printf("Configuration file error 🙊. Maybe you should add %s 🧐\n", yellow("projects"))
			return nil
		}

		fileDomains := make(map[string]string)
		currentDirectory, _ := os.Getwd()
		for _, element := range config.Pomdok.Projects {
			if element.Domain == "" {
				fmt.Printf("Configuration file error 🙊. One of the project has empty/no %s 🧐\n", yellow("domain"))
				return nil
			}
			if element.Path == "" {
				fmt.Printf("Configuration file error 🙊. One of the project has empty/no %s 🧐\n", yellow("path"))
				return nil
			}

			fullPath := currentDirectory + element.Path
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				fmt.Printf("Configuration file error 🙊. %s path is not found 🧐\n", bold(fullPath))
				return nil
			}

			if _, ok := fileDomains[element.Domain]; ok {
				fmt.Printf("Configuration file error 🙊. Domain %s is used more than one time 🧐\n", yellow(element.Domain))
				return nil
			}
			fileDomains[element.Domain] = fullPath
		}

		symfonyJsonData := SymfonyJsonProxy{
			Tld:     config.Pomdok.Tld,
			Port:    7080,
			Domains: fileDomains,
		}
		symfonyJson, _ := json.MarshalIndent(symfonyJsonData, "", "  ")

		user, _ := user.Current()
		ioutil.WriteFile(fmt.Sprintf("%s/.symfony/proxy.json", user.HomeDir), symfonyJson, 0644)
		fmt.Printf("Project setup done ✔\n")

		return nil
	},
}
