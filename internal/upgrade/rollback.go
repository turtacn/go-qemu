package upgrade

import "fmt"

type RollbackHandler struct{}

func (r *RollbackHandler) Rollback(vmid int) error {
	return fmt.Errorf("rollback handled for vm %d", vmid)
}
