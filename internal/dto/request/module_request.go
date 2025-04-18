package request

import "fmt"

type GetModuleRequest struct {
	ModuleID int64 `json:"module_id"`
}

func (req *GetModuleRequest) ValidateRequest() (err error) {
	if req.ModuleID <= 0 {
		return fmt.Errorf("module_id не может быть меньше 1")
	}

	return nil
}
