package google

import "testing"

func TestParseLoggingSinkId(t *testing.T) {
	tests := []struct {
		val             string
		defResourceType string
		defResourceId   string
		out             *LoggingSinkId
		errExpected     bool
	}{
		{
			val:             "projects/my-project/sinks/my-sink",
			defResourceType: "projects",
			defResourceId:   "my-project",
			out:             &LoggingSinkId{"projects", "my-project", "my-sink"},
			errExpected:     false,
		}, {
			val:             "folders/foofolder/sinks/woo",
			defResourceType: "invalidResourceType", // ignored as its in the val
			defResourceId:   "InvalidResourceId",   // ignored as its in the val
			out:             &LoggingSinkId{"folders", "foofolder", "woo"},
			errExpected:     false},
		{
			val: "kitchens/the-big-one/sinks/second-from-the-left",
			nil, true},
		// TODO add more tests!!!!

	}

	for _, test := range tests {
		out, err := parseLoggingSinkId(test.val)
		if err != nil {
			if !test.errExpected {
				t.Errorf("Got error with val %#v: error = %#v", test.val, err)
			}
		} else {
			if *out != *test.out {
				t.Errorf("Mismatch on val %#v: expected %#v but got %#v", test.val, test.out, out)
			}
		}
	}
}

func TestLoggingSinkId(t *testing.T) {
	tests := []struct {
		val         LoggingSinkId
		canonicalId string
		parent      string
	}{
		{
			val:         LoggingSinkId{"projects", "my-project", "my-sink"},
			canonicalId: "projects/my-project/sinks/my-sink",
			parent:      "projects/my-project",
		}, {
			val:         LoggingSinkId{"folders", "foofolder", "woo"},
			canonicalId: "folders/foofolder/sinks/woo",
			parent:      "folders/foofolder",
		},
	}

	for _, test := range tests {
		canonicalId := test.val.canonicalId()

		if canonicalId != test.canonicalId {
			t.Errorf("canonicalId mismatch on val %#v: expected %#v but got %#v", test.val, test.canonicalId, canonicalId)
		}

		parent := test.val.parent()

		if parent != test.parent {
			t.Errorf("parent mismatch on val %#v: expected %#v but got %#v", test.val, test.parent, parent)
		}
	}
}
