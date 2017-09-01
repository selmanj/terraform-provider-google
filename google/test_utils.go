package google

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// To update:
// bigquery_dataset.go
// bigquery_table
// compute_disk
// compute_image
// compute_instance
// container_cluster
// container_node_pool
// spanner_instance

func TestAccCheckStructHasLabel(strct interface{}, labelKey, labelValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return CheckStructHasLabel(strct, labelKey, labelValue)
	}
}

func CheckStructHasLabel(strct interface{}, labelKey, labelValue string) error {
	labels, err := ReflectMapStringStringFromStruct(strct, "Labels")
	if err != nil {
		return err
	}
	v, ok := labels[labelKey]
	if !ok {
		return fmt.Errorf("Expected label key '%s' with value '%s' not found: labels = %#v", labelKey, labelValue, labels)
	}
	if v != labelValue {
		return fmt.Errorf("Label value for key '%s' does not match: expected '%s', but found '%s'", labelKey, labelValue, v)
	}
	return nil
}

func ReflectMapStringStringFromStruct(strct interface{}, fieldName string) (map[string]string, error) {
	strctVal := reflect.ValueOf(strct)
	// If we got a ptr, go ahead and deref
	if strctVal.Kind() == reflect.Ptr {
		strctVal.Pointer()
	}

	if strctVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Expected to reflect on value of type struct, but got %#v", kindToString[strctVal.Kind()])
	}

	fieldVal := strctVal.FieldByName(fieldName)
	if fieldVal.Kind() != reflect.Map {
		return nil, fmt.Errorf("Expected struct's field %#v to have type map, but instead found %#v", fieldName, fieldVal.Kind())
	}

	mp, ok := fieldVal.Interface().(map[string]string)
	if !ok {
		return nil, fmt.Errorf("Expected type map[string]string for field %s but it's a different map type", fieldName)
	}

	return mp, nil
}

// This is fine.
var kindToString = map[reflect.Kind]string{
	reflect.Invalid:       "Invalid",
	reflect.Bool:          "Bool",
	reflect.Int:           "Int",
	reflect.Int8:          "Int8",
	reflect.Int16:         "Int16",
	reflect.Int32:         "Int32",
	reflect.Int64:         "Int64",
	reflect.Uint:          "Uint",
	reflect.Uint8:         "Uint8",
	reflect.Uint16:        "Uint16",
	reflect.Uint32:        "Uint32",
	reflect.Uint64:        "Uint64",
	reflect.Uintptr:       "Uintptr",
	reflect.Float32:       "Float32",
	reflect.Float64:       "Float64",
	reflect.Complex64:     "Complex64",
	reflect.Complex128:    "Complex128",
	reflect.Array:         "Array",
	reflect.Chan:          "Chan",
	reflect.Func:          "Func",
	reflect.Interface:     "Interface",
	reflect.Map:           "Map",
	reflect.Ptr:           "Ptr",
	reflect.Slice:         "Slice",
	reflect.String:        "String",
	reflect.Struct:        "Struct",
	reflect.UnsafePointer: "UnsafePointer",
}
