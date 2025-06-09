package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

// baseURL/interface/method/version?parameters
const location = "https://api.steampowered.com/IGameServersService/"
const version = "v1"

// Steam base type
type Steam struct {
	apiKey string
}

// New initializes a new Steam instance with given api key
func New(apiKey string) *Steam {
	if apiKey == "" {
		log.Fatal(fmt.Errorf("missing steam api key"))
	}

	return &Steam{
		apiKey: apiKey,
	}
}

// Steam returns a JSON { response: } object, which wraps all return values.
type steamResponse struct {
	Response json.RawMessage `json:"response"`
}

// FML
type serversResponse struct {
	Servers []Account `json:"servers"`
}

// Account is an abstraction around LoginToken, for use with SteamCMD dedicated servers.
type Account struct {
	SteamID    string `json:"steamid,omitempty"`
	AppID      uint16 `json:"appid,omitempty"`
	LoginToken string `json:"login_token,omitempty"`
	Memo       string `json:"memo,omitempty"`
	IsDeleted  bool   `json:"is_deleted,omitempty"`
	IsExpired  bool   `json:"is_expired,omitempty"`
	LastLogon  int    `json:"rt_last_logon,omitempty"`
}

// Remove the { response: data } wrapper, and return inner json as byte array.
func unwrapResponse(response *[]byte) error {
	resp := steamResponse{}
	if err := json.Unmarshal(*response, &resp); err != nil {
		return err
	}
	*response = ([]byte)(resp.Response)
	return nil
}

// Wraps requests for Steam Web API, to generalize insertion of API key,
// and handling of Response Header.
func (s *Steam) querySteam(command string, method string, params map[string]string) (data []byte, err error) {
	// Prep request
	req, err := http.NewRequest(method, location+command+"/"+version, nil)
	if err != nil {
		return nil, err
	}

	// Add API Key and extra parameters
	q := url.Values{}
	q.Add("key", s.apiKey)
	for key, value := range params {
		q.Add(key, value)
	}
	// Encode parameters and append them to the url
	req.URL.RawQuery = q.Encode()

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Drop if Error Header present
	if respErrState := resp.Header.Get("X-error_message"); respErrState != "" {
		return nil, errors.New(respErrState)
	}

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Remove wrapper
	if err = unwrapResponse(&body); err != nil {
		return nil, err
	}

	return body, nil
}

// CreateAccount creates an account with a token, for use with SteamCMD dedicated servers.
func (s *Steam) CreateAccount(appID int, memo string) (account Account, err error) {
	// Build query string
	params := make(map[string]string)
	params["appid"] = strconv.Itoa(appID)
	params["memo"] = memo

	// Execute request
	data, err := s.querySteam("CreateAccount", "POST", params)
	if err != nil {
		return account, err
	}

	// Decode response
	if err := json.Unmarshal(data, &account); err != nil {
		return account, err
	}

	return account, nil
}

// GetAccountList returns a list of all accounts.
func (s *Steam) GetAccountList() (accounts []Account, err error) {
	data, err := s.querySteam("GetAccountList", "GET", nil)
	if err != nil {
		return accounts, err
	}

	var list serversResponse

	if err := json.Unmarshal(data, &list); err != nil {
		return accounts, err
	}

	accounts = list.Servers

	return accounts, nil
}

// DeleteAccount deletes an account, immediately expiring its LoginToken.
func (s *Steam) DeleteAccount(steamID string) (err error) {
	params := make(map[string]string)
	params["steamid"] = steamID

	_, err = s.querySteam("DeleteAccount", "POST", params)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAllAccounts deletes all accounts registered by the user.
func (s *Steam) DeleteAllAccounts() (err error) {
	accounts, err := s.GetAccountList()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err = s.DeleteAccount(account.SteamID)
		if err != nil {
			slog.Error("Error deleting account", "error", err)
		}
	}

	return nil
}

// ResetLoginToken generates a new LoginToken on an existing steamID.
func (s *Steam) ResetLoginToken(steamID string) (account Account, err error) {
	params := make(map[string]string)
	params["steamID"] = steamID

	data, err := s.querySteam("ResetLoginToken", "POST", params)
	if err != nil {
		return account, err
	}

	if err := json.Unmarshal(data, &account); err != nil {
		return account, err
	}

	return account, nil
}
