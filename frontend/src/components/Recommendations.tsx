import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Lightbulb, TrendingUp, Star } from "lucide-react";

interface RecommendationsProps {
  recommendations: string[];
}

const Recommendations = ({ recommendations }: RecommendationsProps) => {
  const getTokenIcon = (tokenSymbol: string) => {
    // Simple logic to assign icons based on token type
    if (tokenSymbol.includes('ETH')) return '⟠';
    if (tokenSymbol.includes('USDC') || tokenSymbol.includes('USDT')) return '💵';
    if (tokenSymbol.includes('UNI')) return '🦄';
    if (tokenSymbol.includes('BTC')) return '₿';
    return '🪙';
  };

  const getRecommendationType = (recommendation: string) => {
    if (recommendation.includes('through') || recommendation.includes('Aave') || recommendation.includes('Compound')) {
      return { type: 'DeFi Protocol', color: 'bg-blue-100 text-blue-800' };
    }
    if (recommendation.includes('ETH') || recommendation.includes('BTC')) {
      return { type: 'Blue Chip', color: 'bg-green-100 text-green-800' };
    }
    if (recommendation.includes('USDC') || recommendation.includes('USDT')) {
      return { type: 'Stablecoin', color: 'bg-gray-100 text-gray-800' };
    }
    return { type: 'Token', color: 'bg-purple-100 text-purple-800' };
  };

  return (
    <Card className="shadow-card">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Lightbulb className="h-5 w-5" />
          AI Recommendations
        </CardTitle>
      </CardHeader>
      <CardContent>
        {recommendations.length > 0 ? (
          <div className="space-y-4">
            <div className="p-3 bg-primary/5 border border-primary/20 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <Star className="h-4 w-4 text-primary" />
                <span className="text-sm font-medium text-primary">Portfolio Optimization</span>
              </div>
              <p className="text-sm text-muted-foreground">
                Based on your current portfolio composition and risk profile, consider these recommendations to improve diversification and reduce volatility.
              </p>
            </div>

            <div className="space-y-3">
              {recommendations.map((recommendation, index) => {
                const recType = getRecommendationType(recommendation);
                const mainToken = recommendation.split(' ')[0];
                
                return (
                  <div
                    key={index}
                    className="flex items-start gap-3 p-3 border rounded-lg hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex-shrink-0">
                      <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
                        <span className="text-lg">{getTokenIcon(recommendation)}</span>
                      </div>
                    </div>
                    
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <h4 className="font-semibold">{recommendation}</h4>
                        <span className={`px-2 py-1 rounded-full text-xs font-medium bg-primary/20 text-primary border border-primary/30`}>
                          {recType.type}
                        </span>
                      </div>
                      
                      <div className="flex items-center gap-2">
                        <TrendingUp className="h-3 w-3 text-success" />
                        <span className="text-xs text-muted-foreground">
                          Recommended for portfolio balance
                        </span>
                      </div>
                    </div>

                    <div className="flex-shrink-0">
                      <Badge variant="outline" className="text-xs">
                        #{index + 1}
                      </Badge>
                    </div>
                  </div>
                );
              })}
            </div>

            <div className="pt-4 border-t">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Total Recommendations</span>
                <Badge variant="secondary">{recommendations.length}</Badge>
              </div>
            </div>
          </div>
        ) : (
          <div className="text-center py-8 text-muted-foreground">
            <Lightbulb className="h-8 w-8 mx-auto mb-2" />
            <p>No recommendations available</p>
            <p className="text-xs">Analyze a wallet to get AI-powered recommendations</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default Recommendations;