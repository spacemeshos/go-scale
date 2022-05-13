SCALE
===


```
go test ./transactions/ -run=xx -bench=BenchmarkSelfSpawn -v -benchmem
goos: linux
goarch: amd64
pkg: github.com/spacemeshos/go-scale/transactions
cpu: 12th Gen Intel(R) Core(TM) i9-12900KF
BenchmarkSelfSpawn
BenchmarkSelfSpawn/Encode
BenchmarkSelfSpawn/Encode-24         	 3286465	       356.6 ns/op	     304 B/op	       4 allocs/op
BenchmarkSelfSpawn/EncodeReuse
BenchmarkSelfSpawn/EncodeReuse-24    	24663522	        46.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkSelfSpawn/Decode
BenchmarkSelfSpawn/Decode-24         	 8302658	       149.1 ns/op	     160 B/op	       1 allocs/op
BenchmarkSelfSpawn/EncodeXDR
BenchmarkSelfSpawn/EncodeXDR-24      	  793075	      1512 ns/op	     776 B/op	      25 allocs/op
BenchmarkSelfSpawn/DecodeXDR
BenchmarkSelfSpawn/DecodeXDR-24      	  712395	      1844 ns/op	    1016 B/op	      27 allocs/op
PASS
ok  	github.com/spacemeshos/go-scale/transactions	6.678s
```

