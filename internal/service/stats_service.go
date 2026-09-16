package service

import (
	"math"
	"sort"

	"webhook/internal/model"
)

// DeliveryOverviewStats 投递总览统计。
type DeliveryOverviewStats struct {
	TotalCount      int     `json:"total_count"`
	DeliveredCount  int     `json:"delivered_count"`
	FailedCount     int     `json:"failed_count"`
	PendingCount    int     `json:"pending_count"`
	RetryingCount   int     `json:"retrying_count"`
	SuccessRate     float64 `json:"success_rate"`
}

// EndpointPerformance 端点性能统计。
type EndpointPerformance struct {
	EndpointID   string  `json:"endpoint_id"`
	TotalCount   int     `json:"total_count"`
	SuccessCount int     `json:"success_count"`
	SuccessRate  float64 `json:"success_rate"`
	AvgLatencyMs float64 `json:"avg_latency_ms"`
	P95LatencyMs float64 `json:"p95_latency_ms"`
}

// EventTypeDistribution 事件类型分布。
type EventTypeDistribution struct {
	EventType string `json:"event_type"`
	Count     int    `json:"count"`
}

// DeliveryTimelineEntry 投递时间线条目。
type DeliveryTimelineEntry struct {
	Date        string `json:"date"`
	Total       int    `json:"total"`
	Delivered   int    `json:"delivered"`
	Failed      int    `json:"failed"`
}

// SnapshotExport 全量快照导出。
type SnapshotExport struct {
	Endpoints        []*model.Endpoint        `json:"endpoints"`
	Subscriptions    []*model.Subscription    `json:"subscriptions"`
	EventTypes       []*model.EventType       `json:"event_types"`
	Events           []*model.Event           `json:"events"`
	Deliveries       []*model.Delivery        `json:"deliveries"`
	DeliveryAttempts []*model.DeliveryAttempt `json:"delivery_attempts"`
	RetryPolicies    []*model.RetryPolicy     `json:"retry_policies"`
	SigningKeys      []*model.SigningKey      `json:"signing_keys"`
	Filters          []*model.Filter          `json:"filters"`
	AuditLogs        []*model.AuditLog        `json:"audit_logs"`
	WebhookLogs      []*model.WebhookLog      `json:"webhook_logs"`
	EndpointGroups   []*model.EndpointGroup   `json:"endpoint_groups"`
	MessageTemplates []*model.MessageTemplate `json:"message_templates"`
}

func (s *Service) GetDeliveryOverviewStats() DeliveryOverviewStats {
	all := s.store.ListDeliveries()
	stats := DeliveryOverviewStats{TotalCount: len(all)}
	for _, d := range all {
		switch d.Status {
		case model.DeliveryDelivered:
			stats.DeliveredCount++
		case model.DeliveryFailed:
			stats.FailedCount++
		case model.DeliveryPending:
			stats.PendingCount++
		case model.DeliveryRetrying:
			stats.RetryingCount++
		}
	}
	if stats.TotalCount > 0 {
		stats.SuccessRate = math.Round(float64(stats.DeliveredCount)/float64(stats.TotalCount)*10000) / 100
	}
	return stats
}

func (s *Service) GetEndpointPerformance() []EndpointPerformance {
	allDeliveries := s.store.ListDeliveries()
	allAttempts := s.store.ListDeliveryAttempts()
	endpointMap := make(map[string]*EndpointPerformance)
	for _, d := range allDeliveries {
		if _, ok := endpointMap[d.EndpointID]; !ok {
			endpointMap[d.EndpointID] = &EndpointPerformance{EndpointID: d.EndpointID}
		}
		endpointMap[d.EndpointID].TotalCount++
		if d.Status == model.DeliveryDelivered {
			endpointMap[d.EndpointID].SuccessCount++
		}
	}
	attemptMap := make(map[string][]int64)
	for _, da := range allAttempts {
		attemptMap[da.DeliveryID] = append(attemptMap[da.DeliveryID], da.DurationMs)
	}
	for _, ep := range endpointMap {
		if ep.TotalCount > 0 {
			ep.SuccessRate = math.Round(float64(ep.SuccessCount)/float64(ep.TotalCount)*10000) / 100
		}
		var allDurations []int64
		for _, d := range allDeliveries {
			if d.EndpointID == ep.EndpointID {
				if durations, ok := attemptMap[d.ID]; ok {
					allDurations = append(allDurations, durations...)
				}
			}
		}
		if len(allDurations) > 0 {
			var sum int64
			for _, v := range allDurations {
				sum += v
			}
			ep.AvgLatencyMs = math.Round(float64(sum)/float64(len(allDurations))*100) / 100
			sort.Slice(allDurations, func(i, j int) bool { return allDurations[i] < allDurations[j] })
			idx := int(math.Ceil(float64(len(allDurations))*0.95)) - 1
			if idx < 0 {
				idx = 0
			}
			if idx >= len(allDurations) {
				idx = len(allDurations) - 1
			}
			ep.P95LatencyMs = float64(allDurations[idx])
		}
	}
	result := make([]EndpointPerformance, 0, len(endpointMap))
	for _, v := range endpointMap {
		result = append(result, *v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalCount > result[j].TotalCount
	})
	return result
}

func (s *Service) GetEventTypeDistribution() []EventTypeDistribution {
	allEvents := s.store.ListEvents()
	counts := make(map[string]int)
	for _, e := range allEvents {
		counts[e.Type]++
	}
	result := make([]EventTypeDistribution, 0, len(counts))
	for k, v := range counts {
		result = append(result, EventTypeDistribution{EventType: k, Count: v})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}

func (s *Service) GetDeliveryTimeline() []DeliveryTimelineEntry {
	allDeliveries := s.store.ListDeliveries()
	dateMap := make(map[string]*DeliveryTimelineEntry)
	for _, d := range allDeliveries {
		date := d.CreatedAt.Format("2006-01-02")
		if _, ok := dateMap[date]; !ok {
			dateMap[date] = &DeliveryTimelineEntry{Date: date}
		}
		dateMap[date].Total++
		if d.Status == model.DeliveryDelivered {
			dateMap[date].Delivered++
		} else if d.Status == model.DeliveryFailed {
			dateMap[date].Failed++
		}
	}
	result := make([]DeliveryTimelineEntry, 0, len(dateMap))
	for _, v := range dateMap {
		result = append(result, *v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})
	return result
}

func (s *Service) ExportSnapshot() SnapshotExport {
	return SnapshotExport{
		Endpoints:        s.store.ListEndpoints(),
		Subscriptions:    s.store.ListSubscriptions(),
		EventTypes:       s.store.ListEventTypes(),
		Events:           s.store.ListEvents(),
		Deliveries:       s.store.ListDeliveries(),
		DeliveryAttempts: s.store.ListDeliveryAttempts(),
		RetryPolicies:    s.store.ListRetryPolicies(),
		SigningKeys:      s.store.ListSigningKeys(),
		Filters:          s.store.ListFilters(),
		AuditLogs:        s.store.ListAuditLogs(),
		WebhookLogs:      s.store.ListWebhookLogs(),
		EndpointGroups:   s.store.ListEndpointGroups(),
		MessageTemplates: s.store.ListMessageTemplates(),
	}
}

// EndpointGroupStats 端点分组统计。
type EndpointGroupStats struct {
	GroupID     string `json:"group_id"`
	GroupName   string `json:"group_name"`
	EndpointCnt int    `json:"endpoint_cnt"`
}

// WebhookLogResultStats Webhook 日志结果分布。
type WebhookLogResultStats struct {
	Result string `json:"result"`
	Count  int    `json:"count"`
}

// GetWebhookLogResultStats 按结果统计 Webhook 日志分布。
func (s *Service) GetWebhookLogResultStats() []WebhookLogResultStats {
	all := s.store.ListWebhookLogs()
	counts := make(map[string]int)
	for _, wl := range all {
		counts[wl.Result]++
	}
	result := make([]WebhookLogResultStats, 0, len(counts))
	for k, v := range counts {
		result = append(result, WebhookLogResultStats{Result: k, Count: v})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}

// TemplateUsageStats 模板使用统计。
type TemplateUsageStats struct {
	EventType string `json:"event_type"`
	Count     int    `json:"count"`
}

// GetTemplateUsageStats 按事件类型统计模板数量。
func (s *Service) GetTemplateUsageStats() []TemplateUsageStats {
	all := s.store.ListMessageTemplates()
	counts := make(map[string]int)
	for _, mt := range all {
		counts[mt.EventType]++
	}
	result := make([]TemplateUsageStats, 0, len(counts))
	for k, v := range counts {
		result = append(result, TemplateUsageStats{EventType: k, Count: v})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}

func (s *Service) GetLatencyStats() (avgMs, p95Ms float64) {
	all := s.store.ListDeliveryAttempts()
	if len(all) == 0 {
		return 0, 0
	}
	durations := make([]int64, len(all))
	var sum int64
	for i, da := range all {
		durations[i] = da.DurationMs
		sum += da.DurationMs
	}
	avgMs = math.Round(float64(sum)/float64(len(durations))*100) / 100
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	idx := int(math.Ceil(float64(len(durations))*0.95)) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(durations) {
		idx = len(durations) - 1
	}
	p95Ms = float64(durations[idx])
	return avgMs, p95Ms
}
