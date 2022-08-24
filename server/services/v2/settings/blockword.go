package settings

import (
	"net/http"
	"strings"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// BlockWord takes a list of words to block comments containing these words. These words are added to the existing list
func (s *Service) BlockWord(r *http.Request, args *commentapi.BlockWordArgs, reply *commentapi.BlockWordRespose) error {
	if len(args.Words) == 0 {
		return api.StatusError{Err: errors.Err("words to block %s must exist", args.Words)}
	}
	creatorChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignatureAndTS(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}

	settings, err := helper.FindOrCreateSettings(creatorChannel)
	if err != nil {
		return err
	}
	var existingWords []string
	if !settings.MutedWords.IsZero() {
		existingWords = strings.Split(settings.MutedWords.String, ",")
	}
	wordsToAdd := strings.Split(args.Words, ",")

	existingWords = append(existingWords, wordsToAdd...)
	settings.MutedWords.SetValid(strings.Join(existingWords, ","))
	err = settings.Update(db.RW, boil.Infer())
	if err != nil {
		return errors.Err(err)
	}
	reply.WordList = existingWords
	return nil
}

// UnBlockWord takes a list of words to remove from the list of blocked words if they exist.
func (s *Service) UnBlockWord(r *http.Request, args *commentapi.UnBlockWordArgs, reply *commentapi.BlockWordRespose) error {
	creatorChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	if creatorChannel == nil {
		return api.StatusError{Err: errors.Err("could not find channel %s with channel id %s", args.ChannelName, args.ChannelID), Status: http.StatusBadRequest}
	}
	err = lbry.ValidateSignatureAndTS(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}

	settings, err := helper.FindOrCreateSettings(creatorChannel)
	if err != nil {
		return err
	}
	existingWords := strings.Split(settings.MutedWords.String, ",")
	wordsToRemove := strings.Split(args.Words, ",")
	remainingWords := make([]string, 0)
skip:
	for _, word := range existingWords {
		for _, wordToRemove := range wordsToRemove {
			if wordToRemove == word {
				continue skip
			}
		}
		remainingWords = append(remainingWords, word)
	}
	settings.MutedWords = null.String{String: "", Valid: false}
	if len(remainingWords) > 0 {
		settings.MutedWords.SetValid(strings.Join(remainingWords, ","))
	}

	err = settings.Update(db.RW, boil.Infer())
	if err != nil {
		return errors.Err(err)
	}
	reply.WordList = remainingWords
	return nil
}

// ListBlockedWords returns a list of all the current blocked words for a channel.
func (s *Service) ListBlockedWords(r *http.Request, args *commentapi.ListBlockedWordsArgs, reply *commentapi.BlockWordRespose) error {
	creatorChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignatureAndTS(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}

	settings, err := helper.FindOrCreateSettings(creatorChannel)
	if err != nil {
		return err
	}
	existingWords := strings.Split(settings.MutedWords.String, ",")
	reply.WordList = existingWords
	return nil
}
