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

Types
---

golang      | notes
------------|-------------------------------------------------------------------
[]byte      | length prefixed byte array with length as u32 compact integer
string      | same as []byte
[...]byte   | appended to the result
bool        | 1 byte, 0 for false, 1 for true
Object{}    | concatenation of fields
*Object{}   | Option. 0 for nil, 1 for Object{}. if 1 - decode Object{}
uint8       | compact u8 [TODO no need for compact u8]
uint16      | compact u16
uint32      | compact u32
uint32      | compact u64
[...]Object | array with objects. encoded by consecutively encoding every object
[]Object    | slice with objects. prefixed with compact u32

Not implemented:
- pointers to arrays and slices
- slices with pointers
- enumerations
- fixed width integers

Code generation
---

```
go install ./scalegen
```

And see examples across different modules, e.g:

```
//go:generate scalegen -pkg examples -file ex1_scale.go -types Ex1 -imports github.com/spacemeshos/go-scale/examples
```