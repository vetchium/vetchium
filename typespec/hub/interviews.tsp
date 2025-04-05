import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/interviews.tsp";
using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model GetHubInterviewsByCandidacyRequest {
  @doc("The Candidacy ID to get Interviews for")
  candidacy_id: string;

  @doc("If provided, will return only the Interviews in the given states. Otherwise all Interviews will be returned.")
  states?: InterviewState[];
}

model HubInterviewer {
  name: string;
  rsvp_status: RSVPStatus;
}

model HubInterview {
  @doc("The Interview ID")
  interview_id: string;

  interview_state: InterviewState;
  start_time: utcDateTime;
  end_time: utcDateTime;
  interview_type: InterviewType;
  description?: string;
  candidate_rsvp_status: RSVPStatus;
  interviewers?: HubInterviewer[];
}

model HubRSVPInterviewRequest {
  interview_id: string;
  rsvp_status: RSVPStatus;
}

@route("/hub/get-interviews-by-candidacy")
interface GetHubInterviewsByCandidacy {
  @tag("Interviews")
  @doc("Get Interviews by Candidacy")
  @post
  @useAuth(HubAuth)
  getHubInterviewsByCandidacy(
    @body request: GetHubInterviewsByCandidacyRequest,
  ): {
    @statusCode statusCode: 200;
    @body response: HubInterview[];
  } | {
    @statusCode statusCode: 400;
    @body error: ValidationErrors;
  };
}

@route("/hub/rsvp-interview")
interface HubRSVPInterview {
  @tag("Interviews")
  @doc("The OrgUser doing this must be an Interviewer in the Interview")
  @post
  @useAuth(HubAuth)
  hubRSVPInterview(@body request: RSVPInterviewRequest):
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
        @doc("The interview_id is invalid")
        @statusCode
        statusCode: 404;
      };
}
