import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { TrendingUp, Shield, Coins, Building2 } from "lucide-react";

interface PortfolioOverviewProps {
  totalValue: number;
  riskScore: number;
  tokenCount: number;
  defiProtocolCount: number;
}

const PortfolioOverview = ({
  totalValue,
  riskScore,
  tokenCount,
  defiProtocolCount,
}: PortfolioOverviewProps) => {
  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(value);
  };

  const getRiskBadgeVariant = (score: number) => {
    if (score <= 0.3) return "success";
    if (score <= 0.6) return "warning";
    return "destructive";
  };

  const getRiskLabel = (score: number) => {
    if (score <= 0.3) return "Low Risk";
    if (score <= 0.6) return "Medium Risk";
    return "High Risk";
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <Card className="shadow-card hover:shadow-card-hover transition-all duration-200">
        <CardHeader className="pb-3">
          <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
            <TrendingUp className="h-4 w-4" />
            Portfolio Value
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatCurrency(totalValue)}</div>
        </CardContent>
      </Card>

      <Card className="shadow-card hover:shadow-card-hover transition-all duration-200">
        <CardHeader className="pb-3">
          <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
            <Shield className="h-4 w-4" />
            Risk Score
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-3">
            <div className="text-2xl font-bold">{(riskScore * 100).toFixed(0)}%</div>
            <Badge variant={getRiskBadgeVariant(riskScore)}>
              {getRiskLabel(riskScore)}
            </Badge>
          </div>
        </CardContent>
      </Card>

      <Card className="shadow-card hover:shadow-card-hover transition-all duration-200">
        <CardHeader className="pb-3">
          <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
            <Coins className="h-4 w-4" />
            Token Holdings
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{tokenCount}</div>
          <div className="text-sm text-muted-foreground">Different tokens</div>
        </CardContent>
      </Card>

      <Card className="shadow-card hover:shadow-card-hover transition-all duration-200">
        <CardHeader className="pb-3">
          <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
            <Building2 className="h-4 w-4" />
            DeFi Protocols
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{defiProtocolCount}</div>
          <div className="text-sm text-muted-foreground">Active protocols</div>
        </CardContent>
      </Card>
    </div>
  );
};

export default PortfolioOverview;