# docker-reverse-proxy

A simple reverse proxy to balance connections between containers based on hostname

![GitHub release (latest by date)](https://img.shields.io/github/v/release/izanagi1995/reverse-proxy?label=Latest%20release)

## Installation instruction

* Install Golang from [the official website](https://golang.org/)
* Run `go get -u github.com/izanagi1995/docker-reverse-proxy`

OR

* Download the latest binaries here : [Access to latest release](https://github.com/izanagi1995/docker-reverse-proxy/releases/latest)

## Running instruction

* Edit `config.yml`
  * Each entry has the hostname to listen to as key
  * The value is an array containing the URLs to the destination
  * Every URL (key and values) should be in format `http://<URL>:<PORT>`
* Start `docker-reverse-proxy`

## Testing

This repository includes a `docker-compose.yml` file and related `config.yml` for testing. Run `docker compose up -d` to start the stack and modify your /etc/hosts file to add demo1.local and demo2.local.

You can then test the software with `curl http://demo1.local:8080` or `curl http://demo2.local:8080`.

## Forcing the destination

By sending a `RP-Target` cookie, you can force the container you want to target, like so : `curl -b 'RP-Target=cnt2:80' http://demo1.local:8080`
