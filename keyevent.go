package main

import (
	"github.com/nsf/termbox-go"
)

type keyEvent struct {
	key termbox.Key
	ch rune
}

func keyEventLoop(keyCh chan keyEvent) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			
			keyEv := keyEvent{}
			keyEv.key = ev.Key
			keyEv.ch = ev.Ch
			if ev.Key == termbox.KeySpace {
				keyEv.ch = 32
			}

			keyCh <- keyEv
		}
	}
}