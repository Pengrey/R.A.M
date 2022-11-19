# R.A.M (Remote Anamnestic Mapper)

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



## Changelog

Still in development phase. No changelog yet.

## License

[MIT](https://github.com/Pengrey/R.A.M/blob/main/LICENSE)