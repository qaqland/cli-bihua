package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	// 键盘监听
	"github.com/eiannone/keyboard"
	// 终端动态显示
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
	next := 0
	input := ""
	p := []int{}
	writer := uilive.New()
	writer.Start()

	// 不知道这个10有啥用，好像是协程 //
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Fprintf(writer, "Press ESC to quit\n")

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		// fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
		if event.Key == keyboard.KeyEsc {
			break
		}
		if event.Rune == '7' {
			fmt.Fprintf(writer, bihua[4*next+0].character+"\n")
			writer.Stop()
			break
		}
		if event.Rune == '8' {
			fmt.Fprintf(writer, bihua[4*next+1].character+"\n")
			writer.Stop()
			break
		}
		if event.Rune == '9' {
			fmt.Fprintf(writer, bihua[4*next+2].character+"\n")
			writer.Stop()
			break
		}
		if event.Rune == '0' {
			fmt.Fprintf(writer, bihua[4*next+3].character+"\n")
			writer.Stop()
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
		case '-':
			if next != 0 {
				next--
			}
		case '=':
			next++
		case '\x00': // Backspace
			if len(input) == 0 {
				continue
			}
			input = input[:len(input)-1]
		}

		p = bihua.match(input)
		fmt.Fprintf(writer, "%s\n", 笔画(input))
		if p == nil {
			fmt.Fprintf(writer.Newline(), "===\n")
			continue
		}
		for j, i := range p {
			k := j + 7
			if k == 10 {
				k = 0
			}
			fmt.Fprintf(writer.Newline(), "%d %s\n", k, bihua[i-1].character)
			if k == 0 {
				break
			}

		}
		time.Sleep(time.Millisecond * 100)
	}
}

// 首先要读取字的数据，给结构体
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

// 匹配有两种，一种是完全输入，一种是不完全输入只匹配前几位
func (bihua *bh) match(input string) []int {
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

// 12345转化为横竖撇点折
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
