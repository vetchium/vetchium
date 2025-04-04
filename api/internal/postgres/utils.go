package postgres

import (
	"fmt"

	"github.com/vetchium/vetchium/typespec/common"
)

func (p *PG) convertToOrgUserRoles(
	dbRoles []string,
) ([]common.OrgUserRole, error) {
	var roles []common.OrgUserRole
	for _, str := range dbRoles {
		role := common.OrgUserRole(str)
		switch role {
		case common.Admin,
			common.ApplicationsCRUD,
			common.ApplicationsViewer,
			common.CostCentersCRUD,
			common.CostCentersViewer,
			common.LocationsCRUD,
			common.LocationsViewer,
			common.OpeningsCRUD,
			common.OpeningsViewer,
			common.OrgUsersCRUD,
			common.OrgUsersViewer:
			roles = append(roles, role)
		default:
			p.log.Err("invalid role in the database", "role", str)
			return nil, fmt.Errorf("invalid role: %s", str)
		}
	}
	return roles, nil
}
