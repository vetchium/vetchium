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
