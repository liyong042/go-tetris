package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func draw(ch rune) {
	w, h := termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(w/2, h/2, ch, termbox.ColorRed, termbox.ColorDefault)

	termbox.Flush()
}

func main() {
	termbox.Init()
	defer termbox.Close()

	for {
		//ev := termbox.PollEvent()
		//draw(ev.Ch)
	}
	fmt.Println("Stop")

}
