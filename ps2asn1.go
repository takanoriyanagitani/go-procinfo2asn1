package ps1asn1

import (
	"bufio"
	"encoding/asn1"
	"encoding/json"
	"io"
	"iter"
	"strconv"
	"strings"
)

type SimpleProcessInfo struct {
	ProcessId  int64  `json:"pid"`
	Rss        int64  `json:"rss"`
	CpuPercent int64  `json:"cpu"`
	MemPercent int64  `json:"mem"`
	State      string `asn1:"utf8" json:"state"`
	Time       string `asn1:"utf8" json:"time"`
}

type LtsvPair struct {
	Label string
	Value string
}

type LtsvLabel string

func (l LtsvLabel) SetProcId(value string, spi *SimpleProcessInfo) error {
	i, e := strconv.Atoi(value)
	spi.ProcessId = int64(i)
	return e
}

func (l LtsvLabel) SetRss(value string, spi *SimpleProcessInfo) error {
	i, e := strconv.Atoi(value)
	spi.Rss = int64(i)
	return e
}

func (l LtsvLabel) SetCpu(value string, spi *SimpleProcessInfo) error {
	i, e := strconv.ParseFloat(value, 64)
	spi.CpuPercent = int64(i)
	return e
}

func (l LtsvLabel) SetMem(value string, spi *SimpleProcessInfo) error {
	i, e := strconv.ParseFloat(value, 64)
	spi.MemPercent = int64(i)
	return e
}

func (l LtsvLabel) SetState(value string, spi *SimpleProcessInfo) error {
	spi.State = value
	return nil
}

func (l LtsvLabel) SetTime(value string, spi *SimpleProcessInfo) error {
	spi.Time = value
	return nil
}

type LtsvLine string

func (l LtsvLine) ToPairs() []string {
	return strings.Split(string(l), "	")
}

func (l LtsvLine) ToSpi() (SimpleProcessInfo, error) {
	var sp []string = l.ToPairs()
	var pairs iter.Seq[LtsvPair] = func(
		yield func(LtsvPair) bool,
	) {
		for _, pair := range sp {
			if !yield(PairString(pair).ToPair()) {
				return
			}
		}
	}

	var spi SimpleProcessInfo
	for pair := range pairs {
		var label string = pair.Label
		var val string = pair.Value

		var err error

		switch label {
		case "pid":
			err = LtsvLabel(label).SetProcId(val, &spi)
		case "rss":
			err = LtsvLabel(label).SetRss(val, &spi)
		case "cpu":
			err = LtsvLabel(label).SetCpu(val, &spi)
		case "mem":
			err = LtsvLabel(label).SetMem(val, &spi)
		case "state":
			err = LtsvLabel(label).SetState(val, &spi)
		case "time":
			err = LtsvLabel(label).SetTime(val, &spi)
		default:
		}

		if nil != err {
			return spi, err
		}
	}

	return spi, nil
}

type LtsvLines iter.Seq[string]

func (l LtsvLines) ToSimpleInfo() iter.Seq2[SimpleProcessInfo, error] {
	return func(yield func(SimpleProcessInfo, error) bool) {
		for line := range l {
			spi, e := LtsvLine(line).ToSpi()
			if !yield(spi, e) {
				return
			}
		}
	}
}

type PairString string

func (p PairString) ToPair() LtsvPair {
	var empty LtsvPair
	var splited []string = strings.SplitN(string(p), ":", 2)
	switch len(splited) {
	case 2:
		return LtsvPair{
			Label: splited[0],
			Value: splited[1],
		}
	default:
		return empty
	}
}

func (s SimpleProcessInfo) ToDer() ([]byte, error) {
	return asn1.Marshal(s)
}

type SpiJson []byte

func (j SpiJson) Parse() (SimpleProcessInfo, error) {
	var ret SimpleProcessInfo
	e := json.Unmarshal(j, &ret)
	return ret, e
}

type SpiArray []SimpleProcessInfo

func (a SpiArray) ToDer() ([]byte, error) {
	return asn1.Marshal(a)
}

func ReaderToLines(rdr io.Reader) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		var s *bufio.Scanner = bufio.NewScanner(rdr)
		for s.Scan() {
			var line []byte = s.Bytes()
			if !yield(line) {
				return
			}
		}
	}
}

func ReaderToStrings(rdr io.Reader) iter.Seq[string] {
	return func(yield func(string) bool) {
		var s *bufio.Scanner = bufio.NewScanner(rdr)
		for s.Scan() {
			var line string = s.Text()
			if !yield(line) {
				return
			}
		}
	}
}

func JsonLinesToProcInfo(
	jsons iter.Seq[[]byte],
) iter.Seq2[SimpleProcessInfo, error] {
	return func(yield func(SimpleProcessInfo, error) bool) {
		for line := range jsons {
			spi, e := SpiJson(line).Parse()
			if !yield(spi, e) {
				return
			}
		}
	}
}

func ProcsToDer(procs iter.Seq2[SimpleProcessInfo, error]) ([]byte, error) {
	var arr []SimpleProcessInfo
	for spi, e := range procs {
		if nil != e {
			return nil, e
		}
		arr = append(arr, spi)
	}
	return SpiArray(arr).ToDer()
}
