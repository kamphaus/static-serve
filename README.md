# static-serve
Static-serve is a very simple and efficient static file server written in Go.

```
Usage:
	-p: ports to serve on (default: 8100)
	-d: the directories of static files to host (default: ./)
	-e: the files to serve in case of error 404 (- to disable error404 handler)
	-l: log access requests
	-hport: the port on which /health and /ready endpoints should be served
	-r	log request/response headers
	-v: verbose logging (e.g. when handling error 404)
Ports, directories and error404 flags can be specified multiple times,
but need to be specified the same amount of times.
```

Static-serve does not show directory listings, it only serves files.

Since static-serve uses Go's `http.FileServer` we have the following features
out of the box:
* basic mime type detection
* caching headers
* range requests
