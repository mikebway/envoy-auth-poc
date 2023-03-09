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

## Running Envoy with a customized configuration

At the time of writing, the default demo Envoy configuration could be found at [envoy-demo.yaml](https://github.com/envoyproxy/envoy/blob/main/configs/envoy-demo.yaml).
This was downloaded and modified to create the [envoy-demo/envoy.yaml](envoy-demo/envoy.yaml) configuration in this 
repository. 

Clone the repo, make [envoy-demo](envoy-demo) your working directory, then edit [envoy-demo/envoy.yaml](envoy-demo/envoy.yaml) 
to replace all instances of `mikebroadway.com` and `mikebroadway_com` to match your preferred proxy target. 

If, as is likely, your target site requires HTTPS, you will also need to change the `port_value` on the last line 
from `80` to `443` and append the following lines with `transport_socket` at the same indentation as `load_assignment`:

```yaml
      transport_socket:
        name: envoy.transport_sockets.tls
        typed_config:
          "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
          sni: www.envoyproxy.io
```

In the [Dockerfile](envoy-demo/Dockerfile), modify the Envoy version referenced to match the version that you pulled:

```dockerfile
FROM envoyproxy/envoy:v1.25-latest
COPY envoy.yaml /etc/envoy/envoy.yaml
```

Then, run the following to build a local Docker image:

```shell
docker build -t envoy-demo:v1 .
```

Finally, give it a try with (not running detached so that you can follow the logs):

```shell
docker run --rm -p 10000:10000 envoy-demo:v1
```

To stop the Envoy instance, just `Ctrl-C` in the shell window running the container.

## Clearing 301 permanent redirects in the Chrome browser

If Google or some other chosen target site sets up a permanent redirect such that [localhost:1000](http://localhost:10000/)
always goes to an old target site even after you have reverted the [envoy.yaml](envoy-demo/envoy.yaml) you can clear
Chrome's cache of 301 redirects as follows:

* Open the developer tools (`F12` or View menu)
* Right click on the refresh button (the circular arrow to te left of the address bar), and select **Empty cache and
  hard reload**. This menu only shows when the developer tools are open.
