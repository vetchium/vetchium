import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/interviews.tsp";
using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model AddInterviewerRequest {
    interview_id: string;
    org_user_email: string;
}

model RemoveInterviewerRequest {
    interview_id: string;
    org_user_email: string;
}

model GetEmployerInterviewsByOpeningRequest {
    opening_id: string;

    @doc("If provided, will return only the Interviews in the given states. Otherwise all Interviews will be returned.")
    states?: InterviewState[];

    @doc("If provided, will return only the Interviews after the given pagination_key. Otherwise all Interviews will be returned. The interviews will be ordered by their start_time, and then the interview_id. Pass the last interview_id fetched to get the next set of Interviews.")
    pagination_key?: string;

    @doc("Number of Interviews to return. Defaults to 40.")
    @maxValue(100)
    limit?: integer;
}

model GetEmployerInterviewsByCandidacyRequest {
    candidacy_id: string;

    @doc("If provided, will return only the Interviews in the given states. Otherwise all Interviews will be returned.")
    states?: InterviewState[];
}

model GetInterviewDetailsRequest {
    interview_id: string;
}

model PutAssessmentRequest {
    interview_id: string;
    decision?: InterviewersDecision;
    positives?: string;
    negatives?: string;
    overall_assessment?: string;
    feedback_to_candidate?: string;

    @doc("If true, will mark the interview as completed")
    mark_interview_completed?: boolean;
}

model Assessment {
    interview_id: string;
    decision?: InterviewersDecision;
    positives?: string;
    negatives?: string;
    overall_assessment?: string;
    feedback_to_candidate?: string;
    feedback_submitted_by?: string;
    feedback_submitted_at?: string;
}

model GetAssessmentRequest {
    interview_id: string;
}

@route("/employer/rsvp-interview")
interface EmployerRSVPInterview {
    @tag("Interviews")
    @doc("The OrgUser doing this must be an Interviewer in the Interview")
    @post
    @useAuth(EmployerAuth)
    employerRSVPInterview(@body request: RSVPInterviewRequest):
        | {
              @statusCode statusCode: 200;
          }
        | {
              @statusCode statusCode: 400;
              @body error: ValidationErrors;
          }
        | {
              @doc("The Interview is not in the correct state for this action")
              @statusCode
              statusCode: 422;
          }
        | {
              @doc("Either the interview_id or the org_user_email is invalid")
              @statusCode
              statusCode: 404;
          };
}

@route("/employer/add-interviewer")
interface EmployerAddInterviewer {
    @tag("Interviews")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    @useAuth(EmployerAuth)
    employerAddInterviewer(@body request: AddInterviewerRequest):
        | {
              @statusCode statusCode: 200;
          }
        | {
              @doc("The interview_id is invalid")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("The interview is not in the correct state for this action")
              @statusCode
              statusCode: 422;
          }
        | {
              @doc("The user corresponding to org_user_email cannot become an interviewer, possibly because they are no longer a valid employee and disabled")
              @statusCode
              statusCode: 403;
          };
}

@route("/employer/remove-interviewer")
interface EmployerRemoveInterviewer {
    @tag("Interviews")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    @useAuth(EmployerAuth)
    employerRemoveInterviewer(@body request: RemoveInterviewerRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Passed orgUserEmail is not an interviewer")
        @statusCode
        statusCode: 404;
    } | {
        @doc("Interview is not in the correct state for this action")
        @statusCode
        statusCode: 422;
    };
}

#deprecated "This endpoint is deprecated"
@route("/employer/get-interviews-by-opening")
interface GetInterviewsByOpening {
    @tag("Interviews")
    @post
    getInterviewsByOpening(
        @body request: GetEmployerInterviewsByOpeningRequest,
    ): {
        @doc("All Interviews of an Opening are returned, ordered by their start_time")
        @statusCode
        statusCode: 200;

        @body response: EmployerInterview[];
    };
}

@route("/employer/get-interviews-by-candidacy")
interface GetInterviewsByCandidacy {
    @tag("Interviews")
    @post
    getInterviewsByCandidacy(
        @body request: GetEmployerInterviewsByCandidacyRequest,
    ): {
        @doc("All Interviews of a Candidacy are returned, ordered by their start_time")
        @statusCode
        statusCode: 200;

        @body response: EmployerInterview[];
    };
}

@route("/employer/get-interview-details")
interface GetInterviewDetails {
    @tag("Interviews")
    @post
    getInterviewDetails(@body request: GetInterviewDetailsRequest): {
        @statusCode statusCode: 200;
        @body response: EmployerInterview;
    };
}

@route("/employer/put-assessment")
interface PutAssesment {
    @doc("Any old feedback will be overwritten")
    @tag("Interviews")
    @post
    putAssesment(@body request: PutAssessmentRequest):
        | {
              @doc("Assessment updated successfully")
              @statusCode
              statusCode: 200;
          }
        | {
              @doc("Interview not found")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("The OrgUser is not an interviewer on this interview. As of now even Admin is forbidden (intentionally) from updating the assessment. This may change in the future.")
              @statusCode
              statusCode: 403;
          }
        | {
              @doc("Interview is not in the correct state for this action")
              @statusCode
              statusCode: 422;
          };
}

#deprecated "This endpoint is deprecated"
@route("/employer/get-assessment")
interface GetAssessment {
    @tag("Interviews")
    @post
    getAssessment(@body request: GetAssessmentRequest): {
        @statusCode statusCode: 200;
        @body response: Assessment;
    };
}
