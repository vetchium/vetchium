import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model GetOnboardStatusRequest {
    client_id: string;
}

enum OnboardStatus {
    DomainNotVerified: "DOMAIN_NOT_VERIFIED",
    DomainVerifiedOnboardPending: "DOMAIN_VERIFIED_ONBOARD_PENDING",
    DomainOnboarded: "DOMAIN_ONBOARDED",
}

model GetOnboardStatusResponse {
    status: OnboardStatus;
}

model SetOnboardPasswordRequest {
    client_id: string;
    password: Password;
    token: string;
}

model EmployerSigninRequest {
    client_id: string;
    email: EmailAddress;
    password: Password;
}

model EmployerSigninResponse {
    token: string;
}

model EmployerTFARequest {
    tfa_code: string;
    tfa_token: string;
    remember_me?: boolean;
}

model EmployerTFAResponse {
    session_token: string;
}

@route("/employer/get-onboard-status")
interface GetOnboardStatus {
    @tag("Employer Auth")
    @post
    GetOnboardStatus(
        @body request: GetOnboardStatusRequest,
    ): GetOnboardStatusResponse;
}

@route("/employer/set-onboard-password")
interface SetOnboardPassword {
    @tag("Employer Auth")
    @post
    SetOnboardPassword(@body request: SetOnboardPasswordRequest): void;
}

@route("/employer/signin")
interface EmployerSignin {
    @tag("Employer Auth")
    @post
    EmployerSignin(
        @body request: EmployerSigninRequest,
    ): EmployerSigninResponse;
}

@route("/employer/tfa")
interface EmployerTFA {
    @tag("Employer Auth")
    @post
    EmployerTFA(@body request: EmployerTFARequest): EmployerTFAResponse;
}
