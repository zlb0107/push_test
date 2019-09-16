package random_list

import (
	"bufio"
	//	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type RandomListInfo struct {
	List [][]uint
}

var RandomList RandomListInfo

func init() {
	RandomList.GetMap("conf/random.map")
}
func (this *RandomListInfo) GetMap(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		tempStrs := strings.Split(line, " ")
		tempList := make([]uint, len(tempStrs))
		for idx, tempStr := range tempStrs {
			tempNum, _ := strconv.ParseInt(tempStr, 10, 10)
			tempList[idx] = uint(tempNum)
		}
		if this.List == nil {
			this.List = make([][]uint, 0)
		}
		this.List = append(this.List, tempList)
	}
}
func (this *RandomListInfo) GetRandomList() []uint {
	rand.Seed(time.Now().UnixNano())
	line := rand.Intn(len(this.List))
	col := rand.Intn(len(this.List[line]))
	return append(this.List[line][col:], this.List[line][:col]...)
}

func (this *RandomListInfo) RandomShuffle(list []int, randomNum int) []int {
	randomList := this.GetRandomList()
	tempList := make([]int, 0, randomNum)
	posMap := make(map[int]bool)
	if len(list) > randomNum {
		for i := 0; i < randomNum; i++ {
			pos := int(randomList[i])
			if pos >= len(list) {
				pos = pos % len(list)
			}

			_, isIn := posMap[pos]
			for isIn {
				pos++
				pos = pos % len(list)
				_, isIn = posMap[pos]
			}

			tempList = append(tempList, list[pos])
			posMap[pos] = true
		}
	} else {
		tempList = list
	}

	return tempList
}
