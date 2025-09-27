import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Loader2, Wallet, TrendingUp, Shield, RefreshCw, Link } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { useWallet } from "@/hooks/useWallet";
import AnimatedBackground from "@/components/AnimatedBackground";
import PortfolioOverview from "@/components/PortfolioOverview";
import TokenHoldings from "@/components/TokenHoldings";
import DeFiPositions from "@/components/DeFiPositions";
import RiskAnalysis from "@/components/RiskAnalysis";
import Recommendations from "@/components/Recommendations";

interface AnalysisData {
  recommended_tokens: string[];
  risk_score: number;
  reasoning: string[];
  token_balances: {
    total_balance_usd: number;
    by_token: Array<{
      token_address: string;
      symbol: string;
      name: string;
      decimals: number;
      price: number;
      balance: number;
      balance_usd: number;
      balance_raw: string;
      network: {
        name: string;
        slug: string;
      };
      img_url_v2?: string;
    }>;
  };
  app_balances: {
    by_app: Array<{
      app: {
        display_name: string;
        slug: string;
      };
      network: {
        name: string;
        slug: string;
      };
      balances: Array<{
        address: string;
        balance_usd: number;
        tokens: any[];
        display_props: {
          label: string;
          images: string[];
        };
      }>;
    }>;
  };
}

const Index = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<AnalysisData | null>(null);
  const [error, setError] = useState<string | null>(null);
  const { toast } = useToast();
  const { address, isConnected, isConnecting, connectWallet, disconnectWallet, error: walletError } = useWallet();

  const analyzeWallet = async () => {
    if (!address) {
      toast({
        title: "Error",
        description: "Please connect your wallet first",
        variant: "destructive",
      });
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`http://localhost:8080/analyze?address=${encodeURIComponent(address)}`);
      
      if (!response.ok) {
        throw new Error(`Analysis failed: ${response.statusText}`);
      }

      const analysisData = await response.json();
      setData(analysisData);
      toast({
        title: "Analysis Complete",
        description: "Wallet analysis has been successfully completed",
      });
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to analyze wallet";
      setError(errorMessage);
      toast({
        title: "Analysis Failed",
        description: errorMessage,
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const refreshAnalysis = () => {
    if (address) {
      analyzeWallet();
    }
  };

  return (
    <div className="min-h-screen bg-background relative">
      <AnimatedBackground />
      {/* Header */}
      <div className="bg-gradient-primary border-b relative z-10">
        <div className="container mx-auto px-4 py-8">
          <div className="flex items-center gap-3 mb-6">
            <Wallet className="h-8 w-8 text-primary-foreground" />
            <h1 className="text-3xl font-bold text-primary-foreground">DeFi Vision</h1>
          </div>

          {/* Wallet Connection */}
          <Card className="bg-card/95 backdrop-blur-sm shadow-card">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Wallet className="h-5 w-5" />
                {isConnected ? 'Connected Wallet' : 'Connect Wallet'}
              </CardTitle>
            </CardHeader>
            <CardContent>
              {!isConnected ? (
                <div className="text-center space-y-4">
                  <p className="text-muted-foreground">Connect your MetaMask wallet to analyze your portfolio</p>
                  <Button onClick={connectWallet} disabled={isConnecting} className="px-8">
                    {isConnecting ? (
                      <>
                        <Loader2 className="h-4 w-4 animate-spin mr-2" />
                        Connecting...
                      </>
                    ) : (
                      <>
                        <Link className="h-4 w-4 mr-2" />
                        Connect MetaMask
                      </>
                    )}
                  </Button>
                  {walletError && (
                    <p className="text-sm text-destructive">{walletError}</p>
                  )}
                </div>
              ) : (
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-muted-foreground">Connected Address</p>
                      <p className="font-mono text-sm">{address?.slice(0, 6)}...{address?.slice(-4)}</p>
                    </div>
                    <Button variant="outline" onClick={disconnectWallet} size="sm">
                      Disconnect
                    </Button>
                  </div>
                  <div className="flex gap-3">
                    <Button onClick={analyzeWallet} disabled={loading} className="flex-1">
                      {loading ? (
                        <>
                          <Loader2 className="h-4 w-4 animate-spin mr-2" />
                          Analyzing
                        </>
                      ) : (
                        'Analyze Portfolio'
                      )}
                    </Button>
                    {data && (
                      <Button variant="outline" onClick={refreshAnalysis} disabled={loading}>
                        <RefreshCw className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Product Info Cards */}
      {!isConnected && (
        <div className="container mx-auto px-4 py-8 relative z-10">
          {/* Large Feature Card */}
          <div className="mb-8">
            <Card className="bg-card/95 backdrop-blur-sm shadow-lg border-primary/20">
              <CardContent className="p-8 text-center">
                <h2 className="text-2xl font-bold mb-4 text-foreground">Advanced DeFi Portfolio Intelligence</h2>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
                  <div className="space-y-2">
                    <div className="text-3xl font-bold text-primary">10K+</div>
                    <div className="text-sm font-medium text-foreground">Token Holdings Tracked</div>
                    <div className="text-xs text-muted-foreground">Comprehensive token analysis</div>
                  </div>
                  <div className="space-y-2">
                    <div className="text-3xl font-bold text-primary">500+</div>
                    <div className="text-sm font-medium text-foreground">DEX & DeFi Protocols</div>
                    <div className="text-xs text-muted-foreground">Cross-protocol position tracking</div>
                  </div>
                  <div className="space-y-2">
                    <div className="text-3xl font-bold text-primary">50+</div>
                    <div className="text-sm font-medium text-foreground">Blockchain Networks</div>
                    <div className="text-xs text-muted-foreground">Multi-chain ecosystem support</div>
                  </div>
                </div>
                <p className="text-foreground/80 max-w-2xl mx-auto">
                  Real-time analysis of your complete DeFi footprint across all major chains and protocols
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Feature Cards Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card className="bg-card/95 backdrop-blur-sm shadow-card border-border/50">
              <CardContent className="p-6">
                <h3 className="font-semibold mb-2 text-foreground">Risk-Aware Portfolio Assistant</h3>
                <p className="text-sm text-muted-foreground">Provides comprehensive DeFi portfolio analysis that works directly with Ethereum addresses.</p>
              </CardContent>
            </Card>
            
            <Card className="bg-card/95 backdrop-blur-sm shadow-card border-border/50">
              <CardContent className="p-6">
                <h3 className="font-semibold mb-2 text-foreground">AI Risk Advisor</h3>
                <p className="text-sm text-muted-foreground">Analyzes wallet balances, DEX LP positions, and futures positions with AI-powered insights.</p>
              </CardContent>
            </Card>
            
            <Card className="bg-card/95 backdrop-blur-sm shadow-card border-border/50">
              <CardContent className="p-6">
                <h3 className="font-semibold mb-2 text-foreground">Smart Recommendations</h3>
                <p className="text-sm text-muted-foreground">Returns risk scores, recommended tokens for diversification, and detailed reasoning.</p>
              </CardContent>
            </Card>
            
            <Card className="bg-card/95 backdrop-blur-sm shadow-card border-border/50">
              <CardContent className="p-6">
                <h3 className="font-semibold mb-2 text-foreground">Backend Integration</h3>
                <p className="text-sm text-muted-foreground">Seamlessly integrates with wallets, portfolio dashboards, and dApps for execution-ready insights.</p>
              </CardContent>
            </Card>
          </div>
        </div>
      )}

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8 relative z-10">
        {error && (
          <Card className="mb-6 border-destructive bg-destructive/5">
            <CardContent className="pt-6">
              <div className="flex items-center gap-2 text-destructive">
                <Shield className="h-5 w-5" />
                <span className="font-medium">Analysis Error</span>
              </div>
              <p className="mt-2 text-sm text-muted-foreground">{error}</p>
            </CardContent>
          </Card>
        )}

        {loading && (
          <div className="flex items-center justify-center py-12">
            <div className="text-center">
              <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-primary" />
              <p className="text-muted-foreground">Analyzing wallet portfolio...</p>
            </div>
          </div>
        )}

        {data && !loading && (
          <div className="space-y-6">
            {/* Portfolio Overview */}
            <PortfolioOverview 
              totalValue={data.token_balances.total_balance_usd}
              riskScore={data.risk_score}
              tokenCount={data.token_balances.by_token.length}
              defiProtocolCount={data.app_balances.by_app.length}
            />

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Risk Analysis */}
              <RiskAnalysis 
                riskScore={data.risk_score}
                reasoning={data.reasoning}
              />

              {/* Recommendations */}
              <Recommendations 
                recommendations={data.recommended_tokens}
              />
            </div>

            {/* Token Holdings */}
            <TokenHoldings tokens={data.token_balances.by_token} />

            {/* DeFi Positions */}
            <DeFiPositions positions={data.app_balances.by_app} />
          </div>
        )}

        {!data && !loading && !error && (
          <div className="text-center py-12">
            <TrendingUp className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
            <h3 className="text-lg font-semibold mb-2">Ready to Analyze</h3>
            <p className="text-muted-foreground">Enter a wallet address above to get started with portfolio analysis</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Index;