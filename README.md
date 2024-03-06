# Put a Map on Me!

Website source for <https://mapon.me>

The code here builds both the static assets and the API runtime for this site.

## Building locally

The development build script is **_dev.sh_** and the entry point is http://localhost:8080/

Your container will require that a Go compiler is installed ([download](https://go.dev/dl/)) and your JS packages updated:

```sh
$ npm install
```

## Testing locally

Running the API functional tests:

```sh
$ go test *[^v].go --test.short
```

Running the API integration tests:

```sh
$ pip install testDependencies.pip
$ pytest test/api.py
```

Running the browser tests:

```sh
$ pip install testDependencies.pip
$ py test/firefox.py --watch
```

## Deploying

The production build script is **_deploy.sh_**

If you choose to launch a modified version of this source, it is recommended to start by using the following route mappings behind a reverse proxy:

- https://hostname/
  - /dist/static/index/
- https://hostname/s/
  - /dist/static/
- https://hostname/api/
  - /dist/apiServer
