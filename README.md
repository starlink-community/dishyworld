# Welcome

This project aims to make it super easy for Starlink users to discover how their Starlink installation is performing. 

![dishy dashboard](../media/dishy-dashboard.png?raw=true)

Features:
  * Hourly speed tests using Speedtest.net
  * Highest latency remote servers from the wifi router
  * Service wide aggregation of latency
  * Failsafes for networks with more than WAN
  * Minimal polling of your Starlink infrastructure
  * Written in golang, so very portable across Windows, Mac, and Linux
  * Utilizes Starlink native GRPC APIs

## Quickstart

The following will bring up the starlink-exporter agent, prometheus, and grafana with a default dashboard. 

```
$ docker-compose -f configs/docker-compose/compose.yml up
```

## Viewing Grafana

Navigate to `localhost:3000`. The default grafana username and password is `admin`. From there, click on the dashboard "Dishy". 

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
