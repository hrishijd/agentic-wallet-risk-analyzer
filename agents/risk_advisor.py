from typing import List, Tuple
from uagents import Agent, Context
from models import RiskRequest, RiskResponse, TokenBalance, TokenPosition
from hyperon import MeTTa
import os
from dotenv import load_dotenv
import logging
from dataclasses import dataclass

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Load environment variables
load_dotenv()

# Enhanced MeTTa rules with fixes for parsing and calculations
RISK_RULES = """
; Helper to match and collect expressions
(= (collect ?pattern) (match &self ?pattern ?pattern))

; Get all asset holdings (wallet tokens + supplied + locked + claimable as assets)
(= (get-all-holdings ?addr)
   (collect (or (token-holding ?addr ?token ?symbol ?v)
                (supplied ?addr ?token ?symbol ?v)
                (locked ?addr ?token ?symbol ?v)
                (claimable ?addr ?token ?symbol ?v))))

; Get all borrowed positions (liabilities)
(= (get-borrowed ?addr)
   (collect (borrowed ?addr ?token ?symbol ?v)))

; Get all locked positions (for illiquidity)
(= (get-locked ?addr)
   (collect (locked ?addr ?token ?symbol ?v)))

; Calculate total assets value
(= (total-assets ?addr)
   (sum (map (lambda $x (case $x ((_ _ _ _ $v) $v)))
             (get-all-holdings ?addr))))

; Calculate total liabilities (borrowed)
(= (total-liabilities ?addr)
   (sum (map (lambda $x (case $x ((_ _ _ _ $v) $v)))
             (get-borrowed ?addr))))

; Calculate total locked value
(= (total-locked ?addr)
   (sum (map (lambda $x (case $x ((_ _ _ _ $v) $v)))
             (get-locked ?addr))))

; Net worth (assets - liabilities)
(= (net-worth ?addr)
   (- (total-assets ?addr) (total-liabilities ?addr)))

; Concentration risk using Herfindahl-Hirschman Index (HHI)
; HHI = sum ( (v_i / total_assets)^2 ) * 10000
; Low diversification if HHI > 2500, high if >5000
(= (hhi ?addr)
   (if (== (total-assets ?addr) 0) 0
       (* 10000
          (sum (map (lambda $x (case $x ((_ _ _ _ $v)
                                        (let $share (/ $v (total-assets ?addr))
                                             (* $share $share))))
                    (get-all-holdings ?addr))))))

(= (risk-factor-concentration ?addr)
   (let $h (hhi ?addr)
        (if (> $h 5000) 0.4
            (if (> $h 2500) 0.2 0.0))))

; Leverage risk: liabilities / assets
; High if >0.5, medium if >0.2
(= (leverage-ratio ?addr)
   (if (== (total-assets ?addr) 0) 0
       (/ (total-liabilities ?addr) (total-assets ?addr))))

(= (risk-factor-leverage ?addr)
   (let $l (leverage-ratio ?addr)
        (if (> $l 0.5) 0.3
            (if (> $l 0.2) 0.15 0.0))))

; Illiquidity risk: locked / assets
; High if >0.5, medium if >0.3
(= (illiquidity-ratio ?addr)
   (if (== (total-assets ?addr) 0) 0
       (/ (total-locked ?addr) (total-assets ?addr))))

(= (risk-factor-illiquidity ?addr)
   (let $i (illiquidity-ratio ?addr)
        (if (> $i 0.5) 0.25
            (if (> $i 0.3) 0.1 0.0))))

; Overall risk score: sum of all risk factors, normalized to 0-1
(= (risk-score ?addr)
   (let* ($conc (risk-factor-concentration ?addr)
          $lev (risk-factor-leverage ?addr)
          $ill (risk-factor-illiquidity ?addr)
          $total (+ $conc $lev $ill))
         (min 1.0 $total)))

; Recommendations
; High concentration: recommend stablecoins like USDC
(= (recommend ?addr USDC "High concentration risk - consider diversifying into stablecoins")
   (> (hhi ?addr) 2500))

; High leverage: recommend BTC or ETH as hedge
(= (recommend ?addr BTC "High leverage detected - consider hedging with volatile assets like BTC")
   (> (leverage-ratio ?addr) 0.2))

(= (recommend ?addr ETH "High leverage detected - consider hedging with volatile assets like ETH")
   (> (leverage-ratio ?addr) 0.2))

; High illiquidity: recommend liquid tokens
(= (recommend ?addr USDC "High illiquidity from locked positions - add more liquid assets")
   (> (illiquidity-ratio ?addr) 0.3))

; General low risk recommendation
(= (recommend ?addr NONE "No specific recommendations - portfolio is balanced")
   (and (< (hhi ?addr) 2500)
        (< (leverage-ratio ?addr) 0.2)
        (< (illiquidity-ratio ?addr) 0.3)))
"""

@dataclass
class AnalysisResult:
    """Data class to hold risk analysis results."""
    risk_score: float
    recommended_tokens: List[str]
    reasoning: List[str]

class RiskAnalyzer:
    """Class to handle risk analysis using an enhanced MeTTa knowledge graph."""
    
    def __init__(self):
        self.metta = MeTTa()
        self.metta.run(RISK_RULES)
        logger.info("MeTTa engine initialized with enhanced risk rules")

    def build_atoms(self, req: RiskRequest) -> Tuple[List[str], float]:
        """
        Builds MeTTa atoms from RiskRequest and calculates total assets (for reference).
        
        Args:
            req: RiskRequest object containing portfolio data
            
        Returns:
            Tuple containing list of atoms and total assets value
        """
        try:
            atoms = []
            addr = req.address
            total_assets = 0.0  # Initialize to 0 to avoid double-counting

            # Add token holdings (wallet) and accumulate total_assets
            for token in req.token_balances.by_token:
                atoms.append(f'(token-holding "{addr}" "{token.token_address}" "{token.symbol}" {token.balance_usd})')
                total_assets += token.balance_usd

            # Add app positions and accumulate total_assets for assets
            for app_balance in req.app_balances.by_app:
                for contract_pos in app_balance.balances:
                    for token_pos in contract_pos.tokens:
                        meta_type = token_pos.meta_type.lower()
                        if meta_type in ["borrowed", "locked", "supplied", "claimable"]:
                            atoms.append(f'({meta_type} "{addr}" "{token_pos.token.token_address}" '
                                         f'"{token_pos.token.symbol}" {token_pos.token.balance_usd})')
                            if meta_type in ["supplied", "locked", "claimable"]:
                                total_assets += token_pos.token.balance_usd
                            # Borrowed is liability, not added to assets

            return atoms, total_assets
        except Exception as e:
            logger.error(f"Error building atoms: {str(e)}")
            raise

    def analyze(self, req: RiskRequest) -> RiskResponse:
        """
        Analyzes risk for a given portfolio using enhanced MeTTa rules.
        
        Args:
            req: RiskRequest object containing portfolio data
            
        Returns:
            RiskResponse object with analysis results
        """
        try:
            # Input validation
            if not req.address or not req.token_balances:
                raise ValueError("Invalid RiskRequest: address and token balances are required")

            # Build atoms
            atoms, python_total_assets = self.build_atoms(req)
            
            # Inject atoms into MeTTa
            for atom in atoms:
                self.metta.run(atom)

            # Debug logging for key metrics
            ta = self.metta.run(f'(total-assets "{req.address}")')
            logger.info(f"Calculated total assets: {ta}")

            tl = self.metta.run(f'(total-liabilities "{req.address}")')
            logger.info(f"Calculated total liabilities: {tl}")

            t_locked = self.metta.run(f'(total-locked "{req.address}")')
            logger.info(f"Calculated total locked: {t_locked}")

            h = self.metta.run(f'(hhi "{req.address}")')
            logger.info(f"Calculated HHI: {h}")

            # Calculate risk score
            rs = self.metta.run(f'(risk-score "{req.address}")')
            risk_score = float(rs[0]) if rs and rs[0] else 0.0

            # Get recommendations
            recs = self.metta.run(f'(recommend "{req.address}" ?token ?reason)')
            recommended_tokens, reasoning = [], []
            for rec in recs:
                if rec:
                    parts = str(rec).strip("()").split(maxsplit=3)
                    if len(parts) >= 4:
                        token = parts[2]
                        if token != "NONE":
                            recommended_tokens.append(token)
                        reasoning.append(parts[3].strip('"'))

            # Default reasoning if no recommendations
            if not reasoning:
                reasoning = ["Portfolio appears well balanced with low risk factors."]

            return RiskResponse(
                recommended_tokens=recommended_tokens,
                risk_score=risk_score,
                reasoning=reasoning
            )
        except Exception as e:
            logger.error(f"Error in risk analysis: {str(e)}")
            return RiskResponse(
                recommended_tokens=[],
                risk_score=0.0,
                reasoning=[f"Analysis failed: {str(e)}"]
            )

# Agent Definition
agent = Agent(
    name="risk_advisor",
    seed=os.getenv("SEED_PHRASE", "default-seed-phrase"),
    port=8000,
    endpoint=["http://localhost:8000/submit"],
)

# Initialize risk analyzer
risk_analyzer = RiskAnalyzer()

@agent.on_message(model=RiskRequest, replies=RiskResponse)
async def handle_risk_request(ctx: Context, sender: str, req: RiskRequest):
    """
    Handles incoming risk analysis requests via agent messaging.
    
    Args:
        ctx: Agent context
        sender: Sender address
        req: RiskRequest object
    """
    logger.info(f"Received risk request from {sender} for address {req.address}")
    response = risk_analyzer.analyze(req)
    await ctx.send(sender, response)

@agent.on_rest_post("/api/analyze", RiskRequest, RiskResponse)
async def handle_rest_risk_request(ctx: Context, req: RiskRequest):
    """
    Handles incoming risk analysis requests via REST API.
    
    Args:
        ctx: Agent context
        req: RiskRequest object
        
    Returns:
        RiskResponse object
    """
    logger.info(f"Received REST risk request for address {req.address}")
    return risk_analyzer.analyze(req)

if __name__ == "__main__":
    try:
        logger.info("Starting risk advisor agent...")
        agent.run()
    except Exception as e:
        logger.error(f"Failed to start agent: {str(e)}")