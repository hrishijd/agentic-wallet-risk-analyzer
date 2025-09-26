from uagents import Agent, Context
from models import RiskRequest, RiskResponse
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()


RULES = """
; Risk factor: concentration
(= (risk-factor (holding ?addr ?token ?v) ?total)
   (if (> ?v (* 0.5 ?total)) 0.3 0.0))

; Risk factor: leverage
(= (risk-factor (futures ?addr ?m ?lev ?val))
   (if (> ?lev 2) 0.2 0.0))

; Default
(= (risk-factor _ _) 0.0)

; Risk score: sum factors
(= (risk-score ?addr ?total)
   (sum (map (lambda $x (risk-factor $x ?total))
             (concat (get-holdings ?addr)
                     (get-dex ?addr)
                     (get-futures ?addr)))))

; Recommendations
(= (recommend ?addr USDC "High concentration in one token")
   (holding ?addr ETH ?v)
   (risk-factor (holding ?addr ETH ?v) ?t)
   (> ?v (* 0.5 ?t)))

(= (recommend ?addr BTC "Leverage exposure detected")
   (futures ?addr ETH-PERP ?lev ?val)
   (> ?lev 2))
"""


# ---------- Helper: Inject Atoms ----------

def build_atoms(req: RiskRequest):
    atoms = []
    addr = req.address

    total_value = sum(h.usd_value for h in req.token_holdings) + \
                  sum(d.usd_value for d in req.dex_positions) + \
                  sum(f.usd_value for f in req.futures_positions)

    for h in req.token_holdings:
        atoms.append(f"(holding {addr} {h.symbol} {h.usd_value})")

    for d in req.dex_positions:
        atoms.append(f"(dex {addr} {d.pool} {d.usd_value})")

    for f in req.futures_positions:
        atoms.append(f"(futures {addr} {f.market} {f.leverage} {f.usd_value})")

    atoms.append(f"(total-portfolio-value {total_value})")
    return atoms, total_value


# ---------- Risk Analysis via MeTTa ----------

def analyze_risk_metta(req: RiskRequest) -> RiskResponse:
    m = MeTTa()
    m.run(RULES)

    atoms, total_value = build_atoms(req)
    for a in atoms:
        m.run(a)

    # Risk score query
    rs = m.run(f"(risk-score {req.address} {total_value})")
    risk_score = float(rs[0]) if rs else 0.0

    # Recommendations
    recs = m.run(f"(recommend {req.address} ?token ?reason)")
    recommendations, reasoning = [], []
    for rec in recs:
        parts = str(rec).strip("()").split()
        if len(parts) >= 4:
            token = parts[2]
            reason = " ".join(parts[3:]).strip('"')
            recommendations.append(token)
            reasoning.append(reason)

    if not recommendations:
        reasoning = ["Portfolio appears well balanced."]

    return RiskResponse(
        recommended_tokens=recommendations,
        risk_score=risk_score,
        reasoning=reasoning,
    )


# ---------- Agent Definition ----------

agent = Agent(
    name="risk_advisor",
    seed=os.getenv("SEED_PHRASE", "default-seed-phrase"),
    port=8000,
    endpoint=["http://localhost:8000/submit"],
)


@agent.on_message(model=RiskRequest, replies=RiskResponse)
async def handle_risk_request(ctx: Context, sender: str, req: RiskRequest):
    ctx.logger.info(f"Received risk request for {sender}")
    response = analyze_risk_metta(req)
    await ctx.send(sender, response)


if __name__ == "__main__":
    agent.run()
