# SCALE

## Types

| golang      | notes                                                              | example                                               |
| ----------- | ------------------------------------------------------------------ | ----------------------------------------------------- |
| []byte      | length prefixed byte array with length as u32 compact integer      |
| string      | same as []byte                                                     |
| [...]byte   | appended to the result                                             |
| bool        | 1 byte, 0 for false, 1 for true                                    |
| Object{}    | concatenation of fields                                            |
| *Object{}   | Option. 0 for nil, 1 for Object{}. if 1 - decode Object{}          |
| uint8       | compact u8 [TODO no need for compact u8]                           |
| uint16      | compact u16                                                        |
| uint32      | compact u32                                                        |
| uint34      | compact u64                                                        |
| *uint8      | Option (0 for nil, 1 otherwise), followed by fixed-width u8        | `&255 -> 01FF`                                        |
| *uint16     | Option (0 for nil, 1 otherwise), followed by compact u16           | `&255 -> 01FD03`                                      |
| *uint32     | Option (0 for nil, 1 otherwise), followed by compact u32           | `&255 -> 01FD03`                                      |
| *uint34     | Option (0 for nil, 1 otherwise), followed by compact u64           | `&255 -> 01FD03`                                      |
| []uint16    | length prefixed (compact u32) followed by compact u16s             | `[4, 15, 23, u16::MAX] -> 10103C5CFEFF0300`           |
| []uint32    | length prefixed (compact u32) followed by compact u32s             | `[4, 15, 23, u32::MAX] -> 10103C5C03FFFFFFFF`         |
| []uint64    | length prefixed (compact u32) followed by compact u64s             | `[4, 15, 23, u64::MAX] -> 10103C5C13FFFFFFFFFFFFFFFF` |
| [...]Object | array with objects. encoded by consecutively encoding every object |
| []Object    | slice with objects. prefixed with compact u32                      |

Not implemented:

- pointers to arrays and slices
- slices with pointers
- enumerations
- fixed width integers

## Code generation

```bash
go install github.com/spacemeshos/go-scale/scalegen
```

`//go:generate scalegen` will discover all struct types and derive EncodeScale/DecodeScale methods. To avoid struct auto-discovery use `-types=U8,U16`.
