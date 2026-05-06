package helix

import (
	"encoding/json"
	"testing"
)

type timeTest struct {
	T Time `json:"started_at"`
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		expectZero    bool
		testStr       string
		datetimeValue string
	}{
		{true, "", ""},
		{true, "null", ""},
		{false, "2018-02-05T08:15:59Z", "2018-02-05 08:15:59 +0000 UTC"},
	}

	for _, testCase := range testCases {
		tme := &timeTest{}
		err := json.Unmarshal([]byte(`{"started_at": "`+testCase.testStr+`"}`), tme)
		if err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if !tme.T.IsZero() {
			if tme.T.String() != testCase.datetimeValue {
				t.Errorf("expected time to be \"%s\", got \"%s\"", testCase.datetimeValue, tme.T.String())
			}
		}

		if testCase.expectZero && !tme.T.IsZero() {
			t.Errorf("expected zero value for time, got \"%s\"", tme.T.String())
		}
	}
}
