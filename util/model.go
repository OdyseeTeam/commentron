package util

import (
	"database/sql"

	"github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/sqlboiler/boil"
)

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
