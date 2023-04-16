# Envoy External Authorization Testbed

This project documents how to run Envoy locally using Docker, serving as a reverse proxy for one or another Docker
hosted services. In particular, how to implement and configure an Envoy [external authorization filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter).

## Installing, configuring, and running a custom Envoy image

Acknowledgement: The guide at [DOCKER : ENVOY - GETTING STARTED](https://www.bogotobogo.com/DevOps/Docker/Docker-Envoy-Getting-Started.php)
is a little out of date but provided the foundation for the notes in this section.

### Find the latest Envoy Docker image

Unfortunately `docker pull envoyproxy/envoy` may not find a `latest` tag. As an alternative, visit the Envoy project 
homepage at [www.envoyproxy.io](https://www.envoyproxy.io/) and look for the bold label declaring: 
**Envoy MAJOR.MINOR.FIX is now available**. For example, at the time of writing, the version label states
**Envoy 1.25.1 is now available**. 

In the examples below, replace every instance of `v1.25-latest`with whatever you find as the current latest MAJOR and 
MINOR versions. For example, you might need to use `v1.34-latest` as your target Envoy version.

### Pull and run Envoy with the default demo configuration

Get the Envoy image:

```shell
docker pull envoyproxy/envoy:v1.25-latest
```

Run Envoy in detached mode, just to confirm that the image works as we expect:

```shell
docker run --rm -d -p 10000:10000 envoyproxy/envoy:v1.25-latest
```

As of Envoy version 1.25.x, using a web browser to hit [localhost:10000](http://localhost:10000/) will then
proxy through to [www.envoyproxy.io](https://www.envoyproxy.io/).

Stop the default Envoy image by listing the running images to find its container ID with:

```shell
docker ps
```

Then stop Envoy with the following, replacing the ID with the one that you found with `docker ps`:

```shell
docker stop c5359bf2a551
```

#### Clearing 301 permanent redirects in the Chrome browser

If Google or some other chosen target site sets up a permanent redirect such that [localhost:10000](http://localhost:10000/)
always goes to an old target site even after you have reverted the [envoy.yaml](envoy-local/envoy.yaml) you can clear
Chrome's cache of 301 redirects as follows:

* Open the developer tools (`F12` or View menu)
* Right click on the refresh button (the circular arrow to te left of the address bar), and select **Empty cache and
  hard reload**. This menu only shows when the developer tools are open.

## Building and running the Testbed

Now you know that you can run Envoy locally, we can move on to setting up a four container demonstration:

* Envoy listening at [localhost:10000](http://localhost:10000/)
* An Envoy external authorization filter ([extauth](extauth)) listening at port 50051
* A simple web service ([authtest](authtest)) listening at [localhost:9090](http://localhost:9090/)
* A "login"/"logout" web service ([login](login)) listening at [localhost:9040](http://localhost:9040/)

Envoy, running at [localhost:10000](http://localhost:10000/), routes all URL paths to the ([authtest](authtest))
service running at [localhost:9090](http://localhost:9090/) except for `/login` and `/logout`. The `/login` and `/logout`
paths are routed to ([login](login)) service listening at [localhost:9040](http://localhost:9040/)

Envoy is configured to route requests for most paths (not all) through the [extauth](extauth) external authorization 
filter. The external authorization filter always adds a `x-extauth-was-here` to mark that the filter was invoked. 
If the [login](login) service has created a "session" cookie with a username in it, the external authorization
service will add a second, `x-extauth-authorization`, header containing a signed JWT that wraps the username found
in the session cookie.

You will be able to reach the web service either through the Envoy reverse proxy at [localhost:10000](http://localhost:10000/),
or bypass the proxy and go direct to [localhost:9090](http://localhost:9090/) to see the difference in results.

The [localhost:9090](http://localhost:9090/) web service simply dumps the contents of the request headers to a text 
response. The only difference that you should see between [localhost:10000](http://localhost:10000/) and
[localhost:9090](http://localhost:9090/) is the possible addition of an `X-Header-Set-By-Extauth` by the external 
authorization filter if Envoy filters the request through the authorization filter.

The [envoy-local/envoy.yaml](envoy-local/envoy.yaml) Envoy configuration shall not route all requests through the 
authorization filter. The following URL path patterns shall be routed directly to the [localhost:9090](http://localhost:9090/)
or [localhost:9040](http://localhost:9040/) services without going through the filter:

* `/` - an exact match; the home page; it is assumed that users would not need to be authorized to view the site's root page.
* `/static` - any URL prefixed with `/static` is assumed to be a request for static content, JavaScript, etc that does
  not require authorization to access.
* `/login` - a prefix; there would not be much point in not allowing anonymous users to identify themselves!
* `/logout` - a prefix; repeatedly asking to logout if you are already logged out should not trigger an "unauthorized" response.

Note that only the `/` pattern requires an exact match. `/loginxyz`, `/login?redirect=/dashboard`, and `/login/otc` 
shall all be treated the same way as `/login`.

* `/graphql` is treated as a special path. Envoy invokes the authorization filter as usual, but also sets a 
  `x-require-auth` **context extension**. A Context extensions is a parameter value passed to the filter that
  can be used to influence the filter's behaviour. `x-require-auth: false` would signal that the the filter
  should set the JWT header if there is a valid session cookie but not throw an error if the cookie does not
  exist or is not valid.

**IMPORTANT:** In no way should this be considered a demonstration of how to secure a web site. The login and session
handling are not close to being good practice. The sole point was to understand and demonstrate how to implement an
Envoy external authorization filter and configure Envoy to invoke it.

### Preparing the `/etc/host` file for the testbed

Before running the custom Envoy image in [envoy-local](envoy-local), you must first add an entry to the `/etc/host`
file. This is in order for the [envoy-local/envoy.yaml](envoy-local/envoy.yaml) Envoy configuration to be able
to locate the external authorization filter outside the Envoy container (localhost / 127.0.0.1 would only route inside 
the container).

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

### Building and starting the Envoy container

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

#### How the [envoy-local Envoy](envoy-local/envoy.yaml) is configured

At the time of writing, the default demo Envoy configuration can be found at [envoy-demo.yaml](https://github.com/envoyproxy/envoy/blob/main/configs/envoy-demo.yaml).
This was downloaded and modified to create the [envoy-local/envoy.yaml](envoy-local/envoy.yaml) configuration in this
repository.

### Building and starting the `authtest`, `login`, and `extauth` containers

All three of these containers are running a Go application. The instructions for building and running the containers
are the same for all three.

To build and run the Go service container images:

1. Ensure that respective `Dockerfile` in each subdirectory references a current version of the Go language image;
   for example: `FROM golang:1.20-alpine`.
2. Open a terminal shell and change to the respective subdirectory
3. Run the following:
   ```shell
   make build
   make run
   ```

Once the container image has been built the first time, you can skip the `make build` step thereafter unless you have
modified the Go source files.

The `make run` command executes Docker to run the containers interactively, i.e. logging to the current shell,
and using the local host network (`--network host`). The container will be removed it is stopped with `ctrl-c`.

## Bonus: Performance testing RS256 signed JWT generation

The command line Go application located under the [rs256](rs256) demonstrates JWT generation and signing with the 
RS256 algorithm. 

A brief [README](rs256/README.md) in that directory explains how the code can be used to assess the time lag 
introduced by the signature calculations, both signing and verifying.
