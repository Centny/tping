# tping
provider a simple tool for testing the network like ping command,but it really connect to server and transfter data between client and server.

### Install

```
go install github.com/Centny/gwf/tping
```

### Usage

Server:

```
tping -m S -p :8080 -l t.log
```

Client

```
tping -h 127.0.0.1:8080 -l t.log
```

### All Options

* `tping -m S` (server)
  * `-l` log file
  * `-p` listen port, default :9100
  * `-a` show all long, default Y
  
	
* `tping [-m C]` (client)
  * `-l` log file
  * `-d` ping delay, default 1s
  * `-t` ping times, default -1
  * `-a` show all long, default Y

### Binary

* Win32
  * [tping.x86.zip](https://raw.githubusercontent.com/Centny/tping/master/bin/tiping.x86.zip)
  * [tping.x64.zip](https://raw.githubusercontent.com/Centny/tping/master/bin/tiping.x64.zip)