package helix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateCustomReward(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *ChannelCustomRewardsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&ChannelCustomRewardsParams{
				BroadcasterID: "145328278",
				Title:         "game analysis 1v1",
				Cost:          50000,
			},
			`{"data": [{"broadcaster_name": "torpedo09","broadcaster_login": "torpedo09","broadcaster_id": "145328278","id": "afaa7e34-6b17-49f0-a19a-d1e76eaaf673","image": null,"background_color": "#00E5CB","is_enabled": true,"cost": 50000,"title": "game analysis 1v1","prompt": "","is_user_input_required": false,"max_per_stream_setting": {"is_enabled": false,"max_per_stream": 0},"max_per_user_per_stream_setting": {"is_enabled": false,"max_per_user_per_stream": 0},"global_cooldown_setting": {"is_enabled": false,"global_cooldown_seconds": 0},"is_paused": false,"is_in_stock": true,"default_image": {"url_1x": "https://static-cdn.jtvnw.net/custom-reward-images/default-1.png","url_2x": "https://static-cdn.jtvnw.net/custom-reward-images/default-2.png","url_4x": "https://static-cdn.jtvnw.net/custom-reward-images/default-4.png"},"should_redemptions_skip_request_queue": false,"redemptions_redeemed_current_stream": null,"cooldown_expires_at": null}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&ChannelCustomRewardsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.CreateCustomReward(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}
	}

	// Test with HTTP Failure
	options := &Options{
		ClientID: "my-client-id",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	}
	c := &Client{
		opts: options,
		ctx:  context.Background(),
	}

	_, err := c.CreateCustomReward(&ChannelCustomRewardsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetCustomRewardsRedemptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *GetCustomRewardsRedemptionsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "gx2pv4208cff0ig9ou7nk3riccffxt"},
			&GetCustomRewardsRedemptionsParams{
				BroadcasterID: "274637212",
				RewardID:      "92af127c-7326-4483-a52b-b0da0be61c01",
				Status:        "CANCELED",
			},
			`{"data":[{"broadcaster_name":"torpedo09","broadcaster_login":"torpedo09","broadcaster_id":"274637212","id":"17fa2df1-ad76-4804-bfa5-a40ef63efe63","user_login":"torpedo09","user_id":"274637212","user_name":"torpedo09","user_input":"","status":"CANCELED","redeemed_at":"2020-07-01T18:37:32Z","reward":{"id":"92af127c-7326-4483-a52b-b0da0be61c01","title":"game analysis","prompt":"","cost":50000}}],"pagination":{"cursor":"eyJiIjpudWxsLCJhIjp7IkN1cnNvciI6Ik1UZG1ZVEprWmpFdFlXUTNOaTAwT0RBMExXSm1ZVFV0WVRRd1pXWTJNMlZtWlRZelgxOHlNREl3TFRBM0xUQXhWREU0T2pNM09qTXlMakl6TXpFeU56RTFOMW89In19"}}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "gx2pv4208cff0ig9ou7nk3riccffxt"},
			&GetCustomRewardsRedemptionsParams{
				BroadcasterID: "",
				RewardID:      "92af127c-7326-4483-a52b-b0da0be61c01",
				Status:        "CANCELED",
			},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetCustomRewardsRedemptions(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if testCase.statusCode == http.StatusOK {
			if len(resp.Data.Redemptions) != 1 {
				t.Errorf("expected 1 redemption, got %d", len(resp.Data.Redemptions))
			}
		}
	}

	// Test with HTTP Failure
	options := &Options{
		ClientID: "gx2pv4208cff0ig9ou7nk3riccffxt",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	}
	c := &Client{
		opts: options,
		ctx:  context.Background(),
	}

	_, err := c.GetCustomRewardsRedemptions(&GetCustomRewardsRedemptionsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestUpdateCustomReward(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *UpdateChannelCustomRewardsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&UpdateChannelCustomRewardsParams{
				ID:            "123",
				BroadcasterID: "145328278",
				Title:         "game analysis 1v1",
				Cost:          50000,
			},
			`{"data": [{"broadcaster_name": "torpedo09","broadcaster_login": "torpedo09","broadcaster_id": "145328278","id": "afaa7e34-6b17-49f0-a19a-d1e76eaaf673","image": null,"background_color": "#00E5CB","is_enabled": true,"cost": 50000,"title": "game analysis 1v1","prompt": "","is_user_input_required": false,"max_per_stream_setting": {"is_enabled": false,"max_per_stream": 0},"max_per_user_per_stream_setting": {"is_enabled": false,"max_per_user_per_stream": 0},"global_cooldown_setting": {"is_enabled": false,"global_cooldown_seconds": 0},"is_paused": false,"is_in_stock": true,"default_image": {"url_1x": "https://static-cdn.jtvnw.net/custom-reward-images/default-1.png","url_2x": "https://static-cdn.jtvnw.net/custom-reward-images/default-2.png","url_4x": "https://static-cdn.jtvnw.net/custom-reward-images/default-4.png"},"should_redemptions_skip_request_queue": false,"redemptions_redeemed_current_stream": null,"cooldown_expires_at": null}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&UpdateChannelCustomRewardsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UpdateCustomReward(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}
	}

	// Test with HTTP Failure
	options := &Options{
		ClientID: "my-client-id",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	}
	c := &Client{
		opts: options,
		ctx:  context.Background(),
	}

	_, err := c.UpdateCustomReward(&UpdateChannelCustomRewardsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestDeleteCustomRewards(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *DeleteCustomRewardsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&DeleteCustomRewardsParams{
				BroadcasterID: "145328278",
				ID:            "84da6b13-efe1-4a82-91d0-25260aeb6a9b",
			},
			``,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&DeleteCustomRewardsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.DeleteCustomRewards(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}
	}

	// Test with HTTP Failure
	options := &Options{
		ClientID: "my-client-id",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	}
	c := &Client{
		opts: options,
		ctx:  context.Background(),
	}

	_, err := c.DeleteCustomRewards(&DeleteCustomRewardsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestGetCustomRewards(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *GetCustomRewardsParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&GetCustomRewardsParams{
				BroadcasterID: "145328278",
			},
			`{"data":[{"broadcaster_name":"Scorfly","broadcaster_login":"scorfly","broadcaster_id":"145328278","id":"7e86e512-f423-45b2-9e4b-2ad99a814ffc","image":null,"background_color":"#421561","is_enabled":false,"cost":50000,"title":"game analysis 1v1","prompt":"","is_user_input_required":false,"max_per_stream_setting":{"is_enabled":false,"max_per_stream":0},"max_per_user_per_stream_setting":{"is_enabled":false,"max_per_user_per_stream":0},"global_cooldown_setting":{"is_enabled":false,"global_cooldown_seconds":0},"is_paused":false,"is_in_stock":true,"default_image":{"url_1x":"https://static-cdn.jtvnw.net/custom-reward-images/default-1.png","url_2x":"https://static-cdn.jtvnw.net/custom-reward-images/default-2.png","url_4x":"https://static-cdn.jtvnw.net/custom-reward-images/default-4.png"},"should_redemptions_skip_request_queue":false,"redemptions_redeemed_current_stream":null,"cooldown_expires_at":null},{"broadcaster_name":"Scorfly","broadcaster_login":"scorfly","broadcaster_id":"145328278","id":"c0d6687c-c543-4a3a-8494-502451e7fa45","image":null,"background_color":"#FFBF00","is_enabled":false,"cost":50000,"title":"game analysis 1v2","prompt":"","is_user_input_required":false,"max_per_stream_setting":{"is_enabled":false,"max_per_stream":0},"max_per_user_per_stream_setting":{"is_enabled":false,"max_per_user_per_stream":0},"global_cooldown_setting":{"is_enabled":false,"global_cooldown_seconds":0},"is_paused":false,"is_in_stock":true,"default_image":{"url_1x":"https://static-cdn.jtvnw.net/custom-reward-images/default-1.png","url_2x":"https://static-cdn.jtvnw.net/custom-reward-images/default-2.png","url_4x":"https://static-cdn.jtvnw.net/custom-reward-images/default-4.png"},"should_redemptions_skip_request_queue":false,"redemptions_redeemed_current_stream":null,"cooldown_expires_at":null}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&GetCustomRewardsParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"broadcaster_id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.GetCustomRewards(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"broadcaster_id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}
	}

	// Test with HTTP Failure
	options := &Options{
		ClientID: "my-client-id",
		HTTPClient: &badMockHTTPClient{
			newMockHandler(0, "", nil),
		},
	}
	c := &Client{
		opts: options,
		ctx:  context.Background(),
	}

	_, err := c.GetCustomRewards(&GetCustomRewardsParams{})
	if err == nil {
		t.Error("expected error but got nil")
	}

	if err.Error() != "failed to execute API request: Oops, that's bad :(" {
		t.Error("expected error does match return error")
	}
}

func TestUpdateCustomRewardRedemptionStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		options    *Options
		params     *UpdateChannelCustomRewardsRedemptionStatusParams
		respBody   string
	}{
		{
			http.StatusOK,
			&Options{ClientID: "my-client-id"},
			&UpdateChannelCustomRewardsRedemptionStatusParams{
				ID:            "17fa2df1-ad76-4804-bfa5-a40ef63efe63",
				BroadcasterID: "274637212",
				RewardID:      "92af127c-7326-4483-a52b-b0da0be61c01",
				Status:        "CANCELLED",
			},
			`{"data": [{"broadcaster_name": "torpedo09", "broadcaster_login": "torpedo09", "broadcaster_id": "274637212", "id": "17fa2df1-ad76-4804-bfa5-a40ef63efe63", "user_id": "274637212", "user_name": "torpedo09", "user_login": "torpedo09", "user_input": "", "status": "CANCELED", "redeemed_at": "2020-07-01T18:37:32Z", "reward": { "id": "92af127c-7326-4483-a52b-b0da0be61c01", "title": "game analysis", "prompt": "", "cost": 50000}}]}`,
		},
		{
			http.StatusBadRequest,
			&Options{ClientID: "my-client-id"},
			&UpdateChannelCustomRewardsRedemptionStatusParams{},
			`{"error":"Bad Request","status":400,"message":"Missing required parameter \"id\""}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient(testCase.options, newMockHandler(testCase.statusCode, testCase.respBody, nil))

		resp, err := c.UpdateChannelCustomRewardsRedemptionStatus(testCase.params)
		if err != nil {
			t.Error(err)
		}

		// Test Bad Request Responses
		if resp.StatusCode == http.StatusBadRequest {
			firstErrStr := "Missing required parameter \"id\""
			if resp.ErrorMessage != firstErrStr {
				t.Errorf("expected error message to be \"%s\", got \"%s\"", firstErrStr, resp.ErrorMessage)
			}
			continue
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code to be \"%d\", got \"%d\"", testCase.statusCode, resp.StatusCode)
		}

		if testCase.statusCode == http.StatusOK {
			numRedemptions := len(resp.Data.Redemptions)
			if numRedemptions != 1 {
				t.Errorf("expected 1 redemption, got %d", numRedemptions)
				continue
			}

			title := resp.Data.Redemptions[0].Reward.Title
			if title != "game analysis" {
				t.Errorf("expected reward title to be \"game analysis\", got \"%s\"", title)
			}
		}
	}
}

// bodyCapturingHTTPClient records the request sent by the client, including
// its body, which mockHTTPClient consumes before the handler can see it.
type bodyCapturingHTTPClient struct {
	req  *http.Request
	body []byte
}

func (c *bodyCapturingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.req = req
	if req.Body != nil {
		defer req.Body.Close() //nolint:errcheck
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		c.body = b
	}

	rr := httptest.NewRecorder()
	rr.WriteHeader(http.StatusOK)
	rr.WriteString(`{"data":[]}`) //nolint:errcheck

	return rr.Result(), nil
}

func newBodyCapturingClient() (*Client, *bodyCapturingHTTPClient) {
	capture := &bodyCapturingHTTPClient{}
	return &Client{
		opts: &Options{ClientID: "my-client-id", APIBaseURL: DefaultAPIBaseURL, HTTPClient: capture},
		ctx:  context.Background(),
	}, capture
}

func boolPtr(b bool) *bool { return &b }

func assertRequestTarget(t *testing.T, req *http.Request, method, url string) {
	t.Helper()
	if req.Method != method {
		t.Errorf("expected method %s, got %s", method, req.Method)
	}
	if got := req.URL.Scheme + "://" + req.URL.Host + req.URL.Path; got != url {
		t.Errorf("expected URL %s, got %s", url, got)
	}
}

func decodeBody(t *testing.T, body []byte) map[string]any {
	t.Helper()
	payload := map[string]any{}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("request body is not valid JSON: %v (%s)", err, body)
	}
	return payload
}

func TestCreateCustomRewardPayload(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		params *ChannelCustomRewardsParams
		want   map[string]any
	}{
		{
			name:   "no body fields set sends an empty object",
			params: &ChannelCustomRewardsParams{BroadcasterID: "145328278"},
			want:   map[string]any{},
		},
		{
			name: "only required fields set sends only those",
			params: &ChannelCustomRewardsParams{
				BroadcasterID: "145328278",
				Title:         "game analysis",
				Cost:          50000,
			},
			want: map[string]any{"title": "game analysis", "cost": float64(50000)},
		},
		{
			name: "explicit false booleans are sent",
			params: &ChannelCustomRewardsParams{
				BroadcasterID:                     "145328278",
				Title:                             "game analysis",
				Cost:                              50000,
				IsEnabled:                         boolPtr(false),
				IsUserInputRequired:               boolPtr(false),
				IsMaxPerStreamEnabled:             boolPtr(false),
				IsMaxPerUserPerStreamEnabled:      boolPtr(false),
				IsGlobalCooldownEnabled:           boolPtr(false),
				ShouldRedemptionsSkipRequestQueue: boolPtr(false),
			},
			want: map[string]any{
				"title":                                 "game analysis",
				"cost":                                  float64(50000),
				"is_enabled":                            false,
				"is_user_input_required":                false,
				"is_max_per_stream_enabled":             false,
				"is_max_per_user_per_stream_enabled":    false,
				"is_global_cooldown_enabled":            false,
				"should_redemptions_skip_request_queue": false,
			},
		},
		{
			name: "every field set is sent",
			params: &ChannelCustomRewardsParams{
				BroadcasterID:                     "145328278",
				Title:                             "game analysis",
				Cost:                              50000,
				Prompt:                            "pick a game",
				IsEnabled:                         boolPtr(true),
				BackgroundColor:                   "#9147FF",
				IsUserInputRequired:               boolPtr(true),
				IsMaxPerStreamEnabled:             boolPtr(true),
				MaxPerStream:                      5,
				IsMaxPerUserPerStreamEnabled:      boolPtr(true),
				MaxPerUserPerStream:               2,
				IsGlobalCooldownEnabled:           boolPtr(true),
				GlobalCooldownSeconds:             60,
				ShouldRedemptionsSkipRequestQueue: boolPtr(true),
			},
			want: map[string]any{
				"title":                                 "game analysis",
				"cost":                                  float64(50000),
				"prompt":                                "pick a game",
				"is_enabled":                            true,
				"background_color":                      "#9147FF",
				"is_user_input_required":                true,
				"is_max_per_stream_enabled":             true,
				"max_per_stream":                        float64(5),
				"is_max_per_user_per_stream_enabled":    true,
				"max_per_user_per_stream":               float64(2),
				"is_global_cooldown_enabled":            true,
				"global_cooldown_seconds":               float64(60),
				"should_redemptions_skip_request_queue": true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, capture := newBodyCapturingClient()
			if _, err := c.CreateCustomReward(tc.params); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertRequestTarget(t, capture.req, http.MethodPost, "https://api.twitch.tv/helix/channel_points/custom_rewards")

			if got := capture.req.URL.Query().Get("broadcaster_id"); got != tc.params.BroadcasterID {
				t.Errorf("expected broadcaster_id query param %q, got %q", tc.params.BroadcasterID, got)
			}

			payload := decodeBody(t, capture.body)
			if !reflect.DeepEqual(payload, tc.want) {
				t.Errorf("unexpected body\n got: %s\nwant: %v", capture.body, tc.want)
			}
		})
	}
}

func TestUpdateCustomRewardPayload(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		params *UpdateChannelCustomRewardsParams
		want   map[string]any
	}{
		{
			name:   "no body fields set sends an empty object",
			params: &UpdateChannelCustomRewardsParams{ID: "123", BroadcasterID: "145328278"},
			want:   map[string]any{},
		},
		{
			name:   "only title set sends only title",
			params: &UpdateChannelCustomRewardsParams{ID: "123", BroadcasterID: "145328278", Title: "new title"},
			want:   map[string]any{"title": "new title"},
		},
		{
			name: "explicit false booleans are sent",
			params: &UpdateChannelCustomRewardsParams{
				ID:                                "123",
				BroadcasterID:                     "145328278",
				IsEnabled:                         boolPtr(false),
				IsUserInputRequired:               boolPtr(false),
				IsMaxPerStreamEnabled:             boolPtr(false),
				IsMaxPerUserPerStreamEnabled:      boolPtr(false),
				IsGlobalCooldownEnabled:           boolPtr(false),
				IsPaused:                          boolPtr(false),
				ShouldRedemptionsSkipRequestQueue: boolPtr(false),
			},
			want: map[string]any{
				"is_enabled":                            false,
				"is_user_input_required":                false,
				"is_max_per_stream_enabled":             false,
				"is_max_per_user_per_stream_enabled":    false,
				"is_global_cooldown_enabled":            false,
				"is_paused":                             false,
				"should_redemptions_skip_request_queue": false,
			},
		},
		{
			name: "every field set is sent",
			params: &UpdateChannelCustomRewardsParams{
				ID:                                "123",
				BroadcasterID:                     "145328278",
				Title:                             "game analysis",
				Cost:                              50000,
				Prompt:                            "pick a game",
				IsEnabled:                         boolPtr(true),
				BackgroundColor:                   "#00E5CB",
				IsUserInputRequired:               boolPtr(true),
				IsMaxPerStreamEnabled:             boolPtr(true),
				MaxPerStream:                      5,
				IsMaxPerUserPerStreamEnabled:      boolPtr(true),
				MaxPerUserPerStream:               2,
				IsGlobalCooldownEnabled:           boolPtr(true),
				GlobalCooldownSeconds:             60,
				IsPaused:                          boolPtr(true),
				ShouldRedemptionsSkipRequestQueue: boolPtr(true),
			},
			want: map[string]any{
				"title":                                 "game analysis",
				"cost":                                  float64(50000),
				"prompt":                                "pick a game",
				"is_enabled":                            true,
				"background_color":                      "#00E5CB",
				"is_user_input_required":                true,
				"is_max_per_stream_enabled":             true,
				"max_per_stream":                        float64(5),
				"is_max_per_user_per_stream_enabled":    true,
				"max_per_user_per_stream":               float64(2),
				"is_global_cooldown_enabled":            true,
				"global_cooldown_seconds":               float64(60),
				"is_paused":                             true,
				"should_redemptions_skip_request_queue": true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, capture := newBodyCapturingClient()
			if _, err := c.UpdateCustomReward(tc.params); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertRequestTarget(t, capture.req, http.MethodPatch, "https://api.twitch.tv/helix/channel_points/custom_rewards")

			query := capture.req.URL.Query()
			if got := query.Get("id"); got != tc.params.ID {
				t.Errorf("expected id query param %q, got %q", tc.params.ID, got)
			}
			if got := query.Get("broadcaster_id"); got != tc.params.BroadcasterID {
				t.Errorf("expected broadcaster_id query param %q, got %q", tc.params.BroadcasterID, got)
			}

			payload := decodeBody(t, capture.body)
			if !reflect.DeepEqual(payload, tc.want) {
				t.Errorf("unexpected body\n got: %s\nwant: %v", capture.body, tc.want)
			}
		})
	}
}

func TestUpdateRedemptionStatusPayload(t *testing.T) {
	t.Parallel()

	c, capture := newBodyCapturingClient()
	params := &UpdateChannelCustomRewardsRedemptionStatusParams{
		ID:            "17fa2df1-ad76-4804-bfa5-a40ef63efe63",
		BroadcasterID: "274637212",
		RewardID:      "92af127c-7326-4483-a52b-b0da0be61c01",
		Status:        "CANCELED",
	}
	if _, err := c.UpdateChannelCustomRewardsRedemptionStatus(params); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertRequestTarget(t, capture.req, http.MethodPatch, "https://api.twitch.tv/helix/channel_points/custom_rewards/redemptions")

	query := capture.req.URL.Query()
	for key, want := range map[string]string{"id": params.ID, "broadcaster_id": params.BroadcasterID, "reward_id": params.RewardID} {
		if got := query.Get(key); got != want {
			t.Errorf("expected %s query param %q, got %q", key, want, got)
		}
	}

	payload := decodeBody(t, capture.body)
	if want := map[string]any{"status": "CANCELED"}; !reflect.DeepEqual(payload, want) {
		t.Errorf("unexpected body\n got: %s\nwant: %v", capture.body, want)
	}
}
