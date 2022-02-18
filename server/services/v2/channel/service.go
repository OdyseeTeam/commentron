package channel

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/auth"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

// Service is the service for the server package "server.*"
type Service struct{}

// Args arguments for the server.Status rpc call
type Args struct {
	MyChannels []channel
}

// Response response for the server.Status rpc call
type Response struct {
	Confirmed   []channel
	UnConfirmed []channel
}

type channel struct {
	commentapi.Authorization
}

// Status shows the status of commentron channels for a user (OAuth Linking)
func (t *Service) Status(r *http.Request, args *Args, reply *Response) error {
	_, user, err := auth.Authenticate(r, nil)
	if err != nil {
		return err
	}
	confirmedChannels := make(map[string]bool)
	if user != nil {
		channels, err := model.Channels(model.ChannelWhere.Sub.EQ(null.StringFrom(user.Sub))).All(db.RO)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}

		for _, c := range channels {
			confirmedChannels[c.ClaimID] = true
		}
	}

	for _, c := range args.MyChannels {
		if !confirmedChannels[c.ChannelID] {
			err := lbry.ValidateSignatureAndTS(c.ChannelID, c.Signature, c.SigningTS, c.ChannelName)
			if err != nil {
				reply.UnConfirmed = append(reply.UnConfirmed, c)
			}
			channel, err := helper.FindOrCreateChannel(c.ChannelID, c.ChannelName)
			if err != nil {
				return err
			}
			if user != nil {
				channel.Sub.SetValid(user.Sub)
				err := channel.Update(db.RW, boil.Infer())
				if err != nil {
					return errors.Err(err)
				}
			}
			continue
		}
		reply.Confirmed = append(reply.Confirmed, c)
	}
	return nil
}
