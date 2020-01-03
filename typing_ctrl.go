package main

import (
	_ "fmt"

	"github.com/nsf/termbox-go"
)

type typingStageData struct {
	articleId string
	title     []rune //string
	desc      []rune //string
	link      string
	trans     []rune
}

type typingGameData struct {
	siteName  string
	siteUrl   string
	stageList []*typingStageData
}

type typingStageStatus struct {
	data               *typingStageData
	titlePos           int
	descPos            int
	elapsedTime        int
	errorCount         int
	translationEnabled bool
	typo               bool
	progress           string
	okCount            int
}

type typingGameStatus struct {
	data               *typingGameData
	totalStageNumber   int
	currentStageNumber int
	currentStage       *typingStageStatus
	stageList          []*typingStageStatus
	isGameOver         bool
}

func packDrawInfo(gameStatus *typingGameStatus) *drawInfo {
	drawInfo := new(drawInfo)
	drawInfo.siteName = gameStatus.data.siteName
	stageStatus := gameStatus.currentStage
	drawInfo.title = stageStatus.data.title
	drawInfo.desc = stageStatus.data.desc
	if stageStatus.translationEnabled {
		drawInfo.trans = stageStatus.data.trans
	}
	drawInfo.totalStageNumber = gameStatus.totalStageNumber
	drawInfo.currentStageNumber = gameStatus.currentStageNumber + 1
	drawInfo.errorCount = stageStatus.errorCount
	drawInfo.elapsedTime = stageStatus.elapsedTime
	drawInfo.speed = 0
	drawInfo.titlePos = stageStatus.titlePos
	drawInfo.descPos = stageStatus.descPos
	drawInfo.isGameOver = gameStatus.isGameOver
	drawInfo.typo = stageStatus.typo
	drawInfo.okCount = stageStatus.okCount
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
	gameStatus.currentStageNumber = 0
	gameStatus.totalStageNumber = len(gameStatus.stageList)
	return gameStatus
}

func controlTypingGame(gameData *typingGameData, keyChan chan keyEvent) *typingGameStatus {
	gameStatus := loadGame(gameData)

	if len(gameStatus.stageList) > 0 {
		gameStatus.currentStage = gameStatus.stageList[0]
	} else {
		gameStatus.gameover()
	}

	drawInfoCh := make(chan drawInfo)
	defer close(drawInfoCh)
	timerCh := make(chan int)
	defer close(timerCh)
	timerOpCh := make(chan int)
	defer close(timerOpCh)

	go drawLoop(drawInfoCh)
	go timerLoop(timerCh, timerOpCh)

	drawInfoCh <- *(packDrawInfo(gameStatus))
	timerOpCh <- 1

	for gameStatus.isGameOver == false {
		select {
		case ev := <-keyChan:
			//fmt.Println(ev.key)
			switch ev.key {
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
				fallthrough
			default:
				if gameStatus.currentStage.beforePlaying() { // start playing stage
					timerOpCh <- 0 // start timer
				}
				if gameStatus.currentStage.update(ev.ch) { // completed stage
					timerOpCh <- 1 // stop timer
				}
			}
		case elapsedTime := <-timerCh:
			gameStatus.currentStage.elapsedTime = elapsedTime
		}
		drawInfoCh <- *(packDrawInfo(gameStatus))
	}
	return gameStatus
}

func (stageStatus *typingStageStatus) beforePlaying() bool {
	if stageStatus.progress == "UNREAD" {
		return true
	} else {
		return false
	}
}

func (stageStatus *typingStageStatus) reset() {
	stageStatus.progress = "UNREAD"
	stageStatus.titlePos = 0
	stageStatus.descPos = -1
	stageStatus.errorCount = 0
	stageStatus.okCount = 0
	stageStatus.elapsedTime = 0
	stageStatus.translationEnabled = false
	stageStatus.typo = false
}

func (stageStatus *typingStageStatus) skip() {
	stageStatus.reset()
	stageStatus.progress = "READ"
}

func (stageStatus *typingStageStatus) showTranslation() {
	stageStatus.translationEnabled = true
}

func (stageStatus *typingStageStatus) update(inputChar rune) bool {
	isCompleted := false

	switch stageStatus.progress {
	case "UNREAD":
		stageStatus.progress = "TYPING_TITLE"
		fallthrough
	case "TYPING_TITLE":
		currentPos := stageStatus.titlePos
		expectedChar := stageStatus.data.title[currentPos]
		if inputChar == expectedChar {
			stageStatus.okCount += 1
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
	case "TYPED_TITLE":
		stageStatus.progress = "TYPING_DESC"
		fallthrough
	case "TYPING_DESC":
		currentPos := stageStatus.descPos
		expectedChar := stageStatus.data.desc[currentPos]
		if inputChar == expectedChar {
			stageStatus.okCount += 1
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
	if (gameStatus.currentStageNumber + 1) >= gameStatus.totalStageNumber {
		gameStatus.isGameOver = true
	} else {
		gameStatus.currentStageNumber += 1
		gameStatus.currentStage = gameStatus.stageList[gameStatus.currentStageNumber]
	}
}
