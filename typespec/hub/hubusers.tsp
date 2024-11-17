import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model LoginRequest {
    email: EmailAddress;
    password: Password;
}

model LoginResponse {
    token: string;
}

model HubTFARequest {
    token: string;
    code: string;
}

model HubTFAResponse {
    session_token: string;
}

@route("/hub/login")
interface Login {
    @tag("HubUsers")
    @post
    login(@body loginRequest: LoginRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/tfa")
interface TFA {
    @tag("HubUsers")
    @post
    tfa(@body hubTFARequest: HubTFARequest): {
        @statusCode statusCode: 200;
        @body hubTFAResponse: HubTFAResponse;
    };
}
