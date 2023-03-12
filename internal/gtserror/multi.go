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

package gtserror

import (
	"errors"
	"fmt"
	"strings"
)

// MultiError allows encapsulating multiple errors under a singular instance,
// which is useful when you only want to log on errors, not return early / bubble up.
type MultiError []string

func (e *MultiError) Append(err error) {
	*e = append(*e, err.Error())
}

func (e *MultiError) Appendf(format string, args ...any) {
	*e = append(*e, fmt.Sprintf(format, args...))
}

// Combine converts this multiError to a singular error instance, returning nil if empty.
func (e MultiError) Combine() error {
	if len(e) == 0 {
		return nil
	}
	return errors.New(`"` + strings.Join(e, `","`) + `"`)
}
