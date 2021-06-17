package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"

	"github.com/mmyj/mmyj-urltopn/top/util"
)

// used to generate test data

var (
	dataSize int
	ndv      int
)

const mb = 1000 * 1000

func flagParse() {
	flag.IntVar(&dataSize, "data size", 100, "data size of MB")
	flag.IntVar(&ndv, "ndv", 1000*1000, "number of distinct value")
}

func main() {
	flagParse()

	dataFile, err := os.Create("data")
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()
	answerFile, err := os.Create("answer")
	if err != nil {
		panic(err)
	}
	defer answerFile.Close()

	buff := bufio.NewWriter(dataFile)
	existUrlMap := make(map[string]int64)
	for w := 0; w < dataSize*mb; {
		r := rand.Intn(ndv)
		url := fmt.Sprintf("https://baidu.com/news/%d", r)
		existUrlMap[url]++
		n, err := buff.WriteString(url + "\n")
		if err != nil {
			panic(err)
		}
		w += n
	}
	err = buff.Flush()
	if err != nil {
		panic(err)
	}

	list := make([]util.Pair, 0, len(existUrlMap))
	for k, v := range existUrlMap {
		list = append(list, util.Pair{
			K: k,
			V: v,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Greater(list[j])
	})
	buff = bufio.NewWriter(answerFile)
	for i := 0; i < 100; i++ {
		if i >= len(list) {
			break
		}
		_, err := buff.WriteString(fmt.Sprintf("%s,%d\n", list[i].K, list[i].V))
		if err != nil {
			panic(err)
		}
	}
	err = buff.Flush()
	if err != nil {
		panic(err)
	}
}
