package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// 键盘监听
	"github.com/eiannone/keyboard"
	// 终端动态显示
	"github.com/gosuri/uilive"
	// 显示在下一行
)

// 词+字
type bihuaPlus struct {
	zi []string
	ci []string
}

var bihua bihuaPlus

func main() {

	// 初始化
	bihua.input()
	next := 0
	input := ""
	p := []int{}
	q := -1
	writer := uilive.New()
	writer.Start()

	// 不知道这个10有啥用，好像是协程
	// 键盘监听
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Fprintf(writer, "按ESC退出\n")

	// 进入死循环
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		// fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
		if event.Key == keyboard.KeyEsc {
			fmt.Fprintf(writer, "\n")
			goto end
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
			input += "6"
		case '7':
			if q == 0 {
				fmt.Fprintf(writer, bihua.zi[p[4*next+0]]+"\n")
				bihua.best(p[4*next+0], q)
				goto end
			} else if q == 1 {
				fmt.Fprintf(writer, bihua.ci[p[4*next+0]]+"\n")
				bihua.best(p[4*next+0], q)
				goto end
			}

		case '8':
			if q == 0 {
				fmt.Fprintf(writer, bihua.zi[p[4*next+1]]+"\n")
				bihua.best(p[4*next+1], q)
				goto end
			} else if q == 1 {
				fmt.Fprintf(writer, bihua.ci[p[4*next+1]]+"\n")
				bihua.best(p[4*next+1], q)
				goto end
			}
		case '9':
			if q == 0 {
				fmt.Fprintf(writer, bihua.zi[p[4*next+2]]+"\n")
				bihua.best(p[4*next+2], q)
				goto end
			} else if q == 1 {
				fmt.Fprintf(writer, bihua.ci[p[4*next+2]]+"\n")
				bihua.best(p[4*next+2], q)
				goto end
			}
		case '0':
			if q == 0 {
				fmt.Fprintf(writer, bihua.zi[p[4*next+3]]+"\n")
				bihua.best(p[4*next+3], q)
				goto end
			} else if q == 1 {
				fmt.Fprintf(writer, bihua.ci[p[4*next+3]]+"\n")
				bihua.best(p[4*next+3], q)
				goto end
			}
		case '-':
			if next != 0 {
				next--
			}
		case '=':
			// p = 1 3 4 5
			// 如果p不合适，就限制了next无限变大
			if len(p) > next*4-1 && len(p) < next*4+5 {
				continue
			}
			next++
		case '\x00': // Backspace
			next = 0
			if len(input) != 0 {
				input = input[:len(input)-1]
			}
		}

		p, q = bihua.match(input)

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
				if q == 1 {
					fmt.Fprintf(writer.Newline(), "%d %s\n", k, bihua.ci[i])
				}
				if q == 0 {
					fmt.Fprintf(writer.Newline(), "%d %s\n", k, bihua.zi[i])
				}
			}
		}
	}
end:
	bihua.output()
	writer.Stop()
}

// 首先要读取字+词的数据，给切片
func (bihua *bihuaPlus) input() {
	file1, err := ioutil.ReadFile(`./L1.ini`)
	if err != nil {
		log.Println(err, "字库文件有误")
		os.Exit(2)
	}
	bihua.zi = strings.Fields(string(file1))
	file2, err := ioutil.ReadFile(`./L2.ini`)
	if err != nil {
		log.Println(err, "词库文件有误")
		os.Exit(2)
	}
	bihua.ci = strings.Fields(string(file2))
}

// 匹配有两种，一种是完全输入，一种是不完全输入只匹配前几位,
// 两个返回值，m为nil时n为-1，n=0字，n=1词
func (bihua *bihuaPlus) match(input string) ([]int, int) {
	m := []int{}
	part := strings.Split(input, "6")
	if input == "" {
		goto nil
	}
	if len(part) == 1 {
		goto zi
	}
	// j = 你好
	for i, j := range bihua.ci {
		k := strings.Split(j, "")
		if len(part) != len(k) {
			continue
		}
		h := 0
		// part = 123 123 123
		for index := range part {
			if part[index] == ZIBH[k[index]] {
				h++
				continue
			}
			if len(part[index]) < len(ZIBH[k[index]]) && part[index] == ZIBH[k[index]][:len(part[index])] {
				h++
			}
		}
		if h == len(part) {
			m = append(m, i)
		}
	}
	if len(m) == 0 {
		goto nil
	}
	return m, 1

zi:
	for i, j := range bihua.zi {
		if input == ZIBH[j] {
			m = append(m, i)
			continue
		}
		if len(input) < len(ZIBH[j]) && input == ZIBH[j][:len(input)] {
			m = append(m, i)
		}
	}
	if len(m) == 0 {
		goto nil
	}
	return m, 0
nil:
	return nil, -1
}

// 匹配成功的数和前面那一个交换顺序
func (bihua *bihuaPlus) best(rank int, q int) {
	if rank < 1 {
		return
	}
	if q == 0 {
		bihua.zi[rank], bihua.zi[rank-1] = bihua.zi[rank-1], bihua.zi[rank]
	} else {
		bihua.ci[rank], bihua.ci[rank-1] = bihua.ci[rank-1], bihua.ci[rank]
	}
}

// 保存
func (bihua *bihuaPlus) output() {
	file1 := strings.Join(bihua.zi, "\n")
	err := ioutil.WriteFile(`./L1.ini`, []byte(file1), 0666)
	if err != nil {
		log.Println(err, "保存字库有误")
		os.Exit(2)
	}
	file2 := strings.Join(bihua.ci, "\n")
	err = ioutil.WriteFile(`./L2.ini`, []byte(file2), 0666)
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
		case '6':
			m += "·"
		}
	}
	return m
}

// 本来想放在另外一个包的，迫于没学懂，先放在这里
