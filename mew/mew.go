package mew

import (
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gotelebot"
	"github.com/eternnoir/gotelebot/types"
	"github.com/eternnoir/mew_gif_bot/mew/handlers"
	"github.com/eternnoir/mew_gif_bot/mew/utils"
	"github.com/eternnoir/mew_gif_bot/mew/utils/localfile"
)

type Config struct {
	TeleToken     string
	AwsAccessKey  string
	AwsSecretKey  string
	BucketName    string
	LoaclFilePath string
	ServerUrl     string
}

type MewGif struct {
	Token     string
	telebot   *gotelebot.TeleBot
	fileSer   utils.FileService
	DebugMode bool
}

func NewMewGif(config Config, debug bool) *MewGif {
	ret := MewGif{}
	log.Infof("Create MewGif with config file: %#v", config)
	ret.DebugMode = debug
	ret.Token = config.TeleToken
	ret.telebot = gotelebot.InitTeleBot(ret.Token)
	ret.fileSer = &localfile.LocalFS{Path: config.LoaclFilePath, ServerUrl: config.ServerUrl}
	return &ret
}

func (mew *MewGif) StartPolling() {
	bot := mew.telebot
	go mew.telebot.StartPolling(true, 60)
	go func() {
		newMsgChan := bot.Messages
		for {
			m := <-newMsgChan // Get new messaage, when new message arrive.
			log.Debugf("Get Message:%#v \n", m)
			go mew.HandleNewMessage(m)
		}
	}()
	func() {
		newQuery := bot.InlineQuerys
		for {
			q := <-newQuery
			log.Debugf("Get NewInlineQuery:%#v \n", q)
			go mew.HandleNewInlineQuery(q)
		}
	}()
}

func (mew *MewGif) HandleNewMessage(message *types.Message) {
	log.Infof("Process new message: %#v", message)
	err := handlers.ProcessMessage(mew.telebot, message, mew.fileSer)
	if err != nil {
		log.Errorf("ProcessTextMessage error: %s", err)
	}
}

func (mew *MewGif) HandleNewInlineQuery(inlinequery *types.InlineQuery) {
	log.Infof("Process new inlinequery: %#v", inlinequery)
}
