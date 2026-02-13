package db

import (
	"testing"

	"github.com/OdyseeTeam/commentron/env"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

// Note from Johnny
//
// It's not really our style to add unit tests but,
//
// 1) I'm less familiar with this code base still and this is a lot of change; and,
// 2) I really want to make sure we don't break anything by turning off the soft delete
//    functionality of SQLBoiler.
//
// If this test passes, soft deletes work *without* modification.
func TestSoftDelete(t *testing.T) {
	conf, err := env.NewWithEnvVars()
	if err != nil {
		logrus.Panic(err)
	}

	err = Init(conf.MySQLDsnRO, conf.MySQLDsnRW, true)
	assert.NoError(t, err)

	// Create a test channel
	channel := &model.Channel{ClaimID: "soft-delete-test"}
	assert.NoError(t, channel.Upsert(RW, boil.Infer(), boil.Infer()))

	// Create two comments
	comment1 := &model.Comment{
		CommentID: "comment1",
		ChannelID: null.StringFrom(channel.ClaimID),
	}
	comment2 := &model.Comment{
		CommentID: "comment2",
		ChannelID: null.StringFrom(channel.ClaimID),
	}
	assert.NoError(t, comment1.Upsert(RW, boil.Infer(), boil.Infer()))
	assert.NoError(t, comment2.Upsert(RW, boil.Infer(), boil.Infer()))

	// Get the comments back from the db
	comments, err := model.Comments(
		model.CommentWhere.ChannelID.EQ(null.StringFrom(channel.ClaimID)),
	).All(RW)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)

	// Delete the first comment
	assert.NoError(t, comment1.Delete(RW, false))
	comments, err = model.Comments(
		model.CommentWhere.ChannelID.EQ(null.StringFrom(channel.ClaimID)),
	).All(RW)

	// Check default query
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, comments[0].CommentID, comment2.CommentID)

	// Check for soft deletes
	comments, err = model.Comments(
		model.CommentWhere.ChannelID.EQ(null.StringFrom(channel.ClaimID)),
		qm.WithDeleted(),
	).All(RW)

	assert.NoError(t, err)
	assert.Len(t, comments, 2)
}
