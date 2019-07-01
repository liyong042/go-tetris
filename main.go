package main

import (
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

//随机产生方块移动
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
	// 0000 0
	// 0000 0
	// 0110 6
	// 0110 6
	brickMap = [][]int{ //各种方块定义 十进制每一位表示一行,数据结构中矩阵
		{66, 66, 66, 66},   //田型方块
		{27, 131, 72, 232}, //T型方块
		{36, 231, 36, 231}, //Z型方块
		{63, 132, 63, 132}, //倒Z型方块
		{311, 17, 223, 74}, //倒L型方块
		{322, 71, 113, 47}, //L型方块
		{1111, 9, 1111, 9}, //-型方块 可以考虑使用16进制
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
func drawBrick(x, y int, brick *Brick) {
	for i := 0; i < brickSize; i++ {
		drawBlock(x+brick[i].x, y+brick[i].y, termbox.ColorRed)
	}
}

//根据方块图创建一个方块
func createBrick(t int) (bk Brick) {
	cnt := 0
	horizontal := t == 9 // 这里针对 长条做了特殊处理
	for i := 0; i <= 3; i++ {
		p := int(math.Pow(10, float64(3-i))) //取位整数
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
	drawBackGround(backGround, 1, 0)       //画游戏地图
	drawBrick(curPosX, curPosY, &curBrick) //画方块
	termbox.Flush()
}

//产生一个随机方块
func createRandBrick() {
	curPosX = 6
	curPosY = 0
	curBrickType = rand.Intn(len(brickMap))
	curBrickIndex = 0
	curBrick = createBrick(brickMap[curBrickType][curBrickIndex])
}

//向下移动
func moveDown() {
	curPosY = curPosY + 1
	if curPosY > 16 {
		createRandBrick()
	}
	draw()
}

//向左
func moveLeft(x int) {
	curPosX += x
	if curPosX > 10 {
		curPosX = 10
	}
	if curPosX <= 2 {
		curPosX = 2
	}
	draw()
}

//向上
func moveUp() {
	curBrickIndex++
	if curBrickIndex >= brickSize {
		curBrickIndex = 0
	}
	curBrick = createBrick(brickMap[curBrickType][curBrickIndex])
	draw()
}

//
func main() {
	//初始界面
	runtime.LockOSThread()
	termbox.Init()
	defer termbox.Close()
	rand.Seed(time.Now().Unix())
	//初始数据
	createRandBrick()
	draw()
	//定时
	ticker := time.NewTicker(time.Millisecond * 1000)
	eventChan := make(chan termbox.Event)
	go func() {
		for {
			eventChan <- termbox.PollEvent()
		}
	}()
	//增加键盘实践
	for {
		select {
		case ev := <-eventChan:
			if ev.Type != termbox.EventKey {
				continue
			}
			switch ev.Key {
			case termbox.KeyArrowDown:
				moveDown()
			case termbox.KeyArrowLeft:
				moveLeft(-1)
			case termbox.KeyArrowRight:
				moveLeft(1)
			case termbox.KeyArrowUp:
				moveUp()
			}
		case <-ticker.C:
			moveDown()
		}
	}
}
