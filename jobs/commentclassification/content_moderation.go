package commentclassification

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/metrics"
	"github.com/OdyseeTeam/commentron/model"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/sirupsen/logrus"
)

// PollAndClassifyNewComments looks for new comments and updates them by calling classification api
//
// I'm using the existing jobs code to make organization obvious. But its called much more frequently
// than the other jobs, so it's a bit of a kludge. Related to this, this function quits early
// to prevent two competing poll jobs from running at once (which could lead to some non-fatal
// but annoying data processing errors), this function will quit early. This makes the
// poll_and_classify_new_comments` job metric a misleading under some circumstances.
//
// The commentron::moderation::comments_classified counter is much clearer for grafana stats.
func PollAndClassifyNewComments() {
	if jobAlreadyInProgress := startJob(); jobAlreadyInProgress {
		return
	}
	defer endJob()

	startTime := time.Now()
	defer metrics.Job(startTime, "poll_and_classify_new_comments")

	batchSize := defaultBatchSize

	lastKnownClassificationTimestamp, err := getLastKnownClassificationTimestamp()
	if err != nil {
		logrus.Error("Error getting last known classified comment: ", err)
		return
	}

	for {
		metrics.PollingCallsForClassifierJob.Inc()

		toClassify, err := queryCommentBatch(lastKnownClassificationTimestamp, batchSize)
		if err != nil {
			logrus.Error("Error getting last known classified comment: ", err)
			return
		} else if len(toClassify) == 0 {
			return
		}

		// call inference server with classifications
		classifications, err := inferCommentClassifications(toClassify)
		if err != nil {
			logrus.Error("Error getting comment classifications: ", err)
			return
		}

		// insert all the classifications
		//
		// this is a little annoying with sqlboiler because it doesn't support bulk inserts
		// (because it wants the primary key back to update the struct)
		for _, classification := range classifications {
			// This must be an upsert because comments can be edited.
			err = classification.Upsert(db.RW, boil.Infer(), boil.Infer())
			if err != nil {
				logrus.Error("Error inserting comment classification: ", err)
				// DO NOT RETURN: keep trying others since we have classifications
				// for them. If they can proceed, the system doesn't get stuck.
			} else {
				metrics.CommentsClassified.Inc()
			}
		}

		// Let next cron job take over if we "hit the end"
		if len(toClassify) != batchSize {
			break
		}
	}
}

func inferCommentClassifications(comments model.CommentSlice) (model.CommentClassificationSlice, error) {
	var classifications model.CommentClassificationSlice

	// Package up the request as a json list of dicts
	reqItems := make([]map[string]string, len(comments))
	for i, comment := range comments {
		reqItems[i] = map[string]string{
			"id":      comment.CommentID,
			"comment": comment.Body,
		}
	}
	reqBytes, err := json.Marshal(reqItems)
	if err != nil {
		return nil, err
	}

	// Make the request
	client := http.Client{
		Timeout: time.Second * 20,
	}
	resp, err := client.Post(inferenceServiceURI, "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Check for 200 status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("inference service returned status %d", resp.StatusCode)
	}

	// Parse it
	var classificationResp classificationResp
	err = json.NewDecoder(resp.Body).Decode(&classificationResp)
	if err != nil {
		return nil, err
	}

	if len(classificationResp.Classifications) == 0 {
		return nil, nil
	}

	modelIdent := null.StringFrom(classificationResp.ModelIdent)

	lookupTable := make(map[string]*classification, len(classificationResp.Classifications))
	for _, classification := range classificationResp.Classifications {
		lookupTable[classification.ID] = classification
	}

	for _, comment := range comments {
		classification := lookupTable[comment.CommentID]
		if classification == nil {
			continue
		}

		classifications = append(classifications, &model.CommentClassification{
			CommentID: comment.CommentID,
			Timestamp: comment.Timestamp,

			Toxicity:       classification.Toxicity,
			SevereToxicity: classification.SevereToxicity,
			Obscene:        classification.Obscene,
			IdentityAttack: classification.IdentityAttack,
			Insult:         classification.Insult,
			Threat:         classification.Threat,
			SexualExplicit: classification.SexualExplicit,

			ModelIdent: modelIdent,
		})
	}

	return classifications, nil
}

type classificationResp struct {
	ModelIdent      string            `json:"model_ident"`
	Classifications []*classification `json:"classifications"`
}

type classification struct {
	ID             string  `json:"id"`
	Toxicity       float32 `json:"toxicity"`
	SevereToxicity float32 `json:"severe_toxicity"`
	Obscene        float32 `json:"obscene"`
	IdentityAttack float32 `json:"identity_attack"`
	Insult         float32 `json:"insult"`
	Threat         float32 `json:"threat"`
	SexualExplicit float32 `json:"sexual_explicit"`
}

// Get comments which occurred after the last classification and before now.
//
// Note: I'm using the RO connection because we don't really need transactions,
// as this operation is effectively idempotent.
//
// WARNING: This assumes that comments have new timestamps on edit.
// as per v1/comments/edit.go they do. If this changes, this will break
// in that a user can post a clean comment; wait the 5 minutes; then edit
// it to be toxic and it will not be updated.
func queryCommentBatch(lastKnownClassificationTimestamp, batchSize int) (model.CommentSlice, error) {
	commentTbl := model.TableNames.Comment
	commentTimestampCol := commentTbl + "." + model.CommentColumns.Timestamp
	commentIDCol := commentTbl + "." + model.CommentColumns.CommentID
	classificationTbl := model.TableNames.CommentClassification
	classificationCommentIDCol := classificationTbl + "." + model.CommentClassificationColumns.CommentID

	comments, err := model.Comments(
		// To ensure none are missed, use = instead of > but with
		// an outer join and nullity check to skip duplicates.
		model.CommentWhere.Timestamp.GTE(lastKnownClassificationTimestamp),
		model.CommentWhere.Timestamp.LT(int(time.Now().Unix())),

		qm.LeftOuterJoin(fmt.Sprintf("%s ON %s = %s", classificationTbl, commentIDCol, classificationCommentIDCol)),
		qm.Where(classificationCommentIDCol+" IS NULL"),

		// Poll in chronological order.
		qm.OrderBy(commentTimestampCol+" ASC"),

		// But don't overwhelm the remote inference server.
		qm.Limit(batchSize),
	).All(db.RO)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return comments, nil
}

// Get the last known timestamp for the comment Classifications.
//
// WARNING: this assumes the `created_at` column replicates the `timestamp`.
// WARNING: depending on timestamp granularity, this can miss some comments?
func getLastKnownClassificationTimestamp() (int, error) {
	cc, err := model.CommentClassifications(qm.OrderBy(model.CommentClassificationColumns.Timestamp + " DESC")).One(db.RO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no Classifications exist, start from five minutes ago.
			return int(time.Now().Add(-5 * time.Minute).Unix()), nil
		}
		return -1, err
	}

	return cc.Timestamp, nil
}

func startJob() (jobInProgressAlready bool) {
	commentClassificationMutex.Lock()
	defer commentClassificationMutex.Unlock()

	if commentClassificationInProgress {
		return true
	}

	commentClassificationInProgress = true
	return false
}

func endJob() {
	commentClassificationMutex.Lock()
	defer commentClassificationMutex.Unlock()

	commentClassificationInProgress = false
}
