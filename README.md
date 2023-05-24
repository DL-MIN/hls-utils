# HLS Utils

A lightweight helper that adds some functionality to the [RTMP module](https://github.com/arut/nginx-rtmp-module) of nginx. It handles the authentication of streaming endpoints and generates current audience counts for a customized player.


## Features

- Authentication of streaming endpoints via API token
- Current audience counts per streaming endpoint


## Requirements

- golang >= 1.19


## Build

```shell
go mod tidy
go build -a -buildmode=exe -trimpath
```


## Usage

Adjust the given configuration file `config.yaml` to your desired values. The search paths for the configuration file are `./config.yaml` and `/etc/hls-utils/config.yaml`, otherwise define the path with the `-config` flag.

```shell
./hls-utils -config /path/to/config.yaml
```


## Nginx configuration

Add an additional `log_format` to your `http` section:  
```nginx
http {
    log_format hls-utils "$request_filename";
}
```

Activate `access_log` in your `server` section:  
```nginx
server {
    location ~* "\.ts$" {
        access_log /path-log/hls.log hls-utils;
    }
}
```

Set `hls_fragment_naming` in your RTMP/HLS section to `sequential`.

```nginx
rtmp {
    server {
        application hls {
            hls_fragment_naming sequential;
        }
    }
}
```
