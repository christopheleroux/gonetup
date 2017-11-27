package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/christopheleroux/gonetup/icons"
	"github.com/getlantern/systray"
)

// TODO channels
// TODO disp ip in menu
// TODO Quit button

const checkInterval time.Duration = 5 * 1e9

//TODO load labels from config file
const (
	menuTitleStopped = "Start"
	menuTitleStarted = "Stop"
)

//Global vars
var startStopCmd *systray.MenuItem
var config *goVpnConf

func main() {
	config = readConf()

	onExit := func() {
		fmt.Println("Starting onExit")
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
		fmt.Println("Finished onExit")
		os.Exit(0)
	}
	startStopCmd = systray.AddMenuItem(menuTitleStopped, menuTitleStopped)

	go monitor(config)
	systray.Run(onReady, onExit)
}

func monitor(config *goVpnConf) {
	connected := false
	systray.SetIcon(icons.Down)

	for true {
		up := netUp(config)
		if up && !connected {
			systray.SetIcon(icons.Up)
			startStopCmd.SetTitle(menuTitleStarted)
			connected = true
		} else if !up && connected {
			systray.SetIcon(icons.Down)
			connected = false
			startStopCmd.SetTitle(menuTitleStopped)
		}
		time.Sleep(checkInterval)
	}
}

func netUp(config *goVpnConf) bool {
	l, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, f := range l {
		ifaceMatched, _ := regexp.MatchString(config.IfaceTemplate, f.Name)
		if ifaceMatched {
			adrs, _ := f.Addrs()
			for _, a := range adrs {
				adrMatch, _ := regexp.MatchString(config.IPTemplate, a.String())
				if adrMatch {
					return true
				}
			}
		}
	}
	return false
}

func onReady() {
	systray.SetTitle("VPN Monitor")
	systray.SetTooltip("VPN Monitor")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	go menuHandler()
}

func menuHandler() {
	for {
		select {
		case <-startStopCmd.ClickedCh:
			if netUp(config) {
				fmt.Println("Stop : " + config.StopCommand)
				go execCmd("x-terminal-emulator -e " + config.StopCommand)
			} else {
				fmt.Println("Start : " + config.StartCommand)
				go execCmd("x-terminal-emulator -e " + config.StartCommand)
			}
		}
	}
}

func execCmd(cmd string) ([]byte, error) {
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]
	return exec.Command(head, parts...).Output()
}
