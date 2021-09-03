package blockedlists

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/sqlboiler/boil"
)

func update(_ *http.Request, args *commentapi.SharedBlockedListUpdateArgs, reply *commentapi.SharedBlockedList) error {
	ownerChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignature(ownerChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}

	list, err := model.BlockedLists(model.BlockedListWhere.ChannelID.EQ(ownerChannel.ClaimID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if args.Remove {
		err := list.BlockedListInvites().DeleteAll(db.RW)
		if err != nil {
			return errors.Err(err)
		}
		err = list.Delete(db.RW)
		if err != nil {
			return errors.Err(err)
		}
	}
	var created bool
	if list == nil {
		if args.Name == nil {
			return api.StatusError{Err: errors.Err("a name must be specified if a new list will get created")}
		}
		list = &model.BlockedList{ChannelID: ownerChannel.ClaimID, Name: *args.Name}
		err := list.Insert(db.RW, boil.Infer())
		if err != nil {
			return errors.Err(err)
		}
		created = true
	}

	if args.Name != nil {
		list.Name = *args.Name
	}
	if args.Category != nil {
		list.Category = *args.Category
	}
	if args.Description != nil {
		list.Description = *args.Description
	}
	if args.MemberInviteEnabled != nil {
		list.MemberInviteEnabled.SetValid(*args.MemberInviteEnabled)
	}
	if args.InviteExpiration != nil {
		list.InviteExpiration.SetValid(*args.InviteExpiration)
	}
	if args.StrikeOne != nil {
		list.StrikeOne.SetValid(*args.StrikeOne)
	}
	if args.StrikeTwo != nil {
		list.StrikeTwo.SetValid(*args.StrikeTwo)
	}
	if args.StrikeThree != nil {
		list.StrikeThree.SetValid(*args.StrikeThree)
	}
	if args.CurseJarAmount != nil {
		list.CurseJarAmount.SetValid(*args.CurseJarAmount)
	}

	err = list.Update(db.RW, boil.Infer())
	if err != nil {
		return errors.Err(err)
	}

	if created {
		blockedList := map[string]interface{}{model.BlockedEntryColumns.BlockedListID: list.ID}
		err := ownerChannel.CreatorChannelBlockedEntries().UpdateAll(db.RW, blockedList)
		if err != nil {
			return errors.Err(err)
		}
		ownerChannel.BlockedListID.SetValid(list.ID)
		ownerChannel.BlockedListInviteID.SetValid(list.ID)
		err = ownerChannel.Update(db.RW, boil.Infer())
		if err != nil {
			return errors.Err(err)
		}
	}

	return populateSharedBlockedList(reply, list)
}

func populateSharedBlockedList(list *commentapi.SharedBlockedList, modelList *model.BlockedList) error {
	list.ID = modelList.ID
	list.Name = &modelList.Name
	list.Description = &modelList.Description
	list.Category = &modelList.Category

	if modelList.MemberInviteEnabled.Valid {
		list.MemberInviteEnabled = &modelList.MemberInviteEnabled.Bool
	}
	if modelList.InviteExpiration.Valid {
		list.InviteExpiration = &modelList.InviteExpiration.Uint64
	}
	if modelList.StrikeOne.Valid {
		list.StrikeOne = &modelList.StrikeOne.Uint64
	}
	if modelList.StrikeTwo.Valid {
		list.StrikeTwo = &modelList.StrikeTwo.Uint64
	}
	if modelList.StrikeThree.Valid {
		list.StrikeThree = &modelList.StrikeThree.Uint64
	}
	if modelList.CurseJarAmount.Valid {
		list.CurseJarAmount = &modelList.CurseJarAmount.Uint64
	}

	return nil
}
