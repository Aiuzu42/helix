package helix

type EntitlementCodeStatus string

const (
	SUCCESSFULLY_REDEEMED EntitlementCodeStatus = "SUCCESSFULLY_REDEEMED"
	ALREADY_CLAIMED       EntitlementCodeStatus = "ALREADY_CLAIMED"
	EXPIRED               EntitlementCodeStatus = "EXPIRED"
	USER_NOT_ELIGIBLE     EntitlementCodeStatus = "USER_NOT_ELIGIBLE"
	NOT_FOUND             EntitlementCodeStatus = "NOT_FOUND"
	INACTIVE              EntitlementCodeStatus = "INACTIVE"
	UNUSED                EntitlementCodeStatus = "UNUSED"
	INCORRECT_FORMAT      EntitlementCodeStatus = "INCORRECT_FORMAT"
	INTERNAL_ERROR        EntitlementCodeStatus = "INTERNAL_ERROR"
)

type CodesParams struct {
	// One of the below
	UserID string   `query:"user_id"`
	Codes  []string `query:"code"` // Limit 20
}

type CodeStatus struct {
	Code   string                `json:"code"`
	Status EntitlementCodeStatus `json:"status"`
}

type ManyCodes struct {
	Codes []CodeStatus `json:"data"`
}

type CodeResponse struct {
	ResponseCommon
	Data ManyCodes
}

// GetEntitlementCodeStatus
// Per https://dev.twitch.tv/docs/api/reference#get-code-status
// Access is controlled via an app access token on the calling service. The client ID associated with the app access token must be approved by Twitch as part of a contracted arrangement.
// Callers with an app access token are authorized to redeem codes on behalf of any Twitch user account.
func (c *Client) GetEntitlementCodeStatus(params *CodesParams) (*CodeResponse, error) {
	resp, err := c.get("/entitlements/codes", &ManyCodes{}, params)
	if err != nil {
		return nil, err
	}

	codes := &CodeResponse{}
	resp.HydrateResponseCommon(&codes.ResponseCommon)
	codes.Data.Codes = resp.Data.(*ManyCodes).Codes

	return codes, nil
}

// RedeemEntitlementCode
// Per https://dev.twitch.tv/docs/api/reference/#redeem-code
// Access is controlled via an app access token on the calling service. The client ID associated with the app access token must be approved by Twitch.
// Callers with an app access token are authorized to redeem codes on behalf of any Twitch user account.
func (c *Client) RedeemEntitlementCode(params *CodesParams) (*CodeResponse, error) {
	resp, err := c.post("/entitlements/code", &ManyCodes{}, params)
	if err != nil {
		return nil, err
	}

	codes := &CodeResponse{}
	resp.HydrateResponseCommon(&codes.ResponseCommon)
	codes.Data.Codes = resp.Data.(*ManyCodes).Codes

	return codes, nil
}
