package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/go-yaml/yaml"
)

const confFile string = ".gonetup.yml"

type goVpnConf struct {
	IfaceTemplate string `yaml:"ifacetemplate"`
	IPTemplate    string `yaml:"iptemplate"`
	StartCommand  string `yaml:"startcommand"`
	StopCommand   string `yaml:"stopcommand"`
}

// TODO configurable icons

func readConf() *goVpnConf {
	userConf := getUserHome() + "/" + confFile
	if _, err := os.Stat(userConf); os.IsNotExist(err) {
		fmt.Println("Config file doesn't exists :", userConf)
		os.Exit(1)
	}

	//read yaml Config file
	b, _ := ioutil.ReadFile(userConf)
	config := goVpnConf{}
	yaml.Unmarshal(b, &config)
	// TODO check validity
	// http://ghodss.com/2014/the-right-way-to-handle-yaml-in-golang/
	return &config
}

func getUserHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
