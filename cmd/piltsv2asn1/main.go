package main

import (
	"context"
	"iter"
	"log"
	"os"

	pa "github.com/takanoriyanagitani/go-procinfo2asn1"
	. "github.com/takanoriyanagitani/go-procinfo2asn1/util"
)

var bytes2stdout func([]byte) IO[Void] = func(dat []byte) IO[Void] {
	return func(_ context.Context) (Void, error) {
		_, e := os.Stdout.Write(dat)
		return Empty, e
	}
}

var stdin2ltsv iter.Seq[string] = pa.ReaderToStrings(os.Stdin)

var spiArray iter.Seq2[pa.SimpleProcessInfo, error] = pa.LtsvLines(stdin2ltsv).
	ToSimpleInfo()

var procDerBytes IO[[]byte] = Bind(
	Of(spiArray),
	Lift(pa.ProcsToDer),
)

var stdin2ltsvs2der2stdout IO[Void] = Bind(procDerBytes, bytes2stdout)

func main() {
	_, e := stdin2ltsvs2der2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
