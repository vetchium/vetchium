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
    addOfficialEmail(@body request: AddOfficialEmailRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/verify-official-email")
interface VerifyOfficialEmail {
    @tag("HubProfile")
    @post
    @useAuth(HubAuth)
    verifyOfficialEmail(@body request: VerifyOfficialEmailRequest): {};
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
