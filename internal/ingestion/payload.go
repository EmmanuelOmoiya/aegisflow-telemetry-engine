package ingestion

import (
	"errors"
	"strings"
)

type TelemetryEvent struct {
	DeviceID  string                 `json:"device_id"`
	Timestamp int64                  `json:"timestamp"`
	Source    string                 `json:"source"`
	EventType string                 `json:"event_type"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func (e *TelemetryEvent) Validate() error {
	if strings.TrimSpace(e.DeviceID) == "" {
		return errors.New("validation failed: device_id cannot be blank")
	}
	if e.Timestamp <= 0 {
		return errors.New("validation failed: timestamp must be a positive epoch integer")
	}
	if strings.TrimSpace(e.Source) == "" {
		return errors.New("validation failed: source identifier required")
	}
	if strings.TrimSpace(e.EventType) == "" {
		return errors.New("validation failed: event_type identifier required")
	}
	if e.Metadata == nil {
		return errors.New("validation failed: metadata object payload cannot be empty")
	}
	return nil
}
