package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type idolResult struct {
	name  string
	rank  int
	votes int
}

type resultType int

const (
	equalRank resultType = iota
	dereGreater
	mobaGreater
)

type sendResult struct {
	name  string
	types resultType
	diff  int
}

func readFile(filepath string) []idolResult {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	res := []idolResult{}

	scanner := bufio.NewScanner(file)
	rankcount := 0
	for scanner.Scan() {
		rankcount++
		line := scanner.Text()
		strs := strings.Split(line, ",")
		votes, _ := strconv.Atoi(strs[1])
		res = append(res, idolResult{strs[0], rankcount, votes})
	}
	return res
}

func compareRank(idol idolResult, datas []idolResult, resChan chan sendResult) {
	for _, data := range datas {
		if idol.name == data.name {
			if idol.rank > data.rank {
				resChan <- sendResult{idol.name, mobaGreater, idol.rank - data.rank}
			} else if idol.rank == data.rank {
				resChan <- sendResult{idol.name, equalRank, 0}
			} else {
				resChan <- sendResult{idol.name, dereGreater, data.rank - idol.rank}
			}
			return
		}
	}
}

func sorter(sl []sendResult) {
	sort.Slice(sl, func(i, j int) bool {
		return sl[i].diff > sl[j].diff
	})
}

func main() {
	var wg sync.WaitGroup
	dere := readFile("./star_distinct.txt")
	moba := readFile("./moba_distinct.txt")
	//all := readFile("./all_distinct.txt")

	maxc := make(chan int, 30)
	resultChan := make(chan sendResult, len(dere))
	for _, v := range dere {
		wg.Add(1)
		go func(data idolResult) {
			maxc <- 1
			defer func() {
				<-maxc
				wg.Done()
			}()
			compareRank(data, moba, resultChan)
		}(v)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	equals := []sendResult{}
	deres := []sendResult{}
	mobas := []sendResult{}

	for res := range resultChan {
		switch res.types {
		case equalRank:
			equals = append(equals, res)
		case dereGreater:
			deres = append(deres, res)
		case mobaGreater:
			mobas = append(mobas, res)
		}
	}

	sorter(equals)
	sorter(deres)
	sorter(mobas)

	fmt.Println("equal")
	for _, v := range equals {
		fmt.Println(v.name)
	}
	fmt.Println()

	fmt.Println("dere")
	for i := 0; i < len(dere); i++ {
		if deres[i].diff < 10 {
			break
		}
		fmt.Println(deres[i].name, deres[i].diff)
	}
	fmt.Println()

	fmt.Println("moba")
	for i := 0; i < len(moba); i++ {
		if mobas[i].diff < 10 {
			break
		}
		fmt.Println(mobas[i].name, mobas[i].diff)
	}

}
