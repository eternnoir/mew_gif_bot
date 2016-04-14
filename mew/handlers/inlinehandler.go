package handlers

import (
	"encoding/hex"
	_ "fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gotelebot"
	"github.com/eternnoir/gotelebot/types"
	"github.com/eternnoir/mew_gif_bot/mew/utils"
	"strings"
)

func ProcessInlineQuery(bot *gotelebot.TeleBot, q *types.InlineQuery, fs utils.FileStore) error {
	log.Infof("Get New InlineQuery %#v", q)
	if strings.TrimSpace(q.Query) == "" {
		return SendAllGifs(bot, q, fs)
	}
	return SendByQuery(bot, q, fs)
}

func SendByQuery(bot *gotelebot.TeleBot, q *types.InlineQuery, fs utils.FileStore) error {

	Enc := hex.EncodeToString([]byte(q.Query))
	fileids, err := fs.Get(Enc)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("Get %d files, query %s", len(fileids), Enc)
	ir := genInlineGifQuery(fileids)

	opt := gotelebot.AnswerInlineQueryOptional{}
	cache := 1
	opt.CacheTime = &cache
	result, err := bot.AnswerInlineQuery(q.Id, ir, &opt)
	if err != nil {
		return err
	}
	log.Infof("Ansered %#v", result)
	return nil
}

func SendAllGifs(bot *gotelebot.TeleBot, q *types.InlineQuery, fs utils.FileStore) error {
	fileids, err := fs.GetAll()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("Get %d files, query %s", len(fileids), q.Query)
	ir := genInlineGifQuery(fileids)

	opt := gotelebot.AnswerInlineQueryOptional{}
	cache := 1
	opt.CacheTime = &cache
	result, err := bot.AnswerInlineQuery(q.Id, ir, &opt)
	if err != nil {
		return err
	}
	log.Infof("Ansered %#v", result)
	return nil
}

func genInlineGifQuery(fileIds []string) []interface{} {
	ret := []interface{}{}
	if len(fileIds) > 50 {
		fileIds = fileIds[:49]
	}

	for _, fid := range fileIds {
		s := strings.Split(fid, "::")
		tfid := s[1]
		iq := types.NewInlineQueryResultCachedMpeg4Gif()
		iq.Id = fid
		iq.Mpeg4FileId = tfid
		ret = append(ret, iq)
	}
	return ret
}
