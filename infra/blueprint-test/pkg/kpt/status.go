package kpt

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	// Individual apply events can have either of resourceApplied or resourceFailed status.
	// https://github.com/GoogleContainerTools/kpt/blob/2a817f60cf7132c88fd2e526c02b800cf927c048/thirdparty/cli-utils/printers/json/formatter.go#L31
	ApplyType                = "apply"
	ResourceAppliedEventType = "resourceApplied"
	ResourceFailedEventType  = "resourceFailed"
	// Unchanged operation represents a resource that remained unchanged.
	ResourceOperationUnchanged = "Unchanged"
	// Group event of type apply has completed status
	CompletedEventType = "completed"
)

type ResourceApplyStatus struct {
	EventType string `json:"eventType"`
	Group     string `json:"group,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Operation string `json:"operation"`
	Type      string `json:"type"`
}

type GroupApplyStatus struct {
	EventType       string `json:"eventType"`
	Count           int    `json:"count"`
	CreatedCount    int    `json:"createdCount"`
	UnchangedCount  int    `json:"unchangedCount"`
	ConfiguredCount int    `json:"configuredCount"`
	FailedCount     int    `json:"failedCount"`
	ServerSideCount int    `json:"serverSideCount"`
	Operation       string `json:"operation"`
	Type            string `json:"type"`
}

// GetPkgApplyResourcesStatus finds individual kpt apply statuses from newline seperated string of apply statuses
// and converts them into a slice of ResourceApplyStatus.
func GetPkgApplyResourcesStatus(jsonStatus string) ([]ResourceApplyStatus, error) {
	var statuses []ResourceApplyStatus
	resourceStatus := strings.Split(jsonStatus, "\n")
	for _, status := range resourceStatus {
		var s ResourceApplyStatus
		err := json.Unmarshal([]byte(status), &s)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling %s: %v", status, err)
		}
		if s.Type == ApplyType && (s.EventType == ResourceAppliedEventType || s.EventType == ResourceFailedEventType) {
			statuses = append(statuses, s)
		}

	}
	return statuses, nil
}

// GetPkgApplyGroupStatus finds the first group kpt apply status from newline seperated string of apply statuses
// and converts it into a GroupApplyStatus.
func GetPkgApplyGroupStatus(jsonStatus string) (GroupApplyStatus, error) {
	var s GroupApplyStatus
	resourceStatuses := strings.Split(jsonStatus, "\n")
	// search in reverse as group status is usually the last status
	for i := len(resourceStatuses) - 1; i >= 0; i-- {
		err := json.Unmarshal([]byte(resourceStatuses[i]), &s)
		if err != nil {
			return s, fmt.Errorf("error unmarshalling %s: %v", resourceStatuses[i], err)
		}
		if s.Type == ApplyType && (s.EventType == CompletedEventType) {
			return s, nil
		}
	}
	return s, fmt.Errorf("unable to find group status in json status %s", jsonStatus)
}
