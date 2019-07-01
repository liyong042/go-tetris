package main

import (
	"runtime"
	"time"
	"github.com/nsf/termbox-go"
)

func main() {
	//初始界面
	runtime.LockOSThread()
	termbox.Init()
	defer termbox.Close()

	initGame()

	//定时
	ticker := time.NewTicker(time.Millisecond * 1000)
	eventChan := make(chan termbox.Event)
	go func() {
		for {
			eventChan <- termbox.PollEvent()
		}
	}()
	//增加键盘事件
	for {
		drawWindow()
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
			case termbox.KeySpace:
				moveEnd()
			}
		case <-ticker.C:
			moveDown()
		}
	}
}
