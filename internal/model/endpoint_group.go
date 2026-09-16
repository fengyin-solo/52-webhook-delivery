package model

import (
	"strings"
	"time"
)

const (
	EndpointGroupActive   = "active"
	EndpointGroupInactive = "inactive"
)

type EndpointGroup struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (eg *EndpointGroup) Validate() error {
	eg.Name = strings.TrimSpace(eg.Name)
	eg.Description = strings.TrimSpace(eg.Description)
	if eg.Name == "" {
		return NewValidationError("name", "分组名称不能为空")
	}
	if len(eg.Name) > 128 {
		return NewValidationError("name", "分组名称不能超过 128 个字符")
	}
	if eg.Status == "" {
		eg.Status = EndpointGroupActive
	}
	if eg.Status != EndpointGroupActive && eg.Status != EndpointGroupInactive {
		return NewValidationError("status", "分组状态不合法")
	}
	return nil
}

type EndpointGroupFilter struct {
	Status  string
	Keyword string
}

func (f EndpointGroupFilter) Match(eg *EndpointGroup) bool {
	if f.Status != "" && eg.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(eg.Name), k) &&
			!strings.Contains(strings.ToLower(eg.Description), k) {
			return false
		}
	}
	return true
}
