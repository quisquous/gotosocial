/*
	GoToSocial
	Copyright (C) GoToSocial Authors admin@gotosocial.org
	SPDX-License-Identifier: AGPL-3.0-or-later

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

import React from "react";

import { useGetAccountQuery } from "../../../../lib/query/admin";
import FormWithData from "../../../../lib/form/form-with-data";
import FakeProfile from "../../../../components/fake-profile";
import { AdminAccount } from "../../../../lib/types/account";
import { HandleSignup } from "./handlesignup";
import { AccountActions } from "./actions";
import { useParams } from "wouter";

export default function AccountDetail() {
	const params: { accountID: string } = useParams();
	
	return (
		<div className="account-detail">
			<h1>Account Details</h1>
			<FormWithData
				dataQuery={useGetAccountQuery}
				queryArg={params.accountID}
				DataForm={AccountDetailForm}
			/>
		</div>
	);
}

interface AccountDetailFormProps {
	backLocation: string,
	data: AdminAccount,
}

function AccountDetailForm({ data: adminAcct, backLocation }: AccountDetailFormProps) {	
	let yesOrNo = (b: boolean) => {
		return b ? "yes" : "no";
	};

	let created = new Date(adminAcct.created_at).toDateString();
	let lastPosted = "never";
	if (adminAcct.account.last_status_at) {
		lastPosted = new Date(adminAcct.account.last_status_at).toDateString();
	}
	const local = !adminAcct.domain;

	return (
		<>
			<FakeProfile {...adminAcct.account} />
			<h3>General Account Details</h3>
			{ adminAcct.suspended &&
				<div className="info">
					<i className="fa fa-fw fa-info-circle" aria-hidden="true"></i>
					<b>Account is suspended.</b>
				</div>
			}
			<dl className="info-list">
				{ !local &&
				<div className="info-list-entry">
					<dt>Domain</dt>
					<dd>{adminAcct.domain}</dd>
				</div>}
				<div className="info-list-entry">
					<dt>Created</dt>
					<dd><time dateTime={adminAcct.created_at}>{created}</time></dd>
				</div>
				<div className="info-list-entry">
					<dt>Last posted</dt>
					<dd>{lastPosted}</dd>
				</div>
				<div className="info-list-entry">
					<dt>Suspended</dt>
					<dd>{yesOrNo(adminAcct.suspended)}</dd>
				</div>
				<div className="info-list-entry">
					<dt>Silenced</dt>
					<dd>{yesOrNo(adminAcct.silenced)}</dd>
				</div>
				<div className="info-list-entry">
					<dt>Statuses</dt>
					<dd>{adminAcct.account.statuses_count}</dd>
				</div>
				<div className="info-list-entry">
					<dt>Followers</dt>
					<dd>{adminAcct.account.followers_count}</dd>
				</div>
				<div className="info-list-entry">
					<dt>Following</dt>
					<dd>{adminAcct.account.following_count}</dd>
				</div>
			</dl>
			{ local &&
			// Only show local account details
			// if this is a local account!
			<>
				<h3>Local Account Details</h3>
				{ !adminAcct.approved &&
					<div className="info">
						<i className="fa fa-fw fa-info-circle" aria-hidden="true"></i>
						<b>Account is pending.</b>
					</div>
				}
				{ !adminAcct.confirmed && 
					<div className="info">
						<i className="fa fa-fw fa-info-circle" aria-hidden="true"></i>
						<b>Account email not yet confirmed.</b>
					</div>
				}
				<dl className="info-list">
					<div className="info-list-entry">
						<dt>Email</dt>
						<dd>{adminAcct.email} {<b>{adminAcct.confirmed ? "(confirmed)" : "(not confirmed)"}</b> }</dd>
					</div>
					<div className="info-list-entry">
						<dt>Disabled</dt>
						<dd>{yesOrNo(adminAcct.disabled)}</dd>
					</div>
					<div className="info-list-entry">
						<dt>Approved</dt>
						<dd>{yesOrNo(adminAcct.approved)}</dd>
					</div>
					<div className="info-list-entry">
						<dt>Sign-Up Reason</dt>
						<dd>{adminAcct.invite_request ?? <i>none provided</i>}</dd>
					</div>
					{ (adminAcct.ip && adminAcct.ip !== "0.0.0.0") &&
					<div className="info-list-entry">
						<dt>Sign-Up IP</dt>
						<dd>{adminAcct.ip}</dd>
					</div> }
					{ adminAcct.locale &&
					<div className="info-list-entry">
						<dt>Locale</dt>
						<dd>{adminAcct.locale}</dd>
					</div> }
				</dl>
			</> }
			{ local && !adminAcct.approved
				?
				<HandleSignup
					account={adminAcct}
					backLocation={backLocation}
				/>
				:
				<AccountActions account={adminAcct} />
			}
		</>
	);
}
