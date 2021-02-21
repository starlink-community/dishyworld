# Welcome

This project aims to make it super easy for Starlink users to discover how their Starlink installation is performing. It is intended for Starlink customers that would like to customize their own dashboard. You will need to be comfortable with prometheus metrics and grafana in order to customize the dashboard.

![dishy dashboard](../media/dishy-dashboard.png?raw=true)

Features:
  * Hourly speed tests using Speedtest.net
  * Easily customizable with Grafana
  * Prometheus available for ad-hoc queries
  * Highest latency remote servers from the wifi router
  * Service wide aggregation of latency
  * Failsafes for networks with more than WAN
  * Minimal polling of your Starlink infrastructure
  * Written in golang, so very portable across Windows, Mac, and Linux
  * Utilizes Starlink native GRPC APIs

If you would prefer to just have a dashboard, without all the customization features, the official dashboard can be found in the Starlink App or on the dish itself at http://192.168.100.1/support/statistics

## Quickstart

The following will bring up the starlink-exporter agent, prometheus, and grafana with a default dashboard. 

```
$ docker-compose -f configs/docker-compose/compose.yml up
```

## Viewing Grafana

The default grafana username and password is `admin`. The dashboard will be available at http://localhost:3000/d/dishy/dishy

## Prometheus

Prometheus is available on `localhost:9090`. 

## Data sharing

By default, the quickstart setup will share anonymized metrics back to the project. The goal is to be able to provide a global view of how Starlink is performing, and allow end users to see how thier setup compares to the global deployment. This can be disabled by using the development instructions. 

# Development in docker

If you would like build and test the agent, you can run: 

```
$ docker-compose -f configs/docker-compose/dev.yml up
```

This will run a docker build before setting up the environment. 

## Development locally

You can compile and run the agent itself with standard go tools: 

```
$ git clone https://github.com/starlink-community/dishyworld.git
$ cd dishyworld 
$ go run cmd/starlink-exporter/main.go
```

# Running just the exporter, BYO Prometheus

If you would like to monitor with your own prometheus instance, you can run just the exporter by following the development instructions, then running the agent:

```
$ go build -o starlink-exporter cmd/starlink-exporter/main.go
$ ./starlink-exporter -h
Usage of ./starlink-exporter:
  -dish_addr string
    	Dishy's address (default "192.168.100.1:9200")
  -history_duration duration
    	Polls history this often, then replays it. This means the current metrics from history will be delayed by this amount because of the history replay, but allows us to poll less frequently. Dishy DVR! (default 2m0s)
  -metrics_addr string
    	/metrics address (default "127.0.0.1:2112")
  -ping_interval duration
    	Ping metrics polling interval. (default 1m0s)
  -status_interval duration
    	Status metrics polling interval. (default 4m0s)
  -wifi_addr string
    	Wifi address (default "192.168.1.1:9000")
...
```

By default `http://localhost:2112/metrics` will be available for your promethues to poll. 

## Running on Raspberry Pi

Currently we do not have pre-packaged binaries, but it is very easy to compile to run on Raspberry Pi.

```
$ GOOS=linux GOARCH=arm go build -o starlink-exporter cmd/starlink-exporter/main.go
```

The binary is now available at `./starlink-exporter` to be copied to your pi. 

## systemd

Raspberry Pi and many other systems commonly use systemd, so here is an example if you are using your own compiled binaries. 

```
$ cat /etc/systemd/system/starlink-exporter.service 
[Unit]
Description=starlink-exporter
After=network-online.target

[Service]
# choose your own path
ExecStart=/home/pi/bin/starlink-exporter 
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable the exporter to start at boot, and start it the first time. 

```
$ systemctl enable starlink-exporter
$ systemctl start starlink-exporter
$ systemctl status starlink-exporter
● starlink-exporter.service - starlink-exporter
   Loaded: loaded (/home/pi/systemd/starlink-exporter.service; enabled; vendor preset: enabled)
   Active: active (running) since Sun 2021-02-21 04:01:16 UTC; 11min ago
 Main PID: 18907 (starlink-export)
   CGroup: /system.slice/starlink-exporter.service
           └─18907 /home/pi/bin/starlink-exporter

Feb 21 04:01:16 piaware systemd[1]: Started starlink-exporter.
...
```
