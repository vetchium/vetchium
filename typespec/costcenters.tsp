import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "./common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

enum CostCenterState {
    ActiveCC: "ACTIVE_CC",
    DefunctCC: "DEFUNCT_CC",
}

model CostCenter {
    @minLength(3)
    @maxLength(64)
    name: string;

    @maxLength(1024)
    notes?: string;

    state: CostCenterState;
}

model AddCostCenterRequest {
    @minLength(3)
    @maxLength(64)
    name: string;

    @maxLength(1024)
    notes?: string;
}

model DefunctCostCenterRequest {
    @minLength(3)
    @maxLength(32)
    name: string;
}

model GetCostCenterRequest {
    @minLength(3)
    @maxLength(32)
    name: string;
}

model GetCostCentersRequest {
    @minLength(3)
    @maxLength(32)
    pagination_key?: string;

    @maxValue(100)
    @doc("Number of Cost Centers to return; 100 is the default if not specified")
    limit?: integer;
}

model RenameCostCenterRequest {
    @minLength(3)
    @maxLength(32)
    old_name: string;

    @minLength(3)
    @maxLength(32)
    new_name: string;
}

model UpdateCostCenterRequest {
    @minLength(3)
    @maxLength(64)
    name: string;

    @maxLength(1024)
    notes: string;
}

@route("/employer/add-cost-center")
interface AddCostCenter {
    @tag("Cost Centers")
    @doc("Requires any of ${Admin}, ${CostCentersCRUD} roles")
    @post
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
    renameCostCenter(@body renameCostCenterRequest: RenameCostCenterRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/update-cost-center")
interface UpdateCostCenter {
    @tag("Cost Centers")
    @doc("Requires any of ${Admin}, ${CostCentersCRUD} roles")
    @post
    updateCostCenter(@body updateCostCenterRequest: UpdateCostCenterRequest): {
        @statusCode statusCode: 200;
    };
}
