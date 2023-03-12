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

package accounts

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form/v4"
	apimodel "github.com/superseriousbusiness/gotosocial/internal/api/model"
	apiutil "github.com/superseriousbusiness/gotosocial/internal/api/util"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
)

// AccountUpdateCredentialsPATCHHandler swagger:operation PATCH /api/v1/accounts/update_credentials accountUpdate
//
// Update your account.
//
//	---
//	tags:
//	- accounts
//
//	consumes:
//	- multipart/form-data
//
//	produces:
//	- application/json
//
//	parameters:
//	-
//		name: discoverable
//		in: formData
//		description: Account should be made discoverable and shown in the profile directory (if enabled).
//		type: boolean
//	-
//		name: bot
//		in: formData
//		description: Account is flagged as a bot.
//		type: boolean
//	-
//		name: display_name
//		in: formData
//		description: The display name to use for the account.
//		type: string
//		allowEmptyValue: true
//	-
//		name: note
//		in: formData
//		description: Bio/description of this account.
//		type: string
//		allowEmptyValue: true
//	-
//		name: avatar
//		in: formData
//		description: Avatar of the user.
//		type: file
//	-
//		name: header
//		in: formData
//		description: Header of the user.
//		type: file
//	-
//		name: locked
//		in: formData
//		description: Require manual approval of follow requests.
//		type: boolean
//	-
//		name: source[privacy]
//		in: formData
//		description: Default post privacy for authored statuses.
//		type: string
//	-
//		name: source[sensitive]
//		in: formData
//		description: Mark authored statuses as sensitive by default.
//		type: boolean
//	-
//		name: source[language]
//		in: formData
//		description: Default language to use for authored statuses (ISO 6391).
//		type: string
//	-
//		name: source[status_content_type]
//		in: formData
//		description: Default content type to use for authored statuses (text/plain or text/markdown).
//		type: string
//	-
//		name: custom_css
//		in: formData
//		description: >-
//			Custom CSS to use when rendering this account's profile or statuses.
//			String must be no more than 5,000 characters (~5kb).
//		type: string
//	-
//		name: enable_rss
//		in: formData
//		description: Enable RSS feed for this account's Public posts at `/[username]/feed.rss`
//		type: boolean
//	-
//		name: fields_attributes
//		in: formData
//		description: Profile fields to be added to this account's profile
//		type: array
//		items:
//			type: object
//
//	security:
//	- OAuth2 Bearer:
//		- write:accounts
//
//	responses:
//		'200':
//			description: "The newly updated account."
//			schema:
//				"$ref": "#/definitions/account"
//		'400':
//			description: bad request
//		'401':
//			description: unauthorized
//		'404':
//			description: not found
//		'406':
//			description: not acceptable
//		'500':
//			description: internal server error
func (m *Module) AccountUpdateCredentialsPATCHHandler(c *gin.Context) {
	authed, err := oauth.Authed(c, true, true, true, true)
	if err != nil {
		apiutil.ErrorHandler(c, gtserror.NewErrorUnauthorized(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	if _, err := apiutil.NegotiateAccept(c, apiutil.JSONAcceptHeaders...); err != nil {
		apiutil.ErrorHandler(c, gtserror.NewErrorNotAcceptable(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	form, err := parseUpdateAccountForm(c)
	if err != nil {
		apiutil.ErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	acctSensitive, errWithCode := m.processor.Account().Update(c.Request.Context(), authed.Account, form)
	if errWithCode != nil {
		apiutil.ErrorHandler(c, errWithCode, m.processor.InstanceGetV1)
		return
	}

	c.JSON(http.StatusOK, acctSensitive)
}

type fieldAttributesBinding struct{}

func (fieldAttributesBinding) Name() string {
	return "FieldAttributes"
}

func (fieldAttributesBinding) Bind(req *http.Request, obj any) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	decoder := form.NewDecoder()
	// change default namespace prefix and suffix to allow correct parsing of the field attributes
	decoder.SetNamespacePrefix("[")
	decoder.SetNamespaceSuffix("]")
	if err := decoder.Decode(obj, req.Form); err != nil {
		return err
	}

	return nil
}

func parseUpdateAccountForm(c *gin.Context) (*apimodel.UpdateCredentialsRequest, error) {
	form := &apimodel.UpdateCredentialsRequest{
		Source: &apimodel.UpdateSource{},
	}

	if err := c.ShouldBind(&form); err != nil {
		return nil, fmt.Errorf("could not parse form from request: %s", err)
	}

	// use custom form binding to support field attributes in the form data
	if err := c.ShouldBindWith(&form, fieldAttributesBinding{}); err != nil {
		return nil, fmt.Errorf("could not parse form from request: %s", err)
	}

	// parse source field-by-field
	sourceMap := c.PostFormMap("source")

	if privacy, ok := sourceMap["privacy"]; ok {
		form.Source.Privacy = &privacy
	}

	if sensitive, ok := sourceMap["sensitive"]; ok {
		sensitiveBool, err := strconv.ParseBool(sensitive)
		if err != nil {
			return nil, fmt.Errorf("error parsing form source[sensitive]: %s", err)
		}
		form.Source.Sensitive = &sensitiveBool
	}

	if language, ok := sourceMap["language"]; ok {
		form.Source.Language = &language
	}

	if statusContentType, ok := sourceMap["status_content_type"]; ok {
		form.Source.StatusContentType = &statusContentType
	}

	if form == nil ||
		(form.Discoverable == nil &&
			form.Bot == nil &&
			form.DisplayName == nil &&
			form.Note == nil &&
			form.Avatar == nil &&
			form.Header == nil &&
			form.Locked == nil &&
			form.Source.Privacy == nil &&
			form.Source.Sensitive == nil &&
			form.Source.Language == nil &&
			form.Source.StatusContentType == nil &&
			form.FieldsAttributes == nil &&
			form.CustomCSS == nil &&
			form.EnableRSS == nil) {
		return nil, errors.New("empty form submitted")
	}

	return form, nil
}
