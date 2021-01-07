package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	// 键盘监听
	"github.com/eiannone/keyboard"
	// 终端动态显示
	"github.com/gosuri/uilive"
)

// 字，这个地方要设置大小
type bh [20887]struct {
	rank      int
	character string
	strokes   string
}

// 词+字
type bihuaPlus struct {
	bh
	ci []string
}

var bihua bihuaPlus

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

	fmt.Fprintf(writer, "ESC退出\n")

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
			fmt.Fprintf(writer, bihua.bh[p[4*next+0]-1].character+"\n")
			bihua.best(p[4*next+0] - 1)
			break
		}
		if event.Rune == '8' {
			fmt.Fprintf(writer, bihua.bh[p[4*next+1]-1].character+"\n")
			bihua.best(p[4*next+1] - 1)
			break
		}
		if event.Rune == '9' {
			fmt.Fprintf(writer, bihua.bh[p[4*next+2]-1].character+"\n")
			bihua.best(p[4*next+2] - 1)
			break
		}
		if event.Rune == '0' {
			fmt.Fprintf(writer, bihua.bh[p[4*next+3]-1].character+"\n")
			bihua.best(p[4*next+3] - 1)
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
			// p = 0 1 3 4 5
			if len(p) > next*4-1 && len(p) < (next+1)*4+1 {
				continue
			}
			next++
		case '\x00': // Backspace
			if len(input) == 0 {
				continue
			}
			next = 0
			input = input[:len(input)-1]
		}

		p = bihua.match(input)
		fmt.Fprintf(writer, "%s\n", 笔画(input))
		fmt.Fprintf(writer.Newline(), "next:%d\n", next)
		if p == nil {
			fmt.Fprintf(writer.Newline(), "===\n")
			continue
		}
		for j, i := range p {
			k := j + 7 - next*4
			if k == 10 {
				k = 0
			}
			if j > next*4-1 && j < (next+1)*4 {
				fmt.Fprintf(writer.Newline(), "%d %s\n", k, bihua.bh[i-1].character)
			}
		}
		// time.Sleep(time.Millisecond * 100)
	}
	writer.Stop()
	bihua.output()
}

// 首先要读取字+词的数据，给结构体
func (bihua *bihuaPlus) input() {
	file1, err := ioutil.ReadFile(`./bihua.dat`)
	if err != nil {
		log.Println(err, "字库文件有误")
		os.Exit(2)
	}
	arr := strings.Fields(string(file1))
	for n, m := range arr {
		switch n % 3 {
		case 0:
			bihua.bh[n/3].rank, _ = strconv.Atoi(m)
		case 1:
			bihua.bh[n/3].character = m
		case 2:
			bihua.bh[n/3].strokes = m
		}
	}
	file2, err := ioutil.ReadFile(`./bihuaPlus.dat`)
	if err != nil {
		log.Println(err, "词库文件有误")
		os.Exit(2)
	}
	brr := strings.Fields(string(file2))
	for n, m := range brr {
		bihua.ci[n] = m
	}
}

// 匹配有两种，一种是完全输入，一种是不完全输入只匹配前几位
func (bihua *bihuaPlus) match(input string) []int {
	m := []int{}
	if input == "" {
		return nil
	}
	for _, j := range bihua.bh {
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
// 词的还没做
func (bihua *bihuaPlus) best(rank int) {
	bihua.bh[rank].strokes, bihua.bh[rank-1].strokes = bihua.bh[rank-1].strokes, bihua.bh[rank].strokes
	bihua.bh[rank].character, bihua.bh[rank-1].character = bihua.bh[rank-1].character, bihua.bh[rank].character
}

// 保存
func (bihua *bihuaPlus) output() {
	file1 := ""
	for n := range bihua.bh {
		file1 += strconv.Itoa(bihua.bh[n].rank) + " " + bihua.bh[n].character + " " + bihua.bh[n].strokes + "\n"
	}
	err := ioutil.WriteFile(`./bihua.dat`, []byte(file1), 0666)
	if err != nil {
		log.Println(err, "保存字库有误")
		os.Exit(2)
	}
	file2 := ""
	for _, m := range bihua.ci {
		file2 += m + "\n"
	}
	err = ioutil.WriteFile(`./bihuaPlus.dat`, []byte(file2), 0666)
	if err != nil {
		log.Println(err, "保存词库有误")
		os.Exit(2)
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
