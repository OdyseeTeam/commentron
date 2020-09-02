package lbry

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lbryio/commentron/util"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/sirupsen/logrus"
)

// APIToken is the token allowed to access the api used for internal-apis
var APIToken string

// APIURL is the url for internal-apis to be used by commentron
var APIURL string

// CommentResponse is the response structure from internal-apis for the comment event api
type CommentResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
	Data    string      `json:"data"`
}

// NotifyOptions Are the options used to construct the comment event api signature.
type NotifyOptions struct {
	ActionType string
	CommentID  string
	ChannelID  *string
	ParentID   *string
	Comment    *string
	ClaimID    string
}

// Notify notifies internal-apis of a new comment when one is recieved.
func Notify(options NotifyOptions) {
	err := notify(options)
	if err != nil {
		logrus.Error("API Notification: ", err)
	}
}

func notify(options NotifyOptions) error {
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
	defer util.CloseBody(response.Body)
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Err(err)
	}
	var me CommentResponse
	err = json.Unmarshal(b, &me)
	if err != nil {
		return errors.Err(err)
	}
	if response.StatusCode > 200 {
		if response.StatusCode <= 300 {
			logrus.Warning("Notification Failure[Status - ", response.StatusCode, "] : ", string(b))
		} else {
			logrus.Error("Notification Failure[Status - ", response.StatusCode, "] : ", string(b))
		}
	}
	return nil
}
