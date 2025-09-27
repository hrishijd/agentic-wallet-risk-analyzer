from uagents import Model
from typing import List, Optional


class Network(Model):
    name: str
    slug: str


class TokenBalance(Model):
    token_address: str
    symbol: str
    name: str
    decimals: float
    price: float
    balance: float
    balance_usd: float
    balance_raw: str
    network: Network
    img_url_v2: Optional[str] = None


class TokenBalances(Model):
    total_balance_usd: float
    by_token: List[TokenBalance]


class App(Model):
    display_name: str
    slug: str


class TokenPosition(Model):
    meta_type: str  # SUPPLIED, BORROWED, CLAIMABLE, VESTING, LOCKED, NFT, WALLET
    token: TokenBalance


class ContractPosition(Model):
    address: str
    balance_usd: float
    tokens: List[TokenPosition]
    display_props: dict  # label, images


class AppBalance(Model):
    app: App
    network: Network
    balances: List[ContractPosition]


class AppBalances(Model):
    by_app: List[AppBalance]


class RiskRequest(Model):
    address: str
    token_balances: TokenBalances
    app_balances: AppBalances


class RiskResponse(Model):
    recommended_tokens: List[str]
    risk_score: float
    reasoning: List[str]