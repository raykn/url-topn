package solution2

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"unsafe"

	"github.com/mmyj/mmyj-urltopn/top/util"
)

// baseSolution 类似mapreduce处理
func baseSolution(mc *util.MemConsumer, batch int64) []util.Pair {
	dataFile, err := os.Open("data")
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	const (
		buffSize = 4096
	)
	tempFiles := make([]*os.File, batch)
	tempFileBuffWrites := make([]*bufio.Writer, batch)
	mc.Consume(int64(unsafe.Sizeof(os.File{})) * batch)
	mc.Consume(int64(unsafe.Sizeof(bufio.Writer{})) * batch)
	mc.Consume(buffSize * batch)

	defer func() {
		mc.Consume(-int64(unsafe.Sizeof(os.File{})) * batch)
		mc.Consume(-int64(unsafe.Sizeof(bufio.Writer{})) * batch)
		mc.Consume(-buffSize * batch)
	}()

	for i := 0; i < int(batch); i++ {
		tempFile := util.TempDir + fmt.Sprintf("temp-%d", i)
		tempFiles[i], err = os.Create(tempFile)
		if err != nil {
			panic(err)
		}
		defer func(id int) {
			tempFiles[id].Close()
			os.Remove(tempFile)
		}(i)
		tempFileBuffWrites[i] = bufio.NewWriterSize(tempFiles[i], buffSize)
	}

	rBuff := bufio.NewReaderSize(dataFile, util.DefaultBufSize)
	mc.Consume(util.DefaultBufSize)
	defer mc.Consume(-util.DefaultBufSize)
	for {
		line, _, err := rBuff.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		sumBytes := md5.Sum(line)
		sum := binary.BigEndian.Uint32(sumBytes[:4])
		buff := tempFileBuffWrites[int64(sum)%batch]
		_, err = buff.WriteString(string(line) + "\n")
		if err != nil {
			panic(err)
		}
	}
	for i := 0; i < int(batch); i++ {
		err = tempFileBuffWrites[i].Flush()
		if err != nil {
			panic(err)
		}
	}

	myHeap := util.NewPairHeap(util.TopN)
	for i := 0; i < int(batch); i++ {
		_, err = tempFiles[i].Seek(0, io.SeekStart)
		if err != nil {
			panic(err)
		}
		rBuff.Reset(tempFiles[i])
		existUrlMap := make(map[string]int64)
		existUrlMapMem := int64(0)
		for {
			line, _, err := rBuff.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			str := string(line)
			if _, ok := existUrlMap[str]; !ok {
				mem := int64(len(str)) + 8
				mc.Consume(mem)
				existUrlMapMem += mem
			}
			existUrlMap[str]++
		}
		for k, v := range existUrlMap {
			p := util.Pair{K: k, V: v}
			myHeap.TryPush(p)
		}
		mc.Consume(-existUrlMapMem)
	}
	list := myHeap.Data()
	sort.Slice(list, func(i, j int) bool {
		return list[i].Greater(list[j])
	})
	mc.Consume(int64(unsafe.Sizeof(util.Pair{})) * int64(len(list)))
	return list
}

func SolutionBatch256(mc *util.MemConsumer) []util.Pair {
	return baseSolution(mc, 256)
}

func SolutionBatch1024(mc *util.MemConsumer) []util.Pair {
	return baseSolution(mc, 1024)
}
