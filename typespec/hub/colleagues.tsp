import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model HubUserShort {
    handle: Handle;
    name: string;
    short_bio: string;
}

model ConnectColleagueRequest {
    handle: Handle;
}

model UnlinkColleagueRequest {
    handle: Handle;
}

model MyColleagueApprovalsRequest {
    pagination_key?: string;

    @minValue(1)
    @maxValue(100)
    limit?: numeric;
}

model MyColleagueApprovals {
    approvals: HubUserShort[];
    pagination_key?: string;
}

model MyColleagueSeeksRequest {
    pagination_key?: string;

    @minValue(1)
    @maxValue(100)
    limit?: numeric;
}

model MyColleagueSeeks {
    seeks: HubUserShort[];
    pagination_key?: string;
}

model ApproveColleagueRequest {
    handle: Handle;
}

model RejectColleagueRequest {
    handle: Handle;
}

model FilterColleaguesRequest {
    @doc("The prefix to match either the name or the handle of the colleagues of the logged in user")
    @minLength(1)
    @maxLength(12)
    prefix: string;

    @doc("The limit on the number of colleagues to return")
    @minValue(1)
    @maxValue(6)
    limit: numeric;
}

union EndorsementState {
    SoughtEndorsement: "SOUGHT_ENDORSEMENT",
    Endorsed: "ENDORSED",
    DeclinedEndorsement: "DECLINED_ENDORSEMENT",
    // TODO: We need a new state to migrate to, when an application is rejected/withdrawn/shortlisted but the Endorser has still not responded. Open Question: Should we hold off showing the Application to the Employer until the Endorser responds ?
}

model MyEndorseApprovalsRequest {
    pagination_key?: string;

    @doc("Filter and return only Endorsements in this state. By default only Endorsements in the SoughtEndorsement state are returned.")
    @maxItems(3)
    state: EndorsementState[];

    @doc("The limit on the number of endorsements to return")
    @minValue(1)
    @maxValue(40)
    limit?: numeric;
}

model MyEndorseApproval {
    application_id: string;
    applicant_handle: Handle;
    applicant_name: string;
    applicant_short_bio: string;
    employer_name: string;
    employer_domain: string;
    opening_title: string;
    opening_url: string;
    application_status: string;
    application_created_at: utcDateTime;
    endorsement_status: EndorsementState;
}

model MyEndorseApprovalsResponse {
    endorsements: MyEndorseApproval[];
    pagination_key?: string;
}

model EndorseApplicationRequest {
    application_id: string;
}

model RejectEndorsementRequest {
    application_id: string;
}

@route("/hub/connect-colleague")
interface ConnectColleague {
    @tag("HubProfile")
    @doc("Request the handle to connect as a colleague with the logged in user. The invitation will be notified to the handle. The invitation will be valid for 60 days. The handle may either accept or reject the invitation before the 60 days. If the handle rejects the invitation, the logged in user will not be able to send an invitation again to the same handle.")
    @post
    @useAuth(HubAuth)
    connectColleague(@body request: ConnectColleagueRequest): {};
}

@route("/hub/unlink-colleague")
interface UnlinkColleague {
    @tag("HubProfile")
    @doc("Remove the colleague relationship of the handle with the logged in user. The user who has removed the colleague relationship can only seek for a new connection again. The other party (who was removed) will not be able to send an invitation again to the one who removed the relationship.")
    @post
    @useAuth(HubAuth)
    unlinkColleague(@body request: UnlinkColleagueRequest): {};
}

@route("/hub/my-colleague-approvals")
interface MyColleagueApprovalsAPI {
    @tag("HubProfile")
    @doc("Get the list of handles of all likely colleagues that have sought to connect with the logged in user, in the last 60 days.")
    @post
    @useAuth(HubAuth)
    myColleagueApprovals(@body request: MyColleagueApprovalsRequest): {
        @statusCode statusCode: 200;
        @body response: MyColleagueApprovals;
    };
}

// TODO: In future we can add a state filter to the request to get only seeks with a certain state
@route("/hub/my-colleague-seeks")
interface MyColleagueSeeksAPI {
    @tag("HubProfile")
    @doc("Get the list of all handles to whom the logged in user has sought a colleague connection, and not yet approved or rejected by the requested handle.")
    @post
    @useAuth(HubAuth)
    myColleagueSeeks(@body request: MyColleagueSeeksRequest): {
        @statusCode statusCode: 200;
        @body response: MyColleagueSeeks;
    };
}

@route("/hub/approve-colleague")
interface ApproveColleague {
    @tag("HubProfile")
    @doc("Logged in user approves a request from the handle to connect as a colleague")
    @post
    @useAuth(HubAuth)
    approveColleague(@body request: ApproveColleagueRequest): {};
}

@route("/hub/reject-colleague")
interface RejectColleague {
    @tag("HubProfile")
    @doc("Logged in user rejects a request from the handle to connect as a colleague")
    @post
    @useAuth(HubAuth)
    rejectColleague(@body request: RejectColleagueRequest): {};
}

@route("/hub/filter-colleagues")
interface FilterColleagues {
    @tag("HubProfile")
    @doc("Filter the list of colleagues of the logged in user")
    @post
    @useAuth(HubAuth)
    filterColleagues(@body request: FilterColleaguesRequest): {
        @statusCode statusCode: 200;
        @body response: HubUserShort[];
    };
}

@route("/hub/my-endorse-approvals")
interface MyEndorseApprovals {
    @tag("HubProfile")
    @doc("Get the list of endorsement requests that the logged in user can either approve or reject. The list is paginated. The items are ordered primarily by the timestamp of the endorsement request, and then by the id of the endorsement request for items with the same timestamp.")
    @post
    @useAuth(HubAuth)
    myEndorseApprovals(@body request: MyEndorseApprovalsRequest): {
        @statusCode statusCode: 200;
        @body response: MyEndorseApprovalsResponse;
    };
}

@route("/hub/endorse-application")
interface EndorseApplication {
    @tag("HubProfile")
    @doc("Endorse an application for a colleague")
    @post
    @useAuth(HubAuth)
    endorseApplication(@body request: EndorseApplicationRequest): {};
}

@route("/hub/reject-endorsement")
interface RejectEndorsement {
    @tag("HubProfile")
    @doc("Reject an endorsement for a colleague")
    @post
    @useAuth(HubAuth)
    rejectEndorsement(@body request: RejectEndorsementRequest): {};
}
