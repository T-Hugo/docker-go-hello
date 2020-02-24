# HTTP server in Go, printing the current hostname for load balancer testing

![Docker Build Status](https://img.shields.io/docker/build/hugobin/go-hello?logo=docker)
![MicroBadger Size](https://img.shields.io/microbadger/image-size/hugobin/go-hello?logo=docker)
![GitHub](https://img.shields.io/github/license/t-hugo/docker-go-hello)

>HTTP server printing the hostname for the load balancer testing, about ~1.67 MB in size!

Program written in Go which spins up a web server and prints its current hostname (i.e. the docker container ID), its IP address and port as well as the request URI and the local time of the web server.

Build via a multi-stage build and compress with [UPX](https://github.com/upx/upx) to make the image as small as possible **~1.67MB**


## How to use this image

``` shell script
docker run --rm -it -p 8080:80 hugobin/go-hello
```

``` shell script
$ curl http://localhost:8080/somepath?foo=bar
Hello from b2a20cada094

Server address: 172.17.0.2:80
Server name: b2a20cada094
Date: 2020-02-23T23:30:33Z
URI: /somepath?foo=bar
```

This image was created to be used as a simple backends for various load balancing testing.