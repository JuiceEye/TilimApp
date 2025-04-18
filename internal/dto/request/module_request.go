package request

import "fmt"

type GetMainPageModuleRequest struct {
	ModuleID int64 `json:"module_id"`
}

func (req *GetMainPageModuleRequest) ValidateRequest() (err error) {
	if req.ModuleID <= 0 {
		return fmt.Errorf("module_id не может быть меньше 1")
	}

	return nil
}
