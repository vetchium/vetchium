import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";
import "../common/interviews.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model GetCandidaciesInfoRequest {
    @doc("If provided, will return only the candidacies for the given recruiter. If not provided, all Candidacies under the given employer will be returned.")
    recruiter_id?: string;

    @doc("If provided, will return only the Candidacies in the given state. Otherwise all Candidacies will be returned.")
    state?: CandidacyState;

    @doc("If provided, will return only the Candidacies after the given pagination_key. Otherwise all Candidacies will be returned. The Candidacies will be returned in reverse chronological order. If two Candidacies are created at the same time, the Candidacies will be further ordered by the candidacy_id. Pass the last candidacy_id fetched to get the next set of Candidacies.")
    pagination_key?: string;

    @doc("Number of Candidacies to return. Defaults to 40.")
    @maxValue(40)
    limit?: integer;
}

model Candidacy {
    candidacy_id: string;
    opening_id: string;
    opening_title: string;
    opening_description: string;
    candidacy_state: CandidacyState;
    applicant_name: string;
    applicant_handle: string;
}

model AddEmployerCandidacyCommentRequest {
    candidacy_id: string;

    @maxLength(2048)
    comment: string;
}

union InterviewType {
    InPerson: "IN_PERSON",
    VideoCall: "VIDEO_CALL",
    TakeHome: "TAKE_HOME",
    Unspecified: "UNSPECIFIED",
}

model AddInterviewRequest {
    candidacy_id: string;
    start_time: utcDateTime;
    end_time: utcDateTime;
    interview_type: InterviewType;

    @doc("Should contain details like location, meeting URLs, etc.")
    @maxLength(2048)
    description: string;

    @doc("List of email addresses of the OrgUsers who will be interviewing the candidate")
    interviewers?: EmailAddress[];
}

model Interview {
    interview_id: string;
    interview_state: InterviewState;
    start_time: utcDateTime;
    end_time: utcDateTime;
    interview_type: InterviewType;
    description?: string;
    interviewers?: EmailAddress[];
    interviewers_decision?: InterviewersDecision;

    @maxLength(4096)
    interviewers_assessment?: string;

    @maxLength(4096)
    feedback_to_candidate?: string;

    created_at: utcDateTime;
}

model GetInterviewsRequest {
    @doc("If provided, will return only the interviews for the given candidacy. If not provided, all Interviews under the given employer will be returned.")
    candidacy_id?: string;

    @doc("If provided, will return only the interviews for the given opening. If not provided, all Interviews under the given employer will be returned. If a mismatching opening_id and candidacy_id are provided, the results are unpredictable. Use only one of the two parameters at a time.")
    opening_id?: string;

    @doc("Results will be ordered based on the created_at timestamp. If more than one interview is created at the sametime, then the interview_id is used for further ordering. This should match the last interview_id fetched. Interviews that are after the value of this parameter will be returned.")
    pagination_key?: string;

    @doc("Number of interviews to return. Defaults to 40.")
    @maxValue(100)
    limit?: integer;
}

@route("/employer/get-candidacies-info")
interface GetCandidaciesInfo {
    @post
    getCandidaciesInfo(@body request: GetCandidaciesInfoRequest): {
        @statusCode statusCode: 200;
        @body response: Candidacy[];
    };
}

@route("/employer/add-candidacy-comment")
interface AddEmployerCandidacyComment {
    @post
    addComment(@body request: AddEmployerCandidacyCommentRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/add-interview")
interface AddInterview {
    @doc("Need one of $OrgUserRole.Admin or $OrgUserRole.ApplicationsCRUD")
    @post
    addInterview(@body request: AddInterviewRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/get-interviews")
interface GetInterviews {
    @post
    getInterviews(@body request: GetInterviewsRequest): {
        @statusCode statusCode: 200;
        @body response: Interview[];
    };
}
