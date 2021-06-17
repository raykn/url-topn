package main

import (
	"testing"

	"github.com/mmyj/mmyj-urltopn/top/solution1"
	"github.com/mmyj/mmyj-urltopn/top/solution2"
	"github.com/mmyj/mmyj-urltopn/top/util"
)

//util.go:140: total: 23.15(mb), max: 58.74(mb)
func TestSolution1(t *testing.T) {
	fn := util.SolutionFunc{Fn: solution1.Solution, Title: "solution1"}
	util.AssertSolutionIsRight(t, fn)
	util.PrintMemUsed(t, fn)
}

//cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
//BenchmarkSolution1
//BenchmarkSolution1-8   	       1	1305060876 ns/op	254068256 B/op	 3384023 allocs/op
func BenchmarkSolution1(b *testing.B) {
	fn := util.SolutionFunc{Fn: solution1.Solution, Title: "solution1"}
	util.Benchmark(b, fn)
}

//util.go:147: SolutionFunc-SolutionBatch256 total: 0.00(mb), max: 1.22(mb)
func TestSolution2_256(t *testing.T) {
	fn := util.SolutionFunc{Fn: solution2.SolutionBatch256, Title: "SolutionBatch256"}
	util.AssertSolutionIsRight(t, fn)
	util.PrintMemUsed(t, fn)
}

//util.go:147: SolutionFunc-SolutionBatch1024 total: 0.00(mb), max: 4.31(mb)
func TestSolution2_1024(t *testing.T) {
	fn := util.SolutionFunc{Fn: solution2.SolutionBatch1024, Title: "SolutionBatch1024"}
	util.AssertSolutionIsRight(t, fn)
	util.PrintMemUsed(t, fn)
}

//cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
//BenchmarkSolution2
//BenchmarkSolution2-8   	       1	1526941482 ns/op	233154144 B/op	 3386936 allocs/op
func BenchmarkSolution2_256(b *testing.B) {
	fn := util.SolutionFunc{Fn: solution2.SolutionBatch256, Title: "SolutionBatch256"}
	util.Benchmark(b, fn)
}

//cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
//BenchmarkSolution2_1024
//BenchmarkSolution2_1024-8   	       1	1828375873 ns/op	236878016 B/op	 3397730 allocs/op
func BenchmarkSolution2_1024(b *testing.B) {
	fn := util.SolutionFunc{Fn: solution2.SolutionBatch1024, Title: "SolutionBatch1024"}
	util.Benchmark(b, fn)
}
