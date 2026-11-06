# Market Research

Consolidated market research for the Proxy0x proxy reselling business. Originally produced as three separate reports with supporting research notes.

---

## 1. Proxy Reselling Market Overview

O mercado de revenda de proxies residenciais movimenta cerca de **USD 122 milhoes em 2025**, com previsao de atingir USD 148 milhoes ate 2030 a um CAGR de 3,98% ([Mordor Intelligence](https://www.mordorintelligence.com/industry-reports/residential-proxy-market)). Mais de **50 novos fornecedores surgiram so em 2025**, a maioria revendendo redes residenciais ou montando farms de SIM cards ([Proxyway 2026 Market Report](https://www.newswire.com/news/proxyway-releases-2026-proxy-server-market-report-22788334)).

### Margins and Unit Economics

Materiais promocionais prometem margens liquidas de 30-50%, mas analise independente aponta que **o custo de aquisicao de clientes consome 40-60% da receita** mesmo em escala, e a margem liquida real pode ser tao baixa quanto **5%** ([Scraping Club](https://www.scraping.club/p/the-unit-economics-of-proxy-providers)). Os precos de proxies residenciais cairam ate **75% entre 2023 e 2025** ([Proxyway 2026](https://proxyway.com/research/proxy-market-research-2026)).

Fato revelador: **nao foram encontrados relatos verificados de primeira mao de revendedores lucrativos em nenhuma plataforma** (Reddit, BlackHatWorld, YouTube, Quora).

### Reseller Programs

| Provider | Min Deposit | Wholesale $/GB | Whitelabel | Notes |
|----------|-------------|----------------|------------|-------|
| DataImpulse | US$ 50 | ~US$ 1.00 | DNS only | API included, no dashboard |
| Geonode | US$ 500/mo | US$ 0.50 (10TB+) | Full | Most robust, high minimum |
| Proxy-Cheap | US$ 5 | ~US$ 1.50 | Full + WHMCS | Opaque company |
| IPRoyal | US$ 200 | ~US$ 1.75 | Sub-user API | No branding |
| CatProxies | Any | ~US$ 0.75 (est.) | Hosted storefront | Young, high risk |
| Bright Data | Free (affiliate) | N/A | No | 50% revenue share |

### Business Models

1. **Pure API reselling** — Buy wholesale, resell under own brand. Lowest barrier.
2. **Whitelabel** — Deeper brand integration with custom dashboard/domain.
3. **Mobile SIM farm** — Own hardware (3-5 modems ~R$ 750-1.500). Proxies sold at US$ 1.70-10/day.
4. **Hybrid SaaS** — Proxies + scraping APIs, web unlockers. Market trending here.
5. **Full infrastructure** — Own IP pools, rotation servers, platform. Capital intensive.
6. **Aggregation/panel** — Resell multiple providers through single account (~40% cheaper).

### Critical Risk: NetNut/FBI Case (July 2026)

FBI desmantelou a NetNut por ligacoes com botnet que inscreveu **2+ milhoes de dispositivos como exit nodes sem consentimento**. Marcas whitelabel que usavam a infraestrutura foram afetadas. **Due diligence no fornecedor upstream e questao de sobrevivencia** ([HackerNoon](https://hackernoon.com/netnut-shut-down-by-the-fbi-heres-what-happened-and-what-to-do-next)).

---

## 2. Provider Comparison (DataImpulse vs IPRoyal vs Proxy-Cheap vs CatProxies)

### Cost Analysis

| Provider | Deposit | Residential/GB | Volume Discount | Expiry |
|----------|---------|----------------|-----------------|--------|
| DataImpulse | US$ 50 | US$ 1.00 (US$ 0.80 at 1TB+) | Yes | Never |
| Proxy-Cheap | US$ 5 | ~US$ 1.50 (US$ 1.20 at 20%) | 5-20% tiered | Never |
| IPRoyal | US$ 200 | ~US$ 1.75 (high vol) | Negotiable | Never |
| CatProxies | Any | ~US$ 0.75-1.25 (est.) | "Up to 70%" (unverified) | Never |

### Financial Projections

**At US$ 200/mo revenue (~67 GB):**
- DataImpulse: **US$ 32 profit (16%)** — only profitable option
- Proxy-Cheap: breakeven
- IPRoyal: -US$ 35 loss
- CatProxies: ~US$ 15 (8%, estimated)

**At US$ 2,000/mo revenue (~667 GB):**
- DataImpulse: **US$ 1,055 (53%)**
- CatProxies: ~US$ 921 (46%)
- Proxy-Cheap: US$ 755 (38%)
- IPRoyal: US$ 488 (24%)

**At US$ 10,000/mo revenue (~3,333 GB):**
- DataImpulse & CatProxies: ~**US$ 6,263 (63%)**
- Proxy-Cheap: US$ 4,930 (49%)
- IPRoyal: US$ 3,930 (39%)

### Brazil/LATAM IP Quality

- **IPRoyal**: 205,985 unique Brazilian residential IPs tested, 97.70% success rate (Proxyway verified)
- **DataImpulse**: Brazil gateway available, but mobile accuracy only 29.69% for Brazil
- **Proxy-Cheap**: No independent benchmarks for Brazil
- **CatProxies**: Grade F (15/100) on proxies.fi, 72% overall success

### Red Flags by Provider

- **DataImpulse**: Trustpilot removed fake reviews; ~15% server IPs in US pool; small team (~20)
- **IPRoyal**: 1M+ Comcast IPs without permission; 48 downtime incidents (455h); discontinued Premium pool
- **Proxy-Cheap**: No founders identified, no legal entity published; not tested by Proxyway; API was "beta invite only"
- **CatProxies**: ~1.5 years old; anonymous team; "high risk" classification (Synthient); 72% success rate

### Recommendation

**Start with DataImpulse (US$ 50)**, build own platform with their API, price at ~US$ 3/GB. Add CatProxies/Proxy-Cheap as secondary at US$ 500/mo. Use IPRoyal for high-quality Brazilian IPs when needed. Design with provider abstraction layer from day 1.

---

## 3. Brazilian Competitors Analysis

### Competitor Map (12 players)

| Competitor | Residential /GB | CNPJ | Status | Threat Level |
|------------|----------------|------|--------|-------------|
| **ProxyBR** | Hidden (registration required) | 48.832.656/0001-05 | Active | **High** |
| **RH7 Proxys** | R$ 4.90 | 54.843.242/0001-76 | Active | **High** |
| **Allan Proxys** | R$ 6.50-7.69 | 56.928.569/0001-76 | **Inapta** (Oct/2026) | **High** |
| GM Proxys | R$ 4.90 | Not disclosed | Unknown | Unknown |
| Luma Proxys | R$ 5.94-8.04 | Not disclosed | Unknown | Unknown |
| CPA Proxys | R$ 9.90 | Not disclosed | Unknown | Low |
| FastProxy | Per-IP model | Not disclosed | Unknown | Unknown |
| GTI Proxy | R$ 40/mo (unlimited) | 18.915.849/0001-50 | Active | Low |
| ShieldProx | Per-IP model | Not verified | Unknown | Low |
| Proxy Brasil | R$ 25.00 (vol: R$ 18) | Not disclosed | Unknown | Low |
| Proxy Roque | R$ 28.80 | 42.020.889/0001-19 | Active | Low |
| Red Proxys | Unknown | Not disclosed | Unknown | None |

### Top 3 Threats

1. **ProxyBR** — AW Solucoes Digitais Ltda, R$ 150K capital, 73K-160K visits/mo, 4.97/5 Reclame Aqui. Strengths: Web Scraping API, iOS app, loyalty program. Weaknesses: hidden pricing, no mobile.

2. **RH7 Proxys** — Axion Tecnologia Ltda. Price leader at R$ 4.90/GB. Broadest portfolio among price leaders. Weakness: at DataImpulse cost of R$ 5.10/GB, mathematically impossible to match their price with positive margin.

3. **Allan Proxys** — Most complete feature set: REST API, carrier selection (Vivo/Claro/TIM), OS fingerprint, SMS service. **Critical red flag: CNPJ marked "Inapta" by Receita Federal since Oct/2026.**

### Pricing Strategy

**Recommended hybrid subscription model (monthly + overage):**

| Plan | Price/mo | GB Included | Effective R$/GB | Gross Margin |
|------|----------|-------------|-----------------|-------------|
| Free Trial | R$ 0 | 500 MB | — | — |
| Basic | R$ 49.90 | 5 GB | R$ 9.98 | 46.4% |
| Pro | R$ 149.90 | 20 GB | R$ 7.50 | 29.5% |
| Enterprise | R$ 699.90 | 100 GB | R$ 7.00 | 24.6% |

**Competitive advantage**: PIX nativo, dashboard em portugues, suporte WhatsApp, nota fiscal — no competitor offers all of these together.

### Revenue Projections (BRL)

| Scale | Clients | Monthly Revenue | Monthly Profit | Margin |
|-------|---------|-----------------|----------------|--------|
| Validation | 10 | R$ 699 | R$ 219 | 31.3% |
| Stable | 50 | R$ 7,245 | R$ 2,674 | 36.9% |
| Scale | 200 | R$ 28,980 | R$ 16,332 | 56.4% |

Break-even at 2-3 clients. PIX advantage: international proxy at US$ 5/GB costs ~R$ 32/GB after IOF + card spread vs ~R$ 10/GB via local PIX provider ([RH7 Blog](https://rh7proxys.com/blog/quanto-custa-proxy-residencial-brasil-2026)).

### Go-to-Market Strategy

**Phase 1 (Months 1-3)**: Launch with 3 tiers + free 500MB trial. Differentiate on experience (PT-BR dashboard, PIX, WhatsApp support, nota fiscal), not price.

**Phase 2 (Months 4-8)**: Content + SEO. Terms like "proxy residencial brasil", "comprar proxy PIX" have low competition. Video tutorials for AdsPower, Multilogin, Dolphin Anty integrations.

**Phase 3 (Month 9+)**: Scale with affiliate program, diversify upstream (Geonode at US$ 0.50/GB for 10TB+), add SMS verification and Web Scraping API for higher ARPU.

---

## Sources

All sources cited inline. Key references:
- [Proxyway 2026 Market Report](https://proxyway.com/research/proxy-market-research-2026)
- [Scraping Club — Unit Economics](https://www.scraping.club/p/the-unit-economics-of-proxy-providers)
- [DataImpulse Reseller Program](https://dataimpulse.com/reseller-program/)
- [HackerNoon — NetNut/FBI](https://hackernoon.com/netnut-shut-down-by-the-fbi-heres-what-happened-and-what-to-do-next)
- [BlackHatWorld — Reselling Discussions](https://www.blackhatworld.com/seo/thoughts-on-reselling-as-a-business-model.1682622/)

*Research conducted circa 2026. Prices and market conditions subject to change.*
