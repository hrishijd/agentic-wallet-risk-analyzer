import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Building2, ExternalLink } from "lucide-react";

interface DeFiPosition {
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
}

interface DeFiPositionsProps {
  positions: DeFiPosition[];
}

const DeFiPositions = ({ positions }: DeFiPositionsProps) => {
  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(value);
  };

  const getProtocolCategory = (slug: string) => {
    const categories: { [key: string]: string } = {
      'DEX': 'DEX',
      'DeFi': 'Lending',
      'Payments': 'Streaming',
      'Privacy': 'Privacy',
      'Social': 'Social',
      '': 'Other'
    };
    return categories[slug] || 'DeFi';
  };

  const getCategoryColor = (category: string) => {
    const colors: { [key: string]: string } = {
      'DEX': 'bg-blue-100 text-blue-800',
      'Lending': 'bg-green-100 text-green-800',
      'Streaming': 'bg-purple-100 text-purple-800',
      'Privacy': 'bg-gray-100 text-gray-800',
      'Social': 'bg-pink-100 text-pink-800',
      'Other': 'bg-orange-100 text-orange-800'
    };
    return colors[category] || 'bg-gray-100 text-gray-800';
  };

  // Filter out positions with zero balance
  const activePositions = positions.filter(position => 
    position.balances.some(balance => Math.abs(balance.balance_usd) > 0.01)
  );

  // Sort by total value
  const sortedPositions = activePositions.sort((a, b) => {
    const totalA = a.balances.reduce((sum, balance) => sum + Math.abs(balance.balance_usd), 0);
    const totalB = b.balances.reduce((sum, balance) => sum + Math.abs(balance.balance_usd), 0);
    return totalB - totalA;
  });

  return (
    <Card className="shadow-card">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Building2 className="h-5 w-5" />
          DeFi Positions
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {sortedPositions.map((position, index) => {
            const totalValue = position.balances.reduce((sum, balance) => sum + balance.balance_usd, 0);
            const category = getProtocolCategory(position.app.slug);
            
            return (
              <div
                key={`${position.app.display_name}-${index}`}
                className="p-4 rounded-lg border bg-card hover:bg-accent/50 transition-colors"
              >
                {/* Protocol Header */}
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
                      <Building2 className="h-5 w-5" />
                    </div>
                    <div>
                      <h4 className="font-semibold">{position.app.display_name}</h4>
                      <div className="flex items-center gap-2">
                        <Badge variant="secondary" className="text-xs">
                          {position.network.name}
                        </Badge>
                        <span className={`px-2 py-1 rounded-full text-xs font-medium ${getCategoryColor(category)}`}>
                          {category}
                        </span>
                      </div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="font-semibold">{formatCurrency(totalValue)}</div>
                    <ExternalLink className="h-4 w-4 text-muted-foreground hover:text-foreground cursor-pointer ml-auto mt-1" />
                  </div>
                </div>

                {/* Position Details */}
                <div className="space-y-2">
                  {position.balances.map((balance, balanceIndex) => (
                    <div
                      key={balanceIndex}
                      className="flex items-center justify-between py-2 px-3 bg-muted/30 rounded-md"
                    >
                      <div className="flex items-center gap-3">
                        {/* Position Images */}
                        <div className="flex -space-x-2">
                          {balance.display_props.images.slice(0, 3).map((img, imgIndex) => (
                            <img
                              key={imgIndex}
                              src={img}
                              alt=""
                              className="w-6 h-6 rounded-full border-2 border-background"
                              onError={(e) => {
                                const target = e.target as HTMLImageElement;
                                target.style.display = 'none';
                              }}
                            />
                          ))}
                        </div>
                        <div>
                          <div className="text-sm font-medium">{balance.display_props.label}</div>
                          <div className="text-xs text-muted-foreground capitalize">{balance.address}</div>
                        </div>
                      </div>
                      <div className="font-medium">
                        {formatCurrency(balance.balance_usd)}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            );
          })}
        </div>

        {activePositions.length === 0 && (
          <div className="text-center py-8 text-muted-foreground">
            <Building2 className="h-8 w-8 mx-auto mb-2" />
            <p>No active DeFi positions found</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default DeFiPositions;