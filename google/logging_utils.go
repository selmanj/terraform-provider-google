package google

import (
	"fmt"
	"regexp"
	"strings"
)

// loggingSinkResourceTypes contains all the possible Stackdriver Logging resource types. Used to parse ids safely.
var loggingSinkResourceTypes = []string{
	"billingAccount",
	"folders",
	"organizations",
	"projects",
}

// LoggingSinkId represents the parts that make up the canonical id used within terraform for a logging resource.
type LoggingSinkId struct {
	resourceType string
	resourceId   string
	name         string
}

// loggingSinkIdRegex matches valid logging sink canonical ids
var loggingSinkIdRegex = regexp.MustCompile("(.+)/(.+)/sinks/(.+)")

// canonicalId returns the LoggingSinkId as the canonical id used within terraform.
func (l LoggingSinkId) canonicalId() string {
	return fmt.Sprintf("%s/%s/sinks/%s", l.resourceType, l.resourceId, l.name)
}

// parent returns the "parent-level" resource that the sink is in (e.g. `folders/foo` for id `folders/foo/sinks/bar`)
func (l LoggingSinkId) parent() string {
	return fmt.Sprintf("%s/%s", l.resourceType, l.resourceId)
}

// parseLoggingSinkId parses an id into a LoggingSinkId, or returns an error on failure. You can provide an id in four
// different patterns:
//
// * [SINK_NAME]
// * sinks/[SINK_NAME]
// * [RESOURCE_NAME]/sinks/[SINK_NAME]
// * [RESOURCE_TYPE]/[RESOURCE_NAME]/sinks/[SINK_NAME]
//
// You must provide a default resourceType and resourceName to use if either are not found in the id.
func parseLoggingSinkId(id, defaultResourceName, defaultResourceType string) (*LoggingSinkId, error) {
	if id == "" {
		return nil, fmt.Errorf("Id \"\" is not a valid logging id")
	}

	// Set resourceType + resourceId to default values; we will overwrite them if necessary
	resourceType := defaultResourceType
	resourceId := defaultResourceName

	parts := strings.Split(id, "/")
	// SinkName is always the last element in the parts list
	sinkName := parts[len(parts)-1]
	// Verify that the second to last element is 'sinks'
	if len(parts) >= 2 && parts[len(parts)-2] != "sinks" {
		return nil, fmt.Errorf("Invalid Logging sink id: %s", id)
	}
	if len(parts) >= 3 {
		resourceId = parts[len(parts)-3]
	}
	if len(parts) == 4 {
		resourceType = parts[len(parts)-4]
	}
	if len(parts) > 4 {
		return nil, fmt.Errorf("Invalid Logging sink id: %s", id)
	}

	// If our resourceType is not a valid logging sink resource type, complain loudly
	validLoggingSinkResourceType := false
	for _, v := range loggingSinkResourceTypes {
		if v == resourceType {
			validLoggingSinkResourceType = true
			break
		}
	}

	if !validLoggingSinkResourceType {
		return nil, fmt.Errorf("Logging resource type %s is not valid. Valid resource types: %#v", parts[1],
			loggingSinkResourceTypes)
	}
	return &LoggingSinkId{
		resourceType: resourceType,
		resourceId:   resourceId,
		name:         sinkName,
	}, nil
}
