package service

import (
	"webhook/internal/model"
)

func (s *Service) BatchCreateEvents(req model.BatchCreateEventRequest) (*model.BatchResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	result := &model.BatchResult{Success: 0, Failed: 0}
	for _, input := range req.Events {
		_, err := s.CreateEvent(input)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Success++
		}
	}
	return result, nil
}

func (s *Service) BatchDisableEndpoints(req model.BatchDisableEndpointRequest) (*model.BatchResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	result := &model.BatchResult{Success: 0, Failed: 0}
	for _, id := range req.IDs {
		existing, err := s.store.GetEndpoint(id)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, "端点 "+id+": "+err.Error())
			continue
		}
		if !model.EndpointCanTransition(existing.Status, model.EndpointDisabled) {
			result.Failed++
			result.Errors = append(result.Errors, "端点 "+id+": 状态不允许停用")
			continue
		}
		existing.Status = model.EndpointDisabled
		if err := s.store.UpdateEndpoint(existing); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, "端点 "+id+": "+err.Error())
			continue
		}
		result.Success++
	}
	return result, nil
}
