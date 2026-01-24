package hmrc

import (
	"bytes" // Required for refilling the body
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/rhysmcneill/dividr/internal/config"
)

// BusinessDetailsResponse matches the union of possible HMRC responses
type BusinessDetailsResponse struct {
	// Format A: The one we want (contains the MTD ID)
	TaxPayerDisplayResponse []struct {
		MtdID string `json:"mtdId"`
	} `json:"taxPayerDisplayResponse"`

	// Format B: The one you received in logs (Income Sources only)
	ListOfBusinesses []struct {
		BusinessID string `json:"businessId"`
		Type       string `json:"typeOfBusiness"`
	} `json:"listOfBusinesses"`
}

// FetchMTDIdentifier swaps a NINO for an MTD ID
func (s *AuthService) FetchMTDIdentifier(accessToken string, nino string, cfg *config.Config) (string, error) {
	// 1. Determine Base URL
	baseURL := BaseURLSandbox
	if cfg.AppEnv == "production" {
		baseURL = BaseURLProd
	}

	// 2. Construct the URL
	// We use the /list endpoint as discovered
	endpoint, err := url.JoinPath(baseURL, "individuals", "business", "details", nino, "list")
	if err != nil {
		return "", err
	}

	// 3. Create Request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	// 4. Set Headers
	// Version 2.0 is usually required for this endpoint
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.hmrc.2.0+json")

	// 5. Execute
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			slog.Warn("failed to close response body", "error", cerr)
		}
	}()

	// --- DEBUG LOGGING START ---
	// Read the body for logging
	bodyBytes, _ := io.ReadAll(resp.Body)
	slog.Info("DEBUG: Business Details Response", "body", string(bodyBytes))

	// CRITICAL FIX: "Rewind" the body so the JSON decoder can read it again
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// --- DEBUG LOGGING END ---

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HMRC business details failed: status %d", resp.StatusCode)
	}

	// 6. Parse JSON
	var details BusinessDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %w", err)
	}

	// 7. Extract ID (Handle multiple formats)

	// Priority 1: Check for the explicit MTD ID (TaxPayerDisplayResponse)
	if len(details.TaxPayerDisplayResponse) > 0 {
		return details.TaxPayerDisplayResponse[0].MtdID, nil
	}

	// Priority 2: Check if we got the "ListOfBusinesses" format (Format B)
	// This format proves the user exists, but it DOES NOT contain the Global MTD ID (XPIT...).
	// It only contains specific Business IDs (XBIS...).
	if len(details.ListOfBusinesses) > 0 {
		// We found businesses, but not the MTD ID we need.
		// We return an error so your Handler's "Dev Bypass" can kick in and inject the XPIT ID.
		return "", fmt.Errorf("found valid business records (e.g. %s) but response is missing the Global MTD ID", details.ListOfBusinesses[0].BusinessID)
	}

	return "", fmt.Errorf("response parsed successfully but contained no business data")
}
