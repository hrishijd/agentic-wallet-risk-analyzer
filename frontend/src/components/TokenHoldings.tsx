import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Coins, TrendingUp, ExternalLink } from "lucide-react";

interface Token {
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
}

interface TokenHoldingsProps {
  tokens: Token[];
}

const TokenHoldings = ({ tokens }: TokenHoldingsProps) => {
  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(value);
  };

  const formatTokenAmount = (value: number) => {
    if (value < 0.01) return value.toExponential(2);
    if (value < 1000) return value.toFixed(6);
    if (value < 1000000) return (value / 1000).toFixed(2) + 'K';
    return (value / 1000000).toFixed(2) + 'M';
  };

  const sortedTokens = [...tokens].sort((a, b) => b.balance_usd - a.balance_usd);

  return (
    <Card className="shadow-card">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Coins className="h-5 w-5" />
          Token Holdings
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {sortedTokens.map((token) => (
            <div
              key={token.token_address}
              className="flex items-center gap-4 p-4 rounded-lg border bg-card hover:bg-accent/50 transition-colors"
            >
              {/* Token Icon */}
              <div className="flex-shrink-0">
                {token.img_url_v2 ? (
                  <img
                    src={token.img_url_v2}
                    alt={token.symbol}
                    className="w-10 h-10 rounded-full"
                    onError={(e) => {
                      const target = e.target as HTMLImageElement;
                      target.style.display = 'none';
                      target.nextElementSibling?.classList.remove('hidden');
                    }}
                  />
                ) : null}
                <div className={`w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center ${token.img_url_v2 ? 'hidden' : ''}`}>
                  <span className="text-sm font-medium">{token.symbol.slice(0, 2)}</span>
                </div>
              </div>

              {/* Token Info */}
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <h4 className="font-semibold truncate">{token.symbol}</h4>
                  <Badge variant="secondary" className="text-xs">
                    {token.network.name}
                  </Badge>
                </div>
                <p className="text-sm text-muted-foreground truncate">{token.name}</p>
              </div>

              {/* Balance Info */}
              <div className="text-right">
                <div className="font-semibold">{formatCurrency(token.balance_usd)}</div>
                <div className="text-sm text-muted-foreground">
                  {formatTokenAmount(token.balance)} {token.symbol}
                </div>
                <div className="text-xs text-muted-foreground">
                  @ {formatCurrency(token.price)}
                </div>
              </div>

              {/* External Link */}
              <div className="flex-shrink-0">
                <ExternalLink className="h-4 w-4 text-muted-foreground hover:text-foreground cursor-pointer" />
              </div>
            </div>
          ))}
        </div>

        {tokens.length === 0 && (
          <div className="text-center py-8 text-muted-foreground">
            <Coins className="h-8 w-8 mx-auto mb-2" />
            <p>No token holdings found</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default TokenHoldings;