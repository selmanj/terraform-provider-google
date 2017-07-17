package google

import (
	"bytes"
	"fmt"
	"reflect"
)

// ConfigBuilder is a helper class for generating Terraform config strings for use in tests.
type ConfigBuilder struct {
	// ResourceName is the name of the resource (e.g. the 'foo' in 'resource "google_compute_instance" "foo" {').
	ResourceName string
	// ResourceType is the type of the resource (e.g. the 'google_compute_instance' in 'resource
	// "google_compute_instance" "foo" {').
	ResourceType string
	// Attributes contains a mapping between all key/value pairs.
	Attributes map[string]interface{}
}

// NewResourceConfigBuilder creates a ConfigBuilder for a resource with the provided type and name.
func NewResourceConfigBuilder(typ, name string) *ConfigBuilder {
	return &ConfigBuilder{ResourceName: name, ResourceType: typ, Attributes: map[string]interface{}{}}
}

// NewNestedConfig is used for nesting maps (e.g. if you wanted to add a set of key/values for labels, you'd do
// something like:
//
// x := NewResourceConfigBuilder("google_container_cluster", "cluster-" + acctest.RandString(10)).
//          WithAttribute("labels", NewNestedConfig().
//              WithAttribute("my_label", "my_value"))
func NewNestedConfig() *ConfigBuilder {
	return &ConfigBuilder{Attributes: map[string]interface{}{}}
}

// WithResourceName sets the Terraform resource name.
func (rb *ConfigBuilder) WithResourceName(name string) *ConfigBuilder {
	rb.ResourceName = name
	return rb
}

// WithResourceType sets the Terraform resource type.
func (rb *ConfigBuilder) WithResourceType(typ string) *ConfigBuilder {
	rb.ResourceType = typ
	return rb
}

// WithAttribute sets an attribute on the resource. Anything that implements the Stringer interface or is a primitive
// can be used here. See NewNestedConfig() as well for an example on how to embed an additional map structure.
func (rb *ConfigBuilder) WithAttribute(key string, obj interface{}, objs ...interface{}) *ConfigBuilder {
	// Check if the value is to be treated as an array
	if objs != nil {
		arr := make([]interface{}, len(objs)+1)
		arr[0] = obj
		for idx := range objs {
			arr[idx+1] = objs[idx]
		}

		return rb.WithAttribute(key, arr)
	}
	rb.Attributes[key] = obj
	return rb
}

// Name returns the "name" attribute (commonly used in GCP resources).
func (rb ConfigBuilder) Name() string {
	return rb.Attributes["name"].(string)
}

// WithName sets the "name" attribute (commonly used in GCP resources).
func (rb *ConfigBuilder) WithName(name string) *ConfigBuilder {
	rb.Attributes["name"] = name
	return rb
}

// Zone returns the "zone" attribute (commonly used in GCP resources).
func (rb ConfigBuilder) Zone() string {
	return rb.Attributes["zone"].(string)
}

// WithZone sets the "zone" attribute (commonly used in GCP resources).
func (rb *ConfigBuilder) WithZone(zone string) *ConfigBuilder {
	rb.Attributes["zone"] = zone
	return rb
}

// String returns a pretty-printed string of the config.
func (rb ConfigBuilder) String() string {
	return rb.StringWithIndent(0, 4, false)
}

type StringWithIndenter interface {
	// StringWithIndent is like String, but allows for control of multiline resources. 'indent' represents how much to
	// indent every line. 'indentLen' controls how much indenting to add when adding additional indentation. 'embedded'
	// represents whether or not the produced string is embedded in a larger structure, in which case the leading
	// indentation on the first line is suppressed.
	StringWithIndent(indent, indentLen int, embedded bool) string
}

func (rb ConfigBuilder) StringWithIndent(indent, indentLen int, embedded bool) string {
	var buf bytes.Buffer

	if !embedded {
		buf.WriteString(spacesOfLength(indent))
	}
	if rb.ResourceName != "" && rb.ResourceType != "" {
		buf.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" ", rb.ResourceType, rb.ResourceName))
	}
	buf.WriteString("{\n")
	for k, v := range rb.Attributes {
		buf.WriteString(spacesOfLength(indent + indentLen))
		buf.WriteString(fmt.Sprintf("%s ", k))

		// Only emit a '= ' if the type is anything but a nested ConfigBuilder
		confBuilderType := reflect.TypeOf((*ConfigBuilder)(nil)).Elem()
		if reflect.TypeOf(v) != confBuilderType {
			buf.WriteString("= ")
		}
		buf.WriteString(asSimpleYAML(v, indent, indentLen))
		buf.WriteString("\n")
	}

	buf.WriteString(spacesOfLength(indent) + "}\n")
	return buf.String()
}

func asSimpleYAML(v interface{}, indent, indentLen int) string {
	switch v.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case []interface{}:
		var buf bytes.Buffer

		buf.WriteString(fmt.Sprintf("["))
		for idx, x := range v.([]interface{}) {
			if idx != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(asSimpleYAML(x, indent, indentLen))
		}
		buf.WriteString("]")
		return buf.String()
	case StringWithIndenter:
		return v.(StringWithIndenter).StringWithIndent(indent+indentLen, indentLen, true)
	case fmt.Stringer:
		return v.(fmt.Stringer).String()
	default:
		return "**UNKNOWN**"
	}
}

// spacesOfLength is a helper function for generating a string consisting of just spaces.
func spacesOfLength(len int) string {
	sp := make([]byte, len)
	for idx := range sp {
		sp[idx] = byte(' ')
	}
	return string(sp)
}
