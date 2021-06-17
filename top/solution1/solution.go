package solution1

import (
	"bufio"
	"io"
	"os"
	"sort"
	"unsafe"

	"github.com/mmyj/mmyj-urltopn/top/util"
)

// Solution 顺序执行，不限制内存大小，最笨方法
func Solution(mc *util.MemConsumer) []util.Pair {
	dataFile, err := os.Open("data")
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	existUrlMap := make(map[string]int64)
	rBuff := bufio.NewReaderSize(dataFile, util.DefaultBufSize)
	mc.Consume(util.DefaultBufSize)
	defer mc.Consume(-util.DefaultBufSize)
	existUrlMapMem := int64(0)
	for {
		line, _, err := rBuff.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err.Error())
		}
		str := string(line)
		if _, ok := existUrlMap[str]; !ok {
			mem := int64(len(str)) + 8
			mc.Consume(mem)
			existUrlMapMem += mem
		}
		existUrlMap[str]++
	}

	list := make([]util.Pair, 0)
	mc.Consume(int64(unsafe.Sizeof(util.Pair{})) * int64(len(existUrlMap)))
	for k, v := range existUrlMap {
		list = append(list, util.Pair{
			K: k,
			V: v,
		})
	}
	mc.Consume(-existUrlMapMem)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Greater(list[j])
	})
	return list
}
