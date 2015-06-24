# Go-Http

A simple command-line static file server written in the Go Programming Language.

An alternative to `python -m SimpleHTTPServer`

# Installation

It's as simple as 

```sh
$ go get -u github.com/GokulSrinivas/go-http/src/go-http/
```
# Usage 

To serve the current directory on `localhost:8080`, simply type 

```sh
$ go-http
```

Options :- 

	 -p   : Specify the port number
	 -d   : Specify subdirectory

## Specifying Port Number

```sh
$ go-http -p=1234
```

This serves the current directory to port 1234 of localhost

## Specifying subdirectory

```sh
$ go-http -d=sample
```

This serves the subdirectory `sample` to port 8080 of localhost

## Example Usage

```sh
$ go-http -p=4141 -d=sample
```

This serves the subdirectory `sample` to port 4141 of localhost

# Contribute

If you want to add features, improve them, or report issues, feel free to send a pull request!

# LICENSE

![GPL V3](https://raw.githubusercontent.com/GokulSrinivas/go-http/master/gpl.png)