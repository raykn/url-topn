package solution3

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"unsafe"

	"github.com/mmyj/mmyj-urltopn/top/util"
)

// baseSolution 在solution2的基础上，开个routine去异步读取文件，临时文件统计的时候扇出统计
func baseSolution(mc *util.MemConsumer, batch int64, workerCount int) []util.Pair {
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

	lineChan := make(chan string, buffSize)
	fetchLine := func(ch chan<- string) {
		for {
			line, _, err := rBuff.ReadLine()
			if err != nil {
				if err == io.EOF {
					close(ch)
					return
				}
				panic(err)
			}
			ch <- string(line)
		}
	}
	go fetchLine(lineChan)

	for line := range lineChan {
		sumBytes := md5.Sum([]byte(line))
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

	pairChan := make(chan util.Pair, buffSize)
	taskChan := make(chan int, batch/2)
	var wg sync.WaitGroup
	masterFunc := func(taskCh chan<- int) {
		for i := 0; i < int(batch); i++ {
			_, err = tempFiles[i].Seek(0, io.SeekStart)
			if err != nil {
				panic(err)
			}
			taskCh <- i
		}
		close(taskCh)
		wg.Wait()
		close(pairChan)
	}
	workerFunc := func(workerMc *util.MemConsumer, taskCh <-chan int, pairCh chan<- util.Pair) {
		var workerBuff *bufio.Reader
		for task := range taskCh {
			if workerBuff == nil {
				workerBuff = bufio.NewReaderSize(tempFiles[task], buffSize)
			} else {
				workerBuff.Reset(tempFiles[task])
			}
			existUrlMap := make(map[string]int64)
			existUrlMapMem := int64(0)
			for {
				line, _, err := workerBuff.ReadLine()
				if err != nil {
					if err == io.EOF {
						break
					}
					panic(err)
				}
				str := string(line)
				if _, ok := existUrlMap[str]; !ok {
					mem := int64(len(str)) + 8
					workerMc.Consume(mem)
					existUrlMapMem += mem
				}
				existUrlMap[str]++
			}
			for k, v := range existUrlMap {
				p := util.Pair{K: k, V: v}
				pairCh <- p
			}
			workerMc.Consume(-existUrlMapMem)
		}
		wg.Done()
	}

	wg.Add(workerCount)
	workerMc := make([]util.MemConsumer, workerCount)
	mc.Consume(int64(unsafe.Sizeof(util.MemConsumer{})) * int64(workerCount))
	defer func() {
		mc.Consume(-int64(unsafe.Sizeof(util.MemConsumer{})) * int64(workerCount))
	}()
	for i := 0; i < workerCount; i++ {
		go workerFunc(&workerMc[i], taskChan, pairChan)
	}
	go masterFunc(taskChan)

	myHeap := util.NewPairHeap(util.TopN)
	for p := range pairChan {
		myHeap.TryPush(p)
	}
	// add max mem
	for i := 0; i < workerCount; i++ {
		mc.Consume(workerMc[i].Max())
		mc.Consume(-workerMc[i].Max())
	}

	list := myHeap.Data()
	sort.Slice(list, func(i, j int) bool {
		return list[i].Greater(list[j])
	})
	mc.Consume(int64(unsafe.Sizeof(util.Pair{})) * int64(len(list)))
	return list
}

func SolutionWorker16(mc *util.MemConsumer) []util.Pair {
	return baseSolution(mc, 256, 16)
}
