/*
   GoToSocial
   Copyright (C) 2021 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package federatingdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/sirupsen/logrus"
	"github.com/superseriousbusiness/gotosocial/internal/ap"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/messages"
	"github.com/superseriousbusiness/gotosocial/internal/util"
)

func (f *federatingDB) Accept(ctx context.Context, accept vocab.ActivityStreamsAccept) error {
	l := f.log.WithFields(
		logrus.Fields{
			"func":   "Accept",
			"asType": accept.GetTypeName(),
		},
	)
	m, err := streams.Serialize(accept)
	if err != nil {
		return err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	l.Debugf("received ACCEPT asType %s", string(b))

	targetAcctI := ctx.Value(util.APAccount)
	if targetAcctI == nil {
		// If the target account wasn't set on the context, that means this request didn't pass through the
		// API, but came from inside GtS as the result of another activity on this instance. That being so,
		// we can safely just ignore this activity, since we know we've already processed it elsewhere.
		return nil
	}
	targetAcct, ok := targetAcctI.(*gtsmodel.Account)
	if !ok {
		l.Error("ACCEPT: target account was set on context but couldn't be parsed")
		return nil
	}

	fromFederatorChanI := ctx.Value(util.APFromFederatorChanKey)
	if fromFederatorChanI == nil {
		l.Error("ACCEPT: from federator channel wasn't set on context")
		return nil
	}
	fromFederatorChan, ok := fromFederatorChanI.(chan messages.FromFederator)
	if !ok {
		l.Error("ACCEPT: from federator channel was set on context but couldn't be parsed")
		return nil
	}

	acceptObject := accept.GetActivityStreamsObject()
	if acceptObject == nil {
		return errors.New("ACCEPT: no object set on vocab.ActivityStreamsAccept")
	}

	for iter := acceptObject.Begin(); iter != acceptObject.End(); iter = iter.Next() {
		// check if the object is an IRI
		if iter.IsIRI() {
			// we have just the URI of whatever is being accepted, so we need to find out what it is
			acceptedObjectIRI := iter.GetIRI()
			if util.IsFollowPath(acceptedObjectIRI) {
				// ACCEPT FOLLOW
				gtsFollowRequest := &gtsmodel.FollowRequest{}
				if err := f.db.GetWhere(ctx, []db.Where{{Key: "uri", Value: acceptedObjectIRI.String()}}, gtsFollowRequest); err != nil {
					return fmt.Errorf("ACCEPT: couldn't get follow request with id %s from the database: %s", acceptedObjectIRI.String(), err)
				}

				// make sure the addressee of the original follow is the same as whatever inbox this landed in
				if gtsFollowRequest.AccountID != targetAcct.ID {
					return errors.New("ACCEPT: follow object account and inbox account were not the same")
				}
				follow, err := f.db.AcceptFollowRequest(ctx, gtsFollowRequest.AccountID, gtsFollowRequest.TargetAccountID)
				if err != nil {
					return err
				}

				fromFederatorChan <- messages.FromFederator{
					APObjectType:     ap.ActivityFollow,
					APActivityType:   ap.ActivityAccept,
					GTSModel:         follow,
					ReceivingAccount: targetAcct,
				}

				return nil
			}
		}

		// check if iter is an AP object / type
		if iter.GetType() == nil {
			continue
		}
		switch iter.GetType().GetTypeName() {
		// we have the whole object so we can figure out what we're accepting
		case ap.ActivityFollow:
			// ACCEPT FOLLOW
			asFollow, ok := iter.GetType().(vocab.ActivityStreamsFollow)
			if !ok {
				return errors.New("ACCEPT: couldn't parse follow into vocab.ActivityStreamsFollow")
			}
			// convert the follow to something we can understand
			gtsFollow, err := f.typeConverter.ASFollowToFollow(ctx, asFollow)
			if err != nil {
				return fmt.Errorf("ACCEPT: error converting asfollow to gtsfollow: %s", err)
			}
			// make sure the addressee of the original follow is the same as whatever inbox this landed in
			if gtsFollow.AccountID != targetAcct.ID {
				return errors.New("ACCEPT: follow object account and inbox account were not the same")
			}
			follow, err := f.db.AcceptFollowRequest(ctx, gtsFollow.AccountID, gtsFollow.TargetAccountID)
			if err != nil {
				return err
			}

			fromFederatorChan <- messages.FromFederator{
				APObjectType:     ap.ActivityFollow,
				APActivityType:   ap.ActivityAccept,
				GTSModel:         follow,
				ReceivingAccount: targetAcct,
			}

			return nil
		}
	}

	return nil
}
