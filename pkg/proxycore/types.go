package proxycore

type ProxyStatus string

const (
	StatusOffline    ProxyStatus = "offline"
	StatusConnecting ProxyStatus = "connecting"
	StatusConnected  ProxyStatus = "connected"
	StatusError      ProxyStatus = "error"
)

type GeoInfo struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	AS          string  `json:"as"`
}

type ProxyEntry struct {
	ID          string      `json:"id"`
	Host        string      `json:"host"`
	Port        string      `json:"port"`
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	Protocol    string      `json:"protocol"`
	Group       string      `json:"group"`
	Status      ProxyStatus `json:"status"`
	CurrentIP   string      `json:"currentIp"`
	Location    *GeoInfo    `json:"location"`
	Latency     int64       `json:"latency"`
	LastCheck   string      `json:"lastCheck"`
	Active      bool        `json:"active"`
	SessionID   string      `json:"sessionId"`
	Order       int         `json:"order"`
	ConnectedAt string      `json:"connectedAt"`
	TotalUptime int64       `json:"totalUptime"`
	Successes   int         `json:"successes,omitempty"`
	Failures    int         `json:"failures,omitempty"`
}

func (p *ProxyEntry) Score() int {
	return p.Successes*2 - p.Failures
}

func (p *ProxyEntry) Address() string {
	return p.Host + ":" + p.Port
}

type LogEntry struct {
	Time    string `json:"time"`
	Action  string `json:"action"`
	ProxyID string `json:"proxyId"`
	Host    string `json:"host"`
	Detail  string `json:"detail"`
}

type PoolStats struct {
	Total      int     `json:"total"`
	Active     int     `json:"active"`
	Healthy    int     `json:"healthy"`
	Errored    int     `json:"errored"`
	AvgLatency float64 `json:"avgLatencyMs"`
}

type FetchResult struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	ProxyID    string            `json:"proxyId"`
	ProxyIP    string            `json:"proxyIp"`
	LatencyMs  int64             `json:"latencyMs"`
	Truncated  bool              `json:"truncated,omitempty"`
}
