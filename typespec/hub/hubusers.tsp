import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

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
    tfa_token: string;
    tfa_code: string;

    @doc("If true, the session will be remembered for a longer duration")
    remember_me: boolean;
}

model HubTFAResponse {
    session_token: string;
}

model ChangePasswordRequest {
    old_password: Password;
    new_password: Password;
}

model ForgotPasswordRequest {
    email: EmailAddress;
}

model ForgotPasswordResponse {
    token: string;
}

model ResetPasswordRequest {
    token: string;
    password: Password;
}

model GetMyHandleResponse {
    handle: string;
}

union HubUserState {
    ActiveHubUserState: "ACTIVE_HUB_USER",
}

@route("/hub/login")
interface Login {
    @tag("HubUsers")
    @post
    login(@body loginRequest: LoginRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("The User account is not in a valid state to login")
        @statusCode
        statusCode: 422;
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

@route("/hub/logout")
interface Logout {
    @tag("HubUsers")
    @post
    @useAuth(HubAuth)
    logout(): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/forgot-password")
interface ForgotPassword {
    @tag("HubUsers")
    @post
    forgotPassword(@body forgotPasswordRequest: ForgotPasswordRequest): {
        @statusCode statusCode: 200;
        @body forgotPasswordResponse: ForgotPasswordResponse;
    };
}

@route("/hub/change-password")
interface ChangePassword {
    @tag("HubUsers")
    @post
    @useAuth(HubAuth)
    changePassword(@body changePasswordRequest: ChangePasswordRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/reset-password")
interface ResetPassword {
    @tag("HubUsers")
    @post
    resetPassword(@body resetPasswordRequest: ResetPasswordRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/get-my-handle")
interface GetMyHandle {
    @tag("HubUsers")
    @get
    @useAuth(HubAuth)
    getMyHandle(): {
        @statusCode statusCode: 200;
        @body getMyHandleResponse: GetMyHandleResponse;
    };
}
