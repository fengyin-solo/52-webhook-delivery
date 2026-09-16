package model

// BatchCreateEventRequest 批量创建事件请求。
type BatchCreateEventRequest struct {
	Events []Event `json:"events"`
}

func (r *BatchCreateEventRequest) Validate() error {
	if len(r.Events) == 0 {
		return NewValidationError("events", "事件列表不能为空")
	}
	if len(r.Events) > 100 {
		return NewValidationError("events", "单次批量创建事件不能超过 100 条")
	}
	for i := range r.Events {
		if err := r.Events[i].Validate(); err != nil {
			return NewValidationError("events["+string(rune('0'+i))+"]", err.Error())
		}
	}
	return nil
}

// BatchDisableEndpointRequest 批量停用端点请求。
type BatchDisableEndpointRequest struct {
	IDs []string `json:"ids"`
}

func (r *BatchDisableEndpointRequest) Validate() error {
	if len(r.IDs) == 0 {
		return NewValidationError("ids", "ID 列表不能为空")
	}
	if len(r.IDs) > 100 {
		return NewValidationError("ids", "单次批量操作不能超过 100 条")
	}
	return nil
}

// BatchResult 批量操作结果。
type BatchResult struct {
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}
