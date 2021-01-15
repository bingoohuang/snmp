package util

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type expandNums struct {
	nums []string
}

func ExpandNums(x []string) []string {
	ox := expandNums{nums: make([]string, 0)}
	for _, xi := range x {
		ox.expand(xi)
	}

	return ox.nums
}

func (ox *expandNums) expand(xi string) {
	xii := strings.Split(xi, ",")
	for _, xij := range xii {
		if !strings.Contains(xij, "-") {
			ox.append(xij)
			continue
		}

		xirange := strings.Split(xij, "-")
		f, err := strconv.Atoi(xirange[0])
		if err != nil {
			log.Printf("W! error x values %v", err)
			continue
		}

		ox.append(xirange[0])

		xirange = xirange[1:]
		if len(xirange) == 0 {
			continue
		}

		ox.expandRange(f, xirange)
	}
}

func (ox *expandNums) append(f string) {
	for _, i := range ox.nums {
		if i == f {
			return
		}
	}

	ox.nums = append(ox.nums, f)
}

func (ox *expandNums) expandRange(f int, to []string) {
	t, err := strconv.Atoi(to[0])
	if err != nil {
		log.Printf("W! error x values %v", err)
		return
	}

	to = to[1:]
	step := 1
	if len(to) > 0 {
		v, err := strconv.Atoi(to[0])
		if err != nil {
			log.Printf("W! error x values %v", err)
			return
		}
		step = v
	}

	for j := f + step; j <= t; j += step {
		ox.append(fmt.Sprintf("%d", j))
	}
}
