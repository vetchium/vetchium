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
    tfa_token: string;
    tfa_code: string;

    @doc("If true, the session will be remembered for a longer duration")
    remember_me: boolean;
}

model HubTFAResponse {
    session_token: string;
}

model InviteUserRequest {
    email: EmailAddress;
}

model ForgotPasswordRequest {
    email: EmailAddress;
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

@route("/hub/logout")
interface Logout {
    @tag("HubUsers")
    @post
    @useAuth(BearerAuth)
    logout(): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/invite-user")
interface InviteUser {
    @tag("HubUsers")
    @post
    @useAuth(BearerAuth)
    inviteUser(@body inviteUserRequest: InviteUserRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/forgot-password")
interface ForgotPassword {
    @tag("HubUsers")
    @post
    forgotPassword(@body forgotPasswordRequest: ForgotPasswordRequest): {
        @statusCode statusCode: 200;
    };
}

model ChangePasswordRequest {
    old_password: Password;
    new_password: Password;
}

@route("/hub/change-password")
interface ChangePassword {
    @tag("HubUsers")
    @post
    @useAuth(BearerAuth)
    changePassword(@body changePasswordRequest: ChangePasswordRequest): {
        @statusCode statusCode: 200;
    };
}

model ResetPasswordRequest {
    token: string;
    password: Password;
}

@route("/hub/reset-password")
interface ResetPassword {
    @tag("HubUsers")
    @post
    resetPassword(@body resetPasswordRequest: ResetPasswordRequest): {
        @statusCode statusCode: 200;
    };
}

model GetMyHandleResponse {
    handle: string;
}

@route("/hub/get-my-handle")
interface GetMyHandle {
    @tag("HubUsers")
    @get
    @useAuth(BearerAuth)
    getMyHandle(): {
        @statusCode statusCode: 200;
        @body getMyHandleResponse: GetMyHandleResponse;
    };
}
