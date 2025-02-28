package lbry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/OdyseeTeam/commentron/helper"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/sirupsen/logrus"
)

var apiToken string
var apiURL string

type apiClient struct{}

// CommentResponse is the response structure from internal-apis for the comment event api
type CommentResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
	Data    string      `json:"data"`
}

// CheckPerkResponse is the response structure from internal-apis for the comment event api
type CheckPerkResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
	Data    struct {
		HasAccess bool `json:"has_access"`
	} `json:"data"`
}

// Notify notifies internal-apis of a new comment when one is recieved.
func (c apiClient) Notify(options NotifyOptions) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("recovered from failed notification to internal-apis: ", r)
		}
	}()
	err := notify(options)
	if err != nil {
		logrus.Error("API Notification: ", errors.FullTrace(err))
	}
}

func (c apiClient) CheckPerk(options CheckPerkOptions) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("recovered from failed notification to internal-apis: ", r)
		}
	}()
	const apiPath = "/membership_perk/check"
	client := http.Client{}
	form := make(url.Values)
	form.Set("auth_token", apiToken)
	form.Set("claim_id", options.ClaimID)
	form.Set("channel_claim_id", options.ChannelClaimID)
	form.Set("type", options.Type)
	if options.Environment != nil {
		form.Set("environment", *options.Environment)
	}

	response, err := client.PostForm(apiURL+apiPath, form)
	if err != nil {
		return false, errors.Err(err)
	}
	if response == nil {
		return false, errors.Err("No response from internal APIs")
	}
	defer helper.CloseBody(response.Body)
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return false, errors.Err(err)
	}
	if response.StatusCode != http.StatusOK {
		return false, errors.Err("API returned non-200 status code: %d", response.StatusCode)
	}
	var perkRes CheckPerkResponse
	err = json.Unmarshal(b, &perkRes)
	if err != nil {
		return false, errors.Err(err)
	}

	return perkRes.Data.HasAccess, nil
}

// ArweavePaymentDetailsResponse is the response structure from internal-apis for the arweave payment details api
type ArweavePaymentDetailsResponse struct {
	Amount         uint64    `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	UserID         uint64    `json:"user_id"`
	ChannelClaimID string    `json:"channel_claim_id"`
	TippedAt       time.Time `json:"tipped_at"`
}
type apiArweavePaymentDetailsResponse struct {
	Success bool                          `json:"success"`
	Error   interface{}                   `json:"error"`
	Data    ArweavePaymentDetailsResponse `json:"data"`
}

func (c apiClient) GetDetailsForTransaction(txID string) (*ArweavePaymentDetailsResponse, error) {
	const apiPath = "/arweave/payment/retrieve"
	client := http.Client{}
	form := make(url.Values)
	form.Set("auth_token", apiToken)
	form.Set("tx_id", txID)

	response, err := client.PostForm(apiURL+apiPath, form)
	if err != nil {
		return nil, errors.Err(err)
	}
	defer helper.CloseBody(response.Body)

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Err(err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.Err("API returned non-200 status code: %d", response.StatusCode)
	}
	var pr apiArweavePaymentDetailsResponse
	err = json.Unmarshal(b, &pr)

	return &pr.Data, errors.Err(err)
}

func notify(options NotifyOptions) error {
	const apiPath = "/event/comment"
	c := http.Client{}
	form := make(url.Values)
	form.Set("auth_token", apiToken)
	form.Set("action_type", options.ActionType)
	form.Set("comment_id", options.CommentID)
	form.Set("claim_id", options.ClaimID)
	form.Set("amount", fmt.Sprintf("%d", options.Amount))
	form.Set("is_fiat", strconv.FormatBool(options.IsFiat))
	form.Set("is_protected", strconv.FormatBool(options.IsProtected))

	if options.Comment != nil {
		form.Set("comment", *options.Comment)
	}

	if options.ChannelID != nil {
		form.Set("channel_id", *options.ChannelID)
	}

	if options.ParentID != nil {
		form.Set("parent_id", *options.ParentID)
	}

	if options.IsFiat && options.Currency != nil {
		form.Set("currency", *options.Currency)
	}

	response, err := c.PostForm(apiURL+apiPath, form)
	if err != nil {
		return errors.Err(err)
	}
	if response == nil {
		return errors.Err("No response from internal APIs")
	}
	defer helper.CloseBody(response.Body)
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Err(err)
	}
	if response.StatusCode > http.StatusOK {
		if response.StatusCode <= http.StatusMultipleChoices {
			logrus.Warning("Notification Failure[Status - ", response.StatusCode, "] : ")
		} else {
			logrus.Error("Notification Failure[Status - ", response.StatusCode, "] : ")
		}
		return errors.Err("API returned non-200 status code: %d - %s", response.StatusCode, string(b))
	}
	//todo: add recursion if necessary

	var me CommentResponse
	err = json.Unmarshal(b, &me)
	if err != nil {
		return errors.Err(err)
	}
	return nil
}
