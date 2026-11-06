# Proxy Reselling Business Models: Comprehensive Analysis (2024-2026)

## What are the different business models for a proxy reselling business?

### Takeaway
There are five main models, ranging from near-zero-capital pure reselling to capital-intensive infrastructure ownership; most newcomers in 2025-2026 entered through white-label reselling or SIM-card farming, with very few building full infrastructure from scratch.

### Cited Findings
- **Pure Reselling (API Reselling):** You purchase proxy bandwidth at wholesale rates from an upstream provider and resell it under your own brand. No hardware is needed. PROXIES.SX's open-source reseller kit demonstrates this: you get an API key, clone a storefront template (Next.js + Stripe + Postgres), set your own prices, and deploy. Wholesale starts at $4/GB with no deposit or monthly fees. -- [PROXIES.SX Reseller Guide](https://www.proxies.sx/blog/how-to-start-a-proxy-reselling-business-2026)
- **White-Label:** Similar to pure reselling but with deeper branding integration. The upstream provider offers a fully branded dashboard, API, and sometimes custom domains so your customers never see the supplier. Multiple providers offer white-label programs including Proxy-Cheap, Proxying.io, and DataImpulse. -- [Proxying.io White Label Program](https://www.proxying.io/white-label-proxy-reseller-program); [Proxy-Cheap White Label](https://www.proxy-cheap.com/white-label-reseller)
- **SIM Card Farm / Own Mobile Infrastructure:** You purchase SIM cards, modems, and servers to build your own mobile proxy network. Proxyway's 2026 report found that among over 50 new proxy vendors established in 2025, the majority were "reselling residential proxy networks or erecting their own SIM card farms." -- [Proxyway 2026 Market Report announcement](https://www.newswire.com/news/proxyway-releases-2026-proxy-server-market-report-22788334)
- **Full Infrastructure Provider:** You own and manage the entire stack: IP pools (residential, datacenter, or mobile), rotation infrastructure, servers, and the customer-facing platform. One forum source describes this as giving "the most control but requires significant upfront investment and ongoing maintenance." -- [BlackHatWorld Forum Discussion](https://www.blackhatworld.com/seo/how-to-start-a-proxy-provider-business.1458349/)
- **Hybrid / SaaS Model:** Some providers combine proxy access with value-added tooling (scraping APIs, SERP APIs, web unlocker tools). This moves the product from raw bandwidth into a SaaS offering with higher margins and stickier customers. Proxyway also benchmarks what it calls "Web Data APIs for AI," a category that blurs the line between proxy and SaaS. -- [Proxyway Web Data APIs Report](https://www.newswire.com/news/proxyway-benchmarks-analyzes-the-market-of-web-data-apis-for-ai-22688581)

### Inferences
- Pure reselling and white-label are by far the most accessible entry points, as evidenced by the flood of new entrants in 2025 using exactly these models.
- The SaaS/hybrid model commands higher margins because it shifts the value proposition from commodity bandwidth to outcomes (successful scrapes, data delivery, etc.), but requires engineering investment.
- SIM card farming represents a middle ground: you control the IP source (which matters for quality and ethics), but costs scale linearly with each SIM and modem.

### Gaps
- No reliable data was found on the number of providers using each model or their relative market share.
- Forum data (BlackHatWorld) was inaccessible due to 403 errors, so practitioner-level startup stories could not be verified.


## What are realistic profit margins for each model?

### Takeaway
Vendor marketing claims 30-50% net margins, but independent analysis suggests real-world margins are much thinner (5-20%) once customer acquisition, churn, refunds, and support are fully accounted for; the model, volume tier, and customer retention are the key drivers.

### Cited Findings
- **Vendor-side claims:** DataImpulse cites 20-40% distributor markup and 30-50% net margins after support and payment costs. -- [DataImpulse Reseller Programs](https://dataimpulse.com/blog/best-proxy-reseller-programs/)
- **PROXIES.SX illustrative model:** A hypothetical reseller selling 600 GB/month at ~$10/GB retail with $3.20/GB wholesale: gross revenue $6,000, wholesale cost $1,920, Stripe fees $186, refunds/chargebacks (~3%) $180, support (15 hrs at $25/hr) $375, net profit ~$3,339 (~56% net margin). The page explicitly states this is "not a representative income, an average, or a guarantee." -- [PROXIES.SX Margin Analysis](https://www.proxies.sx/blog/how-much-do-proxy-resellers-make)
- **Per-order gross margins at different wholesale tiers (selling 25 GB at $300 retail):**
  - At $4.00/GB wholesale: 67% gross margin
  - At $3.20/GB wholesale: 73% gross margin
  - At $2.40/GB wholesale: 80% gross margin
  -- [PROXIES.SX Margin Analysis](https://www.proxies.sx/blog/how-much-do-proxy-resellers-make)
- **Independent critical view (Proxidize example at $1.00/GB sold):** Mobile data and network costs: $0.20 (20%), customer acquisition: $0.50 (50%), infrastructure/support/R&D: $0.25 (25%), remaining margin: $0.05 (5%). -- [Scraping Club - Unit Economics of Proxy Providers](https://www.scraping.club/p/the-unit-economics-of-proxy-providers)
- **Customer acquisition cost can consume 40-60% of revenue** even at scale, according to the same source. -- [Scraping Club - Unit Economics](https://www.scraping.club/p/the-unit-economics-of-proxy-providers)
- **Bandwidth acquisition cost:** ~$0.10-0.20 per GB for residential/mobile traffic, resold at $4-8/GB, representing a theoretical "20x to 80x" markup on raw bandwidth. -- [Scraping Club - Unit Economics](https://www.scraping.club/p/the-unit-economics-of-proxy-providers)
- **Skeptical reseller view:** "Many resellers operate on thin margins. Maybe 10-20% profit, not 300%." -- [PROXIES.SX Blog](https://www.proxies.sx/blog/how-much-do-proxy-resellers-make)

### Inferences
- The gap between gross margins (60-80%) and net margins (5-20% realistically) is explained almost entirely by customer acquisition costs and churn.
- Resellers who can acquire customers organically (SEO, community, existing audience) will have dramatically better unit economics than those relying on paid acquisition.
- Volume tiers create a powerful incentive loop: more sales lower your wholesale cost, which widens margins, which funds more acquisition.
- The Proxidize example (5% net at $1/GB) and the PROXIES.SX example (56% net at $10/GB) show that pricing power is everything. The spread between $1/GB and $10/GB retail is where margins live or die.

### Gaps
- No audited, independently verified profit margin data was found for proxy resellers of any size.
- Datacenter proxy margins (which should be higher due to lower bandwidth costs) were not broken out in any source found.
- No data on SaaS/hybrid model margins specifically was found.


## What is the minimum capital needed to start for each model?

### Takeaway
Pure reselling can start with under $100 (just bandwidth and a domain); white-label programs range from $0 to a few hundred dollars; SIM farming requires $500-5,000+ for hardware; full infrastructure requires tens of thousands or more.

### Cited Findings
- **Pure reselling:** PROXIES.SX advertises "no deposit" and "pay per GB - no monthly fees," meaning you pay only for bandwidth as you sell it. At $4/GB, serving your first 5 customers buying 5 GB each costs $100 in bandwidth. -- [PROXIES.SX Reseller Guide](https://www.proxies.sx/blog/how-to-start-a-proxy-reselling-business-2026)
- **White-label:** One source states that "whitelabeling proxies would require capital to buy the proxies and market them, which is not a lot." Some programs advertise no upfront deposits. -- [Proxying.io Blog](https://www.proxying.io/blog/how-to-start-your-own-proxy-business)
- **Full infrastructure:** "Owning servers gives you the most control but requires significant upfront investment and ongoing maintenance." No dollar figure was provided for proxy-specific infrastructure. -- [BlackHatWorld Forum](https://www.blackhatworld.com/seo/how-to-start-a-proxy-provider-business.1458349/)

### Inferences
- For pure reselling, the real capital need is not bandwidth but marketing. If CAC is 40-60% of revenue, a reseller spending $500/month on ads to generate $1,000 in revenue is barely breaking even.
- SIM card farming capital can be estimated: a 4G modem costs $30-80, a SIM card is $5-15/month, and you need a server ($50-200/month). A 10-modem setup could cost $500-1,500 upfront plus $100-300/month in SIM and hosting costs.
- Full infrastructure providers like Bright Data, Oxylabs, and Smartproxy have raised millions in venture capital, suggesting the capital requirements for competing at scale are substantial.

### Gaps
- No concrete dollar figures for building full proxy infrastructure from scratch were found in any source.
- SIM farming capital estimates are inferred from general hardware prices, not from proxy-specific sources.
- White-label minimum commitments (monthly minimums, setup fees) vary by provider and were not systematically documented.


## How do successful proxy businesses scale their operations?

### Takeaway
Scaling happens along three axes: volume tiers that reduce per-GB cost, expanding the IP pool (geographic coverage and proxy types), and moving up the value chain from raw bandwidth to tooling (APIs, dashboards, integrations).

### Cited Findings
- **Volume tier leverage:** Wholesale rates drop significantly with volume. At PROXIES.SX: $4.00/GB for 1-24 GB, $3.60 for 25-49 GB, $3.20 for 50-99 GB, $2.80 for 100-249 GB, $2.40 for 250+ GB. This means a reseller at 250+ GB/month pays 40% less per GB than one at entry level. -- [PROXIES.SX Reseller Guide](https://www.proxies.sx/blog/how-to-start-a-proxy-reselling-business-2026)
- **Provider evolution:** Proxyway noted that SOAX and IPRoyal, "both scrappy newcomers just a few years ago, today have a well-defined business presence, their own proxy networks." This suggests the path from reseller to infrastructure owner is viable. -- [Proxyway Market Research 2023](https://proxyway.com/research/proxy-market-research-2023)
- **Value-added services:** Proxyway now separately benchmarks "Web Data APIs for AI," indicating that the market rewards providers who add scraping, parsing, and data delivery tooling on top of raw proxy access. -- [Proxyway Web Data APIs Benchmark](https://www.newswire.com/news/proxyway-benchmarks-analyzes-the-market-of-web-data-apis-for-ai-22688581)
- **AI-agent payments:** PROXIES.SX describes an x402 USDC payment flow for AI agents that can purchase proxy access programmatically, opening a new customer channel. -- [PROXIES.SX Reseller Guide](https://www.proxies.sx/blog/how-to-start-a-proxy-reselling-business-2026)

### Inferences
- The most common scaling path appears to be: start as a reseller, build a customer base, reinvest in infrastructure (own SIM farms or peer networks), and then add SaaS features.
- Geographic expansion (adding proxies in more countries) increases total addressable market without changing the core business model.
- Providers that successfully scaled appear to have transitioned from pure bandwidth arbitrage to platform businesses where the proxy is one component of a data-collection stack.

### Gaps
- No detailed case studies of a proxy business scaling from zero to significant revenue were found.
- Revenue milestones (e.g., time to $10K MRR, $100K MRR) are not documented in any public source found.


## What are the recurring costs (bandwidth, servers, support, payment processing)?

### Takeaway
For resellers, the main recurring costs are wholesale bandwidth (40-60% of revenue), payment processing (~3%), support labor, and refunds/chargebacks (~3%); for infrastructure owners, add server hosting, SIM plans, IP leasing, and engineering.

### Cited Findings
- **Wholesale bandwidth:** $2.40-$4.00/GB depending on volume tier (PROXIES.SX); $0.10-0.20/GB raw acquisition cost for providers who source their own residential traffic. -- [PROXIES.SX](https://www.proxies.sx/blog/how-much-do-proxy-resellers-make); [Scraping Club](https://www.scraping.club/p/the-unit-economics-of-proxy-providers)
- **Payment processing:** Stripe charges 2.9% + $0.30 per transaction. On a $60 sale, that is ~$2.04 (3.4%). On a $1,200 sale, it is ~$35.10 (2.9%). -- [PROXIES.SX Margin Analysis](https://www.proxies.sx/blog/how-much-do-proxy-resellers-make)
- **Refunds and chargebacks:** Estimated at ~3% of revenue in the PROXIES.SX illustrative model. -- [PROXIES.SX Margin Analysis](https://www.proxies.sx/blog/how-much-do-proxy-resellers-make)
- **Support:** Estimated at ~15 hours/month at $25/hour ($375/month) for a 40-order/month business. -- [PROXIES.SX Margin Analysis](https://www.proxies.sx/blog/how-much-do-proxy-resellers-make)
- **Infrastructure (for providers):** Proxidize reports $0.20/GB for mobile data and network costs, plus $0.25/GB for infrastructure, servers, support, and R&D combined. -- [Scraping Club](https://www.scraping.club/p/the-unit-economics-of-proxy-providers)
- **Customer acquisition:** 40-60% of revenue at scale, the single largest cost category for proxy providers. -- [Scraping Club](https://www.scraping.club/p/the-unit-economics-of-proxy-providers)

### Inferences
- For pure resellers, the cost structure is simple: bandwidth + payment fees + support + refunds. Marketing is the variable that determines profitability.
- For infrastructure providers, engineering and R&D become significant costs that are hard to quantify from outside.
- Crypto payment options (like x402 USDC) could reduce payment processing fees, though adoption is currently limited.

### Gaps
- Server hosting costs for proxy infrastructure were not quantified in any source.
- SIM plan costs per modem per month were not found in proxy-specific sources.
- Software licensing costs (dashboard, API platform, billing) for white-label setups were not documented.


## How do proxy resellers acquire customers? What marketing channels work?

### Takeaway
The dominant channels are affiliate programs (30-50% recurring commissions are common), SEO/content marketing, proxy comparison and review sites, and community presence on forums; paid advertising exists but the high CAC (40-60% of revenue) makes it risky for smaller operators.

### Cited Findings
- **Affiliate programs pay 30-50% recurring commissions** and offer cookie durations of 30-90 days. Major programs include Oxylabs (50%), Proxy-Seller (30% for 6 months), Bright Data, and others. -- [AffMaven Proxy Affiliate Programs](https://affmaven.com/proxy-affiliate-programs/)
- **Proxy comparison sites** like Proxyway, G2, and various review blogs serve as discovery channels. Proxyway's annual market research report is the most cited independent benchmark in the industry. -- [Proxyway Market Research](https://proxyway.com/research/proxy-market-research-2025)
- **Target customer segments** include scraping teams, social-media operators, ad verification companies, and AI agents. -- [PROXIES.SX Reseller Guide](https://www.proxies.sx/blog/how-to-start-a-proxy-reselling-business-2026)
- **A VPN reseller marketing guide** recommends: SEO and content marketing for those with writing ability, paid ads for those with budget, and affiliate programs for those with connections in the industry. -- [VPN Resellers Marketing Guide](https://vpnresellers.pages.dev/blog/5-proven-marketing-strategies-sell-more-vpn)
- **CAC is 40-60% of revenue** even at scale, according to a proxy industry operator. -- [Scraping Club](https://www.scraping.club/p/the-unit-economics-of-proxy-providers)

### Inferences
- Organic acquisition (SEO, content, community) is almost certainly the only path to healthy margins for a bootstrapped reseller, given the 40-60% CAC figures from paid channels.
- Getting listed on proxy comparison sites (Proxyway, G2, TrustPilot) appears to be a high-leverage activity since buyers in this space actively research and compare.
- The high affiliate commissions (30-50%) tell you something about the lifetime value: providers can afford to give away that much because the remaining LTV justifies it.
- Building an audience in adjacent communities (web scraping, SEO, social media automation, sneaker bots) before launching a proxy business would dramatically reduce acquisition costs.

### Gaps
- Specific conversion rates by channel for proxy businesses were not found.
- Cost per acquisition by channel (SEO vs. paid vs. affiliate) was not broken out in any source.
- No data on the effectiveness of proxy marketplace/aggregator listings for smaller resellers was found.


## What is the typical customer lifetime value and churn rate in this industry?

### Takeaway
Proxy-specific LTV and churn data is extremely scarce; the one data point found (27% overall churn for a proxy business listed on Flippa) suggests high churn, consistent with commodity SaaS benchmarks.

### Cited Findings
- **One proxy business listing on Flippa reported 27% overall churn.** -- [Flippa Proxy Business Listing](https://flippa.com/14062743)
- **General SaaS benchmarks:** Monthly churn of 5% implies ~20-month average customer lifespan. Churn is often front-loaded: one example shows 30% churn in month 1, falling to 15% in month 2. -- [Monday.com Churn Guide](https://monday.com/blog/crm-and-sales/what-is-customer-churn/); [ChurnZero](https://churnzero.com/churnopedia/lifetime-value-ltv-or-customer-lifetime-value-cltv/)
- **LTV formula:** ARPA x Gross Margin / Churn Rate. -- [Wall Street Prep LTV/CAC Guide](https://www.wallstreetprep.com/knowledge/ltv-cac-ratio/)
- **A small reduction in churn compounds into a large increase in LTV.** -- [Stripe CLV Guide](https://stripe.com/resources/more/customer-lifetime-value)

### Inferences
- If a proxy reseller charges $60/month (a 5 GB starter plan) and has 10% monthly churn, the average customer lasts 10 months with an LTV of $600 gross, or ~$400 net after bandwidth costs.
- If monthly churn is 15-20% (which the front-loaded churn data suggests is plausible for a commodity proxy service), LTV drops to $300-400 gross.
- The 27% churn figure from the Flippa listing, if it is annual, is actually reasonable for B2B SaaS. If monthly, it is catastrophic. The listing did not specify the time period.
- Given the commodity nature of proxy bandwidth and low switching costs, proxy churn is likely higher than average SaaS churn. Reducing churn through value-added features (dashboards, API integrations, support quality) would have an outsized impact on profitability.

### Gaps
- No authoritative proxy-industry churn benchmark was found (monthly or annual).
- No LTV figures specific to proxy businesses were found in any source.
- Contract length distribution (month-to-month vs. annual) for proxy customers was not documented.
- Cohort-based churn analysis for proxy services was not available.


## What legal considerations exist for proxy resellers?

### Takeaway
The primary legal risks are IP sourcing ethics (botnet/consent issues), GDPR and data privacy compliance, liability for customers' end uses, and the unsettled legal landscape around web scraping; the July 2026 FBI action against NetNut demonstrates that upstream provider risk can cascade to resellers.

### Cited Findings
- **IP sourcing is the core liability issue.** Providers acquire residential IPs through methods "ranging from transparent SDK opt-in programs to illegal botnet enrollment." Under GDPR Article 28, processors must ensure sub-processors obtain valid consent. -- [PlainProxies Compliance Blog](https://plainproxies.com/blog/?p=1992)
- **FBI/NetNut action (July 2026):** The FBI, Google, Lumen, and Shadowserver disrupted NetNut, a residential proxy provider, over alleged links to the "Popa botnet" which enrolled over 2 million consumer devices as proxy exit nodes without consent. Parent company Alarum Technologies is under investigation. Bloomberg reported the FBI had been investigating for over a year. -- [HackerNoon - NetNut FBI Action](https://hackernoon.com/netnut-shut-down-by-the-fbi-heres-what-happened-and-what-to-do-next)
- **White-label downstream risk:** The HackerNoon article notes that "many residential proxy brands used NetNut's infrastructure as a white-label backend," meaning the takedown affected businesses that were not directly involved in the alleged misconduct. -- [HackerNoon](https://hackernoon.com/netnut-shut-down-by-the-fbi-heres-what-happened-and-what-to-do-next)
- **KYC requirements:** Proxidize states that customers must complete identity checks before platform access, including government ID, email, IP, region, and device identifiers. Non-verified customers are denied access. -- [Proxidize Ethics Blog](https://proxidize.com/blog/laws-and-ethics-of-being-a-proxy-provider/)
- **Scraping law is unsettled:** Multiple sources note that the legal status of web scraping "remains complex" and "laws regarding web scraping, automation, and proxy usage vary by jurisdiction." Violating a site's ToS is not always illegal but can lead to civil liability. -- [DEV Community](https://dev.to/xavier_fok/proxy-compliance-legal-considerations-every-operator-should-know-2jg0)
- **Criminal exposure:** "Proxy servers that have been used for illegal purposes can lead to fines and even criminal liability." -- [LightningProxies Legal Guide](https://lightningproxies.net/blog/are-proxies-legal-businesses-and-individuals)
- **GDPR and EU AI Act compliance** is becoming a factor, especially for providers serving scraping use cases. -- [Evomi Compliance Blog](https://evomi.com/blog/the-compliance-squeeze-scraping-between-gdpr-the-ai-act-and-platform-tos)

### Inferences
- For resellers, the biggest legal risk is upstream: if your provider sources IPs unethically or illegally, you share in the liability and business disruption risk, as the NetNut case demonstrates.
- Due diligence on upstream providers is not optional -- it is a business survival issue, not just a compliance checkbox.
- Implementing KYC for your own customers provides legal cover and reduces abuse risk, though it may reduce conversion rates.
- Operating in the EU requires particular attention to GDPR and the emerging AI Act, especially if serving web scraping customers.

### Gaps
- No specific legal cases resulting in monetary damages or criminal penalties for proxy resellers (as opposed to providers) were found.
- Insurance options for proxy businesses (E&O, cyber liability) were not documented.
- Jurisdiction-specific legal analysis (US vs. EU vs. other) was not found in sufficient detail.


## How does the competitive landscape look -- is it saturated or is there room for newcomers?

### Takeaway
The market is crowded and commoditizing, with over 50 new vendors entering in 2025 alone and prices dropping up to 70% in two years; however, differentiation through niche focus, superior UX, value-added tooling, or AI-focused services still creates openings.

### Cited Findings
- **Over 50 new proxy vendors were established in 2025,** the majority reselling residential proxy networks or erecting their own SIM card farms. -- [Proxyway 2026 Market Report](https://www.newswire.com/news/proxyway-releases-2026-proxy-server-market-report-22788334)
- **Proxyway's 2025 report analyzes 12 leading providers** and finds that "the surge in AI demand, increasing focus on user experience, and commoditization are reshaping the competitive landscape." -- [Proxyway 2025 Market Research Press Release](https://www.newswire.com/news/proxy-market-research-2025-ai-commoditization-shape-the-proxy-server-22569498)
- **Prices have dropped up to 70% in two years** for residential proxies, with "significant performance gaps despite narrowing price differences." -- [JoinMassive / Proxyway data](https://www.joinmassive.com/blog/best-residential-proxy-providers) (vendor source, treat as approximate)
- **Providers are differentiating through customization and user experience** to stay competitive in the commoditized market. -- [Proxyway 2025 Press Release](https://www.newswire.com/news/proxy-market-research-2025-ai-commoditization-shape-the-proxy-server-22569498)
- **Market size estimates:** The broader proxy server market is roughly $2-4 billion in 2025-2026, growing at 7-11% CAGR depending on segment definition. -- [Research and Markets](https://www.researchandmarkets.com/reports/6218627/proxy-servers-market-forecast); [Knowledge Sourcing Intelligence](https://www.knowledge-sourcing.com/report/proxy-servers-market); [DataIntelo](https://dataintelo.com/report/global-proxy-services-and-software-market)
- **Residential proxy market specifically:** USD 122 million in 2025, expected to reach USD 148 million by 2030, at a slower 3.98% CAGR. -- [Mordor Intelligence](https://www.mordorintelligence.com/industry-reports/residential-proxy-server-market)
- **Former newcomers have matured:** SOAX and IPRoyal, "both scrappy newcomers just a few years ago, today have a well-defined business presence, their own proxy networks." -- [Proxyway 2023 Research](https://proxyway.com/research/proxy-market-research-2023)

### Inferences
- The market is unquestionably crowded at the commodity bandwidth layer. A newcomer competing purely on price for undifferentiated residential proxy GB will struggle.
- Viable entry strategies for newcomers in 2025-2026 likely include: (1) niche focus (specific geography, use case like sneakers or social media), (2) superior developer experience and tooling, (3) AI-agent integration as a differentiator, (4) ethically sourced proxies as a selling point post-NetNut, (5) bundled data services rather than raw bandwidth.
- The 70% price drop in residential proxies is a double-edged sword: it makes it harder to profit from pure reselling, but it also increases total market demand as more businesses can afford to use proxies.
- The existence of 50+ new entrants in a single year, despite commoditization, suggests that many are chasing the business without fully understanding the economics. The washout rate is likely high.

### Gaps
- Market share data for individual providers was not publicly available.
- Failure rate or washout statistics for new proxy businesses were not found.
- Specific data on which niches or geographies are underserved was not available.
- Revenue or size data for mid-tier providers (between the top 12 and the long tail) was not found.


## Summary: Financial Viability Assessment for Someone Starting with Limited Capital

### Realistic Scenario for a Bootstrapped Reseller (2025-2026)

**Startup costs:** $100-500 (domain, initial bandwidth, basic marketing)

**Month 1-3 (Building phase):**
- Revenue: $0-500/month
- Focus: Setting up storefront, creating content, building presence in target communities
- Key risk: Zero revenue while investing time

**Month 4-6 (Early traction):**
- Revenue: $500-2,000/month
- Gross margin: 60-70% (at low volume tiers)
- Net margin after all costs: 10-30% (heavily dependent on customer acquisition method)
- Key challenge: Acquiring customers without significant ad spend

**Month 7-12 (Growth phase):**
- Revenue: $2,000-6,000/month (if gaining traction)
- Gross margin: 65-75% (climbing volume tiers)
- Net margin: 20-40% (if organic acquisition is working)
- Key milestone: Reaching volume tier discounts that meaningfully improve margins

**Critical success factors:**
1. Organic customer acquisition ability (content, SEO, community)
2. Low churn through good support and value-added features
3. Niche focus rather than competing on commodity bandwidth
4. Ethical upstream sourcing to avoid regulatory and reputational risk

**Red flags for non-viability:**
1. Relying entirely on paid acquisition (CAC will eat margins)
2. Competing on price alone in a commoditizing market
3. No differentiation from the upstream provider's own retail offering
4. Ignoring legal and compliance requirements

### Sources Note
Most financial figures in this research come from vendor marketing materials and should be treated as optimistic. The Scraping Club unit economics analysis is the most independent source found, and it paints a significantly less rosy picture (5% net margin at Proxidize scale) than vendor-side claims (30-50% net margins). The truth for a small reseller likely falls between these extremes, closer to 10-25% net margins with organic acquisition and active churn management.
