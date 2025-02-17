import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

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

model Bio {
    handle: string;
    full_name: string;

    @doc("The short bio of the user. The one that gets displayed below the user's name.")
    short_bio: string;

    @doc("The long bio of the user. Something like 'About Me' of a user")
    long_bio: string;
}

model UpdateBioRequest {
    @doc("The handle for the logged in user")
    @minLength(1)
    @maxLength(32)
    handle?: string;

    @doc("The full name for the logged in user")
    @minLength(1)
    @maxLength(64)
    full_name?: string;

    @doc("The short bio for the logged inuser")
    @minLength(1)
    @maxLength(64)
    short_bio?: string;

    @doc("The long bio for the logged in user")
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
    addOfficialEmail(@body request: AddOfficialEmailRequest):
        | {
              @statusCode statusCode: 200;
          }
        | {
              @doc("The user has reached the maximum (50) allowed official emails per user.")
              @statusCode
              statusCode: 422;
          }
        | {
              @doc("The domain of the email address is not associated with an employer")
              @statusCode
              statusCode: 412;
          }
        | {
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
    @get
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
