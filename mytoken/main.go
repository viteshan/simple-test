package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {

	type Rank struct {
		Symbol   string `json:"symbol"`
		DevRank  int    `json:"rank"`
		MarkRank int    `json:"market_rank"`
	}
	//coin, err := ioutil.ReadFile("/Users/jie/Documents/vite/src/github.com/viteshan/simple-test/mytoken/coin.json")
	//if err != nil {
	//	panic(err)
	//}

	dev, err := ioutil.ReadFile("/Users/jie/Documents/vite/src/github.com/viteshan/simple-test/mytoken/dev.json")
	if err != nil {
		panic(err)
	}

	var devs []*Rank

	//json.Unmarshal(coin, &coins)
	json.Unmarshal(dev, &devs)
	//
	//coinM := make(map[string]*Rank)
	//devM := make(map[string]*Rank)
	//
	//for _, v := range coins {
	//	fmt.Print(v.Symbol, ":", v.Rank, " ")
	//	coinM[v.Symbol] = v
	//}
	//fmt.Println()
	//
	//for _, v := range devs {
	//	fmt.Print(v.Symbol, ":", v.Rank, " ")
	//	devM[v.Symbol] = v
	//	_, ok := coinM[v.Symbol]
	//	if !ok {
	//		coinM[v.Symbol] = &Rank{Symbol: v.Symbol, Rank: -1}
	//	}
	//}
	//fmt.Println()
	//
	//for k, v := range coinM {
	//	if v.Rank == -1 {
	//		fmt.Println(k, devM[k].Rank)
	//	}
	//}

	for _, v := range devs {
		fmt.Println(v.DevRank, v.Symbol, v.MarkRank, v.DevRank-v.MarkRank)
	}
}
