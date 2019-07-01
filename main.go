package main

import (
	"strings"
	"time"
	"github.com/nsf/termbox-go"
)
//简单画一个 田子 方块
//常量声明
//游戏地图
const (
	backColor  = termbox.ColorBlue
	backGround = `
		WWWWWWWWWWWW  WWWWWW
		WkkkkkkkkkkW  WkkkkW
		WkkkkkkkkkkW  WkkkkW
		WkkkkkkkkkkW  WkkkkW
		WkkkkkkkkkkW  WkkkkW
		WkkkkkkkkkkW  WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WkkkkkkkkkkW  BBBBBB
		WkkkkkkkkkkW  WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WkkkkkkkkkkW  BBBBBB
		WkkkkkkkkkkW  WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW  BBBBBB
		WkkkkkkkkkkW  WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WWWWWWWWWWWW

		kkkkkkkkkkkkkkkkkkkk
		WWWWWWWWWWWWWWWWWWWW
	`
)
//地图颜色
var (
	colorMapping = map[rune]termbox.Attribute{
		'k': termbox.ColorBlack,
		'K': termbox.ColorBlack | termbox.AttrBold,
		'r': termbox.ColorRed,
		'R': termbox.ColorRed | termbox.AttrBold,
		'g': termbox.ColorGreen,
		'G': termbox.ColorGreen | termbox.AttrBold,
		'y': termbox.ColorYellow,
		'Y': termbox.ColorYellow | termbox.AttrBold,
		'b': termbox.ColorBlue,
		'B': termbox.ColorBlue | termbox.AttrBold,
		'm': termbox.ColorMagenta,
		'M': termbox.ColorMagenta | termbox.AttrBold,
		'c': termbox.ColorCyan,
		'C': termbox.ColorCyan | termbox.AttrBold,
		'w': termbox.ColorWhite,
		'W': termbox.ColorWhite | termbox.AttrBold,
	}
	curY = 0
)

//地图字符对应颜色
func getColorByCh(ch rune) termbox.Attribute {
	if c, ok := colorMapping[ch]; ok {
		return c
	}
	return backColor
}

//画游戏地图
func drawBackGround(text string, left, top int) {
	lines := strings.Split(text, "\n")

	for y, line := range lines {
		for x, ch := range line {
			drawBack(left+x, top+y, getColorByCh(ch))
		}
	}
}
//画背景
func drawBack(x, y int, color termbox.Attribute) {
	termbox.SetCell(2*x-1, y, ' ', backColor, color)
	termbox.SetCell(2*x, y, ' ', backColor, color)
}

//画方块格子
func drawBlock(x, y int, color termbox.Attribute) {
	termbox.SetCell(2*x-1, y, ' ', backColor, color)
	termbox.SetCell(2*x, y, ' ', backColor, color)
}

//画图
func draw() {
	termbox.Clear(backColor, backColor)
	drawBackGround(backGround, 1, 0) //画游戏地图
	drawBlockT( 5, curY )  //画方块
	termbox.Flush()
}
//画方块
func drawBlockT( x, y int ){
	drawBlock( x, y, termbox.ColorRed )
	drawBlock( x+1, y, termbox.ColorRed )
	drawBlock( x, y+1, termbox.ColorRed )
	drawBlock( x+1, y+1, termbox.ColorRed )
}

func main() {
	termbox.Init()
	defer termbox.Close()
	ticker := time.NewTicker(time.Millisecond * 1000)

	curY =2
	draw()
	for range ticker.C {
		curY= curY+ 1
		if curY>18{
			curY=2
		}
		draw()
	}
}
