import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model AddOfficialEmailRequest {
    email: EmailAddress;
}

model VerifyOfficialEmailRequest {
    email: EmailAddress;

    @doc("The verification code sent to the email address")
    code: string;
}

model TriggerVerificationRequest {
    email: EmailAddress;
}

model DeleteOfficialEmailRequest {
    email: EmailAddress;
}

model OfficialEmail {
    email: EmailAddress;

    @doc("The date and time when the email address was last verified. If null, the email address was never verified.")
    last_verified_at?: utcDateTime;

    @doc("Whether a verification is in progress for the email address.")
    verify_in_progress: boolean;
}

model GetBioRequest {
    handle: string;
}

@doc("Represents the state of a colleague connection from the perspective of the logged-in user")
union ColleagueConnectionState {
    @doc("No connection exists and the user can be sent a colleague request")
    CAN_SEND_REQUEST: {},

    @doc("No connection exists and the user cannot be sent a colleague request (e.g., no common verified domains)")
    CANNOT_SEND_REQUEST: {},

    @doc("The logged-in user has sent a request that is pending")
    REQUEST_SENT_PENDING: {},

    @doc("The target user has sent a request that is pending the logged-in user's action")
    REQUEST_RECEIVED_PENDING: {},

    @doc("The users are connected as colleagues")
    CONNECTED: {},

    @doc("The logged-in user previously rejected a request from this user")
    REJECTED_BY_ME: {},

    @doc("This user previously rejected a request from the logged-in user")
    REJECTED_BY_THEM: {},

    @doc("The logged-in user previously unlinked a connection with this user")
    UNLINKED_BY_ME: {},

    @doc("This user previously unlinked a connection with the logged-in user")
    UNLINKED_BY_THEM: {},
}

model Bio {
    handle: string;
    full_name: string;

    @doc("The short bio of the user. The one that gets displayed below the user's name.")
    short_bio: string;

    @doc("The long bio of the user. Something like 'About Me' of a user")
    long_bio: string;

    @doc("The list of verified mail domains for the user.")
    verified_mail_domains?: string[];

    @doc("The state of colleague connection with this user from the perspective of the logged-in user")
    colleague_connection_state: ColleagueConnectionState;
}

model UpdateBioRequest {
    @doc("The full name for the logged in user. If not provided, the full name will not be updated.")
    @minLength(1)
    @maxLength(64)
    full_name?: string;

    @doc("The short bio for the logged in user. If not provided, the short bio will not be updated.")
    @minLength(1)
    @maxLength(64)
    short_bio?: string;

    @doc("The long bio for the logged in user. If not provided, the long bio will not be updated.")
    @minLength(1)
    @maxLength(1024)
    long_bio?: string;
}

model UploadProfilePictureRequest {
    @doc("Profile image in JPEG, PNG, or WEBP format (Max: 5MB size, Dimensions: 200x200 to 2048x2048).")
    image: bytes;
}

@route("/hub/add-official-email")
interface AddOfficialEmail {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    addOfficialEmail(@body request: AddOfficialEmailRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("The user has reached the maximum (50) allowed official emails per user.")
        @statusCode
        statusCode: 422;
    } | {
        @doc("The email address is already associated with an account")
        @statusCode
        statusCode: 409;
    };
}

@route("/hub/verify-official-email")
interface VerifyOfficialEmail {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    verifyOfficialEmail(@body request: VerifyOfficialEmailRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Email address could not be verified, possibly because the verification code is invalid or expired")
        @statusCode
        statusCode: 422;
    };
}

@route("/hub/trigger-verification")
interface TriggerVerification {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    triggerVerification(@body request: TriggerVerificationRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Email address verification not triggered, because its prior verification is still valid")
        @statusCode
        statusCode: 422;
    };
}

@route("/hub/delete-official-email")
interface DeleteOfficialEmail {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    deleteOfficialEmail(@body request: DeleteOfficialEmailRequest): {};
}

@route("/hub/my-official-emails")
interface MyOfficialEmails {
    @tag("HubProfile")
    @get
    @useAuth(HubAuth)
    myOfficialEmails(): {
        @statusCode statusCode: 200;
        @body emails: OfficialEmail[];
    };
}

@route("/hub/get-bio")
interface GetBio {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    getBio(@body request: GetBioRequest): {
        @body bio: Bio;
        @statusCode statusCode: 200;
    };
}

@route("/hub/update-bio")
interface UpdateBio {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    updateBio(@body request: UpdateBioRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/upload-profile-picture")
interface UploadProfilePicture {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    uploadProfilePicture(@bodyRoot formData: UploadProfilePictureRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/profile-picture/{handle}")
interface GetProfilePicture {
    @tag("HubProfile")
    @get
    @useAuth(HubAuth)
    @doc("Get a user's profile picture")
    getProfilePicture(@path handle: string): {
        @statusCode statusCode: 200;
        @body image: bytes;
    } | {
        @doc("User not found or user has no profile picture")
        @statusCode
        statusCode: 404;
    };
}

@route("/hub/remove-profile-picture")
interface RemoveProfilePicture {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    @doc("Remove the profile picture of the logged in user")
    removeProfilePicture(): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/my-tier")
interface MyTier {
    @tag("HubProfile")
    @get
    @useAuth(HubAuth)
    myTier(): {
        @statusCode statusCode: 200;
        @body tier: HubUserTier;
    };
}
