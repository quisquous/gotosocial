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

package webfinger_test

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/superseriousbusiness/gotosocial/internal/ap"
	"github.com/superseriousbusiness/gotosocial/internal/api/wellknown/webfinger"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/email"
	"github.com/superseriousbusiness/gotosocial/internal/federation"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/media"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
	"github.com/superseriousbusiness/gotosocial/internal/processing"
	"github.com/superseriousbusiness/gotosocial/internal/state"
	"github.com/superseriousbusiness/gotosocial/internal/storage"
	"github.com/superseriousbusiness/gotosocial/internal/typeutils"
	"github.com/superseriousbusiness/gotosocial/testrig"
)

type WebfingerStandardTestSuite struct {
	// standard suite interfaces
	suite.Suite
	db           db.DB
	state        state.State
	tc           typeutils.TypeConverter
	mediaManager media.Manager
	federator    federation.Federator
	emailSender  email.Sender
	processor    *processing.Processor
	storage      *storage.Driver
	oauthServer  oauth.Server

	// standard suite models
	testTokens       map[string]*gtsmodel.Token
	testClients      map[string]*gtsmodel.Client
	testApplications map[string]*gtsmodel.Application
	testUsers        map[string]*gtsmodel.User
	testAccounts     map[string]*gtsmodel.Account
	testAttachments  map[string]*gtsmodel.MediaAttachment
	testStatuses     map[string]*gtsmodel.Status

	// module being tested
	webfingerModule *webfinger.Module
}

func (suite *WebfingerStandardTestSuite) SetupSuite() {
	suite.testTokens = testrig.NewTestTokens()
	suite.testClients = testrig.NewTestClients()
	suite.testApplications = testrig.NewTestApplications()
	suite.testUsers = testrig.NewTestUsers()
	suite.testAccounts = testrig.NewTestAccounts()
	suite.testAttachments = testrig.NewTestAttachments()
	suite.testStatuses = testrig.NewTestStatuses()
}

func (suite *WebfingerStandardTestSuite) SetupTest() {
	suite.state.Caches.Init()
	testrig.StartWorkers(&suite.state)

	testrig.InitTestLog()
	testrig.InitTestConfig()

	suite.db = testrig.NewTestDB(&suite.state)
	suite.state.DB = suite.db
	suite.tc = testrig.NewTestTypeConverter(suite.db)
	suite.storage = testrig.NewInMemoryStorage()
	suite.state.Storage = suite.storage
	suite.mediaManager = testrig.NewTestMediaManager(&suite.state)
	suite.federator = testrig.NewTestFederator(&suite.state, testrig.NewTestTransportController(&suite.state, testrig.NewMockHTTPClient(nil, "../../../../testrig/media")), suite.mediaManager)
	suite.emailSender = testrig.NewEmailSender("../../../../web/template/", nil)
	suite.processor = testrig.NewTestProcessor(&suite.state, suite.federator, suite.emailSender, suite.mediaManager)
	suite.webfingerModule = webfinger.New(suite.processor)
	suite.oauthServer = testrig.NewTestOauthServer(suite.db)
	testrig.StandardDBSetup(suite.db, suite.testAccounts)
	testrig.StandardStorageSetup(suite.storage, "../../../../testrig/media")

	suite.NoError(suite.processor.Start())
}

func (suite *WebfingerStandardTestSuite) TearDownTest() {
	testrig.StandardDBTeardown(suite.db)
	testrig.StandardStorageTeardown(suite.storage)
	testrig.StopWorkers(&suite.state)
}

func accountDomainAccount() *gtsmodel.Account {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := &privateKey.PublicKey

	acct := &gtsmodel.Account{
		ID:                    "01FG1K8EA7SYHEC7V6XKVNC4ZA",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		Username:              "aaaaa",
		Domain:                "",
		Privacy:               gtsmodel.VisibilityDefault,
		Language:              "en",
		URI:                   "http://gts.example.org/users/aaaaa",
		URL:                   "http://gts.example.org/@aaaaa",
		InboxURI:              "http://gts.example.org/users/aaaaa/inbox",
		OutboxURI:             "http://gts.example.org/users/aaaaa/outbox",
		FollowingURI:          "http://gts.example.org/users/aaaaa/following",
		FollowersURI:          "http://gts.example.org/users/aaaaa/followers",
		FeaturedCollectionURI: "http://gts.example.org/users/aaaaa/collections/featured",
		ActorType:             ap.ActorPerson,
		PrivateKey:            privateKey,
		PublicKey:             publicKey,
		PublicKeyURI:          "http://gts.example.org/users/aaaaa/main-key",
	}

	return acct
}
