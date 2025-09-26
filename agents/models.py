from uagents import Model
from typing import List


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


class RiskRequest(Model):
    address: str
    token_holdings: List[TokenHolding]
    dex_positions: List[DexPosition]
    futures_positions: List[FuturesPosition]


class RiskResponse(Model):
    recommended_tokens: List[str]
    risk_score: float
    reasoning: List[str]
