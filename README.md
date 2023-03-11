# Envoy External Authorization Testbed

This project documents how to run Envoy locally using Docker, serving as a reverse proxy for one or another Docker
hosted services. In particular, how to implement and configure an Envoy [external authorization filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter).

## Installing and configuring a custom Envoy image

Acknowledgement: The guide at [DOCKER : ENVOY - GETTING STARTED](https://www.bogotobogo.com/DevOps/Docker/Docker-Envoy-Getting-Started.php)
is a little out of date but provided the foundation for the notes in this section.

### Find the latest Envoy Docker image

Unfortunately `docker pull envoyproxy/envoy` may not find a `latest` tag. As an alternative, visit the Envoy project 
homepage at [www.envoyproxy.io](https://www.envoyproxy.io/) and look for the bold label declaring: 
**Envoy MAJOR.MINOR.FIX is now available**. For example, at the time of writing, the version label states
**Envoy 1.25.1 is now available**. 

In the examples below, replace every instance of `v1.25-latest`with whatever you find as the current latest MAJOR and 
MINOR versions. For example, you might need to use `v1.34-latest` as your target Envoy version.

## Pull and run Envoy with the default demo configuration

Get the Envoy image:

```shell
docker pull envoyproxy/envoy:v1.25-latest
```

Run Envoy in detached mode, just to confirm that the image works as we expect:

```shell
docker run --rm -d -p 10000:10000 envoyproxy/envoy:v1.25-latest
```

As of Envoy version 1.25.x, using a web browser to hit [localhost:1000](http://localhost:10000/) will then
proxy through to [www.envoyproxy.io](https://www.envoyproxy.io/).

Stop the default Envoy image by listing the running images to find its container ID with:

```shell
docker ps
```

Then stop Envoy with the following, replacing the ID with the one that you found with `docker ps`:

```shell
docker stop c5359bf2a551
```

### Clearing 301 permanent redirects in the Chrome browser

If Google or some other chosen target site sets up a permanent redirect such that [localhost:1000](http://localhost:10000/)
always goes to an old target site even after you have reverted the [envoy.yaml](envoy-local/envoy.yaml) you can clear
Chrome's cache of 301 redirects as follows:

* Open the developer tools (`F12` or View menu)
* Right click on the refresh button (the circular arrow to te left of the address bar), and select **Empty cache and
  hard reload**. This menu only shows when the developer tools are open.

## Preparing the `/etc/host` file for the testbed

Before running the custom Envoy image in [envoy-local](envoy-local), you must first add an entry to the `/etc/host`
file. 

1. Determine the IP address of you local machine. If you are using a Mac, this can be done with the following commands:
   ```shell
   ipconfig getifaddr en0
   ipconfig getifaddr en1
   ```
   Typically, only one of those two commands will display a result; most likely for the `en0` interface. Only if you
   have both Ethernet and Wi-Fi will they both yield a result (which is which ¯\_(ツ)_/¯).
   
   Let's imagine that the address you found was `192.168.1.15`.
2. Edit the `/etc/host` file:
   ```shell
   sudo vi /etc/host
   ```
3. Add the following line at the end of the file, replacing the address with the value that you found in step 1.
   ```text
   192.168.1.15    thishost
   ```
   Save the file and quit `vi`

## Building and starting the Envoy container

To build and run the Envoy container image: 

1. Ensure that [envoy-local/Dockerfile](envoy-local/Dockerfile) references a current version of the Envoy image
   (see above). For example: `FROM envoyproxy/envoy:v1.25-latest`.
2. Open a terminal shell and change directory to `.../envoy-local`
3. Run the following:
   ```shell
   make build
   make run
   ```

Once the container image has been built the first time, you can skip the `make build` step thereafter unless you have 
modified the [envoy-local/envoy.yaml](envoy-local/envoy.yaml) file.

The `make run` command executes Docker to run the Envoy container interactively, i.e. logging to the current shell,
and using the local host network (`--network host`). The container will be removed when Envoy is stopped with `ctrl-c`.

### How the [envoy-local Envoy](envoy-local/envoy.yaml) is configured

At the time of writing, the default demo Envoy configuration can be found at [envoy-demo.yaml](https://github.com/envoyproxy/envoy/blob/main/configs/envoy-demo.yaml).
This was downloaded and modified to create the [envoy-local/envoy.yaml](envoy-local/envoy.yaml) configuration in this
repository.

[envoy-local/envoy.yaml](envoy-local/envoy.yaml) was modified as follows:

* Routes requests to http://localhost:1000 to http://thishost:9090
* Stripped out the TLS/HTTPS support

## Building and starting the `authtest` container

`authtest` is a crude Go web service that dumps the contents of the request cookies and `Authorization` header. Its
sole purpose is to illustrate what, if anything, the external authorization filter has done.

To build and run the `authtest` container image:

1. Ensure that [envoy-local/Dockerfile](envoy-local/Dockerfile) references a current version of the Go language image;
   for example: `FROM golang:1.20-alpine`.
2. Open a terminal shell and change directory to `.../authtest`
3. Run the following:
   ```shell
   make build
   make run
   ```

Once the container image has been built the first time, you can skip the `make build` step thereafter unless you have
modified the Go source files under [./authtest](authtest).

The `make run` command executes Docker to run the `authtest` container interactively, i.e. logging to the current shell,
and using the local host network (`--network host`). The container will be removed when Envoy is stopped with `ctrl-c`.

