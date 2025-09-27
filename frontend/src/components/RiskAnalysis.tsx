import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Shield, AlertTriangle, CheckCircle, Info } from "lucide-react";

interface RiskAnalysisProps {
  riskScore: number;
  reasoning: string[];
}

const RiskAnalysis = ({ riskScore, reasoning }: RiskAnalysisProps) => {
  const getRiskLevel = (score: number) => {
    if (score <= 0.3) return { level: "Low", color: "success", icon: CheckCircle };
    if (score <= 0.6) return { level: "Medium", color: "warning", icon: Info };
    return { level: "High", color: "destructive", icon: AlertTriangle };
  };

  const risk = getRiskLevel(riskScore);
  const RiskIcon = risk.icon;

  const getRiskDescription = (score: number) => {
    if (score <= 0.3) return "Your portfolio shows low risk characteristics with good diversification.";
    if (score <= 0.6) return "Your portfolio has moderate risk with some areas for improvement.";
    return "Your portfolio shows high risk characteristics that require attention.";
  };

  return (
    <Card className="shadow-card">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Shield className="h-5 w-5" />
          Risk Analysis
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Risk Score Display */}
        <div className="flex items-center justify-between p-4 bg-muted/30 rounded-lg">
          <div className="flex items-center gap-3">
            <RiskIcon className={`h-6 w-6 ${
              risk.color === 'success' ? 'text-success' :
              risk.color === 'warning' ? 'text-warning' : 'text-destructive'
            }`} />
            <div>
              <div className="font-semibold">Risk Level: {risk.level}</div>
              <div className="text-sm text-muted-foreground">
                {getRiskDescription(riskScore)}
              </div>
            </div>
          </div>
          <div className="text-right">
            <div className="text-2xl font-bold">{(riskScore * 100).toFixed(0)}%</div>
            <Badge variant={risk.color as any}>{risk.level} Risk</Badge>
          </div>
        </div>

        {/* Risk Factors */}
        <div>
          <h4 className="font-semibold mb-3 flex items-center gap-2">
            <AlertTriangle className="h-4 w-4" />
            Risk Factors
          </h4>
          <div className="space-y-3">
            {reasoning.map((reason, index) => (
              <Alert key={index} className="border-l-4 border-l-warning">
                <AlertDescription className="text-sm leading-relaxed">
                  {reason}
                </AlertDescription>
              </Alert>
            ))}
          </div>
        </div>

        {/* Risk Metrics */}
        <div className="grid grid-cols-2 gap-4 pt-4 border-t">
          <div className="text-center">
            <div className="text-2xl font-bold text-muted-foreground">
              {(riskScore * 100).toFixed(1)}%
            </div>
            <div className="text-sm text-muted-foreground">Risk Score</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-muted-foreground">
              {reasoning.length}
            </div>
            <div className="text-sm text-muted-foreground">Risk Factors</div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default RiskAnalysis;