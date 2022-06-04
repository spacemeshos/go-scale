SCALE
===

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

`//go:generate scalegen` will discover all struct types and derive EncodeScale/DecodeScale methods. To avoid structs autodiscovery use `-types=U8,U16`.
