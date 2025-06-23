import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

@doc("Request to score a single application's resume")
model ApplicationSortRequest {
    @doc("The unique identifier for the application")
    application_id: string;

    @doc("S3 path to the resume file in format s3://bucket/key")
    resume_path: string;
}

@doc("Request to score multiple resumes against a job description in a batch")
model SortingHatRequest {
    @doc("The job description to score resumes against")
    job_description: string;

    @doc("List of applications to score")
    application_sort_requests: ApplicationSortRequest[];
}

@doc("Scores for a single application from all models")
model SortingHatScore {
    @doc("The application ID this score relates to")
    application_id: string;

    @doc("Scores from different models")
    model_scores: ModelScore[];
}

@doc("Response containing scores for all applications in the batch")
model SortingHatResponse {
    @doc("List of application scores")
    scores: SortingHatScore[];
}

@route("/score-batch")
interface ScoreBatch {
    @post
    @tag("SortingHat")
    @doc("Score multiple resumes against a job description in a single batch")
    scoreBatch(@body request: SortingHatRequest): {
        @statusCode statusCode: 200;
        @body response: SortingHatResponse;
    } | {
        @doc("Bad request - invalid input data")
        @statusCode
        statusCode: 400;
    } | {
        @doc("Internal server error")
        @statusCode
        statusCode: 500;
    };
}
