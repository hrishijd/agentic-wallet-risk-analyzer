from uagents import Agent, Context
from models import TokenHolding, DexPosition, FuturesPosition, RiskRequest, RiskResponse
from typing import List
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

def analyze_risk(
    holdings: List[TokenHolding],
    dex_positions: List[DexPosition],
    futures_positions: List[FuturesPosition],
) -> RiskResponse:
    total_value = sum(h.usd_value for h in holdings) + \
                  sum(d.usd_value for d in dex_positions)

    reasoning = []

    # --- Portfolio concentration risk
    if holdings:
        largest = max(holdings, key=lambda h: h.usd_value)
        if total_value > 0 and largest.usd_value / total_value > 0.5:
            reasoning.append(
                f"High concentration: {largest.symbol} is more than 50% of portfolio."
            )

    # --- DEX exposure risk
    if dex_positions:
        for d in dex_positions:
            if d.token0 == "ETH" or d.token1 == "ETH":
                reasoning.append(
                    f"DEX LP {d.id} is exposed to ETH volatility."
                )

    # --- Futures leverage risk
    high_leverage = [p for p in futures_positions if p.leverage > 2]
    if high_leverage:
        reasoning.append(
            f"{len(high_leverage)} futures positions use leverage > 2, increasing risk."
        )

    # --- Risk scoring
    risk_score = 0.3  # baseline safe portfolio
    if reasoning:
        risk_score += 0.4
    if len(high_leverage) > 0:
        risk_score += 0.3
    risk_score = min(risk_score, 1.0)

    # --- Recommendations
    recommended = []
    if risk_score > 0.7:
        recommended = ["USDC", "DAI"]  # stable assets for hedging
    elif risk_score > 0.5:
        recommended = ["BTC", "ETH"]   # diversify to majors
    else:
        recommended = ["ETH", "MATIC"] # add some growth assets

    return RiskResponse(
        recommended_tokens=recommended,
        risk_score=risk_score,
        reasoning=reasoning or ["Portfolio appears well balanced."]
    )


# ---------- Agent Definition ----------

agent = Agent(
    name="risk_advisor",
    seed=os.getenv("SEED_PHRASE", "default-seed-phrase"),
)


@agent.on_query(model=RiskRequest, replies=RiskResponse)
async def handle_risk_request(ctx: Context, sender: str, req: RiskRequest):
    ctx.logger.info(f"Received risk request for {req.address}")
    response = analyze_risk(req.token_holdings, req.dex_positions, req.futures_positions)
    await ctx.send(sender, response)


if __name__ == "__main__":
    agent.run()