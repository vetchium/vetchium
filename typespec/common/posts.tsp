import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model Post {
    id: string;
    content: string;
    tags: string[];
    author_name: string;
    author_handle: Handle;
    created_at: string;
    upvotes_count: int32;
    downvotes_count: int32;
    score: int32;
}
