# R.A.M (Remote Anamnestic Mapper)

```
 (                          *     
 )\ )         (           (  `    
(()/(         )\          )\))(   
 /(_))     ((((_)(       ((_)()\  
(_))        )\ _ )\      (_()((_) 
| _ \       (_)_\(_)     |  \/  | 
|   /   _    / _ \    _  | |\/| | 
|_|_\  (_)  /_/ \_\  (_) |_|  |_| 

[Remote Anamnestic Mapper (v0.1 beta)]
```

R.A.M is a fast and simple ram dump retrieval tool. Multiple executable to better fit the needs of the user. Written in Go. R.A.M is mainly used to retrieve the memory dump of a remote machine. It can also be used to retrieve the memory dump of a local machine.

## Installation

#### Server

```bash
$ git clone https://github.com/Pengrey/R.A.M.git
$ cd R.A.M/server
$ go build -o server
$ chmod +x server
```

#### Agent

```bash
$ git clone https://github.com/Pengrey/R.A.M.git
$ cd R.A.M/agent
$ go build -o agent
$ chmod +x agent
```

> 🚧 Warning
>
> The agent relies on the tool `memdump` to retrieve the memory dump. You can download it [here](http://www.porcupine.org/forensics/tct.html)
>
> The agent should also be run with root privileges.

## Usage

#### Server

```bash
$ ./server -h
Usage of ./server:
  -port string
        Port to be used for communication (default "8080")
  -s    Remove prompt from startup
```

#### Agent

```bash
$ ./agent -h
Usage of ./agent:
  -LPORT string
        Port to be used for communication (default "8081")
  -RHOST string
        IP of the server (default "<local IP>")
  -RPORT string
        Port of the server (default "8080")
```

## Demo

https://user-images.githubusercontent.com/55480558/202860749-801b4b22-95a7-4987-ada6-be65e00cffc9.mp4

## Changelog

Still in development phase. No changelog yet.

## License

[MIT](https://github.com/Pengrey/R.A.M/blob/main/LICENSE)
