import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "./locations.tsp";
import "./costcenters.tsp";
using TypeSpec.Http;
using TypeSpec.Rest;

@service({
    title: "Vetchi API",
})
@server("http://example.com", "blah")
namespace Vetchi;
