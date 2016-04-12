package handlers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gotelebot"
	"github.com/eternnoir/gotelebot/types"
	"github.com/eternnoir/mew_gif_bot/mew/utils"
)

func ProcessChosen(telebot *gotelebot.TeleBot, result *types.ChosenInlineResult, fs utils.FileStore) error {
	log.Infof("Get New ChosenInlineResult %#v", result)
	log.Infof("Choise result %s", result.ResultId)
	fs.Hint(result.ResultId)
	return nil
}
