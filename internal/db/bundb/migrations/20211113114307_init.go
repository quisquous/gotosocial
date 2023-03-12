// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package migrations

import (
	"context"

	gtsmodel "github.com/superseriousbusiness/gotosocial/internal/db/bundb/migrations/20211113114307_init"
	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		models := []interface{}{
			&gtsmodel.Account{},
			&gtsmodel.Application{},
			&gtsmodel.Block{},
			&gtsmodel.DomainBlock{},
			&gtsmodel.EmailDomainBlock{},
			&gtsmodel.Follow{},
			&gtsmodel.FollowRequest{},
			&gtsmodel.MediaAttachment{},
			&gtsmodel.Mention{},
			&gtsmodel.Status{},
			&gtsmodel.StatusToEmoji{},
			&gtsmodel.StatusToTag{},
			&gtsmodel.StatusFave{},
			&gtsmodel.StatusBookmark{},
			&gtsmodel.StatusMute{},
			&gtsmodel.Tag{},
			&gtsmodel.User{},
			&gtsmodel.Emoji{},
			&gtsmodel.Instance{},
			&gtsmodel.Notification{},
			&gtsmodel.RouterSession{},
			&gtsmodel.Token{},
			&gtsmodel.Client{},
		}
		for _, i := range models {
			if _, err := db.NewCreateTable().Model(i).IfNotExists().Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	}

	down := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			return nil
		})
	}

	if err := Migrations.Register(up, down); err != nil {
		panic(err)
	}
}
