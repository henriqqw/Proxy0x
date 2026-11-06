export namespace main {
	
	export class GeoInfo {
	    country: string;
	    countryCode: string;
	    region: string;
	    city: string;
	    isp: string;
	    org: string;
	    lat: number;
	    lon: number;
	    timezone: string;
	    as: string;
	
	    static createFrom(source: any = {}) {
	        return new GeoInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.country = source["country"];
	        this.countryCode = source["countryCode"];
	        this.region = source["region"];
	        this.city = source["city"];
	        this.isp = source["isp"];
	        this.org = source["org"];
	        this.lat = source["lat"];
	        this.lon = source["lon"];
	        this.timezone = source["timezone"];
	        this.as = source["as"];
	    }
	}
	export class LogEntry {
	    time: string;
	    action: string;
	    proxyId: string;
	    host: string;
	    detail: string;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.action = source["action"];
	        this.proxyId = source["proxyId"];
	        this.host = source["host"];
	        this.detail = source["detail"];
	    }
	}
	export class MCPStatus {
	    installed: boolean;
	    exePath: string;
	    exeExists: boolean;
	    configPath: string;
	    healthOk: boolean;
	    healthError?: string;
	
	    static createFrom(source: any = {}) {
	        return new MCPStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.installed = source["installed"];
	        this.exePath = source["exePath"];
	        this.exeExists = source["exeExists"];
	        this.configPath = source["configPath"];
	        this.healthOk = source["healthOk"];
	        this.healthError = source["healthError"];
	    }
	}
	export class ProxyEntry {
	    id: string;
	    host: string;
	    port: string;
	    username: string;
	    password: string;
	    protocol: string;
	    group: string;
	    status: string;
	    currentIp: string;
	    location?: GeoInfo;
	    latency: number;
	    lastCheck: string;
	    active: boolean;
	    sessionId: string;
	    order: number;
	    connectedAt: string;
	    totalUptime: number;
	
	    static createFrom(source: any = {}) {
	        return new ProxyEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.protocol = source["protocol"];
	        this.group = source["group"];
	        this.status = source["status"];
	        this.currentIp = source["currentIp"];
	        this.location = this.convertValues(source["location"], GeoInfo);
	        this.latency = source["latency"];
	        this.lastCheck = source["lastCheck"];
	        this.active = source["active"];
	        this.sessionId = source["sessionId"];
	        this.order = source["order"];
	        this.connectedAt = source["connectedAt"];
	        this.totalUptime = source["totalUptime"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

