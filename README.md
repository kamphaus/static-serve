# static-serve
Satic-serve is a very simple and efficient static file server written in Go.

```
Usage:
	-p="8100": port to serve on
	-d=".":    the directory of static files to host
	-e="":     a file to serve in case of error 404
	-l:        log access requests
```

Static-serve does not show directory listings, it only serves files.
