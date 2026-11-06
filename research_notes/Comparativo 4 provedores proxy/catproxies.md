# CatProxies -- Reseller/Whitelabel Platform Research

## What is CatProxies' business model and how does the "no deposit, hosted storefront" work?

### Takeaway
CatProxies is a proxy reseller platform operated by MOTAN DATA SOLUTIONS S.R.L. (Romania, founded ~2025) that lets approved resellers buy proxy bandwidth at wholesale rates, set their own retail prices with no price floor or revenue share, and either integrate via API or use a hosted white-label storefront. The entry-level "Starter" tier accepts any deposit amount (effectively no minimum), making it one of the lowest-barrier reseller programs in the market.

### Cited Findings
- Resellers buy at discounted rates (up to 70% off retail on select plans), set their own retail prices, and keep the full margin. There is no revenue share and no price floor. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Three reseller tiers exist: Starter (any deposit), Plus ($100+ deposit, unlocks datacenter/IPv6/ISP products), and Volume ($500+ deposit, best rates, custom invoicing, higher balance limit). -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Deposits and unused bandwidth do not expire. No monthly commitment or recurring spend on any tier. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- White-label storefront is hosted on CatProxies infrastructure, carrying the reseller's branding. End users do not see CatProxies in authentication, dashboards, or invoices. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Applicants must pass review and KYC; onboarding completes within 1-2 business days. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- The operating entity is MOTAN DATA SOLUTIONS S.R.L., registered in Targoviste, Romania (CUI 51601374, VAT RO51601374). -- [CatProxies Ethics Page](https://dashboard.catproxies.com/ethics)
- A tech directory (Techimply) states CatProxies was founded in 2025 and lists the team as "freelancers" -- this is weak evidence from a low-quality directory. -- [Techimply](https://www.techimply.com/profile/catproxies)

### Inferences
- The "no deposit" claim is technically accurate for the Starter tier: you can deposit any amount and begin reselling. This is genuinely low-barrier compared to competitors like DataImpulse ($50 minimum) or Proxy-Cheap ($5 minimum).
- The company is very young (founded ~2025 in Romania), which explains the limited independent coverage and small review footprint.
- The "freelancer team" descriptor, if accurate, suggests a small operation that could face scaling or continuity challenges.

### Gaps
- No independent confirmation of how the white-label storefront actually looks or functions. No screenshots or demo storefront were found from real resellers.
- Whether the storefront handles payment processing for end customers or if the reseller must handle that themselves is unclear from public documentation.
- The exact branding customization options (custom domain, logo, colors, CSS) are not documented publicly beyond the claim that "end users do not see CatProxies."


## What is the pricing structure? Per-GB rates, volume discounts, commission model

### Takeaway
CatProxies' retail pricing is competitive, especially at higher volumes (residential from $2.50/GB at 1 GB down to $1.26/GB at 250 GB). Reseller wholesale rates are not publicly disclosed but are described as "up to 70% off retail." The model is pure wholesale markup (no commission/revenue share), which gives resellers maximum pricing flexibility but also requires them to manage their own margin strategy.

### Cited Findings
- Standard Residential retail: $2.50/GB (1 GB), $2.00 (5 GB), $1.80 (10 GB), $1.60 (25 GB), $1.50 (50 GB), $1.40 (100 GB), $1.26 (250 GB). -- [CatProxies Pricing](https://catproxies.com/pricing)
- Mobile retail: $4.50/GB (1 GB), $4.25 (5 GB), $4.00 (10 GB), $3.75 (25 GB), $3.50 (50 GB), $3.25 (100 GB). -- [CatProxies Pricing](https://catproxies.com/pricing)
- Mobile Unlimited: $250/day (1 day), $100/day (7 days), $70/day (30 days). -- [CatProxies Pricing](https://catproxies.com/pricing)
- Static ISP: from $1.50/IP/day to $4.00/IP/month. Dedicated ISP: from $5.00/IP (1 IP) to $1.90/IP (100 IPs). -- [CatProxies Pricing](https://catproxies.com/pricing)
- Datacenter: from $4/day; pricing varies by IP count, bandwidth, and plan length. -- [CatProxies Pricing](https://catproxies.com/pricing)
- Reseller dashboard mockup shows an illustrative average margin of 28.4%, but this is not guaranteed. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Market comparison: typical residential average is $3.47/GB (median $3.00); typical mobile average is $4.19/GB. -- [Proxies.fi Benchmarks](https://proxies.fi/provider/catproxies)
- DataImpulse offers $1/GB residential wholesale; IPRoyal offers ~$1.75/GB pay-as-you-go; Geonode offers from $0.50/GB at high volume. -- [DataImpulse Blog](https://dataimpulse.com/blog/best-proxy-reseller-programs/)

### Inferences
- At retail, CatProxies' residential pricing ($2.50/GB entry, $1.26/GB at 250 GB) is competitive but not the cheapest in the market. Budget providers like Geonode and FlameProxies undercut it.
- If "up to 70% off retail" means the best wholesale rate on the $2.50 base, that would be ~$0.75/GB -- extremely competitive if achievable. More likely, 70% off applies only to the highest-priced products or smallest tiers.
- The pure-markup model (no revenue share) is more profitable for high-volume resellers but requires upfront capital to deposit.
- Compared to IPRoyal (~$1.75/GB reseller) and DataImpulse ($1/GB wholesale), CatProxies' wholesale rates would need to be in the $1-1.50/GB range to be competitive for residential. The exact rates are not public.

### Gaps
- Exact reseller wholesale rates per tier are not publicly disclosed; they are shared only after approval.
- Whether the "up to 70% off" applies to all products or only specific ones is unclear.
- No information on how pricing compares at identical volume levels against specific competitors' reseller tiers.


## What storefront/dashboard do they provide? Branding options?

### Takeaway
CatProxies provides a reseller panel with extensive order management features (custom plans, bandwidth adjustments, customer tagging, analytics) and a hosted white-label storefront option. However, the specific branding customization capabilities (custom domains, logo uploads, color themes) are not documented in detail on their public pages.

### Cited Findings
- Reseller panel features: create custom-size plans at any GB amount, add/deduct bandwidth on live plans, delete plans with unused bandwidth refunded, exchange remaining bandwidth between residential and mobile, reset passwords. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Customer management: per-customer pages showing orders, spend, active plans, and transaction history. Customers can be grouped, tagged, filtered, and labeled. Bulk actions available. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Resellers can view any plan as their customer sees it. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Billing page logs purchases, extensions, deposits, and refunds. Business invoices issued in the reseller's company name and downloadable. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Real-time analytics: bandwidth tracking, active orders, plan status, per-sub-user usage. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- White-label hostname options: CNAME records (follows gateway changes automatically) or A records (more private, requires manual IP updates). Available for Custom ISP, Premium Residential, Mobile, and Standard Residential. -- [CatProxies API Docs](https://catproxies.com/docs)

### Inferences
- The reseller panel is surprisingly feature-rich for a young company: bandwidth exchange, automatic refunds on expiration, custom-size plans, and customer-view simulation are advanced features.
- The white-label hostname support via CNAME/A records means resellers can serve proxies from their own domains (e.g., proxy.mybrand.com), which is a strong white-label feature.
- The storefront itself appears to be a pre-built frontend that resellers can deploy with their branding, but no public demo or detailed customization documentation exists.

### Gaps
- No public documentation on storefront visual customization (logo, colors, CSS, layout).
- No public demo or screenshot of the white-label storefront from a real reseller's perspective.
- Whether the storefront handles end-customer payment processing (Stripe, PayPal, crypto) or if the reseller must integrate their own payment gateway is not documented.


## What proxy types and quality? Pool size, geo coverage

### Takeaway
CatProxies offers six proxy types with claimed pool sizes that are typical for mid-tier providers. Independent benchmarks are mixed: editorial reviews rate IP reputation highly (85-90% trust scores), while one automated benchmark (proxies.fi) rated them F with 72% success rate. The discrepancy suggests performance may vary significantly by use case and target.

### Cited Findings
- Residential: 20M+ IPs, 194-195 countries, country/state/city/ASN targeting, rotating and sticky sessions (up to 24h), HTTP and SOCKS5. -- [CatProxies Pricing](https://catproxies.com/pricing)
- Mobile: 5M+ IPs, 150+ countries, 4G/5G, country/state/ISP targeting, HTTP and SOCKS5. -- [CatProxies Pricing](https://catproxies.com/pricing)
- Datacenter: 15,000 IPs in 5 countries, 10 Gbps network. -- [CatProxies Pricing](https://catproxies.com/pricing)
- Static ISP: 680,000+ IPs in 62 countries. -- [CatProxies Pricing](https://catproxies.com/pricing)
- IPv6: 14 /32 subnets, 1 Gbps unlimited plans. -- [CatProxies Pricing](https://catproxies.com/pricing)
- Dedicated ISP: private static IPs in USA, UK, Canada, Australia, Germany, Hong Kong, Japan. -- [CatProxies Pricing](https://catproxies.com/pricing)
- IPFighter tested 10 residential IPs: 9 of 10 scored 80%+ trust, were not flagged as proxies, and were blacklist-clean. One Hong Kong IP was flagged as datacenter and geolocated to Texas. -- [IPFighter Review](https://ipfighter.com/catproxies-review)
- ProxyVero benchmark (Nov 2026) showed 100% success on Amazon, Google, Instagram, Nike, TikTok, and generic HTTPS, with average latency ~2,325 ms. -- [ProxyVero Review](https://www.proxyvero.com/providers/catproxies/review/)
- Proxies.fi grades CatProxies F (15/100) with 72% success rate and ~$2.50/GB. -- [Proxies.fi](https://proxies.fi/provider/catproxies)
- Proxies.fi top performers: Bright Data 99.8%, Massive 99.9%, FlameProxies 99.4%, Vital-Data 99.7%. -- [Proxies.fi](https://proxies.fi/provider/catproxies)

### Inferences
- The extreme divergence between ProxyVero (100% success on 6 targets) and proxies.fi (72% overall) suggests methodology differences. Proxies.fi likely tests a broader range of harder targets at higher volume, while ProxyVero tested 6 popular sites.
- CatProxies likely resells upstream residential pools (Synthient classifies them as such), which means pool quality depends on their upstream provider. This is common among smaller providers.
- The 20M+ residential IP claim is plausible for a reseller pooling multiple upstream sources, but unique IPs available at any given time would be much smaller.
- Speed varies wildly (533 kbps to 217 Mbps in IPFighter tests), which is typical of residential proxies but worse than enterprise-grade providers.

### Gaps
- No reliable independent data on unique IP count vs. total pool size.
- No independent geo-accuracy testing beyond IPFighter's 10-IP sample.
- Which upstream provider(s) CatProxies sources from is not publicly known.
- No independent mobile proxy benchmarks found.


## What is the API like? Documentation, endpoints, integration ease

### Takeaway
The API is well-documented and comprehensive for a young provider, covering order management, proxy generation, targeting, billing, and white-label hostname management. It uses standard REST conventions with Bearer token auth. The 1,000 req/hour rate limit is adequate for small-to-mid-scale operations but could constrain large resellers with many concurrent customers.

### Cited Findings
- Base URL: `https://catproxies.com/api/v1/public`. Auth: Bearer token (keys start with `cp_` + 64 hex chars). -- [CatProxies Docs](https://catproxies.com/docs)
- Rate limits: 1,000 requests/hour per key (general), 60 per 10 minutes for mobile targeting endpoints, 1 per 10 seconds for custom plan creation, 1 per 30 seconds for exchange operations. -- [CatProxies Docs](https://catproxies.com/docs)
- Key endpoints: account info, billing history, store/products listing, order creation/management/renewal, proxy credential retrieval, password reset, IP whitelisting, usage stats, bandwidth adjustment, plan exchange, customer tagging, white-label hostname management, targeting options. -- [CatProxies Docs](https://catproxies.com/docs)
- Error handling includes specific guidance against blind retries on purchases (409/500 responses may indicate payment already processed). -- [CatProxies Docs](https://catproxies.com/docs)
- Products supported via API: Standard Residential (v2resi), Premium Residential (resix), Custom ISP (Isp), Dedicated ISP (IspP), Datacenter (DatacenterP), IPv6 (Ipv6p), Rotating Mobile (RotatingMobile). -- [CatProxies Docs](https://catproxies.com/docs)
- No SDKs documented. Compatible with Selenium, Puppeteer, Playwright, Python, JavaScript. -- [ProxyVero Review](https://www.proxyvero.com/providers/catproxies/review/)
- CatProxies claims "setup in under 5 minutes" and "API integration typically takes a few days for experienced developers." -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)

### Inferences
- The API covers the full lifecycle a reseller needs: provisioning, management, targeting, billing, and white-labeling. This is above average for a provider of this size.
- The lack of SDKs means resellers must write their own HTTP client wrappers, which adds development overhead but is common for smaller providers.
- The 1,000 req/hour limit means roughly 16 requests per minute, which is sufficient for managing dozens of customers but could be tight for automated high-frequency operations.
- The payment-safety warnings in the docs (don't retry purchase calls blindly) suggest the team has encountered real-world issues and documented them, which is a positive sign of operational maturity.

### Gaps
- No webhook or callback support documented (resellers must poll for status changes).
- No sandbox/test environment mentioned for development.
- No API changelog or versioning policy documented.


## Who is behind CatProxies? Company background, reputation

### Takeaway
CatProxies is operated by MOTAN DATA SOLUTIONS S.R.L., a Romanian company registered in Targoviste. The company appears to have been founded in 2025, making it approximately 1-1.5 years old. No founders or team members are publicly named. The company is very small and new, with limited independent verification of its claims. This is itself a significant signal about maturity and risk.

### Cited Findings
- Legal entity: MOTAN DATA SOLUTIONS S.R.L., CUI 51601374, VAT RO51601374, registered in Targoviste, Romania. -- [CatProxies Ethics Page](https://dashboard.catproxies.com/ethics)
- Trustpilot profile claimed July 2024, company listed as based in the United States. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- Techimply directory lists founding year as 2025 and team as "freelancers." -- [Techimply](https://www.techimply.com/profile/catproxies)
- No founders, executives, or team members are publicly named on the website or ethics page. -- [CatProxies Ethics Page](https://dashboard.catproxies.com/ethics)
- Synthient classifies CatProxies as a reseller of upstream residential pools and labels it "high risk." -- [Synthient](https://synthient.com/context/provider/catproxies)
- Ethics page effective date: 6 May 2026, which indicates relatively recent formalization of policies. -- [CatProxies Ethics Page](https://dashboard.catproxies.com/ethics)
- Support hours: Monday to Friday, 9:00-17:00 EET. Channels: email, live chat, Telegram, Discord, X. -- [CatProxies Ethics Page](https://dashboard.catproxies.com/ethics)

### Inferences
- The discrepancy between "based in the United States" (Trustpilot) and Romanian registration is concerning. It could mean the company operates from Romania but markets to a US audience, or the Trustpilot location is inaccurate.
- A founding date of ~2025 means CatProxies has roughly 1-1.5 years of track record. This is very young for a company you'd be building a business dependency on.
- No named team members + "freelancer" team descriptor + Romanian SRL registration = this is likely a very small operation, possibly 1-5 people.
- The Synthient "high risk" classification and the reseller-of-upstream-pools designation mean CatProxies likely doesn't own its own IP infrastructure, adding a layer of dependency.

### Gaps
- No LinkedIn profiles, team bios, or founder information found anywhere.
- Company revenue, funding status, and employee count are unknown.
- No information on what upstream providers CatProxies resells from.
- Whether MOTAN DATA SOLUTIONS S.R.L. has any other products or history is unknown.


## What do real users say? Reviews across platforms

### Takeaway
User review volume is very low (57 Trustpilot, ~9 G2), which is expected for a company this young. The reviews show a polarized pattern: mostly enthusiastic 5-star reviews (88% on Trustpilot) praising speed and support, with a small cluster of 1-star reviews (10%) reporting IP failures, poor support, and alleged review manipulation. The low volume and polarized distribution make it difficult to draw confident conclusions about real-world reliability.

### Cited Findings
- Trustpilot: 4.4/5 TrustScore from 57 reviews. Star distribution: 88% 5-star, 2% 4-star, 0% 3-star, 0% 2-star, 10% 1-star. 18 reviews in the last 12 months. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- G2: 4.7/5 from ~9 verified reviews. One reviewer noted CatProxies costs more than competitors but judged quality higher. -- [G2](https://www.g2.com/sellers/cat-proxies)
- Positive themes: proxy quality and speed, affordable pricing, responsive support, real residential IPs (not flagged as datacenter), easy setup. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- Negative: "proxies dies randomly," "about 50% of IPs failed on target sites," support called "toxic," alleged offer of payment for 5-star reviews, refund refused and account blocked after complaint. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- Company responds to 100% of negative reviews within ~2 weeks, disputing claims and providing public chat transcripts. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- Review authenticity concerns: many reviewers have only 1-2 reviews, several 5-star reviews are very short/generic, clustering of reviews in May 2026, one reviewer mentions "four star review" in a 5-star review suggesting solicitation, only one review carries a "Verified" label. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- BlackHatWorld: vendor promotional thread exists but was blocked from full access. The post lists products and prefers buyers with established forum reputation. -- [BlackHatWorld](https://www.blackhatworld.com/posts/18179051/)
- No Reddit discussions about CatProxies were found in searches. -- Search results
- No Proxyway review of CatProxies was found. -- Search results

### Inferences
- The 88%/0%/0%/0%/10% distribution (almost all 5-star or 1-star, nothing in between) is a pattern sometimes associated with review solicitation/manipulation, though it can also occur naturally with small sample sizes.
- The allegation of being offered payment for reviews is serious but unverified; the company denies it.
- The complete absence from Reddit and major review platforms (Proxyway) indicates CatProxies is not yet on the radar of the broader proxy community.
- The low total review count (66 across Trustpilot and G2 combined) after ~1.5 years of operation suggests a small customer base.

### Gaps
- No Proxyway review exists (Proxyway is considered the most authoritative proxy review site).
- No Reddit threads discussing CatProxies were found.
- Full BlackHatWorld thread content could not be accessed (403 error).
- Capterra reviews could not be verified (multiple regional sites listed but none with accessible reviews).


## What are the limitations? What can't you do vs other providers?

### Takeaway
Key limitations include: restricted product access below $100 deposit, a 1,000 req/hour API rate limit, no SDKs, limited refund policy (datacenter/ISP/unlimited/IPv6 are non-refundable), no city-level mobile targeting, and the fundamental risk of building on a very young, small, unproven company. The hosted storefront's payment handling and customization depth are unclear.

### Cited Findings
- Datacenter, IPv6, and Static ISP products require a minimum $100 deposit (Plus tier). -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- API rate limit: 1,000 requests/hour per key. -- [CatProxies Docs](https://catproxies.com/docs)
- No SDKs documented. -- [CatProxies Docs](https://catproxies.com/docs)
- Datacenter, ISP, unlimited residential, and IPv6 products are non-refundable. -- [IPFighter Review](https://ipfighter.com/catproxies-review)
- No city-level mobile targeting. -- [JoshWP Review](https://joshwp.com/catproxy-review/)
- Datacenter IPs have low trust on hardened targets. -- [JoshWP Review](https://joshwp.com/catproxy-review/)
- Mobile pool (5M+) is smaller than the 10-20M+ pools of the largest specialists (Oxylabs, Decodo, ProxyEmpire). -- [JoshWP Review](https://joshwp.com/catproxy-review/)
- Rotation-control API documentation is thinner than enterprise leaders'. -- [JoshWP Review](https://joshwp.com/catproxy-review/)
- Refund terms are less explicit than some competitors'. -- [JoshWP Review](https://joshwp.com/catproxy-review/)
- Support hours: Monday-Friday, 9:00-17:00 EET only. -- [CatProxies Ethics Page](https://dashboard.catproxies.com/ethics)

### Inferences
- The Monday-Friday support hours are a significant limitation for a reseller serving end customers who expect 24/7 support. If a customer's proxy breaks on Saturday, the reseller is on their own.
- The non-refundable policy on most product types shifts risk to the reseller, who may need to absorb costs from customer complaints.
- The $100 minimum to access datacenter/ISP products means the "no deposit" claim only applies to residential and mobile reselling.
- As a reseller of upstream pools, CatProxies cannot guarantee pool composition, ASN diversity, or geo-accuracy the way first-party providers can.

### Gaps
- Maximum concurrent connections per reseller account is not documented.
- Whether resellers can set usage limits per end customer is unclear.
- SLA terms for resellers (uptime guarantee, compensation) are not publicly documented.


## Is this a serious platform or a small/unknown player? How established are they?

### Takeaway
CatProxies is a small, new player in the proxy market. It is a legitimate registered business with a real product, but by every available metric -- company age (~1.5 years), review count (66 total), team size (likely very small), absence from major review platforms, and "high risk" classification by Synthient -- it is unestablished. The feature set is surprisingly polished for its size, but the business risk is high compared to established providers.

### Cited Findings
- Founded approximately 2025 based on directory listing (weak source). -- [Techimply](https://www.techimply.com/profile/catproxies)
- Trustpilot profile claimed July 2024. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- Total review count: ~57 (Trustpilot) + ~9 (G2) = ~66 reviews across all platforms found. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com); [G2](https://www.g2.com/sellers/cat-proxies)
- Not reviewed by Proxyway (the most authoritative proxy review site). -- Search results
- Not included in DataImpulse's ranking of 8 best reseller programs (only mentioned in passing for white-label support). -- [DataImpulse Blog](https://dataimpulse.com/blog/best-proxy-reseller-programs/)
- Proxies.fi grades it F (15/100), placing it below all 15 providers in their benchmark. -- [Proxies.fi](https://proxies.fi/provider/catproxies)
- Synthient classifies it as high risk and a reseller of upstream pools. -- [Synthient](https://synthient.com/context/provider/catproxies)
- No Reddit community presence found. -- Search results
- Capterra listings exist across multiple regions (India, UK, Australia, etc.) but with no visible reviews. -- Search results

### Inferences
- Comparison to established players: Bright Data (founded 2014), Oxylabs (2015), IPRoyal (2020), Rayobyte (2015) all have years of track record, thousands of reviews, and extensive independent coverage. CatProxies has none of this.
- The polished website, comprehensive API docs, and multi-product catalog suggest competent development work, but this doesn't guarantee business stability or long-term viability.
- The DataImpulse article's dismissive treatment (two passing mentions despite covering 8 reseller programs in depth) reflects CatProxies' marginal market position.
- For a reseller building a business on top of CatProxies, the key risk is dependency on a young, small company that could change terms, experience financial difficulties, or shut down.

### Gaps
- No financial information (revenue, funding, profitability) available.
- No information on total customer count or active reseller count.
- No web archive analysis was performed to track site evolution over time.


## What is the scalability path? Can you grow from small to large?

### Takeaway
CatProxies' tier structure (Starter -> Plus -> Volume) provides a clear growth path with improving rates at higher deposit levels, and the API supports programmatic management. However, scalability is ultimately constrained by CatProxies' own scale as a young, small company, the 1,000 req/hour API limit, and the dependency on unknown upstream providers.

### Cited Findings
- Tier progression: Starter (any deposit) -> Plus ($100+) -> Volume ($500+), with pricing improving at each tier. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Custom pricing available for high-volume teams; contact sales for tailored plans. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Large one-off orders receive same-day quotes through Contact Sales. -- [CatProxies Pricing](https://catproxies.com/pricing)
- Volume tier includes custom invoicing and higher platform balance limit. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- API supports full lifecycle management: order creation, bandwidth adjustment, plan exchange, customer management. -- [CatProxies Docs](https://catproxies.com/docs)

### Inferences
- The low entry barrier (any deposit) and progressive tier structure make it easy to start small and test, then scale up if results are satisfactory.
- However, at true scale (thousands of end customers, TB-level monthly bandwidth), you'd want guarantees about: uptime SLAs, dedicated account management, higher API rate limits, and pool quality guarantees that CatProxies cannot yet credibly provide based on their market position.
- A realistic use case: CatProxies could work for a small reseller doing $500-5,000/month. Beyond that, the risk of relying on a company this young increases significantly.

### Gaps
- No information on CatProxies' largest resellers or maximum volume they've supported.
- No published SLA details for reseller accounts.
- Whether API rate limits can be increased for Volume-tier resellers is not documented.


## Support quality and responsiveness

### Takeaway
Support receives mixed signals: editorial reviews and many Trustpilot reviewers praise responsiveness and helpfulness, while a minority of Trustpilot reviewers describe it as "toxic." Support hours are limited to business hours (Mon-Fri, 9:00-17:00 EET), which is a notable constraint for resellers. Multiple support channels exist (email, chat, Telegram, Discord).

### Cited Findings
- Support channels: email, live chat, Telegram, Discord, X (Twitter). -- [CatProxies Ethics Page](https://dashboard.catproxies.com/ethics)
- Support hours: Monday-Friday, 9:00-17:00 EET. -- [CatProxies Ethics Page](https://dashboard.catproxies.com/ethics)
- Starter tier includes a dedicated account manager and full API access. Plus and Volume tiers include priority support. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- JoshWP review gives docs and support 8.6/10 but includes no first-hand support experience. -- [JoshWP Review](https://joshwp.com/catproxy-review/)
- Trustpilot positive reviews frequently praise responsive, helpful support. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- Two 1-star Trustpilot reviewers called support "toxic" and "bad." One claimed being offered payment for a 5-star review. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- CatProxies replies to 100% of negative Trustpilot reviews within ~2 weeks with detailed rebuttals. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)

### Inferences
- The detailed rebuttals to negative reviews suggest the company takes reputation seriously, but the small team may struggle with tone/diplomacy in contentious cases (hence the "toxic" complaint).
- Business-hours-only support (EET timezone) means no coverage during US afternoon/evening or weekends. This is a real operational risk for resellers with US-based end customers.
- The "dedicated account manager" for all resellers sounds good but likely means a single point of contact on a small team, not a dedicated resource.

### Gaps
- No independent SLA data on response times.
- Whether weekend/holiday support exists for Volume-tier resellers is unclear (one source says "daily including holidays," another says Mon-Fri only).


## Any known issues, complaints, or red flags?

### Takeaway
Several red flags exist: (1) the proxies.fi F grade with 72% success rate; (2) Synthient "high risk" classification; (3) the company's extreme youth and anonymity; (4) review authenticity concerns on Trustpilot; (5) allegations of offering payment for reviews; (6) refund refusal and account blocking reported by a user; and (7) inconsistent information between different sources (US vs. Romania location). None of these are individually disqualifying, but collectively they warrant caution.

### Cited Findings
- Proxies.fi F grade (15/100), 72% success rate: "We cannot recommend CatProxies at its current reliability level." -- [Proxies.fi](https://proxies.fi/provider/catproxies)
- Synthient classification: reseller of upstream pools, high risk. -- [Synthient](https://synthient.com/context/provider/catproxies)
- Trustpilot review (Kaimur Kumar): "about 50% of IPs failed," refund refused, account access blocked after raising concerns. Company disputes this, says proxies tested working and blames misconfiguration. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- Trustpilot review (Janis): alleges being offered payment for 5-star review. Company denies this. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)
- Trustpilot location listed as "United States" while company is registered in Romania. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com); [CatProxies Ethics Page](https://dashboard.catproxies.com/ethics)
- IPFighter test: 1 of 10 residential IPs was incorrectly flagged as datacenter and geolocated to wrong country (requested Hong Kong, got Texas). -- [IPFighter Review](https://ipfighter.com/catproxies-review)
- No team members publicly identified despite being ~1.5 years old. -- Research finding
- Review distribution (88% 5-star, 10% 1-star, nothing in between) and many single-review accounts raise manipulation concerns. -- [Trustpilot](https://www.trustpilot.com/review/catproxies.com)

### Inferences
- The pattern of account blocking after a refund dispute (if true) is a serious red flag for resellers who may need to advocate on behalf of their end customers.
- The US/Romania location discrepancy is not necessarily deceptive (many companies market from one country while being legally registered in another) but it reduces transparency.
- The combination of anonymous leadership + young company + reseller-of-upstream-pools + high-risk classification creates meaningful counterparty risk.
- None of the red flags are proven fraud indicators, but the totality suggests treating CatProxies as a high-risk partner requiring careful due diligence and testing before commitment.

### Gaps
- Whether the review-payment allegation is true cannot be verified from available sources.
- No legal complaints, lawsuits, or regulatory actions were found (though this could be because the company is too new/small to have attracted them).
- Full BlackHatWorld thread could not be accessed to check for user complaints there.


## How does their hosted storefront compare to building your own platform?

### Takeaway
The hosted storefront removes the need to build a frontend, integrate payments, and manage infrastructure, which is a significant time/cost savings. However, the lack of public documentation on customization depth, payment handling, and branding flexibility makes it impossible to compare meaningfully against a custom-built platform. The storefront is best understood as a quick-start option for resellers without development resources, not as a feature-complete alternative to a custom build.

### Cited Findings
- Storefront is hosted on CatProxies infrastructure with reseller's branding. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- End users do not see CatProxies in auth, dashboards, or invoices. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- White-label hostname management available via API for CNAME/A record setup. -- [CatProxies Docs](https://catproxies.com/docs)
- DataImpulse also offers "true own-domain white-label" with $50 minimum deposit. -- [DataImpulse Blog](https://dataimpulse.com/blog/best-proxy-reseller-programs/)
- Proxy-Cheap offers white-label setups as well. -- [DataImpulse Blog](https://dataimpulse.com/blog/best-proxy-reseller-programs/)

### Inferences
- For a reseller without development resources, the hosted storefront is a significant advantage -- you can start selling immediately without building a frontend or payment integration.
- For a reseller with development resources (like the user building proxy_ui), the API is likely more valuable than the storefront, as it allows building a custom frontend with full control over UX, branding, and features.
- The storefront option is most valuable as a way to validate the business model quickly before investing in custom development.
- The competition (DataImpulse, Proxy-Cheap) also offers white-label, so this is not unique to CatProxies, though the "any deposit" entry point may be the lowest barrier.

### Gaps
- No public documentation on what the storefront actually looks like or how much can be customized.
- Whether the storefront supports custom payment gateways (Stripe, PayPal, crypto) is unknown.
- Whether you can switch from storefront to API-only (or vice versa) without disruption is unknown.
- Whether end customers' data (emails, usage patterns) is accessible to the reseller is unknown.


## Do they handle billing/payments for you or do you handle your own?

### Takeaway
CatProxies handles the billing relationship between the reseller and CatProxies (deposits, balance management, invoicing in the reseller's company name). Whether the hosted storefront handles end-customer payment collection is not clearly documented. If using the API to build a custom frontend, the reseller would handle their own end-customer billing.

### Cited Findings
- Billing page in reseller panel logs purchases, extensions, deposits, and refunds. Business invoices issued in the reseller's company name and downloadable. -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- Deposits accepted via crypto (BTC, ETH, USDT) and apparently fiat (Visa, Mastercard, Apple Pay, Google Pay, SBP, Pix). -- [CatProxies Reseller Page](https://catproxies.com/become-reseller); [CatProxies Ethics Page](https://dashboard.catproxies.com/ethics)
- Custom invoicing available on Volume tier ($500+). -- [CatProxies Reseller Page](https://catproxies.com/become-reseller)
- The API includes order management and billing endpoints but not end-customer payment collection endpoints. -- [CatProxies Docs](https://catproxies.com/docs)

### Inferences
- The reseller model appears to be: (1) reseller deposits funds with CatProxies, (2) reseller provisions plans for end customers from their balance, (3) reseller collects payment from end customers separately. This is a standard wholesale model.
- If using the hosted storefront, it's possible that CatProxies handles end-customer payment collection, but this is not confirmed.
- For API-based resellers, you definitely need your own payment processing (Stripe, PayPal, crypto) and customer billing system.
- The "invoices in reseller's company name" feature is useful for B2B resellers but doesn't address the end-customer payment question.

### Gaps
- The most critical gap: whether the hosted storefront collects end-customer payments, and if so, whether funds flow through CatProxies or directly to the reseller.
- What payment methods end customers can use on the hosted storefront is unknown.
- Whether resellers receive detailed billing reports for tax/accounting purposes is unclear.


## Overall Assessment and Comparison Context

### Summary for Report Writer

**What CatProxies is:** A young (~1.5 year old) Romanian company that resells upstream proxy pools under a white-label model. It offers a surprisingly feature-rich reseller program with low entry barriers (any deposit amount), a comprehensive API, and a hosted storefront option.

**Strengths:**
1. Lowest entry barrier found: any deposit amount to start reselling (vs. $50 DataImpulse, $5 Proxy-Cheap)
2. Full white-label: hosted storefront + API + custom hostnames
3. Comprehensive product range: residential, mobile, datacenter, ISP, IPv6
4. No revenue share / no price floor: resellers keep full markup
5. Non-expiring deposits and bandwidth
6. Feature-rich reseller panel (bandwidth exchange, auto-refunds, customer management)
7. Well-documented API covering full order lifecycle

**Weaknesses:**
1. Extremely young and small company with anonymous leadership
2. Reseller of upstream pools (not a first-party network owner)
3. Mixed reliability: F grade from proxies.fi (72% success), though editorial reviews are more positive
4. Very low review volume (66 total across platforms) with authenticity concerns
5. Not covered by Proxyway or included in major reseller program rankings
6. Business-hours-only support (Mon-Fri EET)
7. Non-refundable policy on most product types
8. Synthient "high risk" classification
9. Unclear storefront customization and payment handling

**Comparison to other reseller providers:**
- vs. DataImpulse: DataImpulse has a $50 minimum but $1/GB wholesale residential, a more established reputation, and is ranked #1 in reseller program comparisons. CatProxies' wholesale rates are not public but likely higher.
- vs. IPRoyal: IPRoyal (~$1.75/GB reseller) has a much larger customer base and years of track record. CatProxies' storefront option is unique, but IPRoyal is more battle-tested.
- vs. Bright Data: Bright Data offers a 50% revenue share model (no capital needed) aimed at enterprises. Completely different model and scale; not directly comparable.
- vs. Proxy-Cheap: Both offer white-label, but Proxy-Cheap has $5 minimum and established market position.

**Viability verdict:** CatProxies is a real product with genuine features, but it is a high-risk partner choice. It is best suited for: (1) testing the reseller business model with minimal capital, (2) resellers who value the hosted storefront and don't want to build their own frontend, or (3) as a secondary/backup provider alongside a more established primary. It should NOT be a sole business dependency at this stage of its maturity.
