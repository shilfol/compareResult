package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type idolResult struct {
	name  string
	rank  int
	votes int
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

func compareRank(idol idolResult, datas []idolResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, data := range datas {
		if idol.name == data.name {
			if idol.rank > data.rank {
				fmt.Println(idol.name, ": dere", idol.rank, " < ", data.rank, "moba, diff:", idol.rank-data.rank)
			} else if idol.rank == data.rank {
				fmt.Println(idol.name, ": dere", idol.rank, " = ", data.rank, "moba")
			} else {
				fmt.Println(idol.name, ": dere", idol.rank, " > ", data.rank, "moba, diff:", data.rank-idol.rank)
			}
			return
		}
	}
}

func main() {
	var wg sync.WaitGroup

	dere := readFile("./star_distinct.txt")
	moba := readFile("./moba_distinct.txt")
	//all := readFile("./all_distinct.txt")

	for _, v := range dere {
		wg.Add(1)
		go compareRank(v, moba, &wg)
	}

	wg.Wait()

}
