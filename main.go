package main

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)
//随机产生方块移动, 增加键盘事件
//常量声明
//游戏地图
const (
	backColor  = termbox.ColorBlue
	brickSize  = 4 //方块数量
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

//方块
type Brick [brickSize]struct{ x, y int }

//地图颜色
var (
	curPosX       = 0
	curPosY       = 0
	curBrick      Brick //当前方块
	curBrickIndex = 0   //当前方块方向
	curBrickType  = 0   //当前方块种类
	colorMapping  = map[rune]termbox.Attribute{
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
	brickMap = [][]int{ //各种方块定义 十进制每一位表示一行
		// 0000 0
		// 0000 0
		// 0110 6
		// 0110 6
		[]int{66, 66, 66, 66},
		// 0000 0
		// 0000 0
		// 0010 2
		// 0111 7
		[]int{27, 131, 72, 232},
		// 0000 0
		// 0000 0
		// 0011 3
		// 0110 6
		[]int{36, 231, 36, 231},
		// 0000 0
		// 0000 0
		// 0110 6
		// 0011 3
		[]int{63, 132, 63, 132},
		// 0000 0
		// 0011 3
		// 0001 1
		// 0001 1
		[]int{311, 17, 223, 74},
		// 0000 0
		// 0011 3
		// 0010 1
		// 0010 1
		[]int{322, 71, 113, 47},
		// Special case since 15 can't be used
		// 0001 1
		// 0001 1
		// 0001 1
		// 0001 1
		[]int{1111, 9, 1111, 9}, //可以考虑使用16进制
	}
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

//画一种方块
func drawBrick(x, y int) {
	for i := 0; i < brickSize; i++ {
		drawBlock(x+curBrick[i].x, y+curBrick[i].y, termbox.ColorRed)
	}
}

//根据方块图创建一个方块
func createBrick(t int) (bk Brick) {
	cnt := 0
	horizontal := t == 9 // 这里针对 长条做了特殊处理
	for i := 0; i <= 3; i++ {
		p := int( math.Pow( 10, float64(3-i)) ) //取位整数
		digit := t / p
		t %= p
		for j := 3; j >= 0; j-- { //行转换
			bin := digit % 2
			digit /= 2
			if bin == 1 || (horizontal && i == brickSize-1) { //这里针对 长条做了特殊处理
				bk[cnt].x = j
				bk[cnt].y = i
				cnt++
			}
		}
	}
	return bk
}

//画图
func draw() {
	termbox.Clear(backColor, backColor)
	drawBackGround(backGround, 1, 0) //画游戏地图
	drawBrick(curPosX, curPosY)      //画方块
	termbox.Flush()
}
//产生一个随机方块
func createRandBrick(){
	curPosX = 6
	curPosY = 0
	curBrickType  =  rand.Intn( len(brickMap) )
	curBrickIndex = 0
	curBrick = createBrick(brickMap[curBrickType][curBrickIndex])
}

//
func main() {
	//初始界面
	termbox.Init()
	defer termbox.Close()
	rand.Seed(time.Now().Unix())

	//初始数据
	createRandBrick()
	draw()

	//定时
	ticker := time.NewTicker(time.Millisecond * 1000)
	for range ticker.C {
		curPosY = curPosY + 1
		if curPosY > 16 {
			createRandBrick()
		}
		draw()
	}
}
