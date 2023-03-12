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

package prune

import (
	"context"
	"fmt"

	"github.com/superseriousbusiness/gotosocial/cmd/gotosocial/action"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/log"
)

// Orphaned prunes orphaned media from storage.
var Orphaned action.GTSAction = func(ctx context.Context) error {
	prune, err := setupPrune(ctx)
	if err != nil {
		return err
	}

	dry := config.GetAdminMediaPruneDryRun()

	pruned, err := prune.manager.PruneOrphaned(ctx, dry)
	if err != nil {
		return fmt.Errorf("error pruning: %s", err)
	}

	if dry /* dick heyyoooooo */ {
		log.Infof(ctx, "DRY RUN: %d items are orphaned and eligible to be pruned", pruned)
	} else {
		log.Infof(ctx, "%d orphaned items were pruned", pruned)
	}

	return prune.shutdown(ctx)
}
