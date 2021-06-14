# GATEWAY BLE


## USAGE

Get binary gateway-ble

> chmod +x gateway-ble

> sudo ./gateway-ble -mqtt=127.0.0.1:1883 -db=http://127.0.0.1:8086 -debug

OR

Get Deb packet

> sudo dpkg -i gateway-linux.deb

> gateway-ble -debug



|    Params             | Value                 |   Description             |
| :------------------   | :-------------------  | :----------------         |
|   mqtt                | ex: 127.0.0.1:1883    | broker                    |
|   interval            | ex: 1s                | Interval publish ble      |
|   db                  | http://127.0.0.1:8086 | Influx DB host            |
|   u                   | string                | DB username               |
|   p                   | string                | DB password               |
|   debug               | boolean               | Mode debug                |  
|   du                  | duration              | Duration scanner BLE ex: 5s |
|   dup                 | Boolean               | BLE ducplicate            |
|   device              | String                | Device Ble                |
|   rssi                | String                | Filter rssi ex: 130       |
|   name                | String                | Filter name ex: P%T%EN%80021F ( [%] use to replace [ ] )  |
|   mac                 | String                | Filter max ex: ac:23      |



## WORKFLOW DEV

### GOLANG

#### INSTALL

> cd /usr/local

> sudo wget https://golang.org/dl/go1.16.4.linux-amd64.tar.gz

> sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.4.linux-amd64.tar.gz

> nano $HOME/.profile

Add /usr/local/go/bin to the PATH environment variable

> export PATH=$PATH:/usr/local/go/bin [obselete]

> export PATH=$PATH:$HOME/go/bin

Apply change

>. ~/.profile

> go version

---

#### PROJET

**INSTALL DEPENDANCES**

> go mod tidy

**ENV**

>go env

**RUN GO**

> go run main.go

> ./main

**BUILD GO**

> go build .

>./main

**Clean**

> go clean -i

**Install**

Install app 

> go install .

> sudo ~/go/bin/gateway-ble

---

```bash
make ~/go/src/<project>

# initializing dependencies file (go.mod)
$ go mod init

# installing a module
$ go get github.com/go-shadow/moment

# updating a module
$ go get -u github.com/go-shadow/moment

# removing a module
$ rm -rf $GOPATH/pkg/mod/github.com/go-shadow/moment@v<tag>-<checksum>/

# pruning modules (removing unused modules from dependencies file)
$ go mod tidy

# download modules being used to local vendor directory (equivalent of downloading node_modules locally)
$ go mod vendor
```

---

### HCI BLUETOOTH

```tips
sudo hciconfig
sudo hciconfig hci down

# for raspberry
sudo service bluetooth stop
```

---

### BLE & ROOT

__FOR ROOT (BLE ACCESS):__

> sudo su

> nano $HOME/.profile

Add /usr/local/go/bin to the PATH environment variable

> export PATH=$PATH:/usr/local/go/bin

> . $HOME/.profile

> make dev

### REFERENCES

- __ref__: https://awesomeopensource.com/project/miguelmota/golang-for-nodejs-developers
- __ref__: https://www.beaconzone.co.uk/blog/
- __ref__: https://github.com/golang/go/wiki/GoArm
- __ref__: https://gorm.io/docs/index.html
- __ref__: https://github.com/tevjef/go-runtime-metrics
- __ref__: https://stackoverflow.com/questions/45585589/golang-fatal-error-concurrent-map-read-and-map-write/45585833


### VSCODE EXT

- Docker: ms-azuretools.vscode-docker
- REST Client: humao.rest-client
- Go: golang.go

---