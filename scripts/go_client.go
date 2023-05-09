package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
	fmt.Println("client,reqID,url,timestamp")
	var counter int = 1
	for _, url := range list {
		clean_url := CleanStr(url)
		fmt.Printf("%s,%d,%s,%s\n", hostname, counter, clean_url, GetTimeMs())
		cmd := exec.Command("dnslookup", clean_url, "192.168.56.2:453")
		err := cmd.Start()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		counter += 1
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

func GetTimeMs() string {
	return time.Now().Format(time.StampMilli)
}

func CleanStr(str string) string {
	str = strings.ReplaceAll(str, "\n", "")
	return str
}
