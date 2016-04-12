package handlers

import (
	"encoding/hex"
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gotelebot"
	"github.com/eternnoir/gotelebot/types"
	"github.com/eternnoir/mew_gif_bot/mew/utils"
	"strings"
)

var UserState map[float64]func(bot *gotelebot.TeleBot, message *types.Message, fs utils.FileStore) error
var UserGif map[float64]*types.Document

func ProcessMessage(bot *gotelebot.TeleBot, message *types.Message, fs utils.FileStore) error {
	log.Infof("Get Text message: %#v", message.Text)
	if UserState == nil {
		UserState = make(map[float64]func(bot *gotelebot.TeleBot, message *types.Message, fs utils.FileStore) error)
		UserGif = make(map[float64]*types.Document)
	}

	if strings.HasPrefix(message.Text, "/new") {
		delete(UserState, message.Chat.Id)
		delete(UserGif, message.Chat.Id)
		return CreateNewGif(bot, message, fs)
	}
	if val, ok := UserState[message.Chat.Id]; ok {
		log.Debugf("%d id state found.", int(message.Chat.Id))
		delete(UserState, message.Chat.Id)
		return val(bot, message, fs)
	}
	return nil
}

func CreateNewGif(bot *gotelebot.TeleBot, m *types.Message, fs utils.FileStore) error {
	log.Debug("Add command triggered")
	msg := "Send me gif file you want to add."
	retmsg, err := bot.SendMessage(int(m.Chat.Id), msg, nil)
	log.Debugf("Get return msg: %#v", retmsg)
	UserState[m.Chat.Id] = CreateNewGifStep2
	return err
}

func CreateNewGifStep2(bot *gotelebot.TeleBot, m *types.Message, fs utils.FileStore) error {
	log.Debugf("%d enter New Gif step2", int(m.Chat.Id))
	if m.Document == nil {
		msg := "Please Send Me gif file."
		_, err := bot.SendMessage(int(m.Chat.Id), msg, nil)
		return err
	}
	log.Debugf("Get Document from %d :%#v", int(m.Chat.Id), m.Document)
	msg := "Ok give me a EMOJI about this gif."
	_, err := bot.SendMessage(int(m.Chat.Id), msg, nil)
	if err != nil {
		return err
	}
	UserState[m.Chat.Id] = CreateNewGifStep3
	UserGif[m.Chat.Id] = m.Document
	return nil
}

func CreateNewGifStep3(bot *gotelebot.TeleBot, m *types.Message, fs utils.FileStore) error {
	log.Debugf("%d enter New Gif step3", int(m.Chat.Id))
	if m.Text == "" {
		msg := "Please Send Me EMOJI."
		_, err := bot.SendMessage(int(m.Chat.Id), msg, nil)
		return err
	}
	log.Debugf("Get Document from %d :%#v", int(m.Chat.Id), m.Document)
	if _, ok := UserGif[m.Chat.Id]; !ok {
		msg := "Something wrong, connot get gif file id."
		_, err := bot.SendMessage(int(m.Chat.Id), msg, nil)
		return err
	}

	doc := UserGif[m.Chat.Id]

	Enc := hex.EncodeToString([]byte(m.Text))

	log.Infof("Put fileid to cache.")
	err := fs.Put(Enc, doc.FileId)
	if err != nil {
		msg := "Something wrong, connot get gif file id."
		bot.SendMessage(int(m.Chat.Id), msg, nil)
		return err
	}
	msg := "Good."
	bot.SendMessage(int(m.Chat.Id), msg, nil)
	return nil
}
