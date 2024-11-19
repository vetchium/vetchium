import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

@minLength(3)
@maxLength(64)
scalar CostCenterName extends string;

enum CostCenterState {
    ActiveCC: "ACTIVE_CC",
    DefunctCC: "DEFUNCT_CC",
}

model CostCenter {
    name: CostCenterName;

    @maxLength(1024)
    notes?: string;

    state: CostCenterState;
}

model AddCostCenterRequest {
    name: CostCenterName;

    @maxLength(1024)
    notes?: string;
}

model DefunctCostCenterRequest {
    name: CostCenterName;
}

model GetCostCenterRequest {
    name: CostCenterName;
}

model GetCostCentersRequest {
    pagination_key?: CostCenterName;

    @maxValue(100)
    @doc("Number of Cost Centers to return; 100 is the default if not specified")
    limit?: integer;
}

model RenameCostCenterRequest {
    old_name: CostCenterName;
    new_name: CostCenterName;
}

model UpdateCostCenterRequest {
    name: CostCenterName;

    @maxLength(1024)
    notes: string;
}

@route("/employer/add-cost-center")
interface AddCostCenter {
    @tag("Cost Centers")
    @doc("Requires any of ${Admin}, ${CostCentersCRUD} roles")
    @post
    @useAuth(BearerAuth)
    addCostcenter(@body costcenter: AddCostCenterRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Cost Center already exists")
        @statusCode
        statusCode: 409;
    } | {
        @statusCode statusCode: 400;
        @body error: ValidationErrors;
    };
}

@route("/employer/defunct-cost-center")
interface DefunctCostCenter {
    @tag("Cost Centers")
    @doc("Requires any of ${Admin}, ${CostCentersCRUD} roles")
    @post
    @useAuth(BearerAuth)
    defunctCostCenter(
        @body defunctCostCenterRequest: DefunctCostCenterRequest,
    ): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/get-cost-center")
interface GetCostCenter {
    @tag("Cost Centers")
    @doc("Requires any of ${Admin}, ${CostCentersCRUD}, ${CostCentersViewer} roles")
    @post
    @useAuth(BearerAuth)
    getCostCenter(@body getCostCenterRequest: GetCostCenterRequest): {
        @statusCode statusCode: 200;
        @body costcenter: CostCenter;
    };
}

@route("/employer/get-cost-centers")
interface GetCostCenters {
    @tag("Cost Centers")
    @doc("Requires any of ${Admin}, ${CostCentersCRUD}, ${CostCentersViewer} roles")
    @post
    @useAuth(BearerAuth)
    getCostCenters(@body getCostCentersRequest: GetCostCentersRequest): {
        @statusCode statusCode: 200;
        @body costcenters: CostCenter[];
    };
}

@route("/employer/rename-cost-center")
interface RenameCostCenter {
    @tag("Cost Centers")
    @doc("Requires any of ${Admin}, ${CostCentersCRUD} roles")
    @post
    @useAuth(BearerAuth)
    renameCostCenter(@body renameCostCenterRequest: RenameCostCenterRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/update-cost-center")
interface UpdateCostCenter {
    @tag("Cost Centers")
    @doc("Requires any of ${Admin}, ${CostCentersCRUD} roles")
    @post
    @useAuth(BearerAuth)
    updateCostCenter(@body updateCostCenterRequest: UpdateCostCenterRequest): {
        @statusCode statusCode: 200;
    };
}
