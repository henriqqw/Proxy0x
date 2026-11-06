# GTI Proxy e FastProxy: Analise Competitiva Completa

## GTI Proxy (gtiproxy.com)

### Dados da Empresa

#### Takeaway
GTI Proxy e operada por GTI BRASIL LTDA, uma microempresa registrada em Goiania/GO desde 2013, com CNPJ verificavel e nota fiscal brasileira. E uma operacao enxuta de um unico socio.

#### Cited Findings
- CNPJ 18.915.849/0001-50, razao social GTI BRASIL LTDA, registrada em 19/09/2013 como Sociedade Empresaria Limitada (Microempresa) -- [CNPJ API](https://cnpj-api.com/en/cnpj/18915849000150)
- Capital social de R$ 100.000,00 -- [CNPJ API](https://cnpj-api.com/en/cnpj/18915849000150)
- Optante do Simples Nacional, nao e MEI -- [CNPJ API](https://cnpj-api.com/en/cnpj/18915849000150)
- Socio unico: RAFAEL SANTOS BARBOSA (qualificacao codigo 49) -- [CNPJ API](https://cnpj-api.com/en/cnpj/18915849000150)
- Endereco: Rua 14, 254, Jd Goias, Goiania/GO, CEP 74810-180 -- [CNPJ API](https://cnpj-api.com/en/cnpj/18915849000150)
- CNAE principal: 63.19-4-00 (Portais, provedores de conteudo e outros servicos de informacao na internet) -- [CNPJ API](https://cnpj-api.com/en/cnpj/18915849000150)
- CNAEs secundarios: 62.09-1-00 (Suporte tecnico em TI), 63.11-9-00 (Tratamento de dados e hospedagem), 63.99-2-00 (Outros servicos de informacao) -- [CNPJ API](https://cnpj-api.com/en/cnpj/18915849000150)
- Emite nota fiscal brasileira sob solicitacao no painel -- [GTI Proxy](https://gtiproxy.com/)
- Verificado no Google Meu Negocio -- [GTI Proxy Sobre](https://gtiproxy.com/sobre)

#### Inferences
- Com capital social de apenas R$ 100.000 e socio unico, e uma operacao pequena, provavelmente com equipe reduzida
- A empresa opera desde 2013, mas nao esta claro se sempre atuou no segmento de proxies
- Sede em Goiania mas proxies em SP, RJ e Brasilia sugere infraestrutura distribuida sem datacenter proprio

#### Gaps
- Nao foi possivel determinar quando a empresa comecou a atuar especificamente com proxies
- Numero de funcionarios nao divulgado
- Parcela de receita vinda de proxies vs. outros servicos de TI nao e conhecida

---

### Planos e Precos (BRL)

#### Takeaway
GTI Proxy oferece 5 planos com precos de R$ 15 (passe 24h movel) a R$ 200/mes (movel mensal). O residencial fixo custa R$ 40/mes com banda ilimitada, confirmando o preco reportado. Ha descontos por volume a partir de 5 IPs.

#### Cited Findings

**Proxies Moveis 4G/5G:**

| Plano | Preco | Periodo | Dados | Rotacao |
|---|---|---|---|---|
| Movel 24h | R$ 15 | 24 horas | 1 GB | Automatica |
| Movel Semanal | R$ 50 | 7 dias | 5 GB | Automatica |
| Movel Mensal | R$ 200 | 1 mes | 60 GB | Manual via link |

-- [GTI Proxy](https://gtiproxy.com/)

**Proxies Residenciais:**

| Plano | Preco | Periodo | Detalhes |
|---|---|---|---|
| Residencial Fixo | R$ 40/mes | Mensal | IP dedicado e fixo, banda ilimitada |
| Residencial Fixo IPv6 | R$ 40/mes | Mensal | Pacote de 5 IPs fixos (nao vendidos avulso) |

-- [GTI Proxy](https://gtiproxy.com/)

**Descontos por volume (Residencial Fixo):**

| Quantidade | Preco por IP |
|---|---|
| 1-4 IPs | R$ 40 |
| 5+ IPs | R$ 36 |
| 10+ IPs | R$ 34 |
| 25+ IPs | R$ 32 |
| 50+ IPs | R$ 30 |

-- [GTI Proxy](https://gtiproxy.com/)

- Credito de R$ 15: quem assinar o Mensal em ate 48 horas apos o passe de 24h tem R$ 15 abatidos na primeira mensalidade -- [GTI Proxy](https://gtiproxy.com/)
- Teste gratis: ate 5 testes por CPF ou CNPJ, com 1 hora entre eles, duracao de 15 minutos, liberado pelo IP do cliente. O IP do teste nao e o mesmo do plano contratado -- [GTI Proxy](https://gtiproxy.com/)
- Reembolso: 7 dias corridos (art. 49 CDC). Integral dentro dos limites de avaliacao (50 MB no 24h, 250 MB no Semanal, 3 GB no Mensal, 2 GB nos residenciais). Proporcional acima desses limites -- [GTI Proxy](https://gtiproxy.com/)
- Sem fidelidade. Ativacao em minutos apos confirmacao do pagamento -- [GTI Proxy](https://gtiproxy.com/)

#### Inferences
- O preco por GB no plano movel mensal (R$ 200/60 GB = R$ 3,33/GB) e competitivo para proxies moveis brasileiros
- O residencial fixo a R$ 40/mes com banda ilimitada e atrativo para scraping de alto volume
- A escala de descontos ate R$ 30/IP para 50+ IPs indica que atendem clientes de medio porte com multiplos IPs
- O teste gratis de 15 minutos e limitado mas demonstra confianca no produto

#### Gaps
- Preco do plano mensal movel pode ter mudado; pagina informa revisao em 19/10/2026
- Nao ha informacao sobre precos para periodos superiores a 1 mes (trimestral, semestral, anual)

---

### Tipos de Proxy e Especificacoes Tecnicas

#### Takeaway
GTI Proxy oferece proxies moveis 4G/5G (TIM, Vivo, Claro) e residenciais fixos (fibra optica). Nao oferece datacenter. Pool total de ~12.800 IPs brasileiros. Suporta HTTP e SOCKS5.

#### Cited Findings
- Proxy movel 4G/5G com IPs de TIM, Vivo e Claro em Sao Paulo, Rio de Janeiro e Brasilia -- [GTI Proxy](https://gtiproxy.com/)
- Proxy residencial com IPs de provedores brasileiros de fibra optica em Sao Paulo e Rio de Janeiro -- [GTI Proxy](https://gtiproxy.com/)
- Residencial Fixo IPv6: pacote de 5 IPs para sites com IPv6 -- [GTI Proxy](https://gtiproxy.com/)
- Nao oferece proxies de datacenter -- [GTI Proxy](https://gtiproxy.com/)
- Pool de 12.847 IPs brasileiros disponiveis (apurado em 30/07/2026) -- [GTI Proxy](https://gtiproxy.com/)
- Protocolos: HTTP e SOCKS5, com autenticacao usuario e senha -- [GTI Proxy](https://gtiproxy.com/)
- Residencial: 100+ threads simultaneas e 100+ Mbps -- [GTI Proxy](https://gtiproxy.com/)
- Rotacao: automatica no Movel 24h e Semanal, manual no Mensal, sem rotacao no Residencial Fixo -- [GTI Proxy](https://gtiproxy.com/)

#### Inferences
- Pool de ~12.800 IPs e relativamente pequeno comparado a provedores internacionais (milhoes), mas faz sentido para IPs 100% brasileiros
- A cobertura de 3 operadoras moveis e adequada para o mercado brasileiro
- A restricao do IPv6 residencial "somente para sites com IPv6" e uma limitacao relevante que o concorrente FastProxy tambem enfrenta
- Sem datacenter = foco em qualidade/autenticidade do IP vs. volume

#### Gaps
- Nomes dos provedores de fibra optica nao divulgados
- Nao ha informacao sobre como obtem/gerenciam os IPs moveis (modems fisicos, SIM cards, etc.)
- Latencia e uptime nao divulgados (a empresa declara que nao publica numeros que nao pode medir)

---

### Formas de Pagamento e Fiscal

#### Takeaway
Aceita Pix e cartao de credito em Real. Emite nota fiscal. Boleto nao e mencionado.

#### Cited Findings
- Pagamento via Pix e cartao de credito, em Real (BRL) -- [GTI Proxy](https://gtiproxy.com/)
- Nota fiscal brasileira sob solicitacao no painel -- [GTI Proxy](https://gtiproxy.com/)
- Boleto nao e mencionado na pagina -- [GTI Proxy](https://gtiproxy.com/)
- Selos no rodape: LGPD e Nota Fiscal -- [GTI Proxy](https://gtiproxy.com/)

#### Inferences
- Ausencia de boleto pode ser uma limitacao para empresas que preferem esse metodo
- Emissao de NF por solicitacao (e nao automatica) pode indicar processo manual

#### Gaps
- Nao foi possivel verificar se aceitam boleto em algum fluxo nao visivel na homepage

---

### API e Painel

#### Takeaway
GTI Proxy possui uma pagina de API e um painel de controle, alem de compatibilidade com mais de 10 navegadores antidetect.

#### Cited Findings
- Pagina "API de proxy" disponivel em /api-proxy -- [GTI Proxy](https://gtiproxy.com/)
- Pagina de configuracao em /configurar-proxy -- [GTI Proxy](https://gtiproxy.com/)
- Ferramentas compativeis: AdsPower, Multilogin, Dolphin{anty}, Evolution API, Puppeteer e Playwright -- [GTI Proxy](https://gtiproxy.com/)
- Parceiros antidetect: GeeLark, MoreLogin, BitBrowser, DuoPlus, AdsPower, GoLogin, Incogniton, Kameleo, Multilogin e Octo Browser -- [GTI Proxy](https://gtiproxy.com/)
- Rotacao manual no plano mensal via link (sugere interface web para rotacao) -- [GTI Proxy](https://gtiproxy.com/)

#### Inferences
- A larga compatibilidade com antidetect browsers indica foco no publico de multi-contas e trafego pago
- A integracao com Evolution API sugere foco em automacao de WhatsApp
- Presenca de pagina de API indica que ha algum nivel de automacao disponivel, mas o conteudo nao foi acessivel

#### Gaps
- O conteudo detalhado da API nao foi extraido (a pagina nao apareceu no conteudo principal)
- Nao se sabe se o painel e desenvolvido internamente ou usa solucao de terceiros

---

### Publico-Alvo e Casos de Uso

#### Takeaway
O publico principal e de gestores de trafego, afiliados e operadores de multi-contas em redes sociais, com foco secundario em scraping e automacao.

#### Cited Findings
- Casos de uso citados: operacoes de WhatsApp, Instagram, TikTok, Meta Ads e BMs, contingencia de contas, web scraping, verificacao de anuncios, monitoramento de e-commerce e SEO, publicacao em classificados -- [GTI Proxy](https://gtiproxy.com/)
- Ferramentas destacadas incluem antidetect browsers (Multilogin, AdsPower, GoLogin etc.) e ferramentas de automacao (Puppeteer, Playwright, Evolution API) -- [GTI Proxy](https://gtiproxy.com/)

#### Inferences
- O foco em "contingencia de contas" e "BMs" (Business Managers do Facebook) indica que o publico principal sao gestores de trafego pago
- A mencao a Evolution API e WhatsApp sugere atendimento ao mercado de automacao de mensagens
- Combinacao de movel + residencial indica estrategia de oferecer o IP certo para cada caso de uso

---

### Presenca de Marketing

#### Takeaway
Presenca digital limitada: Instagram (@gtiproxy), LinkedIn, Google Meu Negocio. Sem evidencia de YouTube, blog ativo ou conteudo de SEO significativo.

#### Cited Findings
- Redes sociais: Instagram (@gtiproxy) e LinkedIn (gti-proxy) -- [GTI Proxy](https://gtiproxy.com/)
- Perfil verificado no Google Meu Negocio -- [GTI Proxy Sobre](https://gtiproxy.com/sobre)
- Suporte 24/7 via WhatsApp (+55 11 92620-5014) e email -- [GTI Proxy](https://gtiproxy.com/)

#### Inferences
- A presenca de marketing parece ser minimalista, confiando mais em indicacao e busca organica
- Numero de WhatsApp com DDD 11 (Sao Paulo) apesar da sede em Goiania sugere posicionamento para o mercado paulista
- A ausencia de YouTube e blog pode indicar estrategia de boca a boca ou canais privados (grupos de Telegram/WhatsApp)

#### Gaps
- Nao foi possivel acessar o perfil do Instagram para verificar numero de seguidores e frequencia de postagens
- Nao ha evidencia de canal no YouTube
- Nao ha evidencia de blog no site

---

### Reputacao e Reviews

#### Takeaway
Nenhuma review ou reclamacao foi encontrada em nenhuma plataforma publica (Trustpilot, Reclame Aqui, Google). Isso pode indicar baixo volume de clientes ou publico que nao usa essas plataformas.

#### Cited Findings
- Nenhuma pagina no Reclame Aqui para "GTI Proxy" ou "GTI Brasil" no segmento de proxy -- busca direta no Reclame Aqui
- Nenhuma review no Trustpilot para gtiproxy.com -- busca no Trustpilot
- A empresa declara que nao publica numeros de clientes ou reviews em seu site -- [GTI Proxy Sobre](https://gtiproxy.com/sobre)

#### Inferences
- A ausencia de reclamacoes no Reclame Aqui nao e necessariamente positiva; pode indicar que o servico e de nicho e poucos clientes usam essa plataforma
- A declaracao explicita de que nao publicam metricas de vaidade pode ser interpretada como transparencia ou como falta de social proof

#### Gaps
- Reviews em grupos de Telegram ou WhatsApp nao sao acessiveis via busca publica
- Nao foi possivel verificar o perfil no Google Meu Negocio para ver avaliacoes

---

### Upstream / Origem dos Proxies

#### Takeaway
GTI Proxy nao revela seus fornecedores upstream. Cita apenas as operadoras moveis (TIM, Vivo, Claro) e "provedores brasileiros de fibra optica" sem nomeá-los.

#### Cited Findings
- Operadoras moveis citadas: TIM, Vivo e Claro -- [GTI Proxy](https://gtiproxy.com/)
- Residencial: "provedores brasileiros de fibra optica" sem identificacao -- [GTI Proxy](https://gtiproxy.com/)
- Pontos de presenca: Sao Paulo, Rio de Janeiro e Brasilia -- [GTI Proxy](https://gtiproxy.com/)

#### Inferences
- O modelo mais provavel para proxies moveis e modems fisicos com SIM cards das operadoras, modelo comum no mercado brasileiro
- Para residencial fixo, provavelmente ha acordos com ISPs locais ou uso de IPs residenciais cedidos/alugados
- O pool de ~12.800 IPs e compativel com operacao propria (nao revenda de upstream internacional)

#### Gaps
- Sem confirmacao do modelo de infra (modems fisicos vs. SDKs vs. revenda)
- Nao ha indicacao de que usem provedores upstream internacionais como Bright Data ou SOAX

---

## FastProxy (fastproxy.com.br)

### Dados da Empresa

#### Takeaway
FastProxy nao divulga CNPJ, razao social ou endereco em seu site. A jurisdicao e Sao Paulo/SP conforme os termos de uso. A ausencia de dados corporativos e uma bandeira de atencao.

#### Cited Findings
- CNPJ nao listado no site -- [FastProxy](https://www.fastproxy.com.br/)
- CNPJ nao encontrado em buscas publicas (CNPJ API, MonitorCNPJ) -- buscas multiplas realizadas
- Foro eleito: Comarca de Sao Paulo/SP -- [FastProxy Termos](https://www.fastproxy.com.br/termos)
- Copyright: "(c) 2026 FastProxy. Feito no Brasil" -- [FastProxy](https://www.fastproxy.com.br/)
- Contato: WhatsApp +55 11 93007-0320, Telegram t.me/fastproxy, Instagram @fastproxybr, email suporte@fastproxy.com.br -- [FastProxy](https://www.fastproxy.com.br/)

#### Inferences
- A ausencia de CNPJ no site e nos registros publicos e preocupante; pode operar como MEI, pessoa fisica, ou simplesmente nao divulgar
- O dominio .com.br requer CNPJ ou CPF para registro no registro.br, entao ha uma pessoa ou empresa por tras
- O foro de Sao Paulo e o DDD 11 no WhatsApp sugerem base em SP

#### Gaps
- CNPJ nao encontrado em nenhuma fonte publica
- Razao social desconhecida
- Endereco fisico desconhecido
- WHOIS do registro.br nao foi consultado diretamente

---

### Planos e Precos (BRL)

#### Takeaway
FastProxy oferece 5 planos mensais de R$ 29,90 a R$ 120,00, com desconto de 18% no pagamento anual. Tambem ha pacotes prepagos com descontos de 15% a 45%.

#### Cited Findings

**Planos Mensais:**

| Plano | Mensal | Anual (por proxy/mes) | Total Anual | Economia Anual | Detalhes |
|---|---|---|---|---|---|
| IPv6 (Economico) | R$ 29,90 | R$ 24,67 | R$ 296,01 | R$ 62,79 | HTTP/S e SOCKS5, entrega imediata, reposicao garantida |
| IPv6 Rotativo | R$ 120,00 | R$ 99,00 | R$ 1.188,00 | R$ 252,00 | Rotacao a cada 10 min, pool ilimitado IPv6 Brasil |
| IPv4 Dedicado | R$ 39,90 | R$ 32,92 | R$ 395,01 | R$ 83,79 | IP exclusivo dedicado, banda ilimitada, suporte prioritario |
| IPv6 p/ Facebook Ads | R$ 79,90 | R$ 65,92 | R$ 791,01 | R$ 167,79 | 1 proxy + 10 rotacoes gratis, alta compatibilidade com BM, suporte VIP |
| ISP (IP fixo) | R$ 49,90 | R$ 41,17 | R$ 494,01 | R$ 104,79 | 100 Mbps, pool separado de datacenter, reputacao mais limpa |

-- [FastProxy](https://www.fastproxy.com.br/)

**Pacotes Prepagos (IPv6):**

| IPs | Duracao | Desconto | Label |
|---|---|---|---|
| 5 IPs | 3 meses | 15% | "Mais vendido" |
| 10 IPs | 6 meses | 25% | "Recomendado" |
| 25 IPs | 12 meses | 40% | "Melhor custo" |
| 50 IPs | 12 meses | 45% | - |

-- [FastProxy](https://www.fastproxy.com.br/)

- Os precos dos pacotes prepagos nao foram exibidos na pagina principal (apenas os percentuais de desconto) -- [FastProxy](https://www.fastproxy.com.br/)
- Garantia de 7 dias (exibido como badge) -- [FastProxy](https://www.fastproxy.com.br/)
- Sem fidelidade. Cancelamento a qualquer momento pelo painel ou suporte -- [FastProxy Termos](https://www.fastproxy.com.br/termos)

#### Inferences
- Calculando o pacote prepago de 5 IPs/3 meses com 15% de desconto: R$ 29,90 x 0,85 x 3 = R$ 76,25 por proxy, ou ~R$ 25,42/mes por proxy
- O plano mais caro (IPv6 Rotativo a R$ 120/mes) indica que a rotacao tem um premium significativo
- O plano "Facebook Ads" a R$ 79,90 com apenas 10 rotacoes sugere que rotacoes adicionais podem ter custo extra
- O ISP a R$ 49,90 e o produto premium para quem precisa de IP de provedor real (nao datacenter)

#### Gaps
- Precos exatos dos pacotes prepagos nao foram exibidos
- FAQ colapsado nao revelou detalhes sobre recorrencia de pagamento e cancelamento
- Politica de reembolso detalhada esta em /reembolso (nao acessada)
- Limites de banda para IPv6 e IPv6 Rotativo nao especificados

---

### Tipos de Proxy e Especificacoes Tecnicas

#### Takeaway
FastProxy oferece IPv6 dedicado, IPv6 rotativo, IPv4 dedicado e ISP residencial. Nao oferece proxies moveis 4G/5G (diferencial do GTI Proxy).

#### Cited Findings
- Tipos oferecidos: IPv6 dedicado, IPv6 rotativo, IPv4 dedicado, IPv6 para Facebook Ads, ISP residencial -- [FastProxy](https://www.fastproxy.com.br/)
- Protocolos: HTTP/S e SOCKS5 -- [FastProxy](https://www.fastproxy.com.br/)
- Autenticacao por usuario e senha -- [FastProxy](https://www.fastproxy.com.br/)
- Reposicao automatica em um clique -- [FastProxy](https://www.fastproxy.com.br/)
- ISP: 100 Mbps, pool separado de datacenter generico -- [FastProxy](https://www.fastproxy.com.br/)
- IPv4: banda ilimitada declarada -- [FastProxy](https://www.fastproxy.com.br/)
- IPv6 Rotativo: rotacao a cada 10 minutos, pool ilimitado de IPv6 Brasil -- [FastProxy](https://www.fastproxy.com.br/)

#### Inferences
- A ausencia de proxies moveis 4G/5G e uma diferenca significativa em relacao ao GTI Proxy
- O pool "ilimitado" de IPv6 provavelmente se refere a sub-redes IPv6 proprias ou alugadas (e facil gerar muitos IPs IPv6)
- "Pool separado de datacenter" no ISP sugere que os IPv6/IPv4 regulares PODEM ser de datacenter, enquanto o ISP usa IPs de provedor residencial
- A rotacao a cada 10 minutos e fixa, nao configuravel

#### Gaps
- Tamanho exato do pool de IPs nao divulgado
- Velocidade/banda dos planos IPv6 e IPv6 Rotativo nao especificada
- Localizacao dos servidores nao detalhada (apenas "Brasil")

---

### Formas de Pagamento

#### Takeaway
FastProxy aceita Pix, cartao e boleto -- confirmando a informacao reportada. Isso e uma vantagem sobre o GTI Proxy que nao menciona boleto.

#### Cited Findings
- Pix, cartao e boleto aceitos -- [FastProxy](https://www.fastproxy.com.br/)
- Servicos vendidos como assinaturas recorrentes -- [FastProxy Termos](https://www.fastproxy.com.br/termos)

#### Inferences
- Aceitacao de boleto e um diferencial para clientes empresariais que preferem esse metodo
- A cobranca recorrente sugere integracao com gateway de pagamento (possivelmente Mercado Pago, PagSeguro ou similar)

#### Gaps
- Nao foi verificado se a nota fiscal e emitida (o site nao menciona)
- Gateway de pagamento utilizado nao identificado

---

### Painel e API

#### Takeaway
FastProxy oferece um painel web com funcoes de gerenciamento de proxies, faturas e regeneracao de IP em tempo real. Tambem ha um PWA instalavel. Nenhuma API publica documentada.

#### Cited Findings
- Dashboard para rastreamento de proxies, faturas e regeneracao de IP em tempo real -- [FastProxy](https://www.fastproxy.com.br/)
- PWA (Progressive Web App) instalavel com notificacoes push, badges de expiracao e atalhos rapidos -- [FastProxy](https://www.fastproxy.com.br/)
- Nenhuma API publica documentada na pagina -- [FastProxy](https://www.fastproxy.com.br/)
- Compatibilidade: AdsPower, Dolphin Anty, Multilogin, GoLogin, Vision, Indigo, Kameleo, Scrapy -- [FastProxy](https://www.fastproxy.com.br/)

#### Inferences
- O PWA e um diferencial moderno que nenhum outro concorrente brasileiro parece oferecer
- A ausencia de API publica pode limitar integracao programatica para usuarios avancados
- A compatibilidade com Scrapy indica que atendem tambem o publico de web scraping (nao so multi-contas)

#### Gaps
- Nao foi possivel verificar se ha API nao documentada (area logada)
- Stack tecnologico do painel nao identificado

---

### Publico-Alvo e Casos de Uso

#### Takeaway
FastProxy atende primariamente gestores de trafego pago (Facebook Ads, Google Ads) e afiliados, com foco secundario em scraping, automacao e privacidade.

#### Cited Findings
- Casos de uso listados: web scraping, automacao, e-commerce, redes sociais, ferramentas de SEO, data mining, bots & IA, privacidade -- [FastProxy](https://www.fastproxy.com.br/)
- Plano especifico para Facebook Ads com "alta compatibilidade com BM" -- [FastProxy](https://www.fastproxy.com.br/)
- Depoimentos de gestores de trafego, afiliados, desenvolvedores e scrapers -- [FastProxy](https://www.fastproxy.com.br/)
- "Mais de 500 clientes ativos" -- [FastProxy](https://www.fastproxy.com.br/)
- Nota de 4.9 de 237 reviews -- [FastProxy](https://www.fastproxy.com.br/)

#### Inferences
- O plano "Facebook Ads" como produto destaque indica que o core business e trafego pago
- 500 clientes ativos e um numero modesto, compativel com operacao pequena
- As 237 reviews nao tem fonte verificavel (possivelmente internas ao painel)

#### Gaps
- Onde as 237 reviews foram coletadas nao e especificado (podem ser do proprio painel)
- Volume de receita ou proxies vendidos nao divulgado

---

### Presenca de Marketing

#### Takeaway
FastProxy tem presenca em Instagram (@fastproxybr), Telegram, WhatsApp e um blog no site. Marketing mais ativo que GTI Proxy, mas sem YouTube evidente.

#### Cited Findings
- Instagram: @fastproxybr -- [FastProxy](https://www.fastproxy.com.br/)
- Telegram: t.me/fastproxy -- [FastProxy](https://www.fastproxy.com.br/)
- WhatsApp: +55 11 93007-0320 -- [FastProxy](https://www.fastproxy.com.br/)
- Blog linkado no menu do site -- [FastProxy](https://www.fastproxy.com.br/)
- Pagina de anunciantes (/anunciantes) sugere programa de afiliados ou parcerias -- [FastProxy](https://www.fastproxy.com.br/)

#### Inferences
- A pagina /anunciantes sugere que FastProxy investe em marketing de afiliados
- O Telegram como canal de suporte e comum no mercado de proxies e sugere publico tech-savvy
- Blog no site indica estrategia de SEO/conteudo

#### Gaps
- Numero de seguidores no Instagram nao verificado
- Conteudo do blog nao acessado
- Canal no YouTube nao encontrado

---

### Suporte 24/7

#### Takeaway
FastProxy declara suporte 24/7 com atendimento humano em ate 15 minutos no plano VIP (Facebook Ads). Canais: WhatsApp, Telegram, email.

#### Cited Findings
- Suporte 24/7 declarado -- [FastProxy](https://www.fastproxy.com.br/)
- Atendimento humano em ate 15 minutos no plano VIP (Facebook Ads) -- [FastProxy](https://www.fastproxy.com.br/)
- Canais: WhatsApp, Telegram, Instagram, email -- [FastProxy](https://www.fastproxy.com.br/)

#### Inferences
- O SLA de 15 minutos apenas no plano VIP sugere que nos planos mais baratos o tempo de resposta pode ser maior
- Multiplos canais de atendimento e positivo mas nao ha como verificar a qualidade sem teste

#### Gaps
- SLA de resposta para planos nao-VIP nao especificado
- Nao foi possivel testar o suporte na pratica

---

### Reputacao e Reviews

#### Takeaway
Nenhuma review externa verificavel foi encontrada para FastProxy. As 237 reviews e nota 4.9 citadas no site nao tem fonte identificavel.

#### Cited Findings
- Nota 4.9 de 237 reviews declarada no site -- [FastProxy](https://www.fastproxy.com.br/)
- Nenhuma pagina no Reclame Aqui para FastProxy -- busca direta
- Nenhuma review no Trustpilot para fastproxy.com.br -- busca direta
- Nenhuma review em plataformas de terceiros encontrada -- buscas multiplas

#### Inferences
- A nota de 4.9 sem fonte verificavel deve ser tratada com ceticismo
- A ausencia total de presenca em plataformas de review e similar ao GTI Proxy
- Ambas as empresas operam em um nicho onde os clientes tipicamente nao deixam reviews publicas

#### Gaps
- Fonte das 237 reviews nao identificada
- Nao foi possivel verificar reviews em grupos privados de Telegram ou WhatsApp

---

### Dados de Retencao / Logs

#### Takeaway
FastProxy declara que mantem logs de conexao por 1 ano, citando o Marco Civil da Internet.

#### Cited Findings
- Dados registrados: IP de origem, data e hora da conexao, volume de trafego e host de destino -- [FastProxy Termos](https://www.fastproxy.com.br/termos)
- Periodo de retencao: 1 ano, citando Lei 12.965/2014 (Marco Civil da Internet) -- [FastProxy Termos](https://www.fastproxy.com.br/termos)
- FastProxy declara que nao monitora o conteudo do trafego, exceto quando abuso ameaca a integridade da rede -- [FastProxy Termos](https://www.fastproxy.com.br/termos)

#### Inferences
- A retencao de logs por 1 ano e uma exigencia legal para provedores de aplicacao no Brasil
- O registro de "host de destino" vai alem do minimo legal e pode preocupar usuarios focados em privacidade
- GTI Proxy nao divulga politica de logs de forma tao explicita

---

## Comparativo GTI Proxy vs. FastProxy

| Criterio | GTI Proxy | FastProxy |
|---|---|---|
| **CNPJ** | 18.915.849/0001-50 (verificado) | Nao divulgado |
| **Desde** | 2013 | Nao informado |
| **Sede** | Goiania/GO | Sao Paulo/SP (inferido) |
| **Preco mais baixo** | R$ 15 (24h movel) / R$ 40/mes (residencial) | R$ 29,90/mes (IPv6) |
| **Proxy Movel 4G/5G** | Sim (R$ 15-200) | Nao |
| **Residencial Fixo** | Sim (R$ 40/mes, banda ilimitada) | ISP a R$ 49,90/mes (100 Mbps) |
| **IPv6** | Sim (pacote 5 IPs R$ 40/mes) | Sim (R$ 29,90/mes por IP) |
| **IPv4 Dedicado** | Nao mencionado separadamente | R$ 39,90/mes |
| **Rotativo** | Movel com rotacao automatica | IPv6 Rotativo R$ 120/mes (10 min) |
| **Plano Facebook Ads** | Nao especifico | R$ 79,90/mes (VIP) |
| **Datacenter** | Nao | Nao (mas IPv6 pode ser datacenter) |
| **Pool de IPs** | 12.847 (verificavel) | Nao divulgado |
| **Protocolos** | HTTP, SOCKS5 | HTTP/S, SOCKS5 |
| **Pagamento** | Pix, cartao | Pix, cartao, boleto |
| **Nota Fiscal** | Sim (sob solicitacao) | Nao mencionado |
| **Teste Gratis** | 15 min residencial | Nao mencionado |
| **API** | Pagina existe (/api-proxy) | Nao documentada |
| **PWA/App** | Nao | Sim |
| **Antidetect** | 10+ navegadores | 6+ navegadores |
| **Suporte** | WhatsApp 24/7 | WhatsApp, Telegram, email 24/7 |
| **Reviews externas** | Nenhuma encontrada | Nenhuma encontrada |
| **Reclame Aqui** | Nenhuma reclamacao | Nenhuma reclamacao |
| **Plano anual** | Nao oferecido | 18% desconto |
| **Desconto volume** | Sim (ate R$ 30/IP para 50+) | Pacotes prepagos (ate 45% off) |

---

## Analise de Forcas e Fraquezas

### GTI Proxy - Forcas
- CNPJ verificavel e empresa ativa desde 2013 (maior credibilidade)
- Proxy movel 4G/5G (diferencial unico vs. FastProxy)
- Residencial fixo com banda ilimitada a R$ 40/mes
- Teste gratis de 15 minutos
- Descontos por volume claros e progressivos
- Nota fiscal brasileira
- Pool de IPs quantificado e datado (transparencia)
- Declaracao explicita de que nao publica metricas que nao pode medir (honestidade incomum)

### GTI Proxy - Fraquezas
- Sem boleto como forma de pagamento
- Presenca de marketing minimalista
- Sem plano anual com desconto
- Pool relativamente pequeno (~12.800 IPs)
- Sem plano especifico para Facebook Ads
- IPv6 residencial limitado a sites com IPv6

### FastProxy - Forcas
- Preco de entrada mais baixo (R$ 29,90/mes IPv6)
- Aceita boleto (vantagem para empresas)
- Plano anual com 18% de desconto
- Plano especifico para Facebook Ads com suporte VIP
- PWA instalavel (diferencial tecnologico)
- IPv6 rotativo com pool ilimitado
- Pacotes prepagos com descontos agressivos (ate 45%)
- Mais canais de suporte (WhatsApp + Telegram + email)

### FastProxy - Fraquezas
- CNPJ nao divulgado (risco de credibilidade)
- Nao oferece proxy movel 4G/5G
- ISP residencial mais caro (R$ 49,90 vs. R$ 40 do GTI) e com limite de 100 Mbps
- Reviews de 4.9/237 sem fonte verificavel
- Nao menciona nota fiscal
- Contadores de stats no site mostram "0" (placeholders)
- IPv6 basico pode ser datacenter (nao confirmado como residencial)
- Sem teste gratis evidente
- Logs de "host de destino" podem preocupar usuarios de privacidade

---

## Notas Metodologicas
- Pesquisa realizada em novembro/2026 usando buscas web e acesso direto aos sites
- Precos verificados nas paginas oficiais dos dois servicos
- CNPJ do GTI Proxy verificado via CNPJ API (fonte publica da Receita Federal)
- CNPJ do FastProxy nao encontrado em nenhuma fonte publica consultada
- Nenhuma review de terceiros foi encontrada para nenhum dos dois servicos
- Os precos podem ter mudado desde a data de verificacao
- O conteudo de paginas internas (FAQ colapsado, area logada, paginas de reembolso) nao foi totalmente acessado
