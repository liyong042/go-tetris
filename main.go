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
	backColor   = termbox.ColorBlue
	brickSize   = 4  //方块数量
	brickWidth  = 12 //游戏方块宽
	brickHeight = 21 //游戏方块高
	panelX      = 3  //游戏面板偏移
	panelY      = 1  //游戏面板偏移
	backGround  = `
		WWWWWWWWWWWWWW  WWWWWWWW
		WkkkkkkkkkkkkW  WkkkkkkW
		WkkkkkkkkkkkkW  WkkkkkkW
		WkkkkkkkkkkkkW  WkkkkkkW
		WkkkkkkkkkkkkW  WkkkkkkW
		WkkkkkkkkkkkkW  WkkkkkkW
		WkkkkkkkkkkkkW  WkkkkkkW
		WkkkkkkkkkkkkW  WWWWWWWW
		WkkkkkkkkkkkkW  
		WkkkkkkkkkkkkW
		WkkkkkkkkkkkkW
		WkkkkkkkkkkkkW  BBBBBBBB
		WkkkkkkkkkkkkW  WWWWWWWW  
		WkkkkkkkkkkkkW  
		WkkkkkkkkkkkkW
		WkkkkkkkkkkkkW
		WkkkkkkkkkkkkW  BBBBBBBB
		WkkkkkkkkkkkkW  WWWWWWWW
		WkkkkkkkkkkkkW
		WkkkkkkkkkkkkW
		WkkkkkkkkkkkkW
		WkkkkkkkkkkkkW  BBBBBBBB
		WWWWWWWWWWWWWW  WWWWWWWW
	`
)

//方块
type Brick [brickSize]struct{ x, y int }

//地图颜色
var (
	curPosX       = 0
	curPosY       = 0
	curBrick      Brick //当前方块
	nextBrick     Brick //下一个方块
	curBrickIndex = 0   //当前方块方向
	curBrickType  = 0   //当前方块种类
	curBkColor    = 0   //当前背景颜色
	nextBkColor   = 0   //下一个背景颜色
	nextBrickType = 0   //下一个方块种类
	//地图颜色
	colorMapping = map[rune]termbox.Attribute{
		'k': termbox.ColorBlack,
		'K': termbox.ColorBlack | termbox.AttrBold,
		'b': termbox.ColorBlue,
		'B': termbox.ColorBlue | termbox.AttrBold,
		'w': termbox.ColorWhite,
		'W': termbox.ColorWhite | termbox.AttrBold,
	}
	//方块矩阵, 设计一个2维数组，周边-1，用来做边界判断
	// -1 -1 -1 -1
	// -1  0  0 -1
	// -1  0  0 -1
	// -1 -1 -1 -1
	brickArray = [brickHeight + 2][brickWidth + 2]int{}
	//方块种类
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
	//方块随机颜色
	brickColors = []termbox.Attribute{
		termbox.ColorBlack,
		termbox.ColorRed,
		termbox.ColorGreen,
		termbox.ColorYellow,
		termbox.ColorBlue,
		termbox.ColorMagenta,
		termbox.ColorCyan,
		termbox.ColorWhite,
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
func drawBrick(x, y int, brick *Brick, c int) {
	for i := 0; i < brickSize; i++ {
		drawBlock(x+brick[i].x, y+brick[i].y, brickColors[c])
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
	drawBackGround(backGround, 1, 0) //画游戏地图
	drawBrickArray()
	drawBrick(curPosX+panelX, curPosY+panelY, &curBrick, curBkColor) //画当前方块
	drawBrick(20, 2, &nextBrick, nextBkColor)                        //画下一个方块
	termbox.Flush()
}

//画方块的数组矩阵
func drawBrickArray() {
	for i := 0; i < brickHeight+2; i++ {
		for j := 0; j < brickWidth+2; j++ {
			if t := brickArray[i][j]; t > 0 {
				drawBlock(panelX+j, panelY+i, brickColors[t])
			}
		}
	}
}

//产生一个随机方块
func createRandBrick() {
	curPosX = 4
	curPosY = 0
	curBrickType = nextBrickType
	curBrickIndex = 0

	curBrick = nextBrick
	curBkColor = nextBkColor

	nextBrickType = rand.Intn( len(brickMap) )
	nextBrick = createBrick(brickMap[nextBrickType][0])
	nextBkColor = rand.Intn( len(brickColors) -1  ) + 1
}

//是否可以放得下
func isPut(x, y int, bk *Brick) bool {
	for i := 0; i < brickSize; i++ {
		tY := bk[i].y + y
		tX := bk[i].x + x
		if brickArray[tY][tX] != 0 {
			return false
		}
	}
	return true
}

//检查是否可以删行
func checkFull() {
	for y := brickHeight; y >= 1; y-- {
		removeFull(y)
	}
}

//删除满的行
func removeFull(y int) {
	for x := 1; x <= brickWidth; x++ {
		if brickArray[y][x] == 0 {
			return
		}
	}
	for y = y - 1; y >= 1; y-- {
		for x := 1; x <= brickWidth; x++ {
			brickArray[y+1][x] = brickArray[y][x]
		}

	}
}

//加入方块数组
func addBrickToMap(x, y int, bk *Brick, c int) {
	for i := 0; i < brickSize; i++ {
		tY := bk[i].y + y
		tX := bk[i].x + x
		brickArray[tY][tX] = c
	}
}

//向下移动
func moveDown() {
	if isPut(curPosX, curPosY+1, &curBrick) {
		curPosY += 1
	} else {
		addBrickToMap(curPosX, curPosY, &curBrick, curBkColor)
		createRandBrick()
	}
}

//向左
func moveLeft(x int) {
	if isPut(curPosX+x, curPosY, &curBrick) {
		curPosX += x
	}
}

//向上
func moveUp() {

	t := curBrickIndex +1
	if t >= brickSize {
		t = 0
	}
	bk := createBrick(brickMap[curBrickType][t])
	if isPut(curPosX, curPosY, &bk ) {
		curBrickIndex = t
		curBrick = bk
	}
}

func initGame() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < brickHeight+2; i++ {
		for j := 0; j < brickWidth+2; j++ {
			brickArray[i][j] = 0
			if i == 0 || i == brickHeight+1 || j == 0 || j == brickWidth+1 {
				brickArray[i][j] = -1
			}
		}
	}
}

//
func main() {

	initGame()

	//初始界面
	runtime.LockOSThread()
	termbox.Init()
	defer termbox.Close()

	//初始数据
	createRandBrick()
	createRandBrick()
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
		draw()
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
			checkFull()
		}
	}
}
