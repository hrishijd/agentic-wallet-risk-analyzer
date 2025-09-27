from uagents import Agent, Context, Model
from typing import List

# Replace with the actual on-chain address of your risk_advisor agent
RISK_ADVISOR_ADDRESS = "agent1qf30d46w02h6qjp473a4wtu0xe2a42eva0sxvg3gv4zneuqcly8j5ukwcf4"  

# ---------- Define the client agent ----------

class TokenHolding(Model):
    symbol: str
    amount: float
    usd_value: float


class DexPosition(Model):
    id: str
    pool: str
    token0: str
    token1: str
    liquidity: float
    usd_value: float


class FuturesPosition(Model):
    id: str
    market: str
    amount: float
    leverage: float
    status: str
    usd_value: float


class RiskRequest(Model):
    address: str
    token_holdings: List[TokenHolding]
    dex_positions: List[DexPosition]
    futures_positions: List[FuturesPosition]


class RiskResponse(Model):
    recommended_tokens: List[str]
    risk_score: float
    reasoning: List[str]

client = Agent(
    name="risk_client",
    seed="risk-client-seed-phrase",
    port=8001,
    endpoint=["http://localhost:8001/submit"],
)


@client.on_event("startup")
async def startup(ctx: Context):
    ctx.logger.info("Client agent starting up...")

    # Build a sample portfolio request
    request = RiskRequest(
        address="0x1234abcd...",
        token_holdings=[
            TokenHolding(symbol="ETH", amount=1.2, usd_value=3000),
            TokenHolding(symbol="USDC", amount=1000, usd_value=1000),
        ],
        dex_positions=[
            DexPosition(id="1", pool="ETH-USDC", token0="ETH", token1="USDC", liquidity=2.0, usd_value=500),
        ],
        futures_positions=[
            FuturesPosition(id="1", market="ETH-PERP", amount=2, leverage=3, status="open", usd_value=1200),
        ],
    )

    # Send query to the risk_advisor agent
    await ctx.send(RISK_ADVISOR_ADDRESS, request)
    ctx.logger.info(f"Sent risk analysis request to {RISK_ADVISOR_ADDRESS}")


@client.on_message(model=RiskResponse)
async def handle_response(ctx: Context, sender: str, msg: RiskResponse):
    ctx.logger.info(f"Received response from {sender}")
    ctx.logger.info(f"Recommended tokens: {msg.recommended_tokens}")
    ctx.logger.info(f"Risk score: {msg.risk_score}")
    ctx.logger.info(f"Reasoning: {msg.reasoning}")


if __name__ == "__main__":
    client.run()
