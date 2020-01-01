package main

import (
	//"fmt"
	"github.com/nsf/termbox-go"
)

type typingStageData struct {
	articleId string
	title []rune //string
	desc []rune //string
	link string
	trans []rune
}

type typingGameData struct {
	siteName string
	siteUrl string
	stageList []*typingStageData
}

type typingStageStatus struct {
	data *typingStageData
	titlePos int
	descPos int
	elapsedTime int
	errorCount int
	translationEnabled bool
	typo bool
	progress string
}

type typingGameStatus struct {
	data *typingGameData
	totalStageNumber int
	currentStageNumber int
	currentStage *typingStageStatus
	stageList []*typingStageStatus
	isGameOver bool
}

func packDrawInfo(gameStatus *typingGameStatus) *drawInfo {
	drawInfo := new(drawInfo)
	drawInfo.siteName = gameStatus.data.siteName
	stageStatus := gameStatus.currentStage;
	drawInfo.title = stageStatus.data.title
	drawInfo.desc = stageStatus.data.desc
	if stageStatus.translationEnabled {
		drawInfo.trans = stageStatus.data.trans
	}
	drawInfo.totalStageNumber = gameStatus.totalStageNumber
	drawInfo.currentStageNumber = gameStatus.currentStageNumber
	drawInfo.errorCount = stageStatus.errorCount
	drawInfo.elapsedTime = stageStatus.elapsedTime
	drawInfo.speed = 0
	drawInfo.titlePos = stageStatus.titlePos
	drawInfo.descPos = stageStatus.descPos
	drawInfo.isGameOver = gameStatus.isGameOver
	return drawInfo
}

func loadGame(gameData *typingGameData) *typingGameStatus {
	gameStatus := new(typingGameStatus)
	gameStatus.data = gameData
	for _, stage := range gameData.stageList {
		stageStatus := new(typingStageStatus)
		stageStatus.data = stage
		stageStatus.reset()
		gameStatus.stageList = append(gameStatus.stageList, stageStatus)
	}
	gameStatus.currentStage = gameStatus.stageList[0]
	return gameStatus
}

func controlTypingGame(gameData *typingGameData) *typingGameStatus{
	gameStatus := loadGame(gameData)

	if len(gameStatus.stageList) > 0 {
		gameStatus.currentStage = gameStatus.stageList[0]
	} else {
		gameStatus.gameover()
	}

	drawInfoCh := make(chan drawInfo)
	// timerCh := make(chan int)

	go drawLoop(drawInfoCh)
	// go timerLoop(statusCh, timerCh)

	drawInfoCh <- *(packDrawInfo(gameStatus))

	for gameStatus.isGameOver == false {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				gameStatus.gameover()
			case termbox.KeyTab:
				gameStatus.currentStage.skip()
				gameStatus.nextStage()
			case termbox.KeyEnter:
				if gameStatus.currentStage.isCompleted() {
					gameStatus.nextStage()
				}
			case termbox.KeyCtrlT:
				gameStatus.currentStage.showTranslation()
			case termbox.KeySpace:
				gameStatus.currentStage.update(32)
			default:
				gameStatus.currentStage.update(ev.Ch)
			}
		}
		drawInfoCh <- *(packDrawInfo(gameStatus))
	}
	return gameStatus
}

func (stageStatus *typingStageStatus) reset() {
	stageStatus.progress = "UNREAD"
	stageStatus.titlePos = 0
	stageStatus.descPos = -1
	stageStatus.errorCount = 0
	stageStatus.elapsedTime = 0
	stageStatus.translationEnabled = false
	stageStatus.typo = false
}

func (stageStatus *typingStageStatus) skip(){
	stageStatus.reset()
	stageStatus.progress = "READ"
}

func (stageStatus *typingStageStatus) showTranslation() {
	stageStatus.translationEnabled = true
}

func (stageStatus *typingStageStatus) update(inputChar rune) bool {
	isCompleted := false
	if stageStatus.progress == "UNREAD" {
		currentPos := stageStatus.titlePos
		expectedChar := stageStatus.data.title[currentPos]
		if inputChar == expectedChar {
			stageStatus.typo = false
			stageStatus.titlePos += 1
			if currentPos == (len(stageStatus.data.title) - 1) {
				stageStatus.progress = "TYPED_TITLE"
				stageStatus.descPos = 0
			}
		} else {
			stageStatus.typo = true
			stageStatus.errorCount += 1
		}
	} else if stageStatus.progress == "TYPED_TITLE" {
		currentPos := stageStatus.descPos
		expectedChar := stageStatus.data.desc[currentPos]
		if inputChar == expectedChar {
			stageStatus.typo = false
			stageStatus.descPos += 1
			if currentPos == len(stageStatus.data.desc) {
				stageStatus.progress = "TYPED_ALL"
				isCompleted = true
			} 
		} else {
			stageStatus.typo = true
			stageStatus.errorCount += 1
		}
	} 
	return isCompleted
}

func (stageStatus *typingStageStatus) isCompleted() bool {
	if stageStatus.progress == "TYPED_ALL" {
		return true
	} else {
		return false
	}
}


func (gameStatus *typingGameStatus) gameover() {
	gameStatus.currentStage.reset()
	gameStatus.isGameOver = true
}

func (gameStatus *typingGameStatus) nextStage() {
	if (gameStatus.currentStageNumber + 1) > gameStatus.totalStageNumber {
		gameStatus.isGameOver = true
	} else {
		gameStatus.currentStageNumber += 1
		gameStatus.currentStage = gameStatus.stageList[gameStatus.currentStageNumber]
	}
}










