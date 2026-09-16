package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	FilterActive   = "active"
	FilterDisabled = "disabled"
)

const (
	OpEQ         = "eq"
	OpNE         = "ne"
	OpContains   = "contains"
	OpStartsWith = "starts_with"
)

type Filter struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Field     string    `json:"field"`
	Operator  string    `json:"operator"`
	Value     string    `json:"value"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (f *Filter) Validate() error {
	f.Name = strings.TrimSpace(f.Name)
	f.Field = strings.TrimSpace(f.Field)
	f.Value = strings.TrimSpace(f.Value)
	if f.Name == "" {
		return NewValidationError("name", "过滤器名称不能为空")
	}
	if f.Field == "" {
		return NewValidationError("field", "字段不能为空")
	}
	if f.Operator == "" {
		return NewValidationError("operator", "操作符不能为空")
	}
	if f.Operator != OpEQ && f.Operator != OpNE && f.Operator != OpContains && f.Operator != OpStartsWith {
		return NewValidationError("operator", "操作符不合法")
	}
	if f.Status == "" {
		f.Status = FilterActive
	}
	if f.Status != FilterActive && f.Status != FilterDisabled {
		return NewValidationError("status", "过滤器状态不合法")
	}
	return nil
}

func (f *Filter) MatchPayload(payload string) bool {
	if f.Status != FilterActive {
		return true
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return false
	}
	fieldVal, ok := data[f.Field]
	if !ok {
		return false
	}
	strVal := ""
	switch v := fieldVal.(type) {
	case string:
		strVal = v
	case float64:
		strVal = strings.TrimRight(strings.TrimSpace(fmt.Sprintf("%v", v)), ".0")
	case bool:
		strVal = fmt.Sprintf("%v", v)
	default:
		strVal = fmt.Sprintf("%v", v)
	}
	switch f.Operator {
	case OpEQ:
		return strVal == f.Value
	case OpNE:
		return strVal != f.Value
	case OpContains:
		return strings.Contains(strVal, f.Value)
	case OpStartsWith:
		return strings.HasPrefix(strVal, f.Value)
	default:
		return false
	}
}

type FilterFilter struct {
	Status  string
	Keyword string
}

func (f FilterFilter) Match(fl *Filter) bool {
	if f.Status != "" && fl.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(fl.Name), k) {
			return false
		}
	}
	return true
}
