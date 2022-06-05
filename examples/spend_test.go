package examples

import (
	"bytes"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/spacemeshos/go-scale"
	"github.com/spacemeshos/go-scale/tester"
)

func FuzzSpendConsistency(f *testing.F) {
	tester.FuzzConsistency[Spend](f)
}

func FuzzSpendSafety(f *testing.F) {
	tester.FuzzSafety[Spend](f)
}

func BenchmarkSpend(b *testing.B) {
	fuzzer := fuzz.NewWithSeed(1001)
	spend := Spend{}
	fuzzer.Fuzz(&spend)

	b.Run("Encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf := bytes.NewBuffer(nil)
			encoder := scale.NewEncoder(buf)
			_, err := spend.EncodeScale(encoder)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Decode", func(b *testing.B) {
		buf := bytes.NewBuffer(nil)
		encoder := scale.NewEncoder(buf)
		_, err := spend.EncodeScale(encoder)
		if err != nil {
			b.Fatal(err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := bytes.NewBuffer(buf.Bytes())
			decoder := scale.NewDecoder(buf)
			var spend Spend
			_, err := spend.DecodeScale(decoder)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
