package resources

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	condition "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

const (
	ConditionTypeReady     = "Ready"
	ConditionTypeInstalled = "Installed"
	ConditionTypeHealthy   = "Healthy"
	ConditionTypeSynced    = "Synced"
	ConditionTypeFailed    = "Failed"
)

const (
	ConditionStatusTrue    = "True"
	ConditionStatusFalse   = "False"
	ConditionStatusUnknown = "Unknown"
)

func ExtractVersion(packageStr string) string {
	parts := strings.Split(packageStr, ":")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "latest"
}

func ExtractInstallDate(timestamp time.Time) string {
	if !timestamp.IsZero() {
		return timestamp.Format("2006-01-02")
	}
	return "Unknown"
}

func ExtractStatusFromConditions(conditions interface{}) string {
	if conditions == nil {
		return "Unknown"
	}

	switch conds := conditions.(type) {
	case []condition.Condition:
		return extractStatusFromCrossplaneConditions(conds)
	case []interface{}:
		return extractStatusFromUnstructuredConditions(conds)
	default:
		return extractStatusFromStructuredConditions(conditions)
	}
}

func ExtractDescription(annotations map[string]string, packageStr string, resourceType string) string {
	if annotations != nil {
		if desc, ok := annotations["description"]; ok {
			return desc
		}
		if desc, ok := annotations["meta.crossplane.io/description"]; ok {
			return desc
		}
	}

	if packageStr != "" {
		parts := strings.Split(packageStr, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			if idx := strings.LastIndex(lastPart, ":"); idx > 0 {
				lastPart = lastPart[:idx]
			}
			if lastPart != "" {
				return fmt.Sprintf("Crossplane %s: %s", resourceType, lastPart)
			}
		}
	}

	return fmt.Sprintf("Crossplane %s package", resourceType)
}

func extractStatusFromCrossplaneConditions(conditions []condition.Condition) string {
	if len(conditions) == 0 {
		return "Unknown"
	}

	var ready, installed, healthy, synced bool
	var lastConditionMsg string

	for _, cond := range conditions {
		if cond.Message != "" {
			lastConditionMsg = cond.Message
		}

		switch string(cond.Type) {
		case ConditionTypeReady:
			ready = (string(cond.Status) == ConditionStatusTrue)
		case ConditionTypeInstalled:
			installed = (string(cond.Status) == ConditionStatusTrue)
		case ConditionTypeHealthy:
			healthy = (string(cond.Status) == ConditionStatusTrue)
		case ConditionTypeSynced:
			synced = (string(cond.Status) == ConditionStatusTrue)
		case ConditionTypeFailed:
			if string(cond.Status) == ConditionStatusTrue {
				if cond.Reason != "" {
					return fmt.Sprintf("Failed: %s", string(cond.Reason))
				}
				return "Failed"
			}
		}
	}

	if ready {
		return ConditionTypeReady
	}
	if healthy {
		return ConditionTypeHealthy
	}
	if installed {
		return ConditionTypeInstalled
	}
	if synced {
		return ConditionTypeSynced
	}

	for _, cond := range conditions {
		if string(cond.Status) == ConditionStatusFalse && cond.Reason != "" {
			return string(cond.Reason)
		}
	}

	if lastConditionMsg != "" && !strings.Contains(lastConditionMsg, "successfully") {
		return "Pending"
	}

	return "Unknown"
}

func extractStatusFromUnstructuredConditions(conditions []interface{}) string {
	if len(conditions) == 0 {
		return "Unknown"
	}

	var ready, installed, healthy, synced bool
	var lastConditionMsg string

	for _, cond := range conditions {
		condMap, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}

		condType, _ := condMap["type"].(string)
		condStatus, _ := condMap["status"].(string)
		condReason, _ := condMap["reason"].(string)
		condMessage, _ := condMap["message"].(string)

		if condMessage != "" {
			lastConditionMsg = condMessage
		}

		switch condType {
		case ConditionTypeReady:
			ready = (condStatus == ConditionStatusTrue)
		case ConditionTypeInstalled:
			installed = (condStatus == ConditionStatusTrue)
		case ConditionTypeHealthy:
			healthy = (condStatus == ConditionStatusTrue)
		case ConditionTypeSynced:
			synced = (condStatus == ConditionStatusTrue)
		case ConditionTypeFailed:
			if condStatus == ConditionStatusTrue {
				if condReason != "" {
					return fmt.Sprintf("Failed: %s", condReason)
				}
				return "Failed"
			}
		}
	}

	if ready {
		return ConditionTypeReady
	}
	if healthy {
		return ConditionTypeHealthy
	}
	if installed {
		return ConditionTypeInstalled
	}
	if synced {
		return ConditionTypeSynced
	}

	for _, cond := range conditions {
		condMap, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}
		condStatus, _ := condMap["status"].(string)
		condReason, _ := condMap["reason"].(string)

		if condStatus == ConditionStatusFalse && condReason != "" {
			return condReason
		}
	}

	if lastConditionMsg != "" && !strings.Contains(lastConditionMsg, "successfully") {
		return "Pending"
	}

	return "Unknown"
}

func extractStatusFromStructuredConditions(conditions interface{}) string {
	v := reflect.ValueOf(conditions)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return "Unknown"
	}

	if v.Len() == 0 {
		return "Unknown"
	}

	var ready, installed, healthy, synced bool
	var lastConditionMsg string

	for i := 0; i < v.Len(); i++ {
		cond := v.Index(i)

		condType := getFieldString(cond, "Type")
		condStatus := getFieldString(cond, "Status")
		condReason := getFieldString(cond, "Reason")
		condMessage := getFieldString(cond, "Message")

		if condMessage != "" {
			lastConditionMsg = condMessage
		}

		switch condType {
		case ConditionTypeReady:
			ready = (condStatus == ConditionStatusTrue)
		case ConditionTypeInstalled:
			installed = (condStatus == ConditionStatusTrue)
		case ConditionTypeHealthy:
			healthy = (condStatus == ConditionStatusTrue)
		case ConditionTypeSynced:
			synced = (condStatus == ConditionStatusTrue)
		case ConditionTypeFailed:
			if condStatus == ConditionStatusTrue {
				if condReason != "" {
					return fmt.Sprintf("Failed: %s", condReason)
				}
				return "Failed"
			}
		}
	}

	if ready {
		return ConditionTypeReady
	}
	if healthy {
		return ConditionTypeHealthy
	}
	if installed {
		return ConditionTypeInstalled
	}
	if synced {
		return ConditionTypeSynced
	}

	for i := 0; i < v.Len(); i++ {
		cond := v.Index(i)
		condStatus := getFieldString(cond, "Status")
		condReason := getFieldString(cond, "Reason")

		if condStatus == ConditionStatusFalse && condReason != "" {
			return condReason
		}
	}

	if lastConditionMsg != "" && !strings.Contains(lastConditionMsg, "successfully") {
		return "Pending"
	}

	return "Unknown"
}

func getFieldString(v reflect.Value, fieldName string) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return ""
	}

	if field.Kind() == reflect.String {
		return field.String()
	}

	if field.CanInterface() {
		if str, ok := field.Interface().(string); ok {
			return str
		}
		if stringer, ok := field.Interface().(fmt.Stringer); ok {
			return stringer.String()
		}
	}

	return ""
}
