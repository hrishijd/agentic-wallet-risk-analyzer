from uagents import Agent, Context, Model
from typing import List
from models import RiskRequest, RiskResponse, TokenBalances, AppBalances, TokenBalance, AppBalance, ContractPosition, TokenPosition, App, Network

# Replace with the actual on-chain address of your risk_advisor agent
RISK_ADVISOR_ADDRESS = "agent1qf30d46w02h6qjp473a4wtu0xe2a42eva0sxvg3gv4zneuqcly8j5ukwcf4"

client = Agent(
    name="risk_client",
    seed="risk-client-seed-phrase",
    port=8001,
    endpoint=["http://localhost:8001/submit"],
)


@client.on_event("startup")
async def startup(ctx: Context):
    ctx.logger.info("Client agent starting up...")

    # Build a sample portfolio request using Zapper structure
    request = RiskRequest(
        address="0x1234abcd...",
        token_balances=TokenBalances(
            total_balance_usd=4000.0,
            by_token=[
                TokenBalance(
                    token_address="0x0000000000000000000000000000000000000000",
                    symbol="ETH",
                    name="Ethereum",
                    decimals=18.0,
                    price=2500.0,
                    balance=1.2,
                    balance_usd=3000.0,
                    balance_raw="1200000000000000000",
                    network=Network(name="Ethereum", slug="ethereum"),
                    img_url_v2="https://example.com/eth.png"
                ),
                TokenBalance(
                    token_address="0xA0b86a33E6441e8C4E2C2C1E4e1E4E1E4E1E4E1E",
                    symbol="USDC",
                    name="USD Coin",
                    decimals=6.0,
                    price=1.0,
                    balance=1000.0,
                    balance_usd=1000.0,
                    balance_raw="1000000000",
                    network=Network(name="Ethereum", slug="ethereum"),
                    img_url_v2="https://example.com/usdc.png"
                ),
            ]
        ),
        app_balances=AppBalances(
            by_app=[
                AppBalance(
                    app=App(display_name="Uniswap V3", slug="uniswap-v3"),
                    network=Network(name="Ethereum", slug="ethereum"),
                    balances=[
                        ContractPosition(
                            address="0x1234567890123456789012345678901234567890",
                            balance_usd=500.0,
                            tokens=[
                                TokenPosition(
                                    meta_type="SUPPLIED",
                                    token=TokenBalance(
                                        token_address="0x0000000000000000000000000000000000000000",
                                        symbol="ETH",
                                        name="Ethereum",
                                        decimals=18.0,
                                        price=2500.0,
                                        balance=0.2,
                                        balance_usd=500.0,
                                        balance_raw="200000000000000000",
                                        network=Network(name="Ethereum", slug="ethereum"),
                                        img_url_v2="https://example.com/eth.png"
                                    )
                                )
                            ],
                            display_props={"label": "ETH/USDC Pool", "images": ["https://example.com/eth.png"]}
                        )
                    ]
                )
            ]
        )
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
