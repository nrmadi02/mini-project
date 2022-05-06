package seeds

import (
	"github.com/nrmadi02/mini-project/domain"
	"github.com/nrmadi02/mini-project/internal/role/utils"
)

func (s Seed) RoleSeed() {

	s.db.FirstOrCreate(&domain.Role{}, domain.Role{Name: role.Admin.String(), ID: 1})
	s.db.FirstOrCreate(&domain.Role{}, domain.Role{Name: role.Client.String(), ID: 2})
}
