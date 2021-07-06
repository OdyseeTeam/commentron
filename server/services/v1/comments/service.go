package comments

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	m "github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

/*
'ping': ping,
DONE 'get_claim_comments': handle_get_claim_comments,  # this gets used ( params: claimid, page, page-size, include-replies, order-by )
DONE 'get_channel_from_comment_id': handle_get_channel_from_comment_id,  # this gets used by SDK
DONE 'create_comment': handle_create_comment,  # this gets used
'abandon_comment': handle_abandon_comment,  # this gets used
'edit_comment': handle_edit_comment  # this gets used
NEW APIS
comment count per claim
comment count per parent
comments per parent ( order-by time, rating, page, page-size )
page for comment ( params: page, [size], order-by )

*/

// Service is the service struct defined for the comment package for rpc service "comment.*"
type Service struct{}

// Create creates a comment
func (c *Service) Create(r *http.Request, args *commentapi.CreateArgs, reply *commentapi.CreateResponse) error {
	return create(r, args, reply)
}

// List lists comments based on filters and arguments passed. The returned result is dynamic based on the args passed
func (c *Service) List(r *http.Request, args *commentapi.ListArgs, reply *commentapi.ListResponse) error {
	return list(r, args, reply)
}

// GetChannelFromCommentID gets the channel info for a specific comment, this is really only used by the sdk
func (c *Service) GetChannelFromCommentID(_ *http.Request, args *commentapi.ChannelArgs, reply *commentapi.ChannelResponse) error {
	comment, err := m.Comments(m.CommentWhere.CommentID.EQ(args.CommentID), qm.Load(m.CommentRels.Channel)).OneG()
	if errors.Is(err, sql.ErrNoRows) {
		return api.StatusError{Err: errors.Err("could not find comment for comment id"), Status: http.StatusBadRequest}
	}
	if err != nil {
		return errors.Err(err)
	}
	if comment.R == nil && comment.R.Channel == nil {
		return api.StatusError{Err: errors.Err("could not find channel for comment"), Status: http.StatusBadRequest}
	}
	reply.ChannelID = comment.R.Channel.ClaimID
	reply.ChannelName = comment.R.Channel.Name
	return nil
}

// Abandon deletes a comment
func (c *Service) Abandon(_ *http.Request, args *commentapi.AbandonArgs, reply *commentapi.AbandonResponse) error {
	item, err := abandon(args)
	if err != nil {
		return errors.Err(err)
	}
	reply.CommentItem = item
	reply.Abandoned = true

	go lbry.API.Notify(lbry.NotifyOptions{
		ActionType: "D",
		CommentID:  item.CommentID,
		ChannelID:  &item.ChannelID,
		ParentID:   &item.ParentID,
		Comment:    &item.Comment,
		ClaimID:    item.ClaimID,
	})

	return nil
}

// Edit edits a comment
func (c *Service) Edit(_ *http.Request, args *commentapi.EditArgs, reply *commentapi.EditResponse) error {
	item, err := edit(args)
	if err != nil {
		return errors.Err(err)
	}
	reply.CommentItem = item

	go lbry.API.Notify(lbry.NotifyOptions{
		ActionType: "U",
		CommentID:  item.CommentID,
		ChannelID:  &item.ChannelID,
		ParentID:   &item.ParentID,
		Comment:    &item.Comment,
		ClaimID:    item.ClaimID,
	})
	return nil
}

// ByID returns the comment from the comment id passed in
func (c *Service) ByID(r *http.Request, args *commentapi.ByIDArgs, reply *commentapi.ByIDResponse) error {
	item, ancestors, err := byID(r, args)
	if err != nil {
		return err
	}
	reply.Item = item
	reply.Ancestors = ancestors
	return nil
}

// Pin sets the pinned flag on a comment
func (c *Service) Pin(r *http.Request, args *commentapi.PinArgs, reply *commentapi.PinResponse) error {
	item, err := pin(r, args)
	if err != nil {
		return err
	}
	reply.Item = item
	return nil
}

// SuperChatList returns comments that are super chat only.
func (c *Service) SuperChatList(r *http.Request, args *commentapi.SuperListArgs, reply *commentapi.SuperListResponse) error {
	return superChatList(r, args, reply)
}
