package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
)

const (
	mb             = 1000 * 1000
	TopN           = 100
	DefaultBufSize = 4096
)

type Pair struct {
	K string
	V int64
}

func (p1 Pair) Greater(p2 Pair) bool {
	if p1.V != p2.V {
		return p1.V > p2.V
	}
	return p1.K < p2.K
}

var TempDir = os.TempDir() + "top-url"

type SolutionFunc func(memConsume *MemConsumer) []Pair

type MemConsumer struct {
	max   int64
	total int64
}

func (m *MemConsumer) Consume(d int64) {
	if m != nil {
		m.total = m.total + d
		if m.total > m.max {
			m.max = m.total
		}
	}
}

func SaveAnswer(solutionNo int, list []Pair) (fileName string) {
	path := TempDir + fmt.Sprintf("my-answer-%d", solutionNo)
	answerFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer answerFile.Close()

	wBuff := bufio.NewWriter(answerFile)
	for i := 0; i < 100; i++ {
		if i >= len(list) {
			break
		}
		_, err := wBuff.WriteString(fmt.Sprintf("%s,%d\n", list[i].K, list[i].V))
		if err != nil {
			panic(err)
		}
	}
	wBuff.Flush()
	return path
}

func CheckAnswer(myAnsFileName string) (ret bool, desc string) {
	myAnsFile, err := os.Open(myAnsFileName)
	if err != nil {
		panic(err)
	}
	defer myAnsFile.Close()
	ansFile, err := os.Open("answer")
	if err != nil {
		panic(err)
	}
	defer ansFile.Close()
	var (
		rBuff1 = bufio.NewReader(ansFile)
		rBuff2 = bufio.NewReader(myAnsFile)
	)
	for {
		line1, _, err1 := rBuff1.ReadLine()
		if err1 != nil && err1 != io.EOF {
			panic(err)
		}
		line2, _, err2 := rBuff2.ReadLine()
		if err2 != nil && err2 != io.EOF {
			panic(err)
		}
		if (err1 == io.EOF && err2 != io.EOF) || (err2 == io.EOF && err1 != io.EOF) {
			return false, "line of answer file is not match"
		}
		if err1 == io.EOF && err2 == io.EOF {
			break
		}
		str1 := string(line1)
		str2 := string(line2)
		if str1 != str2 {
			return false, fmt.Sprintf("want %s, but %s", str1, str2)
		}
	}
	return true, ""
}

func AssertSolutionIsRight(t *testing.T, solutionNo int, solutionFn SolutionFunc) {
	myAnsFileName := SaveAnswer(solutionNo, solutionFn(nil))
	defer os.Remove(myAnsFileName)
	ok, desc := CheckAnswer(myAnsFileName)
	if !ok {
		t.Fatal(fmt.Sprintf("SolutionFunc-%d failed, %s", solutionNo, desc))
	}
}

func GetCpuPprof(sol SolutionFunc) {
	fcpu, err := os.Create("cpu.prof")
	defer fcpu.Close()
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(fcpu)
	sol(nil)
	pprof.StopCPUProfile()
}

func GetMemPprof(sol SolutionFunc) {
	runtime.GC()
	fmem, err := os.Create("mem.prof")
	defer fmem.Close()
	if err != nil {
		panic(err)
	}
	pprof.WriteHeapProfile(fmem)
	sol(nil)
}

func PrintMemUsed(t *testing.T, sol SolutionFunc) {
	var mc MemConsumer
	sol(&mc)
	t.Logf("total: %0.2f(mb), max: %0.2f(mb)", float64(mc.total)/mb, float64(mc.max)/mb)
}

func Benchmark(b *testing.B, solutionFn SolutionFunc) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solutionFn(nil)
	}
}
