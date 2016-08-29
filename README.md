# OBNodeMan
OpenBazaar Node Manager

This is a small Go wrapper for the OpenBazaar v1 (python) app - it allows your OpenBazaar shop to be remotely restarted with a single http call if it becomes unresponsive. This project has only been tested on Ubuntu Linux. It:

1. It is a single executable file,

2. Launches OpenBazzar upon startup,

3. Exposes a http API allowing:

a. OpenBazaar python process to be stopped

b. OpenBazaar python process to be started

c. OpenBazaar python process to be re-started


It should only be used by people with deep knowledge of the Go language, OpenBazaar and network security best practices. It is suitable for use within a private network such as a properly configured AWS VPC. The http API has no security applied - do not expose the API port to the internet without making changes to secure the API as doing so will allow malicious parties to shut down your OpenBazaar server at will. 

## Usage Instructions
These instructions target a Windows development environment and will output a binary suited for use on Linux. A level of familiarity with OpenBazaar and Linux is assumed:
- Set up Go on your Windows development PC: [https://golang.org/doc/install](https://golang.org/doc/install)
- Clone this repo to your Windows development PC, make sure it builds Ok:
```
git clone https://github.com/gzook/obnodeman.git
cd obnodeman
go build
```
- Review main.go to, ensure that you are Ok with the default port (3080)
- Review nodeman.go, function Start() to you are Ok with the arguments that will be used to launch openbazaard.py (particulary the -a 0.0.0.0 flag)
- Produce the Linux executable (named obnodeman)
```
build-linux.bat
```
- Clone the OpenBazaar-Server repo onto your Linux server, launch it, make sure OpenBazaar is working Ok, shut OpenBazaar down
- Copy the Linux executable onto your OpenBazaar host server, placing it into the directory that contains the OpenBazaar server launch file openbazaard.py 
- Make the file executable
```
sudo chmod +x obnodeman
```
- Launch OBNodeMan, wait a short period, then confirm your OpenBazaar node is reachable using the starndard OpenBazaar client app
```
./obnodeman
```
* In order to stop your OpenBazaar node perform a http GET (e.g. use your web browser) against http://{yourServerIp}:{port e.g. 3080}/stop
* In order to start your OpenBazaar node perform a http GET against http://{yourServerIp}:{port e.g. 3080}/start
* In order to stop and then start your OpenBazaar node perform a http GET against http://{yourServerIp}:{port e.g. 3080}/restart
- To run OBNodeman as a service using upstart, create the config file:
```
sudo nano /etc/init/obnodeman.conf
``` 
and then insert the text, replacing {path-to-folder} and {userId} with your folder containg obnodeman, e.g. /home/myuser/OpenBazaar-Server and the Linux user name, e.g. myuser:
```
description "OpenBazaar Node Manager"
author "timm@gzook.com"

start on runlevel [2345]
stop on runlevel [!2345]


respawn
setuid {userId}
setgid {userId}
chdir {path-to-folder}
exec ./obnodeman

```