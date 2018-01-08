# Diligent [![Build Status](https://travis-ci.org/senseyeio/diligent.svg?branch=master)](https://travis-ci.org/senseyeio/diligent) [![GoDoc](https://godoc.org/github.com/senseyeio/diligent?status.svg)](https://godoc.org/github.com/senseyeio/diligent)

Get the licenses associated with your software dependencies. Enforce that the open source licenses used by your software are compatible with your business and software licensing.

 - Identify licenses across many languages using a single tool
 - Whitelist licenses or categories of license which are compatible with your software
 - Integrates with continuous integration to stop builds if non-whitelisted licenses are detected
 - Docker image available = super easy to run
 - Multiple output formats

## Usage

### With Docker

```
docker run -v {location of file}:/dep senseyeio/diligent {file name}
```
for instance, if I had a node application at `~/app` which contained a `package.json` file at `~/app/package.json`, run the following command:
```
docker run -v ~/app:/dep senseyeio/diligent package.json
```

### Locally

The following requirements need to be satisfied when running locally:
 - `go` command line toold
 - `GOPATH` defined

The following assumes `$GOPATH/bin` is within your `PATH`:
```
go install github.com/senseyeio/diligent
```

Run the resulting binary as follows:
```
dil {file path}
```

## Package Manager Support

The following languages and dependency managers are supported

 - Go
   - govendor (vendor.json)
 - Node / Javascript
   - NPM (package.json)

## Whitelisting

Whitelisting is possible by specifying license identifiers or categories of licenses.
To see the identifiers and categories available please look at the [license definitions](https://github.com/senseyeio/diligent/blob/master/license.go)

For example, the following would whitelist all permissive licenses and in addition `GPL-3.0`:
```
docker run -v {location of file}:/dep senseyeio/diligent -w GPL-3.0 -w permissive {file name}
```

If licenses are found which do not match the specified whitelist, the application will return a non zero exit code.
This is compatible with most CI solutions and can be used to stop builds if incompatible licenses are discovered.
