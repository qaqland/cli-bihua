package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
)

// 这个地方要设置大小
type bh [20]struct {
	rank      int
	character string
	strokes   string
}

var bihua bh

func main() {
	fmt.Println("hello")
	bihua.input()

	encoding.Register()
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	// time.Sleep(2 * time.Second)
	input := ""
	p := []int{}
	for {
		if input == "" {
			displayy(s, 1, "开始打字")
		}

		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			} else if ev.Key() == tcell.KeyBackspace {
				if len(input) == 0 {
					continue
				}
				input = input[:len(input)-1]
				p = bihua.match(input)
				s.Clear()
				displayy(s, 1, input)
				if p == nil {
					displayy(s, 2, "没有匹配")
				}
				for x, y := range p {
					displayy(s, x+2, bihua[y-1].character)
				}
			} else {
				input += string(ev.Rune())
				p = bihua.match(input)
				s.Clear()
				displayy(s, 1, input)
				if p == nil {
					displayy(s, 2, "没有匹配")
				}
				for x, y := range p {
					displayy(s, x+2, bihua[y-1].character)
				}

			}
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

// 首先要读取数据
func (bihua *bh) input() {
	file, err := ioutil.ReadFile(`./bihua.dat`)
	if err != nil {
		log.Println(err)
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
func (bihua *bh) output() {
	for n := range bihua {
		fmt.Println(bihua[n])
	}
}

///////////////////////////////////////////////////////////////////////////////
func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 { //这里应该是为了避免零宽字符
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}
func displayy(s tcell.Screen, h int, str string) {
	emitStr(s, 2, h, tcell.StyleDefault, str)
	s.Show()
}
