import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";
import "../common/candidacies.tsp";
import "../common/interviews.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model FilterCandidacyInfosRequest {
    @doc("If provided, will return only the Candidacies for the given opening. If not provided, all Candidacies of the Employer across all the Openings will be returned.")
    opening_id?: string;

    @doc("If provided, will return only the candidacies for the given recruiter. If not provided, all Candidacies under the given employer will be returned.")
    recruiter_email?: string;

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

model AddInterviewRequest {
    candidacy_id: string;
    start_time: utcDateTime;
    end_time: utcDateTime;
    interview_type: InterviewType;

    @doc("Should contain details like location, meeting URLs, etc.")
    @maxLength(2048)
    description: string;

    // TODO: Perhaps should not allow more than 25 interviews per Candidacy
}

model AddInterviewResponse {
    interview_id: string;
}

model EmployerInterview {
    interview_id: string;
    interview_state: InterviewState;
    start_time: utcDateTime;
    end_time: utcDateTime;
    interview_type: InterviewType;
    description?: string;
    interviewers?: OrgUserShort[];
    interviewers_decision?: InterviewersDecision;

    @maxLength(4096)
    positives?: string;

    @maxLength(4096)
    negatives?: string;

    @maxLength(4096)
    overall_assessment?: string;

    @maxLength(4096)
    feedback_to_candidate?: string;

    feedback_submitted_by?: OrgUserShort;
    feedback_submitted_at?: utcDateTime;
    created_at: utcDateTime;
}

@route("/employer/filter-candidacy-infos")
interface FilterCandidacyInfos {
    @tag("Candidacies")
    @post
    filterCandidacyInfos(@body request: FilterCandidacyInfosRequest): {
        @statusCode statusCode: 200;
        @body response: Candidacy[];
    };
}

@route("/employer/get-candidacy-info")
interface GetEmployerCandidacyInfo {
    @tag("Candidacies")
    @post
    getEmployerCandidacyInfo(@body request: GetCandidacyInfoRequest): {
        @statusCode statusCode: 200;
        @body response: Candidacy;
    };
}

@route("/employer/add-candidacy-comment")
interface AddEmployerCandidacyComment {
    @tag("Candidacies")
    @post
    addComment(@body request: AddEmployerCandidacyCommentRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Comment cannot be added to the candidacy as it has reached a final state")
        @statusCode
        statusCode: 422;
    } | {
        @doc("User is not authorized to add a comment to the candidacy")
        @statusCode
        statusCode: 403;
    };
}

@route("/employer/add-interview")
interface AddInterview {
    @tag("Interviews")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    addInterview(@body request: AddInterviewRequest): {
        @body response: AddInterviewResponse;
        @statusCode statusCode: 200;
    };
}

@route("/employer/get-candidacy-comments")
interface GetCandidacyComments {
    @tag("Candidacies")
    @post
    getEmployerCandidacyComments(@body request: GetCandidacyCommentsRequest): {
        @doc("If an invalid candidacy ID is provided, an empty list is returned.")
        @statusCode
        statusCode: 200;

        @body response: CandidacyComment[];
    };
}
