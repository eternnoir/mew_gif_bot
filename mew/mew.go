package mew

import (
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gotelebot"
	"github.com/eternnoir/gotelebot/types"
	"github.com/eternnoir/mew_gif_bot/mew/handlers"
	"github.com/eternnoir/mew_gif_bot/mew/utils"
	"github.com/eternnoir/mew_gif_bot/mew/utils/redis"
)

type Config struct {
	TeleToken  string
	RedisAddr  string
	RedisPwd   string
	WelcomeMsg string
	ListenUrl  string
	KeyPath    string
	CerPath    string
}

type MewGif struct {
	Token     string
	telebot   *gotelebot.TeleBot
	storeage  utils.FileStore
	DebugMode bool
}

func NewMewGif(config Config, debug bool) *MewGif {
	ret := MewGif{}
	log.Infof("Create MewGif with config file: %#v", config)
	ret.DebugMode = debug
	ret.Token = config.TeleToken
	ret.telebot = gotelebot.InitTeleBot(ret.Token)
	handlers.WelcomeMsg = config.WelcomeMsg
	fs, err := redis.NewRedisStore(config.RedisAddr, config.RedisPwd)
	if err != nil {
		log.Error(err)
		return nil
	}
	ret.storeage = fs
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
	go func() {
		newQuery := bot.InlineQuerys
		for {
			q := <-newQuery
			log.Debugf("Get NewInlineQuery:%#v \n", q)
			go mew.HandleNewInlineQuery(q)
		}
	}()

	func() {
		newresult := bot.ChosenInlineResults
		for {
			r := <-newresult
			log.Debugf("Get ChosenInlineResult:%#v \n", r)
			go mew.ProcessChosen(r)
		}
	}()

}

func (mew *MewGif) HandleNewMessage(message *types.Message) {
	log.Infof("Process new message: %#v", message)
	err := handlers.ProcessMessage(mew.telebot, message, mew.storeage)
	if err != nil {
		log.Errorf("ProcessTextMessage error: %s", err)
	}
}

func (mew *MewGif) HandleNewInlineQuery(inlinequery *types.InlineQuery) {
	log.Infof("Process new inlinequery: %#v", inlinequery)
	err := handlers.ProcessInlineQuery(mew.telebot, inlinequery, mew.storeage)
	if err != nil {
		log.Errorf("ProcessInlineQuery error: %s", err)
	}
}

func (mew *MewGif) ProcessChosen(chosen *types.ChosenInlineResult) {
	log.Infof("Process new ChosenInlineResult: %#v", chosen)
	err := handlers.ProcessChosen(mew.telebot, chosen, mew.storeage)
	if err != nil {
		log.Errorf("ProcessChosen error: %s", err)
	}
}
