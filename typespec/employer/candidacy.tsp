import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";
import "../common/interviews.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model AddEmployerCandidacyCommentRequest {
    candidacy_id: string;

    @maxLength(2048)
    comment: string;
}

@route("/employer/add-candidacy-comment")
interface AddEmployerCandidacyComment {
    @post
    addComment(@body request: AddEmployerCandidacyCommentRequest): {
        @statusCode statusCode: 200;
    };
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

@route("/employer/add-interview")
interface AddInterview {
    @doc("Need one of $OrgUserRole.Admin or $OrgUserRole.ApplicationsCRUD")
    @post
    addInterview(@body request: AddInterviewRequest): {
        @statusCode statusCode: 200;
    };
}

model Interview {
    interview_id: string;
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

    state: InterviewState;
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

@route("/employer/get-interviews")
interface GetInterviews {
    @post
    getInterviews(@body request: GetInterviewsRequest): {
        @statusCode statusCode: 200;
        @body response: Interview[];
    };
}
