package main

import (
	//	"fmt"
	"github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
)

/*
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
*/

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	keyChan := make(chan keyEvent)
	go keyEventLoop(keyChan)

	gameData := new(typingGameData)
	gameData.siteName = "Test0001"
	gameData.siteUrl = "http://test0001.com"

	stageData := new(typingStageData)
	stageData.articleId = "1"
	stageData.title = []rune("Maxar is selling")
	stageData.desc = []rune("Satellite industry giant Maxar is selling MDA, its subsidiary focused on space robotics")
	stageData.link = "http://test0001.com/123"
	stageData.trans = []rune("衛星産業の巨人Maxarは、宇宙ロボットに特化した子会社であるMDAを10億ドル（約765.23百万ドル）で販売している、とロイターは報告している。購買事業体は、民間投資会社であるノーザンプライベートキャピタルが率いる企業のコンソーシアムであり、MDAのカナダ事業のすべてを買収します。")

	gameData.stageList = append(gameData.stageList, stageData)

	stageData2 := new(typingStageData)
	stageData2.articleId = "2"
	stageData2.title = []rune("Article2")
	stageData2.desc = []rune("Description2")
	stageData2.link = "http://test0002.com/123"
	stageData2.trans = []rune("テスト2")

	gameData.stageList = append(gameData.stageList, stageData2)

	controlTypingGame(gameData, keyChan)

}
