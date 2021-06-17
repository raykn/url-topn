package main

import (
	"testing"

	"github.com/mmyj/mmyj-urltopn/top/solution1"
	"github.com/mmyj/mmyj-urltopn/top/solution2"
	"github.com/mmyj/mmyj-urltopn/top/util"
)

//util.go:140: total: 23.15(mb), max: 58.74(mb)
func TestSolution1(t *testing.T) {
	fn := solution1.Solution
	util.AssertSolutionIsRight(t, 1, fn)
	util.PrintMemUsed(t, fn)
}

//cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
//BenchmarkSolution1
//BenchmarkSolution1-8   	       1	1305060876 ns/op	254068256 B/op	 3384023 allocs/op
func BenchmarkSolution1(b *testing.B) {
	fn := solution1.Solution
	util.Benchmark(b, fn)
}

//util.go:144: total: 0.00(mb), max: 2.26(mb)
func TestSolution2(t *testing.T) {
	fn := solution2.Solution
	util.AssertSolutionIsRight(t, 2, fn)
	util.PrintMemUsed(t, fn)
}

//cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
//BenchmarkSolution2
//BenchmarkSolution2-8   	       1	1526941482 ns/op	233154144 B/op	 3386936 allocs/op
func BenchmarkSolution2(b *testing.B) {
	fn := solution2.Solution
	util.Benchmark(b, fn)
}
