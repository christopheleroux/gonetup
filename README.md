# Go Network monitor : gonetup

Tiny go app that monitors a network interface and displays a systray status icon.
Monitoring is based on interface and ip regex matching.
Systray menu provides a quick access to start and stop command.


I designed this app as my first go lang experiment in order to drive a vpnc connection.
I could have used the gnome network manager to achieve but it was not easy because of routing specific scripts

## Installation

```
# Get sources
go get github.com/christopheleroux/gonetup
# build
cd $GOPATH/src/github.com/christopheleroux/gonetup && ./build.sh

#Copy config file
cp $GOPATH/src/github.com/christopheleroux/gonetup $HOME/.gotnetup.yml

# Configure $HOME/.gotnetup.yml with you favourite editor

#run gonetup
$HOME/go/bin/gonetup

```

## Dependencies


## Example : my local vpnc configuration

```
ifacetemplate: "tun*"
iptemplate: ".*"
startcommand: "sudo vpnc"
stopcommand: "sudo vpnc-disconnect"
```
