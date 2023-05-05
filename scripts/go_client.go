package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	var n int
	n = 250
	rand.Seed(time.Now().UnixNano())
	iHost := Shuffle(NewSlice(1, n, 1))
	iDomain := Shuffle(NewSlice(0, n, 1))
	list := make([]string, 0, 1+(n*n))
	for _, d := range iDomain {
		for _, h := range iHost {
			url := "host" + strconv.Itoa(h) + ".domain" + strconv.Itoa(d) + ".rsx218-dox.cnam.fr\n"
			list = append(list, url)
		}
	}
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
	for _, url := range list {
		fmt.Print(url)
		cmd := exec.Command("dnslookup", url, "192.168.56.2:453")
		err := cmd.Start()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func NewSlice(start, end, step int) []int {
	if step <= 0 || end < start {
		return []int{}
	}
	s := make([]int, 0, 1+(end-start)/step)
	for start <= end {
		s = append(s, start)
		start += step
	}
	return s
}

func Shuffle(vals []int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]int, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}
