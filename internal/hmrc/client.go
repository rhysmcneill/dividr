package hmrc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rhysmcneill/dividr/internal/config"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
	fraudMeta  FraudMetadata
	cfg        *config.Config
}

func NewClient(isProd bool, token string, fraudMeta FraudMetadata, cfg *config.Config) *Client {
	baseURL := BaseURLSandbox
	env := cfg.AppEnv
	if env == "production" {
		baseURL = BaseURLProd
	}

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		token:      token,
		fraudMeta:  fraudMeta,
		cfg:        cfg,
	}
}

// Do wraps the standard HTTP call with Auth + Fraud Headers
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.hmrc.1.0+json")
	AddFraudHeaders(req, c.fraudMeta)
	return c.httpClient.Do(req)
}

// GetObligations fetches open deadlines (Story 2.6.5 requirement)
func (c *Client) GetObligations(mtdID string, fromDate string, toDate string) (*http.Response, error) {
	// Endpoint: /organisations/mtd/income-tax/{mtdId}/obligations
	url := fmt.Sprintf("%s/organisations/mtd/income-tax/%s/obligations?from=%s&to=%s&status=O", c.baseURL, mtdID, fromDate, toDate)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}
