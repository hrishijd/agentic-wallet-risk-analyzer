package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

	// Fetch portfolio data from Zapper API
	riskRequest, err := s.fetchPortfolioFromZapper(address)
	if err != nil {
		http.Error(w, "failed to fetch portfolio data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return positions response (converted from new structure for backward compatibility)
	// Note: This is a simplified conversion for backward compatibility
	var positions []GraphQLPosition
	for _, appBalance := range riskRequest.AppBalances.ByApp {
		for _, contractPos := range appBalance.Balances {
			// Create a simplified GraphQL position for backward compatibility
			if len(contractPos.Tokens) >= 2 {
				pos := GraphQLPosition{
					ID:    contractPos.Address,
					Owner: address,
					Pool: struct {
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
					}{
						ID: contractPos.Address,
						Token0: struct {
							ID       string `json:"id"`
							Symbol   string `json:"symbol"`
							Decimals string `json:"decimals"`
						}{
							ID:       contractPos.Tokens[0].Token.TokenAddress,
							Symbol:   contractPos.Tokens[0].Token.Symbol,
							Decimals: fmt.Sprintf("%.0f", contractPos.Tokens[0].Token.Decimals),
						},
						Token1: struct {
							ID       string `json:"id"`
							Symbol   string `json:"symbol"`
							Decimals string `json:"decimals"`
						}{
							ID:       contractPos.Tokens[1].Token.TokenAddress,
							Symbol:   contractPos.Tokens[1].Token.Symbol,
							Decimals: fmt.Sprintf("%.0f", contractPos.Tokens[1].Token.Decimals),
						},
						FeeTier: "3000", // Default fee tier
					},
					Liquidity: "1.0",
				}
				positions = append(positions, pos)
			}
		}
	}

	response := PositionsResponse{
		Address:   address,
		Positions: positions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// fetchPortfolioFromZapper calls Zapper API to fetch portfolio data
func (s *Server) fetchPortfolioFromZapper(address string) (*RiskRequest, error) {
	// Fetch token balances
	tokenBalances, err := s.fetchTokenBalances(address)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token balances: %w", err)
	}

	// Fetch app balances
	appBalances, err := s.fetchAppBalances(address)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch app balances: %w", err)
	}

	// Combine both responses into RiskRequest
	riskRequest := &RiskRequest{
		Address:       address,
		TokenBalances: *tokenBalances,
		AppBalances:   *appBalances,
	}

	return riskRequest, nil
}

// fetchTokenBalances fetches token balances using the exact curl query
func (s *Server) fetchTokenBalances(address string) (*TokenBalances, error) {
	query := `query TokenBalances($addresses: [Address!]!, $first: Int, $chainIds: [Int!]) {
		portfolioV2(addresses: $addresses, chainIds: $chainIds) {
			tokenBalances {
				totalBalanceUSD
				byToken(first: $first) {
					totalCount
					edges {
						node {
							name
							symbol
							price
							tokenAddress
							imgUrlV2
							decimals
							balanceRaw
							balance
							balanceUSD
							onchainMarketData {
								priceChange24h
								marketCap
							}
						}
					}
				}
			}
		}
	}`

	request := GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"addresses": []string{address},
			"first":     5,
			"chainIds":  []int{8453}, // Base chain
		},
	}

	body, err := s.makeZapperRequest(request)
	if err != nil {
		return nil, err
	}

	tokenBalances, err := s.parseTokenBalancesResponse(address, body)
	if err != nil {
		return nil, err
	}

	return tokenBalances, nil
}

// fetchAppBalances fetches app balances using the exact curl query
func (s *Server) fetchAppBalances(address string) (*AppBalances, error) {
	query := `query AppBalances($addresses: [Address!]!, $first: Int = 10) {
		portfolioV2(addresses: $addresses) {
			appBalances {
				totalBalanceUSD
				byApp(first: $first) {
					totalCount
					edges {
						node {
							balanceUSD
							app {
								displayName
								imgUrl
								description
								category {
									name
								}
							}
							network {
								name
								chainId
							}
							positionBalances(first: 10) {
								edges {
									node {
										... on AppTokenPositionBalance {
											type
											symbol
											balance
											balanceUSD
											price
											groupLabel
											displayProps {
												label
												images
											}
										}
										... on ContractPositionBalance {
											type
											balanceUSD
											groupLabel
											tokens {
												metaType
												token {
													... on BaseTokenPositionBalance {
														symbol
														balance
														balanceUSD
													}
												}
											}
											displayProps {
												label
												images
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`

	request := GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"addresses": []string{address},
			"first":     5,
		},
	}

	body, err := s.makeZapperRequest(request)
	if err != nil {
		return nil, err
	}

	appBalances, err := s.parseAppBalancesResponse(address, body)
	if err != nil {
		return nil, err
	}

	return appBalances, nil
}

// makeZapperRequest makes a request to Zapper API with the provided query
func (s *Server) makeZapperRequest(request GraphQLRequest) ([]byte, error) {
	// Marshal request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Zapper request: %w", err)
	}

	// Create HTTP request to Zapper API
	req, err := http.NewRequest("POST", "https://public.zapper.xyz/graphql", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Get Zapper API key from environment variable
	apiKey := os.Getenv("ZAPPER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ZAPPER_API_KEY environment variable is not set")
	}
	req.Header.Set("x-zapper-api-key", apiKey)

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

	return body, nil
}

// parseTokenBalancesResponse parses the token balances response
func (s *Server) parseTokenBalancesResponse(address string, responseBody []byte) (*TokenBalances, error) {
	// Print the raw JSON response before parsing
	fmt.Printf("Raw Token Balances Response for address %s:\n", address)
	fmt.Printf("%s\n\n", string(responseBody))

	var resp struct {
		Data struct {
			PortfolioV2 struct {
				TokenBalances struct {
					TotalBalanceUSD float64 `json:"totalBalanceUSD"`
					ByToken         struct {
						Edges []struct {
							Node struct {
								Name              string  `json:"name"`
								Symbol            string  `json:"symbol"`
								Price             float64 `json:"price"`
								TokenAddress      string  `json:"tokenAddress"`
								ImgURLV2          *string `json:"imgUrlV2"`
								Decimals          float64 `json:"decimals"`
								BalanceRaw        string  `json:"balanceRaw"`
								Balance           float64 `json:"balance"`
								BalanceUSD        float64 `json:"balanceUSD"`
								OnchainMarketData *struct {
									PriceChange24h *float64 `json:"priceChange24h"`
									MarketCap      *float64 `json:"marketCap"`
								} `json:"onchainMarketData"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"byToken"`
				} `json:"tokenBalances"`
			} `json:"portfolioV2"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token balances response: %w", err)
	}

	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("token balances API errors: %v", resp.Errors)
	}

	// Convert to our TokenBalances structure
	tokenBalances := make([]TokenBalance, 0)
	for _, edge := range resp.Data.PortfolioV2.TokenBalances.ByToken.Edges {
		node := edge.Node
		tokenBalance := TokenBalance{
			TokenAddress: node.TokenAddress,
			Symbol:       node.Symbol,
			Name:         node.Name,
			Decimals:     node.Decimals,
			Price:        node.Price,
			Balance:      node.Balance,
			BalanceUSD:   node.BalanceUSD,
			BalanceRaw:   node.BalanceRaw,
			Network: Network{
				Name: "Base", // Since we're using chainIds: [8453]
				Slug: "base",
			},
			ImgURLV2: node.ImgURLV2,
		}
		tokenBalances = append(tokenBalances, tokenBalance)
	}

	return &TokenBalances{
		TotalBalanceUSD: resp.Data.PortfolioV2.TokenBalances.TotalBalanceUSD,
		ByToken:         tokenBalances,
	}, nil
}

// parseAppBalancesResponse parses the app balances response
func (s *Server) parseAppBalancesResponse(address string, responseBody []byte) (*AppBalances, error) {
	// Print the raw JSON response before parsing
	fmt.Printf("Raw App Balances Response for address %s:\n", address)
	fmt.Printf("%s\n\n", string(responseBody))

	var resp struct {
		Data struct {
			PortfolioV2 struct {
				AppBalances struct {
					ByApp struct {
						Edges []struct {
							Node struct {
								BalanceUSD float64 `json:"balanceUSD"`
								App        struct {
									DisplayName string `json:"displayName"`
									ImgURL      string `json:"imgUrl"`
									Description string `json:"description"`
									Category    struct {
										Name string `json:"name"`
									} `json:"category"`
								} `json:"app"`
								Network struct {
									Name    string `json:"name"`
									ChainID int    `json:"chainId"`
								} `json:"network"`
								PositionBalances struct {
									Edges []struct {
										Node struct {
											Type       string   `json:"type"`
											Symbol     *string  `json:"symbol"`
											Balance    *string  `json:"balance"` // Can be string or number
											BalanceUSD float64  `json:"balanceUSD"`
											Price      *float64 `json:"price"`      // Changed back to float64 as it's a number in response
											GroupLabel *string  `json:"groupLabel"` // Can be null
											Tokens     *[]struct {
												MetaType string `json:"metaType"`
												Token    struct {
													Symbol     string `json:"symbol"`
													Balance    string `json:"balance"`    // String in response
													BalanceUSD string `json:"balanceUSD"` // String in response
												} `json:"symbol"`
											} `json:"tokens"`
											DisplayProps struct {
												Label  string   `json:"label"`
												Images []string `json:"images"`
											} `json:"displayProps"`
										} `json:"node"`
									} `json:"edges"`
								} `json:"positionBalances"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"byApp"`
				} `json:"appBalances"`
			} `json:"portfolioV2"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal app balances response: %w", err)
	}

	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("app balances API errors: %v", resp.Errors)
	}

	// Convert to our AppBalances structure
	appBalances := make([]AppBalance, 0)
	for _, appEdge := range resp.Data.PortfolioV2.AppBalances.ByApp.Edges {
		appNode := appEdge.Node

		contractPositions := make([]ContractPosition, 0)
		for _, posEdge := range appNode.PositionBalances.Edges {
			posNode := posEdge.Node

			tokenPositions := make([]TokenPosition, 0)
			if posNode.Tokens != nil {
				for _, token := range *posNode.Tokens {
					// Parse string values to float64
					balance, _ := strconv.ParseFloat(token.Token.Balance, 64)
					balanceUSD, _ := strconv.ParseFloat(token.Token.BalanceUSD, 64)

					tokenPos := TokenPosition{
						MetaType: token.MetaType,
						Token: TokenBalance{
							Symbol:     token.Token.Symbol,
							Balance:    balance,
							BalanceUSD: balanceUSD,
							// Fill in other required fields with defaults
							TokenAddress: "unknown",
							Name:         token.Token.Symbol,
							Decimals:     18,
							Price:        0,
							BalanceRaw:   token.Token.Balance, // Use the original string value
							Network: Network{
								Name: appNode.Network.Name,
								Slug: "base", // Default to base
							},
						},
					}
					tokenPositions = append(tokenPositions, tokenPos)
				}
			}

			// Handle nullable GroupLabel
			address := "unknown"
			if posNode.GroupLabel != nil {
				address = *posNode.GroupLabel
			}

			contractPos := ContractPosition{
				Address:    address,
				BalanceUSD: posNode.BalanceUSD,
				Tokens:     tokenPositions,
				DisplayProps: DisplayProps{
					Label:  posNode.DisplayProps.Label,
					Images: posNode.DisplayProps.Images,
				},
			}
			contractPositions = append(contractPositions, contractPos)
		}

		appBalance := AppBalance{
			App: App{
				DisplayName: appNode.App.DisplayName,
				Slug:        appNode.App.Category.Name, // Use category as slug
			},
			Network: Network{
				Name: appNode.Network.Name,
				Slug: "base", // Default to base
			},
			Balances: contractPositions,
		}
		appBalances = append(appBalances, appBalance)
	}

	return &AppBalances{
		ByApp: appBalances,
	}, nil
}

// Zapper GraphQL response structures
type ZapperTokenNode struct {
	Symbol       string  `json:"symbol"`
	TokenAddress string  `json:"tokenAddress"`
	Balance      float64 `json:"balance"`
	BalanceUSD   float64 `json:"balanceUSD"`
	Price        float64 `json:"price"`
	ImgURLV2     *string `json:"imgUrlV2"`
	Name         string  `json:"name"`
	Decimals     float64 `json:"decimals"`
	BalanceRaw   string  `json:"balanceRaw"`
	Network      struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"network"`
}

type ZapperTokenEdge struct {
	Node ZapperTokenNode `json:"node"`
}

type ZapperTokenConnection struct {
	Edges []ZapperTokenEdge `json:"edges"`
}

type ZapperTokenBalances struct {
	TotalBalanceUSD float64               `json:"totalBalanceUSD"`
	ByToken         ZapperTokenConnection `json:"byToken"`
}

type ZapperAppToken struct {
	MetaType string `json:"metaType"`
	Token    struct {
		Symbol       string  `json:"symbol"`
		TokenAddress string  `json:"tokenAddress"`
		Balance      float64 `json:"balance"`
		BalanceUSD   float64 `json:"balanceUSD"`
		Price        float64 `json:"price"`
		ImgURLV2     *string `json:"imgUrlV2"`
		Name         string  `json:"name"`
		Decimals     float64 `json:"decimals"`
		BalanceRaw   string  `json:"balanceRaw"`
		Network      struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"network"`
	} `json:"token"`
}

type ZapperContractNode struct {
	Address      string           `json:"address"`
	BalanceUSD   float64          `json:"balanceUSD"`
	Tokens       []ZapperAppToken `json:"tokens"`
	DisplayProps struct {
		Label  string   `json:"label"`
		Images []string `json:"images"`
	} `json:"displayProps"`
}

type ZapperContractEdge struct {
	Node ZapperContractNode `json:"node"`
}

type ZapperContractConnection struct {
	Edges []ZapperContractEdge `json:"edges"`
}

type ZapperAppNode struct {
	App struct {
		DisplayName string `json:"displayName"`
		Slug        string `json:"slug"`
	} `json:"app"`
	Network struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"network"`
	Balances ZapperContractConnection `json:"balances"`
}

type ZapperAppEdge struct {
	Node ZapperAppNode `json:"node"`
}

type ZapperAppConnection struct {
	Edges []ZapperAppEdge `json:"edges"`
}

type ZapperAppBalances struct {
	ByApp ZapperAppConnection `json:"byApp"`
}

type ZapperPortfolioV2 struct {
	TokenBalances ZapperTokenBalances `json:"tokenBalances"`
	AppBalances   ZapperAppBalances   `json:"appBalances"`
}

type ZapperResponse struct {
	Data struct {
		PortfolioV2 ZapperPortfolioV2 `json:"portfolioV2"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// parseZapperResponse converts Zapper API response to RiskRequest format
func (s *Server) parseZapperResponse(address string, responseBody []byte) (*RiskRequest, error) {
	// Print the raw JSON response before parsing
	fmt.Printf("Raw Zapper API Response for address %s:\n", address)
	fmt.Printf("%s\n\n", string(responseBody))

	var zapperResp ZapperResponse
	if err := json.Unmarshal(responseBody, &zapperResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Zapper response: %w", err)
	}

	// Check for errors
	if len(zapperResp.Errors) > 0 {
		return nil, fmt.Errorf("Zapper API errors: %v", zapperResp.Errors)
	}

	portfolio := zapperResp.Data.PortfolioV2

	// Convert token balances
	tokenBalances := make([]TokenBalance, 0)
	for _, edge := range portfolio.TokenBalances.ByToken.Edges {
		node := edge.Node
		tokenBalance := TokenBalance{
			TokenAddress: node.TokenAddress,
			Symbol:       node.Symbol,
			Name:         node.Name,
			Decimals:     node.Decimals,
			Price:        node.Price,
			Balance:      node.Balance,
			BalanceUSD:   node.BalanceUSD,
			BalanceRaw:   node.BalanceRaw,
			Network: Network{
				Name: node.Network.Name,
				Slug: node.Network.Slug,
			},
			ImgURLV2: node.ImgURLV2,
		}
		tokenBalances = append(tokenBalances, tokenBalance)
	}

	// Convert app balances
	appBalances := make([]AppBalance, 0)
	for _, appEdge := range portfolio.AppBalances.ByApp.Edges {
		appNode := appEdge.Node

		contractPositions := make([]ContractPosition, 0)
		for _, contractEdge := range appNode.Balances.Edges {
			contractNode := contractEdge.Node

			tokenPositions := make([]TokenPosition, 0)
			for _, appToken := range contractNode.Tokens {
				tokenPos := TokenPosition{
					MetaType: appToken.MetaType,
					Token: TokenBalance{
						TokenAddress: appToken.Token.TokenAddress,
						Symbol:       appToken.Token.Symbol,
						Name:         appToken.Token.Name,
						Decimals:     appToken.Token.Decimals,
						Price:        appToken.Token.Price,
						Balance:      appToken.Token.Balance,
						BalanceUSD:   appToken.Token.BalanceUSD,
						BalanceRaw:   appToken.Token.BalanceRaw,
						Network: Network{
							Name: appToken.Token.Network.Name,
							Slug: appToken.Token.Network.Slug,
						},
						ImgURLV2: appToken.Token.ImgURLV2,
					},
				}
				tokenPositions = append(tokenPositions, tokenPos)
			}

			contractPos := ContractPosition{
				Address:    contractNode.Address,
				BalanceUSD: contractNode.BalanceUSD,
				Tokens:     tokenPositions,
				DisplayProps: DisplayProps{
					Label:  contractNode.DisplayProps.Label,
					Images: contractNode.DisplayProps.Images,
				},
			}
			contractPositions = append(contractPositions, contractPos)
		}

		appBalance := AppBalance{
			App: App{
				DisplayName: appNode.App.DisplayName,
				Slug:        appNode.App.Slug,
			},
			Network: Network{
				Name: appNode.Network.Name,
				Slug: appNode.Network.Slug,
			},
			Balances: contractPositions,
		}
		appBalances = append(appBalances, appBalance)
	}

	riskRequest := &RiskRequest{
		Address: address,
		TokenBalances: TokenBalances{
			TotalBalanceUSD: portfolio.TokenBalances.TotalBalanceUSD,
			ByToken:         tokenBalances,
		},
		AppBalances: AppBalances{
			ByApp: appBalances,
		},
	}

	return riskRequest, nil
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

// Zapper API and Risk Advisor types
type Network struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type TokenBalance struct {
	TokenAddress string  `json:"token_address"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	Decimals     float64 `json:"decimals"`
	Price        float64 `json:"price"`
	Balance      float64 `json:"balance"`
	BalanceUSD   float64 `json:"balance_usd"`
	BalanceRaw   string  `json:"balance_raw"`
	Network      Network `json:"network"`
	ImgURLV2     *string `json:"img_url_v2,omitempty"`
}

type TokenBalances struct {
	TotalBalanceUSD float64        `json:"total_balance_usd"`
	ByToken         []TokenBalance `json:"by_token"`
}

type App struct {
	DisplayName string `json:"display_name"`
	Slug        string `json:"slug"`
}

type TokenPosition struct {
	MetaType string       `json:"meta_type"`
	Token    TokenBalance `json:"token"`
}

type DisplayProps struct {
	Label  string   `json:"label"`
	Images []string `json:"images"`
}

type ContractPosition struct {
	Address      string          `json:"address"`
	BalanceUSD   float64         `json:"balance_usd"`
	Tokens       []TokenPosition `json:"tokens"`
	DisplayProps DisplayProps    `json:"display_props"`
}

type AppBalance struct {
	App      App                `json:"app"`
	Network  Network            `json:"network"`
	Balances []ContractPosition `json:"balances"`
}

type AppBalances struct {
	ByApp []AppBalance `json:"by_app"`
}

type RiskRequest struct {
	Address       string        `json:"address"`
	TokenBalances TokenBalances `json:"token_balances"`
	AppBalances   AppBalances   `json:"app_balances"`
}

type RiskResponse struct {
	RecommendedTokens []string `json:"recommended_tokens"`
	RiskScore         float64  `json:"risk_score"`
	Reasoning         []string `json:"reasoning"`
}

// AnalyzeWithASI fetches portfolio data from Zapper and calls the risk advisor API for analysis
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

	// Fetch portfolio data from Zapper API
	riskRequest, err := s.fetchPortfolioFromZapper(address)
	if err != nil {
		http.Error(w, "failed to fetch portfolio data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Call the risk advisor API
	riskResponse, err := s.callRiskAdvisorAPI(*riskRequest)
	if err != nil {
		http.Error(w, "failed to get risk analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the risk analysis response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(riskResponse)
}
