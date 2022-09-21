package helper

import (
	"database/sql"
	"time"

	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/model"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/karlseguin/ccache/v2"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var channelCache = ccache.New(ccache.Configure().GetsPerPromote(1).MaxSize(100000))

// FindOrCreateChannel gets the channel from commentron database or creates it and returns it
func FindOrCreateChannel(channelClaimID, channelName string) (*model.Channel, error) {
	item, err := channelCache.Fetch(channelClaimID, 1*time.Hour, func() (interface{}, error) {
		channel, err := getChannel(channelClaimID, channelName)
		if err != nil {
			return nil, err
		}
		return channel, nil
	})
	if err != nil {
		return nil, err
	}
	c, ok := item.Value().(*model.Channel)
	if !ok {
		return nil, errors.Err("could not convert item to channel")
	}
	return c, nil
}

func getChannel(channelClaimID, channelName string) (*model.Channel, error) {
	channel, err := model.Channels(model.ChannelWhere.ClaimID.EQ(channelClaimID)).One(db.RO)
	if errors.Is(err, sql.ErrNoRows) {
		if channelName == "" {
			return nil, errors.Err("channel name cannot be blank")
		}
		channel = &model.Channel{
			ClaimID: channelClaimID,
			Name:    channelName,
		}
		err = nil
		err := channel.Insert(db.RW, boil.Infer())
		if err != nil {
			return nil, errors.Err(err)
		}
	}
	return channel, errors.Err(err)
}

var settingsCache = ccache.New(ccache.Configure().GetsPerPromote(1).MaxSize(100000))

// FindOrCreateSettings gets the settings for the creator from commentron database or creates it and returns it
func FindOrCreateSettings(creatorChannel *model.Channel) (*model.CreatorSetting, error) {
	item, err := settingsCache.Fetch(creatorChannel.ClaimID, 1*time.Hour, func() (interface{}, error) {
		setting, err := getSettings(creatorChannel)
		if err != nil {
			return nil, err
		}
		return setting, nil
	})
	if err != nil {
		return nil, err
	}
	setting, ok := item.Value().(*model.CreatorSetting)
	if !ok {
		return nil, errors.Err("could not convert item to creator setting")
	}
	return setting, nil
}

func getSettings(creatorChannel *model.Channel) (*model.CreatorSetting, error) {
	settings, err := creatorChannel.CreatorChannelCreatorSettings().One(db.RO)
	if errors.Is(err, sql.ErrNoRows) {
		settings = &model.CreatorSetting{CreatorChannelID: creatorChannel.ClaimID, CommentsEnabled: null.BoolFrom(true)}
		err = nil
		err := settings.Insert(db.RW, boil.Infer())
		if err != nil {
			return nil, errors.Err(err)
		}
	}
	return settings, nil
}
