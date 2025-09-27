package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Server struct {
	// No dependencies needed for position fetching
}

// Simple response structure for positions
type PositionsResponse struct {
	Address   string            `json:"address"`
	Positions []GraphQLPosition `json:"positions"`
}

// GraphQL request and response structures
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data struct {
		Positions []GraphQLPosition `json:"positions"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

type GraphQLPosition struct {
	ID                  string `json:"id"`
	Owner               string `json:"owner"`
	Liquidity           string `json:"liquidity"`
	DepositedToken0     string `json:"depositedToken0"`
	DepositedToken1     string `json:"depositedToken1"`
	WithdrawnToken0     string `json:"withdrawnToken0"`
	WithdrawnToken1     string `json:"withdrawnToken1"`
	CollectedFeesToken0 string `json:"collectedFeesToken0"`
	CollectedFeesToken1 string `json:"collectedFeesToken1"`
	Pool                struct {
		ID     string `json:"id"`
		Token0 struct {
			ID       string `json:"id"`
			Symbol   string `json:"symbol"`
			Decimals string `json:"decimals"`
		} `json:"token0"`
		Token1 struct {
			ID       string `json:"id"`
			Symbol   string `json:"symbol"`
			Decimals string `json:"decimals"`
		} `json:"token1"`
		FeeTier string `json:"feeTier"`
	} `json:"pool"`
}

func NewServer() *Server {
	return &Server{}
}

// AnalyzeRequest is the request payload for the /analyze endpoint
type AnalyzeRequest struct {
	Address string `json:"address"`
}

// AnalyzeResponse is the simplified ASI:One analysis response we return
type AnalyzeResponse struct {
	Address         string `json:"address"`
	Recommendations string `json:"recommendations"`
}

func (s *Server) GetPositions(w http.ResponseWriter, r *http.Request) {
	// Parse Ethereum address from query parameters
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address parameter is required", http.StatusBadRequest)
		return
	}

	// Validate Ethereum address format (basic validation)
	if len(address) != 42 || address[:2] != "0x" {
		http.Error(w, "invalid Ethereum address format", http.StatusBadRequest)
		return
	}

	// Fetch position data from The Graph
	positions, err := s.fetchPositionsFromGraph(address)
	if err != nil {
		http.Error(w, "failed to fetch position data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return positions response
	response := PositionsResponse{
		Address:   address,
		Positions: positions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// fetchPositionsFromGraph calls The Graph API to fetch position data
func (s *Server) fetchPositionsFromGraph(address string) ([]GraphQLPosition, error) {
	// GraphQL query
	query := `query PositionsByOwner($owner: String!) {
		positions(first: 100, where: { owner: $owner }) {
			id
			owner
			pool {
				id
				token0 {
					id
					symbol
					decimals
				}
				token1 {
					id
					symbol
					decimals
				}
				feeTier
			}
			liquidity
			depositedToken0
			depositedToken1
			withdrawnToken0
			withdrawnToken1
			collectedFeesToken0
			collectedFeesToken1
		}
	}`

	// Create GraphQL request
	request := GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"owner": address,
		},
	}

	// Marshal request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://gateway.thegraph.com/api/subgraphs/id/5zvR82QoaXYFyDEKLZ9t6v9adgnptxYpKpSbxtgVENFV", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer a89fc492b43ed0cdd300d82e82845b90")

	// Make HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse GraphQL response
	var graphQLResp GraphQLResponse
	if err := json.Unmarshal(body, &graphQLResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GraphQL response: %w", err)
	}

	// Check for GraphQL errors
	if len(graphQLResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", graphQLResp.Errors)
	}

	return graphQLResp.Data.Positions, nil
}

// transformPositionsToRiskRequest converts GraphQL positions to RiskRequest format
func (s *Server) transformPositionsToRiskRequest(address string, positions []GraphQLPosition) RiskRequest {
	// Initialize all arrays as empty slices
	tokenHoldings := make([]TokenHolding, 0)
	dexPositions := make([]DexPosition, 0)
	futuresPositions := make([]FuturesPosition, 0)

	// For now, we'll create sample data based on the positions
	// In a real implementation, you'd fetch actual token balances and USD values

	for _, pos := range positions {
		// Convert GraphQL position to DexPosition
		// Note: This is a simplified conversion - you'd need to calculate actual USD values
		dexPos := DexPosition{
			ID:        pos.ID,
			Pool:      fmt.Sprintf("%s-%s", pos.Pool.Token0.Symbol, pos.Pool.Token1.Symbol),
			Token0:    pos.Pool.Token0.Symbol,
			Token1:    pos.Pool.Token1.Symbol,
			Liquidity: 1.0,    // Simplified - would need actual calculation
			USDValue:  1000.0, // Simplified - would need price lookup
		}
		dexPositions = append(dexPositions, dexPos)

		// Create token holdings from the position tokens
		token0Holding := TokenHolding{
			Symbol:   pos.Pool.Token0.Symbol,
			Amount:   1.0,   // Simplified
			USDValue: 500.0, // Simplified
		}
		tokenHoldings = append(tokenHoldings, token0Holding)

		token1Holding := TokenHolding{
			Symbol:   pos.Pool.Token1.Symbol,
			Amount:   1.0,   // Simplified
			USDValue: 500.0, // Simplified
		}
		tokenHoldings = append(tokenHoldings, token1Holding)
	}

	return RiskRequest{
		Address:          address,
		TokenHoldings:    tokenHoldings,
		DexPositions:     dexPositions,
		FuturesPositions: futuresPositions,
	}
}

// callRiskAdvisorAPI calls the risk advisor API with the given risk request
func (s *Server) callRiskAdvisorAPI(riskRequest RiskRequest) (*RiskResponse, error) {
	// Marshal the risk request to JSON
	requestBody, err := json.Marshal(riskRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal risk request: %w", err)
	}

	// Create HTTP request to the risk advisor API
	req, err := http.NewRequest("POST", "http://localhost:8000/api/analyze", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("risk advisor API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var riskResponse RiskResponse
	if err := json.Unmarshal(body, &riskResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal risk response: %w", err)
	}

	return &riskResponse, nil
}

// Risk Advisor API request/response types
type TokenHolding struct {
	Symbol   string  `json:"symbol"`
	Amount   float64 `json:"amount"`
	USDValue float64 `json:"usd_value"`
}

type DexPosition struct {
	ID        string  `json:"id"`
	Pool      string  `json:"pool"`
	Token0    string  `json:"token0"`
	Token1    string  `json:"token1"`
	Liquidity float64 `json:"liquidity"`
	USDValue  float64 `json:"usd_value"`
}

type FuturesPosition struct {
	ID       string  `json:"id"`
	Market   string  `json:"market"`
	Amount   float64 `json:"amount"`
	Leverage float64 `json:"leverage"`
	Status   string  `json:"status"`
	USDValue float64 `json:"usd_value"`
}

type RiskRequest struct {
	Address          string            `json:"address"`
	TokenHoldings    []TokenHolding    `json:"token_holdings"`
	DexPositions     []DexPosition     `json:"dex_positions"`
	FuturesPositions []FuturesPosition `json:"futures_positions"`
}

type RiskResponse struct {
	RecommendedTokens []string `json:"recommended_tokens"`
	RiskScore         float64  `json:"risk_score"`
	Reasoning         []string `json:"reasoning"`
}

// AnalyzeWithASI fetches positions and calls the risk advisor API for analysis
func (s *Server) AnalyzeWithASI(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address parameter is required", http.StatusBadRequest)
		return
	}
	if len(address) != 42 || address[:2] != "0x" {
		http.Error(w, "invalid Ethereum address format", http.StatusBadRequest)
		return
	}

	// Fetch position data from The Graph
	positions, err := s.fetchPositionsFromGraph(address)
	if err != nil {
		http.Error(w, "failed to fetch position data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Transform positions to RiskRequest format
	riskRequest := s.transformPositionsToRiskRequest(address, positions)

	// Call the risk advisor API
	riskResponse, err := s.callRiskAdvisorAPI(riskRequest)
	if err != nil {
		http.Error(w, "failed to get risk analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the risk analysis response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(riskResponse)
}
