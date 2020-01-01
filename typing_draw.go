package main

import (
	"github.com/nsf/termbox-go"
)

const baseColor = termbox.ColorDefault
const baseBackColor = termbox.ColorDefault
const cursorForeColor = termbox.ColorBlack
const cursorBackColor = termbox.ColorYellow
const untypedColor = termbox.ColorDefault
const typedColor = termbox.ColorRed
const errorCountColor = termbox.ColorRed

const lineLength = 78

/*
const progressOffset = 37
const speedOffset = 20
const errorOffset = 34
const accuracyOffset = 48
const timeOffset = 65
*/

type drawInfo struct {
	siteName string
	title []rune
	desc []rune
	trans []rune
	totalStageNumber int
	currentStageNumber int
	errorCount int
	elapsedTime int
	speed int
	titlePos int
	descPos int
	isGameOver bool
}

func drawHeader(info *drawInfo, lineLength int, lineOffset int, rowOffset int) int {
	foreColor := baseColor
	line := lineOffset
	row := rowOffset

	// draw site info (max 32 characters)
	runes := []rune(info.siteName)
	for i, c := range runes {
		row = i + rowOffset
		termbox.SetCell(row, line, c, foreColor, termbox.ColorDefault)
	}

	// draw progress (current / total)
	// draw error count (or accuracy)
	// draw speed
	// draw elasped time

	return line
}


func drawTitle(info *drawInfo, lineLength int, lineOffset int, rowOffset int) int {
	foreColor := typedColor
	backColor := baseBackColor
	line := lineOffset
	row  := rowOffset
	for i, ch := range info.title {
		if i == info.titlePos {
			foreColor = cursorForeColor
			backColor = cursorBackColor
		} else if i > info.titlePos {
			foreColor = untypedColor
			backColor = baseBackColor
		}
		line = (i / lineLength) + lineOffset
		row =  (i % lineLength) + rowOffset
		termbox.SetCell(row, line, ch, foreColor, backColor)
	}
	return line
}

func drawDesc(info *drawInfo, lineLength int, lineOffset int, rowOffset int) int {
	foreColor := typedColor
	backColor := baseBackColor
	line := lineOffset
	row  := rowOffset
	for i, ch := range info.desc {
		if i == info.descPos {
			foreColor = cursorForeColor
			backColor = cursorBackColor
		} else if i > info.descPos {
			foreColor = untypedColor
			backColor = baseBackColor
		}
		line = (i / lineLength) + lineOffset
		row =  (i % lineLength) + rowOffset
		termbox.SetCell(row, line, ch, foreColor, backColor)
	}
	return line
}

func drawTranslation(info *drawInfo, lineLength int, lineOffset int, rowOffset int) int {
	foreColor := baseColor
	line := lineOffset
	row  := rowOffset
	
	for i, ch := range info.trans {
		line = ((i*2) / lineLength) + lineOffset
		row =  ((i*2) % lineLength) + rowOffset
		termbox.SetCell(row, line, ch, foreColor, termbox.ColorDefault)
	}
	return line
}

func drawLoop(drawInfoChan chan drawInfo) {
	for {
		currentLine := 0
		info := <-drawInfoChan
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		if info.isGameOver {
			break
		}
		currentLine = drawHeader(&info, lineLength, 0, 1)
		currentLine = drawTitle(&info, lineLength, currentLine + 2, 1)
		currentLine = drawDesc(&info, lineLength, currentLine + 2, 1)
		currentLine = drawTranslation(&info, lineLength, currentLine + 2, 1)
		termbox.Flush()
	}
}







