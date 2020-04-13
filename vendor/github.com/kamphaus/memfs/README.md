memfs
=====

Implementation of http.FileSystem where the files stay in memory.<br>
It uses [<b>fsnotify</b>](https://github.com/howeyc/fsnotify) to keep the cache updated.

Example:
<code><pre>
github.com/ranveerkunal/memfs/example $ go build memfs_code.go
github.com/ranveerkunal/memfs/example $ ./memfs_code
</pre></code>

[http://localhost:9999/memfs](http://localhost:9999/memfs)

```go
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ranveerkunal/memfs"
)

func main() {
	path := flag.String("path", "./", "")
	addr := flag.String("addr", ":9999", "")
	verbose := flag.Bool("verbose", true, "")
	flag.Parse()

	fs, err := memfs.New(*path)
	if err != nil {
		log.Fatalf("Failed to create memfs: %s err: %v", *path, err)
	}

	if (*verbose) {
		log.Printf("logging to stderr ...")
		memfs.SetLogger(memfs.Verbose)
	}

	http.Handle("/memfs/", http.StripPrefix("/memfs/", http.FileServer(fs)))

	log.Printf("path: %s addr:%s", *path, *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
```
<pre>
Benchmark on mac: darwin 64
~/ranveerkunal/memfs % go test memfs_test.go -bench=. -cpu=4 -parallel=4
temp dir: /tmp/memfs406771321
writing small file
writing big file
ready to benchmark ...
testing: warning: no tests to run
PASS
BenchmarkNonExistentMemFS-4      5000000               700 ns/op
BenchmarkNonExistentDiskFS-4      500000              3996 ns/op
BenchmarkSmallFileMemFS-4          10000            111634 ns/op
BenchmarkSmallFileDiskFS-4         10000            128475 ns/op
BenchmarkBigFileMemFS-4               20          83455262 ns/op
BenchmarkBigFileDiskFS-4              20          96320175 ns/op
ok      command-line-arguments  26.610s
</pre>

<pre>
Benchmark on linux:
~/kamphaus/memfs % go test memfs_test.go -bench=. -cpu=4 -parallel=4
temp dir: /tmp/memfs018986335
writing small file
writing medium file
writing big file
ready to benchmark ...
goos: linux
goarch: amd64
BenchmarkNonExistentMemFS-4    	 5259802	       223 ns/op
BenchmarkNonExistentDiskFS-4   	  870788	      1422 ns/op
BenchmarkSmallFileMemFS-4      	   21523	     54604 ns/op
BenchmarkSmallFileDiskFS-4     	   18241	     65965 ns/op
BenchmarkMediumFileMemFS-4     	     426	   2790326 ns/op
BenchmarkMediumFileDiskFS-4    	     458	   2603583 ns/op
BenchmarkBigFileMemFS-4        	      24	  46187587 ns/op
BenchmarkBigFileDiskFS-4       	      33	  35086878 ns/op
PASS
ok  	command-line-arguments	19.123s
</pre>

It looks it's a speedup for small or non existent files.
As soon as there are medium sized files (5 MiB) it's more efficient if the kernel takes care of caching the files.

