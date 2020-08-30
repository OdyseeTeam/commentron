package lbry

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lbryio/lbry.go/v2/extras/errors"
)

var APIToken string
var APIURL string

type CommentResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
	Data    string      `json:"data"`
}

type NotifyOptions struct {
	ActionType string
	CommentID  string
	ChannelID  *string
	ParentID   *string
	Comment    *string
	ClaimID    string
}

func Notify(options NotifyOptions) error {
	c := http.Client{}
	form := make(url.Values)
	form.Set("auth_token", APIToken)
	form.Set("action_type", options.ActionType)
	form.Set("comment_id", options.CommentID)
	form.Set("claim_id", options.ClaimID)

	if options.Comment != nil {
		form.Set("comment", *options.Comment)
	}

	if options.ChannelID != nil {
		form.Set("channel_id", *options.ChannelID)
	}

	if options.ParentID != nil {
		form.Set("parent_id", *options.ParentID)
	}

	response, err := c.PostForm(APIURL, form)
	if err != nil {
		return errors.Err(err)
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Err(err)
	}
	var me CommentResponse
	err = json.Unmarshal(b, &me)
	if err != nil {
		return errors.Err(err)
	}
	return nil
}
