# Technical Requirements for a Proxy Reselling Platform (2025-2026)

## What APIs do major proxy providers offer for resellers?

### Takeaway
Most major proxy providers offer a reseller/sub-user API with endpoints for customer provisioning, zone/product management, bandwidth tracking, and IP allocation. The API maturity and documentation quality vary significantly across providers.

### Cited Findings

**Bright Data (Reseller API)**
- Base URL pattern: `https://api.brightdata.com/{partner}/...` where each reseller has a unique partner segment — [Bright Data Reseller API Docs](https://docs.brightdata.com/api-reference/reseller-api/introduction)
- Key endpoints include: `POST /{partner}/signup_shared_isp` (create user), `GET /{partner}/proxies?id={id}&zone={zone}` (generate proxy credentials), `PUT /{partner}/usage_limit` (change zone usage limit), `PUT /{partner}/zone` (create zone with allocated IPs), `GET /{partner}/ips` (get allocated IPs) — [Bright Data Reseller API Docs](https://docs.brightdata.com/api-reference/reseller-api/introduction)
- 25+ reseller endpoints documented including: Allocate IPs, Create/Change Usage Webhook Settings, Disable/Enable Customer, Generate Proxy/SOCKS5 Username, Get Available ASNs/Carriers/Countries, Get Customers, Refresh IPs, Remove Allocated IPs, and multiple sign-up endpoints per product type (Dedicated DC, Dedicated ISP, Residential, Shared DC, Smart Residential, Shared ISP) — [Bright Data Reseller API Docs](https://docs.brightdata.com/api-reference/reseller-api/introduction)
- API keys use five permission levels: Admin, Finance, Ops, Limit, and User, with optional expiration dates — [Bright Data API Introduction](https://docs.brightdata.com/api-reference/introduction)
- Each partner can create up to 50 zones per user; new users default to a 1KB usage limit; proxy access uses superproxy ports 31111-31114 for load balancing — [Bright Data Reseller API Docs](https://docs.brightdata.com/api-reference/reseller-api/introduction)
- Usage limits use the format `{"limit": {"limit": N}}` with unit `bw` (bandwidth in bytes) — [Bright Data Reseller API Docs](https://docs.brightdata.com/api-reference/reseller-api/introduction)

**Oxylabs (Dashboard API)**
- Base URL: `https://api.oxylabs.io/` — [Oxylabs Dashboard API](https://developers.oxylabs.io/dashboard/dashboard-api)
- Authentication: Bearer token in Authorization header; API key issued by Oxylabs Customer Support — [Oxylabs Dashboard API](https://developers.oxylabs.io/dashboard/dashboard-api)
- Endpoints: `GET /stats/v1/filters/instances` (list product instances), `GET /stats/v1/usage` (aggregated usage stats with `traffic_bytes` field) — [Oxylabs Dashboard API](https://developers.oxylabs.io/dashboard/dashboard-api)
- Usage queries support grouping by DAY, TARGET, and other dimensions, limited to 31-day date ranges per request — [Oxylabs Dashboard API](https://developers.oxylabs.io/dashboard/dashboard-api)
- No sub-user management endpoints documented publicly; reseller-specific features require contacting support — [Oxylabs Dashboard API](https://developers.oxylabs.io/dashboard/dashboard-api)

**Smartproxy / Decodo (API)**
- Base URL: `https://api.decodo.com/v1` — [Smartproxy API GitHub](https://github.com/Smartproxy/Smartproxy-API)
- Authentication: `POST /auth/` returns `user_id` and `token` — [Smartproxy API GitHub](https://github.com/Smartproxy/Smartproxy-API)
- Sub-user endpoints: `GET /users/{userId}/sub-users` (list), `POST /users/{userId}/sub-users` (create with username, password, service_type), `PUT /users/{userId}/sub-users` (update traffic_limit/password), `DELETE /users/{userId}/sub-users/{sub_user_id}` — [Smartproxy API GitHub](https://github.com/Smartproxy/Smartproxy-API)
- Traffic tracking: `GET /users/{userId}/sub-users/{username}/traffic` — [Smartproxy API GitHub](https://github.com/Smartproxy/Smartproxy-API)
- Additional endpoints: whitelisted IPs management (list/add/delete), subscription details, endpoint listing by type — [Smartproxy API GitHub](https://github.com/Smartproxy/Smartproxy-API)
- Code examples available in Java, PHP, and C# — [Smartproxy API GitHub](https://github.com/Smartproxy/Smartproxy-API)

**Evomi (Reseller API)**
- Endpoints: `/v2/reseller/my_info`, `/v2/reseller/proxy_settings_v2`, `/v2/reseller/subusers` (list/create), `/v2/reseller/subusers/{username}` (details), `/v2/reseller/subusers/{username}/balance` (add balance), `/v2/reseller/subusers/{username}/usage` (usage history) — [Evomi Reseller Integration Guide](https://proxy-docs.evomi.com/reseller-guide-integration)
- White-label CNAME support: resellers can point a subdomain (e.g., `proxies.yourdomain.com`) to Evomi endpoints via CNAME record — [Evomi Reseller Integration Guide](https://proxy-docs.evomi.com/reseller-guide-integration)
- Dedicated load balancer option (prepaid 6 months, EU or US) for A-record-based custom domains — [Evomi Reseller Integration Guide](https://proxy-docs.evomi.com/reseller-guide-integration)

**IPRoyal**
- Reseller API authorized via `X-Access-Token` header, found/regenerated in dashboard Settings under API section — [IPRoyal API Docs](https://docs.iproyal.com/proxies/isp-static-residential/api)
- Legacy API deprecated as of September 15, 2025 — [IPRoyal API Docs](https://docs.iproyal.com/proxies/isp-static-residential/api)
- API covers ISP static residential, Sneaker DC, and mobile proxy products — [IPRoyal API Docs](https://docs.iproyal.com/proxies/isp-static-residential/api)
- Sneaker DC API activation requires identity verification — [IPRoyal Sneaker DC API](https://docs.iproyal.com/proxies/sneaker-dc/api)

**NetNut**
- Offers a reseller program with sub-user creation, usage management, and client access control to the residential proxy network — [NetNut Resellers](https://netnut.io/resellers/)
- Pricing plans start from $100, no pay-as-you-go — [Proxyway NetNut Review](https://proxyway.com/reviews/netnut-proxies)
- API top-ups reported to take over an hour to process in at least one reseller's experience — [Rampage Networks Note](https://rampagenetworks.notion.site/Netnut-Private-Resi-435e0518edb647a784dab07e3a2b82cf)

### Inferences
- Bright Data and Smartproxy/Decodo have the most feature-complete documented reseller APIs with sub-user management, bandwidth tracking, and IP allocation.
- Evomi's balance-based model (top up sub-user balance after payment) is the simplest integration pattern for new resellers.
- Oxylabs appears enterprise-focused with reseller features gated behind support contact rather than self-service API docs.
- A modular provider adapter layer is advisable since APIs differ significantly in authentication, endpoint structure, and provisioning models.

### Gaps
- Webshare, SOAX, and Infatica reseller API documentation could not be found through web search. Their reseller programs may exist but are not publicly documented.
- Bright Data's full reseller API authentication method is not documented on the introduction page; individual endpoint pages may contain this.
- NetNut's specific API endpoint documentation was not found publicly.


## What does a typical reseller API integration look like?

### Takeaway
The standard integration pattern involves creating sub-users via the provider API after payment, storing credentials, assembling proxy connection strings from fixed endpoints, and tracking usage through provider-reported bandwidth data.

### Cited Findings
- Evomi's recommended flow: (1) customer registers on your platform, (2) wait until customer pays, (3) create sub-user via API, (4) store the sub-user username (not the password), (5) retrieve and display password on demand, (6) confirm account ready — [Evomi Reseller Integration Guide](https://proxy-docs.evomi.com/reseller-guide-integration)
- Proxy credentials are assembled from stored username + password + fixed product endpoints and ports. Example: `username:password@rp.evomi.com:1000` (HTTP) or `:1002` (SOCKS) — [Evomi Reseller Integration Guide](https://proxy-docs.evomi.com/reseller-guide-integration)
- Product endpoint mapping is fixed per provider: Residential (rp.evomi.com:1000/1002), Datacenter (dcp.evomi.com:2000/2002), Mobile (mp.evomi.com:3000/3002) — [Evomi Reseller Integration Guide](https://proxy-docs.evomi.com/reseller-guide-integration)
- Targeting options (country, session type) are handled through password parameters rather than separate endpoints — [Evomi Reseller Integration Guide](https://proxy-docs.evomi.com/reseller-guide-integration)
- Bright Data uses a zone-based model: each product gets a zone with its own credentials; usage limits are set per zone in bytes — [Bright Data Reseller API Docs](https://docs.brightdata.com/api-reference/reseller-api/introduction)
- Smartproxy/Decodo uses a sub-user model with explicit traffic limits: `PUT /users/{userId}/sub-users` to set `traffic_limit` — [Smartproxy API GitHub](https://github.com/Smartproxy/Smartproxy-API)
- Traffic data is separate from the management path: customers connect directly to the provider's proxy servers, while the reseller's platform manages accounts through the provider's API — [Evomi Reseller Guide Overview](https://proxy-docs.evomi.com/reseller-guide-overview)

### Inferences
- There are two main provisioning models: (a) balance-based (Evomi: top up sub-user balance after purchase) and (b) zone/limit-based (Bright Data: create zones with byte-based usage caps).
- The "never store password" pattern from Evomi suggests providers may rotate or regenerate credentials, making on-demand retrieval important.
- A provider-agnostic abstraction layer should normalize these different models into a unified interface for the reseller platform.

### Gaps
- Webhook-based usage notification details for Bright Data (Create/Change Usage Webhook Settings endpoint) were not fully documented in the available pages.
- Real-time usage synchronization frequency and latency across providers was not found.


## Are there ready-made dashboard/panel solutions for proxy reselling?

### Takeaway
Several turnkey white-label reseller programs exist with built-in dashboards, and there is one dedicated WHMCS module for proxy reselling, but fully independent open-source proxy reseller panels are scarce.

### Cited Findings

**White-label reseller programs with dashboards**
- Proxidize: white-label customer dashboard, reseller management portal, usage analytics, dedicated partner success manager, margin control — [Proxidize Reseller Program](https://proxidize.com/partners/reseller-program)
- ProxyWing: reseller API for ordering/renewals/IP whitelisting, dashboard for resellers and clients to manage orders — [ProxyWing Reseller](https://proxywing.com/reseller)
- ColdProxy: dashboard and Reseller API, Customer Plan API for generating proxy lists — [ColdProxy Reseller](https://coldproxy.com/reseller/)
- CatProxies: reseller analytics dashboard, full reseller API access, white-label branding, approval within one business day — [CatProxies Reseller](https://dashboard.catproxies.com/become-reseller)
- MaskProxy: white-label system with customizable control panel, API integration, payment gateway, and domain — [MaskProxy Agent](https://maskproxy.io/become-an-agent.html)
- DataImpulse: Open API for building custom dashboards — [DataImpulse Reseller](https://scriptagc.wasmer.app/https_dataimpulse_com/reseller-program)

**WHMCS Integration**
- Proxy-Cheap Whitelabel Proxy Reseller Module: uses Proxy-Cheap API to sell four proxy types (Datacenter, Static Residential, Residential Rotating, Mobile) through WHMCS — [WHMCS Marketplace](https://marketplace.whmcs.com/product/6889-whitelabel-proxy-reseller-module)
- Latest release v1.2.1 (March 2026), compatible with WHMCS v8.13 — [WHMCS Marketplace](https://marketplace.whmcs.com/product/6889-whitelabel-proxy-reseller-module)
- Deposit as low as $5 to start — [WHMCS Marketplace](https://marketplace.whmcs.com/product/6889-whitelabel-proxy-reseller-module)
- For general reseller billing: WHMPRESS Products Reseller for WHMCS ($10/month or $39/year, open-source at $399 one-time) and ModulesGarden Resellers Center (supports WHMCS v9.0) — [WHMCS Marketplace](https://marketplace.whmcs.com/product/8273-product-reseller-for-whmcs); [ModulesGarden](https://modulesgarden.com/products/whmcs/resellers-center)

**Open-source proxy/VPN panels (adaptable for reselling)**
- Remnawave Panel: TypeScript/NestJS, built on Xray-core, ~5.3k GitHub stars, AGPL-3.0 licensed — [GitTrend Remnawave](https://gittrend.io/repo/remnawave/panel)
- PPanel: Go backend with user management, subscriptions, payments, orders, tickets, and node management; separate frontend — [GitHub PPanel](https://github.com/CrackerCat/PPanel-server); [GitHub PPanel Frontend](https://github.com/perfect-panel/frontend)
- V2Board: PHP/Laravel panel for multiple proxy protocols — [GitHub V2Board](https://github.com/Mintapril/v2board)
- ProxyPanel: web application managing multiple proxy protocols — [GitHub ProxyPanel](https://github.com/ProxyPanel)

**Custom builds**
- A Freelancer listing for a custom reverse-proxy reseller platform with dashboard was priced at $1,500-3,000 — [Freelancer](https://www.freelancer.cl/projects/network-administration/proxy-reseller-reverse-proxy-setup)

### Inferences
- For MVP, using a white-label program (Proxidize, ColdProxy, or MaskProxy) eliminates the need to build a dashboard from scratch.
- For a custom platform, open-source panels like PPanel or Remnawave provide a starting foundation for user management and billing, but would need adaptation for proxy-specific features (sub-user API integration, bandwidth metering, credential generation).
- WHMCS is a viable option only if the provider (Proxy-Cheap) is your upstream; no generic "connect any provider" WHMCS proxy module exists.

### Gaps
- Active maintenance status and community size for the open-source panels (PPanel, V2Board, ProxyPanel) was not verified.
- Whether the open-source panels support multiple upstream proxy providers simultaneously was not determined.


## What payment processing integrations are common?

### Takeaway
Stripe is the dominant payment processor for proxy reseller platforms, with PayPal as secondary and cryptocurrency as an emerging option. Webhook-driven automated provisioning after payment confirmation is the standard pattern.

### Cited Findings
- PPanel (open-source panel) supports payment gateways as a core feature alongside subscriptions and orders — [GitHub PPanel](https://github.com/CrackerCat/PPanel-server)
- A proxy platform's payments overview describes Stripe invoices captured without customer interaction when automatic payment is on, then marked as paid by the API — [PingProxies Payments Overview](https://documentation.pingproxies.com/checkout-and-payments/payments-overview)
- A hosting provider's billing documentation describes withholding automatic provisioning until Stripe manually verifies and clears cryptocurrency funds — [Hostwinds Billing Docs](https://www.hostwinds.com/product-docs/billing/managing-cryptocurrency-payments)
- An npm package `proxy-generics-stripe` provides a generic payment processor adapter for Stripe, configured with API keys — [npm proxy-generics-stripe](https://www.npmjs.com/package/proxy-generics-stripe)
- MaskProxy's white-label system includes payment gateway customization as a core feature — [MaskProxy Agent](https://maskproxy.io/become-an-agent.html)

### Inferences
- The standard payment flow is: (1) customer selects plan, (2) Stripe/PayPal processes payment, (3) webhook fires on payment confirmation, (4) backend calls provider API to create/top-up sub-user, (5) credentials delivered to customer.
- For cryptocurrency, an additional settlement verification step is needed before provisioning due to confirmation delays and chargeback risk.
- Common payment processors in the proxy space: Stripe (primary), PayPal, and crypto gateways (e.g., CoinGate, NOWPayments, BTCPay Server for self-hosted).
- WHMCS natively supports Stripe, PayPal, and various crypto gateways, making it a shortcut for payment integration.

### Gaps
- Specific crypto payment gateway adoption rates among proxy resellers could not be found.
- Chargeback rates and fraud prevention strategies specific to the proxy industry were not found.


## How do resellers handle user authentication and access control?

### Takeaway
Proxy access uses provider-issued credentials (username:password) routed through the provider's proxy servers, while the reseller platform itself uses standard web authentication (JWT/sessions) for dashboard access.

### Cited Findings
- Bright Data API keys have five permission levels: Admin, Finance, Ops, Limit, and User — [Bright Data API Introduction](https://docs.brightdata.com/api-reference/introduction)
- Proxy credentials use the format `brd-customer-[ACCOUNT_ID]-zone-[PROXY_NAME]` as the username — [Bright Data API Introduction](https://docs.brightdata.com/api-reference/introduction)
- Smartproxy/Decodo supports IP whitelisting as an alternative authentication method: endpoints for listing, adding, and deleting whitelisted IPs — [Smartproxy API GitHub](https://github.com/Smartproxy/Smartproxy-API)
- Evomi's guide explicitly says to store the sub-user username but not the password, retrieving it on demand — [Evomi Reseller Integration Guide](https://proxy-docs.evomi.com/reseller-guide-integration)
- Oxylabs API keys are issued by Customer Support with Bearer token authentication — [Oxylabs Dashboard API](https://developers.oxylabs.io/dashboard/dashboard-api)

### Inferences
- Two authentication layers exist: (1) platform authentication (customer logs into reseller dashboard via email/password, OAuth, or API key) and (2) proxy authentication (customer uses provider-issued credentials or IP whitelisting to access proxies).
- IP whitelisting is an important feature for enterprise customers who cannot use username:password authentication in their tools.
- The reseller platform should never store proxy passwords in plaintext; instead, retrieve them on-demand from the provider API.

### Gaps
- Specific OAuth/SSO implementations used by proxy reseller platforms were not found.
- Rate limiting and brute-force protection patterns for reseller dashboard login were not found.


## What monitoring and analytics tools are needed?

### Takeaway
Reseller platforms need real-time bandwidth usage dashboards, per-customer usage tracking, provider API health monitoring, and standard application monitoring (uptime, errors, latency).

### Cited Findings
- Oxylabs provides usage statistics API with `traffic_bytes` field, supporting grouping by DAY and TARGET dimensions, limited to 31-day ranges — [Oxylabs Dashboard API](https://developers.oxylabs.io/dashboard/dashboard-api)
- Bright Data offers usage monitoring for the residential network calculated via bandwidth consumption, displayed in the Control Panel's Statistics view — [Bright Data Usage Monitoring](https://docs.brightdata.com/general/usage-monitoring/Monitoring_Residential_Network_Usage)
- PingProxies describes a per-customer ledger logging usage as decrements and top-ups as increments, described as critical for billing and quota management — [PingProxies Analytics](https://documentation.pingproxies.com/api-explainers/analytics)
- Proxidize includes usage analytics as part of their reseller management portal — [Proxidize Reseller Program](https://proxidize.com/partners/reseller-program)
- Pingproxies offers a React-based dashboard for monitoring proxy and API usage, including sub-user activity — [Pingproxies Dashboard](https://pingproxies.com/features/best-in-class-dashboard)

### Inferences
- Required monitoring layers: (1) Provider usage API polling (bandwidth consumed per sub-user), (2) Platform metrics (request counts, error rates, API call latency), (3) Business metrics (revenue, active users, churn), (4) Infrastructure monitoring (server health, database performance).
- Standard tools: Grafana + Prometheus for infrastructure, custom dashboards for business metrics, provider API polling for usage data.
- Alerting needed for: usage threshold warnings, quota exhaustion, provider API failures, payment failures, abnormal usage patterns.

### Gaps
- Specific monitoring stack recommendations from existing reseller platforms were not found.
- Provider API polling intervals and rate limits for usage endpoints were not documented.


## How is bandwidth metering and usage tracking implemented?

### Takeaway
Bandwidth is metered at the provider's proxy server level and exposed through provider APIs. Reseller platforms poll these APIs and maintain their own usage ledgers for billing, quotas, and customer dashboards.

### Cited Findings
- Bright Data usage is calculated via bandwidth consumption (bytes) at the proxy network level — [Bright Data Usage Monitoring](https://docs.brightdata.com/general/usage-monitoring/Monitoring_Residential_Network_Usage)
- Bright Data supports usage webhooks for real-time notification (endpoint titles: Create or Change Usage Webhook Settings, Get Usage Webhook Info) — [Bright Data Reseller API Docs](https://docs.brightdata.com/api-reference/reseller-api/introduction)
- Bright Data usage limits use byte-based units: `{"limit": {"limit": 50000000}}` with unit `bw` — [Bright Data Reseller API Docs](https://docs.brightdata.com/api-reference/reseller-api/introduction)
- Evomi tracks balance per sub-user; the balance endpoint adds credit, and the usage endpoint returns history — [Evomi Reseller Integration Guide](https://proxy-docs.evomi.com/reseller-guide-integration)
- Smartproxy/Decodo provides per-sub-user traffic endpoint: `GET /users/{userId}/sub-users/{username}/traffic` — [Smartproxy API GitHub](https://github.com/Smartproxy/Smartproxy-API)
- PingProxies implements a per-customer ledger model: usage as decrements, top-ups as increments — [PingProxies Analytics](https://documentation.pingproxies.com/api-explainers/analytics)
- A glossary definition describes metered proxies as tracking bandwidth consumption for each customer account, measuring both inbound and outbound traffic — [CorsProxy.io Glossary](https://corsproxy.io/glossary/metered-proxy/)
- RFC 2227 (1997) defined HTTP hit-metering and usage-limiting for proxy caches as a standard — [RFC 2227](https://www.rfc-editor.org/rfc/rfc2227.txt)

### Inferences
- Two metering models exist: (a) byte-based (Bright Data, Smartproxy) where bandwidth in bytes is the billing unit, and (b) balance-based (Evomi) where monetary credit is deducted as bandwidth is consumed.
- Best practice is to maintain a local usage ledger synchronized with the provider API, rather than relying solely on real-time provider queries. This protects against API downtime and enables faster dashboard rendering.
- Usage webhooks (where available, e.g., Bright Data) are preferred over polling for near-real-time tracking.
- The ledger pattern (increments for top-ups, decrements for usage) is the cleanest implementation for reconciliation.

### Gaps
- Webhook payload formats and delivery guarantees for provider usage notifications were not documented.
- How providers handle bandwidth discrepancies or billing disputes with resellers was not found.


## What are the typical tech stacks used for proxy reselling platforms?

### Takeaway
There is no single standard stack, but the most common patterns involve a modern JavaScript/TypeScript or Python backend with a React frontend, PostgreSQL database, and Stripe for payments.

### Cited Findings
- Remnawave Panel: TypeScript + NestJS backend, AGPL-3.0 — [GitTrend Remnawave](https://gittrend.io/repo/remnawave/panel)
- PPanel: Go backend for server, separate web frontend — [GitHub PPanel](https://github.com/CrackerCat/PPanel-server)
- V2Board: PHP/Laravel — [GitHub V2Board](https://github.com/Mintapril/v2board)
- Freelancer portfolio: Django + FastAPI backend, React frontend, PostgreSQL, admin interface for products/orders/support — [Freelancer Project](https://www.freelancer.cn/projects/api-integration/proxy-seller-site-with-api)
- Another freelancer proposal: PHP-based platform connecting to provider APIs — [Freelancer Project](https://www.freelancer.cl/projects/network-administration/proxy-reseller-reverse-proxy-setup)
- Pingproxies: React-based dashboard — [Pingproxies Dashboard](https://pingproxies.com/features/best-in-class-dashboard)
- Bright Data CLI: `npm install -g @brightdata/cli`, Python SDK: `SyncBrightDataClient`, JavaScript SDK: `bdclient` — [Bright Data API Introduction](https://docs.brightdata.com/api-reference/introduction)

### Inferences
- **MVP Tech Stack**: Next.js or React frontend + Node.js/Express or FastAPI backend + PostgreSQL + Stripe + one provider API (Evomi or Smartproxy for simplest integration).
- **Production Tech Stack**: React/Next.js frontend + NestJS or Django backend + PostgreSQL + Redis (caching, sessions, rate limiting) + Stripe + PayPal + multiple provider adapters + background job queue (Bull/Celery) for usage sync and provisioning.
- **Infrastructure**: Docker containerization, reverse proxy (Nginx/Caddy), CI/CD pipeline, object storage for logs.
- The choice between Node.js/TypeScript and Python/Django depends on team expertise; both are well-represented in the space.

### Gaps
- Traffic volumes and scaling requirements for typical proxy reseller platforms were not found.
- Specific hosting providers commonly used by proxy resellers were not identified.


## How do resellers handle customer support and ticket systems?

### Takeaway
Most reseller platforms use either built-in ticket systems (in open-source panels) or third-party solutions like Freshdesk, Zendesk, or Intercom, with live chat being essential for this industry.

### Cited Findings
- PPanel includes ticket management as a core feature alongside user management, subscriptions, and payments — [GitHub PPanel](https://github.com/CrackerCat/PPanel-server)
- Proxidize offers a dedicated partner success manager for resellers — [Proxidize Reseller Program](https://proxidize.com/partners/reseller-program)
- NetNut received complaints about support quality from resellers on Trustpilot — [Trustpilot NetNut](https://nz.trustpilot.com/review/netnut.io)
- Capterra reviewers for NetNut noted the dashboard and API access are intuitive, suggesting self-service reduces support load — [Capterra NetNut Reviews](https://www.capterra.com/p/10015437/NetNut/reviews/?page=2)

### Inferences
- Support tiers for a proxy reseller: (1) self-service knowledge base and FAQ, (2) ticket system for technical issues, (3) live chat for sales and urgent issues, (4) escalation path to upstream provider for infrastructure issues.
- WHMCS includes a built-in ticket system, making it a two-in-one solution for billing + support.
- For custom platforms, integrating a third-party like Freshdesk (API-based) or building a simple ticket system is standard.
- A good self-service dashboard with real-time usage, credential regeneration, and proxy testing reduces support volume significantly.

### Gaps
- Average support ticket volumes for proxy reseller businesses were not found.
- Specific third-party support tool adoption rates in the proxy industry were not found.


## How does automated provisioning work?

### Takeaway
Automated provisioning is webhook-driven: payment confirmation triggers API calls to the provider to create sub-users and allocate resources, typically completing within seconds for most providers.

### Cited Findings
- Evomi provisioning flow: create sub-user after payment, then add balance via API — sub-user creation is immediate through the API — [Evomi Reseller Integration Guide](https://proxy-docs.evomi.com/reseller-guide-integration)
- Bright Data: new users default to a 1KB usage limit, which the reseller increases after payment via `PUT /{partner}/usage_limit` — [Bright Data Reseller API Docs](https://docs.brightdata.com/api-reference/reseller-api/introduction)
- Bright Data: zone creation with IP allocation done in a single API call: `PUT /{partner}/zone?id={id}&type={type}&ips_amount={n}&zone={name}` — [Bright Data Reseller API Docs](https://docs.brightdata.com/api-reference/reseller-api/introduction)
- CatProxies claims reseller approval within one business day — [CatProxies Reseller](https://dashboard.catproxies.com/become-reseller)
- NetNut API top-ups reported to take over an hour to process in one instance — [Rampage Networks Note](https://rampagenetworks.notion.site/Netnut-Private-Resi-435e0518edb647a784dab07e3a2b82cf)
- A hosting provider withholds automatic provisioning until crypto payments are verified and cleared — [Hostwinds Billing Docs](https://www.hostwinds.com/product-docs/billing/managing-cryptocurrency-payments)

### Inferences
- Standard provisioning pipeline: (1) Stripe/PayPal webhook fires `payment_intent.succeeded`, (2) backend validates webhook signature, (3) backend calls provider API to create/update sub-user, (4) backend generates proxy connection details, (5) customer receives credentials via email and dashboard.
- Provisioning latency varies by provider: most are near-instant via API, but NetNut may have delays.
- For crypto payments, implement a "pending" state with background polling until payment is confirmed.
- Idempotency is critical: webhook retries must not create duplicate sub-users.

### Gaps
- Provider API rate limits for provisioning endpoints were not found.
- SLA guarantees for API availability from providers were not found.


## Are there open-source proxy management panels or billing systems?

### Takeaway
Several open-source panels exist for proxy/VPN protocol management with billing, but none are purpose-built specifically for proxy reselling with upstream provider API integration.

### Cited Findings
- **Remnawave Panel**: TypeScript/NestJS, Xray-core based, ~5.3k GitHub stars, AGPL-3.0 — [GitTrend Remnawave](https://gittrend.io/repo/remnawave/panel)
- **PPanel**: Go backend, covers user management, subscriptions, payments, orders, tickets, and node management; GNU licensed frontend — [GitHub PPanel](https://github.com/CrackerCat/PPanel-server); [GitHub PPanel Frontend](https://github.com/perfect-panel/frontend)
- **V2Board**: PHP/Laravel, multiple proxy protocols, requires modified XrayR backend — [GitHub V2Board](https://github.com/Mintapril/v2board)
- **ProxyPanel**: manages multiple proxy protocols — [GitHub ProxyPanel](https://github.com/ProxyPanel)
- **WHMCS**: not open-source but extensible; the Proxy-Cheap module is the only proxy-specific module found — [WHMCS Marketplace](https://marketplace.whmcs.com/product/6889-whitelabel-proxy-reseller-module)

### Inferences
- The open-source panels (PPanel, Remnawave, V2Board) are designed for operators running their own proxy infrastructure, not for resellers connecting to upstream providers. They would need significant modification to add provider API integration.
- PPanel is the closest to a reseller platform due to its subscription/billing/ticket features, but it still assumes you operate the proxy nodes.
- For a reseller-specific platform, building on a general SaaS billing framework (like Laravel Cashier, Django-Oscar, or a custom Stripe integration) and adding provider adapters is likely more efficient than adapting a proxy panel.

### Gaps
- Whether any of these open-source panels have been forked specifically for reselling purposes could not be determined.
- Active contributor counts and release frequencies for these projects were not verified.


## What security measures are needed?

### Takeaway
Proxy reseller platforms need DDoS protection for the platform itself, abuse prevention to block misuse of proxies, API security with rate limiting, and standard web application security measures.

### Cited Findings
- Proxy operators have primary responsibility for rate limiting and abuse prevention — [Triton RPC Proxy Docs](https://docs.triton.one/rpc-pool/proxying)
- Recommended abuse prevention measures: SSRF protection against private IP ranges, per-user rate limits, and temporary blocks after repeated violations — [Nolxy Abuse Policy](https://nolxy.pages.dev/abuse)
- Cloudflare: proxying DNS records lets DDoS protection stop attacks before they reach the origin server — [Cloudflare DDoS Protection](https://developers.cloudflare.com/learning-paths/prevent-ddos-attacks/baseline/proxy-dns-records/)
- Alibaba Cloud Anti-DDoS uses an AI engine to learn normal traffic patterns and filter application-layer attacks — [Alibaba Cloud Anti-DDoS](https://intl.aliyun.com/help/en/anti-ddos/anti-ddos-pro-and-premium/product-overview/what-are-anti-ddos-pro-and-anti-ddos-premium)
- Commercial proxy providers can see traffic patterns and should be expected to spot abuse — [Rappler/Qurium Report](https://www.rappler.com/technology/ddos-attack-rappler-october-2023-proxy-providers-rayobyte-fineproxy-weaponized-qurium/)
- AWS guidance on securing open proxies recommends regular audits and strict access controls — [AWS Security Blog](https://aws.amazon.com/blogs/security/securing-open-proxies-in-your-aws-environment/)

### Inferences
- **Platform security**: Cloudflare or similar CDN/DDoS protection for the dashboard and API; WAF rules for common attacks; rate limiting on all API endpoints.
- **Abuse prevention**: Terms of service with prohibited use cases; automated monitoring for suspicious patterns (high request volumes, targeting sensitive sites); per-customer bandwidth caps; ability to instantly disable abusive accounts via provider API.
- **API security**: API key rotation, webhook signature verification (for Stripe/provider callbacks), input validation, SQL injection prevention, CSRF protection.
- **Data security**: Encrypt sensitive data at rest (API keys, tokens); never store proxy passwords locally; use secrets management (Vault, AWS Secrets Manager); regular security audits.
- **Compliance**: KYC/identity verification for high-volume customers; maintain audit logs; have an abuse reporting mechanism.

### Gaps
- Industry-specific compliance requirements for proxy resellers (if any exist beyond general terms of service) were not found.
- Specific DDoS attack patterns targeting proxy reseller platforms were not found.
- Insurance or liability protection patterns for proxy resellers were not found.
