package google

import "testing"

func TestReflectMapStringStringFromStruct(t *testing.T) {
	labels := map[string]string{}

	labels["foobar"] = "bazz"
	err := CheckStructHasLabel(struct{ Labels map[string]string }{Labels: labels}, "foobar", "baz")
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
	}
}
