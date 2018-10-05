# Diligent [![Build Status](https://travis-ci.org/senseyeio/diligent.svg?branch=master)](https://travis-ci.org/senseyeio/diligent) [![Docker Build Status](https://img.shields.io/docker/build/senseyeio/diligent.svg)](https://hub.docker.com/r/senseyeio/diligent/) [![GoDoc](https://godoc.org/github.com/senseyeio/diligent?status.svg)](https://godoc.org/github.com/senseyeio/diligent)

Get the licenses associated with your software dependencies. Enforce that the open source licenses used by your software are compatible with your business and software licensing.

 - Identify licenses across many languages using a single tool
 - Whitelist licenses or categories of license which are compatible with your software
 - Integrates with continuous integration to stop builds if non-whitelisted licenses are detected
 - Docker image available = super easy to run
 - Multiple output formats

## Language and Dependency Manager Support

The following languages and dependency managers are supported:

 - Go
   - govendor (vendor.json)
   - dep (Gopkg.lock)
 - Node / Javascript
   - NPM (package.json)

## Usage
The following command demonstrates how to use docker to run diligent:
```
docker run -v {project}:/dep senseyeio/diligent ls {path}
```
For example, if you had a node application at `/app` which contained a `package.json` file, you would run the following command:
```
docker run -v /app:/dep senseyeio/diligent ls .
```
Using diligent without docker is detailed later in the readme.

## Whitelisting

The `check` command can check that your depedencies' licenses match a given license whitelist.
Whitelisting is possible by specifying license identifiers or categories of licenses.
To see the identifiers and categories available please look at the [license definitions](https://github.com/senseyeio/diligent/blob/master/license.go).

For example, the following would whitelist all permissive licenses and in addition `GPL-3.0`:
```
docker run -v {project}:/dep senseyeio/diligent check -w GPL-3.0 -w permissive {path}
```

If licenses are found which do not match the specified whitelist, the application will return a non zero exit code (see exit code section below).
This is compatible with most CI solutions and can be used to stop builds if incompatible licenses are discovered.

To see what licenses you are whitelisting you can call the `whitelist` command:
```
docker run senseyeio/diligent whitelist -w GPL-3.0 -w permissive
```

If no `-w` flags are defined, diligent will always return a non zero exit code.

## Running Locally

The following requirements need to be satisfied when running locally:
 - `go` command line tool
 - `GOPATH` defined

The following assumes `$GOPATH/bin` is within your `PATH`:
```
go install github.com/senseyeio/diligent/cmd/diligent
```

Run the resulting binary as follows:
```
diligent {path}
```

## Exit codes

|Code|Explanation|
| ------------- | ------------- |
| 64  | Failed to determine licenses of one or more dependencies, however, dependencies which were successfully handled were acceptable  |
| 65  | Failed to output  |
| 66  | Failed to load provided file  |
| 67  | Fatal error when trying to determine licenses  |
| 68  | Discovered licenses do not match provided whitelist  |
| 69  | Could not process provided file  |
| 70  | The whitelist provided was invalid  |
