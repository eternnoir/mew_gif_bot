package handlers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gotelebot"
	"github.com/eternnoir/gotelebot/types"
	"github.com/eternnoir/mew_gif_bot/mew/utils"
	"strings"
)

var UserState map[float64]func(bot *gotelebot.TeleBot, message *types.Message, fs utils.FileService) error
var UserGif map[float64]string

func ProcessMessage(bot *gotelebot.TeleBot, message *types.Message, fs utils.FileService) error {
	log.Infof("Get Text message: %#v", message.Text)
	if UserState == nil {
		UserState = make(map[float64]func(bot *gotelebot.TeleBot, message *types.Message, fs utils.FileService) error)
		UserGif = make(map[float64]string)
	}
	if strings.HasPrefix(message.Text, "/new") {
		delete(UserState, message.Chat.Id)
		return CreateNewGif(bot, message, fs)
	}
	if val, ok := UserState[message.Chat.Id]; ok {
		log.Debugf("%d id state found.", int(message.Chat.Id))
		delete(UserState, message.Chat.Id)
		return val(bot, message, fs)
	}
	return nil
}

func CreateNewGif(bot *gotelebot.TeleBot, m *types.Message, fs utils.FileService) error {
	log.Debug("Add command triggered")
	msg := "Send me gif file you want to add."
	retmsg, err := bot.SendMessage(int(m.Chat.Id), msg, nil)
	log.Debugf("Get return msg: %#v", retmsg)
	UserState[m.Chat.Id] = CreateNewGifStep2
	return err
}
func CreateNewGifStep2(bot *gotelebot.TeleBot, m *types.Message, fs utils.FileService) error {
	log.Debugf("%d enter New Gif step2", int(m.Chat.Id))
	if m.Document == nil {
		msg := "Please Send Me gif file."
		_, err := bot.SendMessage(int(m.Chat.Id), msg, nil)
		return err
	}
	log.Debugf("Get Document from %d :%#v", int(m.Chat.Id), m.Document)
	return nil
}
