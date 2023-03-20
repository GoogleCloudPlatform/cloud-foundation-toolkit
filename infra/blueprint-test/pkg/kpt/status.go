package kpt

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	// Resource event of type apply has apply type
	ResourceApplyType = "apply"
	// Status of successful resource event
	ResourceOperationSuccessful = "Successful"
	// Group event of type apply has summary type
	CompletedEventType = "summary"
)

type ResourceApplyStatus struct {
	Group     string `json:"group,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
}

type GroupApplyStatus struct {
	Action     string `json:"action"`
	Count      int    `json:"count"`
	Failed     int    `json:"failed"`
	Skipped    int    `json:"skipped"`
	Status     string `json:"status"`
	Successful int    `json:"successful"`
	Timestamp  string `json:"timestamp"`
	Type       string `json:"type"`
}

// GetPkgApplyResourcesStatus finds individual kpt apply statuses from newline separated string of apply statuses
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
		if s.Type == ResourceApplyType {
			statuses = append(statuses, s)
		}

	}
	return statuses, nil
}

// GetPkgApplyGroupStatus finds the first group kpt apply status from newline separated string of apply statuses
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
		if s.Type == CompletedEventType {
			return s, nil
		}
	}
	return s, fmt.Errorf("unable to find group status in json status %s", jsonStatus)
}
