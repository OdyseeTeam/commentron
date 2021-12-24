package commentapi

type AppealBlockListArgs struct {
	Authorization
}

type AppealBlockListResponse struct {
	Blocks []Appeal `json:"blocks"`
}

type AppealListArgs struct {
	ModAuthorization
}

type AppealListResponse struct {
	Appeals          []Appeal `json:"appeals"`
	ModeratedAppeals []Appeal `json:"moderated_appeals"`
}

type AppealStatus int

const (
	// AppealPending the default value for all appeals
	AppealPending AppealStatus = iota
	// AppealEscalated appeal is escalated to shared block list owner
	AppealEscalated
	// AppealAccepted creator who blocked, has accepted the appeal
	AppealAccepted
	// AppealRejected creator who blocked, has rejected the appeal
	AppealRejected
)

type AppealRequest struct {
	AppealMessage   string       `json:"appeal_message"`
	ResponseMessage string       `json:"response_message"`
	AppealStatus    AppealStatus `json:"status"`
	TxID            string       `json:"tx_id,omitempty"`
}

type Appeal struct {
	BlockedList    SharedBlockedList `json:"blocked_list,omitempty"`
	BlockedChannel BlockedChannel    `json:"blocked_channel"`
	AppealRequest  AppealRequest     `json:"appeal_request,omitempty"`
}

type AppealFileArgs struct {
	Authorization

	SharedBlockedListID  uint64 `json:"blocked_list_id"`
	BlockedByChannelID   string `json:"blocked_by_channel_id"`
	BlockedByChannelName string `json:"blocked_by_channel_name"`
	AppealMessage        string `json:"appeal_message"`
	TxID                 string `json:"tx_id,omitempty"`
}

type AppealCloseArgs struct {
	Authorization

	BlockedChannelID   string       `json:"blocked_channel_id"`
	BlockedChannelName string       `json:"blocked_channel_name"`
	AppealStatus       AppealStatus `json:"status"`
	ResponseMessage    string       `json:"response_message"`
}
