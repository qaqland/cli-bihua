package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
)

// 这个地方要设置大小
type bh [20]struct {
	rank      int
	character string
	strokes   string
}

var bihua bh

func main() {

	// 初始化
	bihua.input()
	input := ""
	p := []int{}
	writer := uilive.New()
	writer.Start()

	// 不知道这个10有啥用，好像是协程
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		// fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
		if event.Key == keyboard.KeyEsc {
			break
		}
		switch event.Rune {
		case '1':
			input += "1"
		case '2':
			input += "2"
		case '3':
			input += "3"
		case '4':
			input += "4"
		case '5':
			input += "5"
		case '6':
		case '7':
		case '8':
		case '9':
		case '0':
		case '-':
		case '=':
		case '\x00':
			if len(input) == 0 {
				continue
			}
			input = input[:len(input)-1]
		}
		p = bihua.match(input)
		fmt.Fprintf(writer, "%s\n", 笔画(input))
		for _, i := range p {
			fmt.Fprintf(writer.Newline(), "%s\n", bihua[i-1].character)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

// 首先要读取数据
func (bihua *bh) input() {
	file, err := ioutil.ReadFile(`./bihua.dat`)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	arr := strings.Fields(string(file))
	for n, m := range arr {
		switch n % 3 {
		case 0:
			bihua[n/3].rank, _ = strconv.Atoi(m)
		case 1:
			bihua[n/3].character = m
		case 2:
			bihua[n/3].strokes = m
		}
	}
}
func (bihua *bh) match(input string) []int {
	n := 0
	m := []int{}
	if input == "" {
		return nil
	}
	for _, j := range bihua {
		if input == j.strokes {
			m = append(m, j.rank)
			continue
		}
		if len(input) < len(j.strokes) && input == j.strokes[:len(input)] {
			m = append(m, j.rank)
			n++
		}
		if n == 6 {
			break
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

// 匹配成功的数和前面那一个交换顺序
func (bihua *bh) best(rank int) {
	bihua[rank].strokes, bihua[rank-1].strokes = bihua[rank-1].strokes, bihua[rank].strokes
	bihua[rank].character, bihua[rank-1].character = bihua[rank-1].character, bihua[rank].character
}
// 没写完
func (bihua *bh) output() {
	for n := range bihua {
		fmt.Println(bihua[n])
	}
}
func 笔画(str string) string {
	m := ""
	for _, n := range str {
		switch n {
		case '1':
			m += "一"
		case '2':
			m += "丨"
		case '3':
			m += "丿"
		case '4':
			m += "丶"
		case '5':
			m += "ㄥ"
		}
	}
	return m
}
