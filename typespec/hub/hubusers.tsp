import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

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

model HubUserInviteRequest {
    email: EmailAddress;
}

@doc("Enum representing the possible tiers for a hub user")
union HubUserTier {
    FreeTier: "FREE_TIER",
    PaidTier: "PAID_TIER",
}

model OnboardHubUserRequest {
    @doc("The token from the invite email")
    token: string;

    @doc("The full name of the user")
    full_name: string;

    @doc("The country code where the user resides")
    resident_country_code: CountryCode;

    @doc("The password for the user account")
    password: Password;

    @doc("The tier selected by the user during onboarding")
    selected_tier: HubUserTier;

    @doc("Optional preferred language for the user")
    preferred_language?: string;

    @doc("Optional short bio for the user profile")
    short_bio?: string;

    @doc("Optional long bio for the user profile")
    long_bio?: string;
}

model OnboardHubUserResponse {
    @doc("The session token for the newly created user")
    session_token: string;

    @doc("The auto-generated handle for the user")
    generated_handle: string;
}

model CheckHandleAvailabilityRequest {
    @doc("The handle to check for availability")
    handle: Handle;
}

model CheckHandleAvailabilityResponse {
    @doc("Whether the handle is available for use")
    is_available: boolean;

    @doc("Suggested alternative handles if the requested one is not available")
    suggested_alternatives?: string[];
}

model SetHandleRequest {
    @doc("The handle sought by the logged inuser")
    handle: Handle;
}

model SignupHubUserRequest {
    email: EmailAddress;
}

model ChangeEmailAddressRequest {
    email: EmailAddress;
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

@route("/hub/invite-hub-user")
interface HubUserInvite {
    @tag("HubUsers")
    @post
    @useAuth(HubAuth)
    hubUserInvite(@body inviteHubUserRequest: HubUserInviteRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/onboard-user")
interface OnboardHubUser {
    @tag("HubUsers")
    @post
    onboardHubUser(@body onboardHubUserRequest: OnboardHubUserRequest): {
        @statusCode statusCode: 200;
        @body onboardHubUserResponse: OnboardHubUserResponse;
    };
}

@route("/hub/check-handle-availability")
interface CheckHandleAvailability {
    @tag("HubUsers")
    @post
    @useAuth(HubAuth)
    checkHandleAvailability(
        @body checkHandleAvailabilityRequest: CheckHandleAvailabilityRequest,
    ): {
        @statusCode statusCode: 200;
        @body checkHandleAvailabilityResponse: CheckHandleAvailabilityResponse;
    };
}

@route("/hub/set-handle")
interface SetHandle {
    @tag("HubUsers")
    @post
    @useAuth(HubAuth)
    setHandle(@body setHandleRequest: SetHandleRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("The user is not eligible to set a handle (not a paid tier user)")
        @statusCode
        statusCode: 403;
    } | {
        @doc("The requested handle is not available")
        @statusCode
        statusCode: 409;
    };
}

@route("/hub/signup")
interface SignupHubUser {
    @tag("HubUsers")
    @post
    signupHubUser(@body signupHubUserRequest: SignupHubUserRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("The domain of the email address is not supported.")
        @statusCode
        statusCode: 460;
    } | {
        @doc("The user is already a member or has been invited.")
        @statusCode
        statusCode: 461;
    };
}

@route("/hub/change-email-address")
interface ChangeEmailAddress {
    @tag("HubUsers")
    @post
    @useAuth(HubAuth)
    changeEmailAddress(@body changeEmailAddressRequest: ChangeEmailAddressRequest): {
        @statusCode statusCode: 200;
    };
}