# HLS Subscriber Stats

To get meaningful statistics of a HTTP Live Stream (HLS), such as the current number of subscribers, it is necessary to analyze a *modified* `access.log`.


## Requirements

- golang


## Build

```shell
go build hls-subscriber-stats.go
```


## Usage (nginx)

1.  Add an additional `log_format` to your `http` section:  
    ```nginx
    http {
        log_format hls-subscriber-stats "$msec $request_filename";
    }
    ```

2.  Activate `access_log` in your `server` section:  
    ```nginx
    server {
        access_log /path-log/hls-access.log hls-subscriber-stats;
    }
    ```

3.  Execute *Go* program in your `application` section:
    ```nginx
    application stream {
        live on;
        hls  on;
        exec_push /path-bin/hls-subscriber-stats -input /path-log/hls-access.log -output "/path-www/$name.json" -name "$name" -interval 10 -segments 3
    }
    ```
