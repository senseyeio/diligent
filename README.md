# Diligent

Get the licenses associated with your software dependencies.

## Package Manager Support

The following languages and dependency managers are supported

### Go

 - govendor (Requires `go` command line tool to be available and `GOPATH` to be defined)

### Node / Javascript

 - npm

## Usage

### With Docker

```
docker run -v {location of file}:/dep senseyeio/diligent {file name}
```
for instance, if I had a node application at `~/app` which contained a `package.json` file at `~/app/package.json`, I could run the following command:
```
docker run -v ~/app:/dep senseyeio/diligent package.json
```

### Locally

Checkout diligent to your go path. Go get and build the application. Run the resulting binary as follows:

```
dil {file path}
```
Given the example above, run the following command:
```
dil ~/app/package.json
```