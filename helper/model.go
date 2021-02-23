package helper

import (
	"database/sql"

	"github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/sqlboiler/boil"
)

// FindOrCreateChannel gets the channel from commentron database or creates it and returns it
func FindOrCreateChannel(channelClaimID, channelName string) (*model.Channel, error) {
	channel, err := model.Channels(model.ChannelWhere.ClaimID.EQ(channelClaimID)).OneG()
	if errors.Is(err, sql.ErrNoRows) {
		channel = &model.Channel{
			ClaimID: channelClaimID,
			Name:    channelName,
		}
		err = nil
		err := channel.InsertG(boil.Infer())
		if err != nil {
			return nil, errors.Err(err)
		}
	}
	return channel, errors.Err(err)
}

// FindOrCreateSettings gets the settings for the creator from commentron database or creates it and returns it
func FindOrCreateSettings(creatorChannel *model.Channel) (*model.CreatorSetting, error) {
	settings, err := creatorChannel.CreatorChannelCreatorSettings().OneG()
	if errors.Is(err, sql.ErrNoRows) {
		settings = &model.CreatorSetting{CreatorChannelID: creatorChannel.ClaimID}
		err = nil
		err := settings.InsertG(boil.Infer())
		if err != nil {
			return nil, errors.Err(err)
		}
	}
	return settings, errors.Err(err)
}
