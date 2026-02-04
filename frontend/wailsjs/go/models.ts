export namespace app {
	
	export class SystemInfo {
	    os: string;
	    arch: string;
	    appVersion: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.arch = source["arch"];
	        this.appVersion = source["appVersion"];
	    }
	}

}

export namespace model {
	
	export class SpeedTestConfig {
	    timeout: number;
	    url: string;
	    pingUrl: string;
	    concurrent: number;
	
	    static createFrom(source: any = {}) {
	        return new SpeedTestConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timeout = source["timeout"];
	        this.url = source["url"];
	        this.pingUrl = source["pingUrl"];
	        this.concurrent = source["concurrent"];
	    }
	}
	export class UIConfig {
	    theme: string;
	    language: string;
	    fontFamily: string;
	    fontSize: number;
	    autoHideOnStart: boolean;
	    closeToTray: boolean;
	    showInDock: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.language = source["language"];
	        this.fontFamily = source["fontFamily"];
	        this.fontSize = source["fontSize"];
	        this.autoHideOnStart = source["autoHideOnStart"];
	        this.closeToTray = source["closeToTray"];
	        this.showInDock = source["showInDock"];
	    }
	}
	export class SystemProxyConfig {
	    exceptions: string;
	    notProxyLocal: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SystemProxyConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.exceptions = source["exceptions"];
	        this.notProxyLocal = source["notProxyLocal"];
	    }
	}
	export class TUNConfig {
	    enabled: boolean;
	    autoRoute: boolean;
	    strictRoute: boolean;
	    stack: string;
	    mtu: number;
	    enableIPv6: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TUNConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.autoRoute = source["autoRoute"];
	        this.strictRoute = source["strictRoute"];
	        this.stack = source["stack"];
	        this.mtu = source["mtu"];
	        this.enableIPv6 = source["enableIPv6"];
	    }
	}
	export class InboundConfig {
	    protocol: string;
	    listenAddr: string;
	    port: number;
	    udpEnabled: boolean;
	    sniffingEnabled: boolean;
	    allowLAN: boolean;
	
	    static createFrom(source: any = {}) {
	        return new InboundConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.protocol = source["protocol"];
	        this.listenAddr = source["listenAddr"];
	        this.port = source["port"];
	        this.udpEnabled = source["udpEnabled"];
	        this.sniffingEnabled = source["sniffingEnabled"];
	        this.allowLAN = source["allowLAN"];
	    }
	}
	export class CoreBasicConfig {
	    logEnabled: boolean;
	    logLevel: string;
	    muxEnabled: boolean;
	    allowInsecure: boolean;
	    fingerprint: string;
	    enableFragment: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CoreBasicConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.logEnabled = source["logEnabled"];
	        this.logLevel = source["logLevel"];
	        this.muxEnabled = source["muxEnabled"];
	        this.allowInsecure = source["allowInsecure"];
	        this.fingerprint = source["fingerprint"];
	        this.enableFragment = source["enableFragment"];
	    }
	}
	export class Config {
	    activeProfileId: string;
	    activeRoutingId: string;
	    activeDnsPreset: string;
	    coreBasic: CoreBasicConfig;
	    inbounds: InboundConfig[];
	    proxyMode: number;
	    tun: TUNConfig;
	    systemProxy: SystemProxyConfig;
	    ui: UIConfig;
	    speedTest: SpeedTestConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.activeProfileId = source["activeProfileId"];
	        this.activeRoutingId = source["activeRoutingId"];
	        this.activeDnsPreset = source["activeDnsPreset"];
	        this.coreBasic = this.convertValues(source["coreBasic"], CoreBasicConfig);
	        this.inbounds = this.convertValues(source["inbounds"], InboundConfig);
	        this.proxyMode = source["proxyMode"];
	        this.tun = this.convertValues(source["tun"], TUNConfig);
	        this.systemProxy = this.convertValues(source["systemProxy"], SystemProxyConfig);
	        this.ui = this.convertValues(source["ui"], UIConfig);
	        this.speedTest = this.convertValues(source["speedTest"], SpeedTestConfig);
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
	
	export class CoreStatus {
	    running: boolean;
	    coreType: number;
	    version: string;
	    startTime?: number;
	    pid?: number;
	    profile?: string;
	
	    static createFrom(source: any = {}) {
	        return new CoreStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.coreType = source["coreType"];
	        this.version = source["version"];
	        this.startTime = source["startTime"];
	        this.pid = source["pid"];
	        this.profile = source["profile"];
	    }
	}
	export class DNSItem {
	    remoteDns: string;
	    directDns: string;
	    bootstrapDns: string;
	    useSystemHosts: boolean;
	    fakeIP: boolean;
	    hosts: string;
	    domainStrategy: string;
	
	    static createFrom(source: any = {}) {
	        return new DNSItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.remoteDns = source["remoteDns"];
	        this.directDns = source["directDns"];
	        this.bootstrapDns = source["bootstrapDns"];
	        this.useSystemHosts = source["useSystemHosts"];
	        this.fakeIP = source["fakeIP"];
	        this.hosts = source["hosts"];
	        this.domainStrategy = source["domainStrategy"];
	    }
	}
	
	export class ProfileItem {
	    id: string;
	    configType: number;
	    remarks: string;
	    subId: string;
	    shareUri: string;
	    sort: number;
	    address: string;
	    port: number;
	    ports?: string;
	    uuid: string;
	    alterId?: number;
	    security: string;
	    flow?: string;
	    network: string;
	    headerType?: string;
	    host?: string;
	    path?: string;
	    streamSecurity: string;
	    allowInsecure: boolean;
	    sni?: string;
	    alpn?: string;
	    fingerprint?: string;
	    publicKey?: string;
	    shortId?: string;
	    spiderX?: string;
	    coreType?: number;
	    extra?: string;
	    muxEnabled?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProfileItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.configType = source["configType"];
	        this.remarks = source["remarks"];
	        this.subId = source["subId"];
	        this.shareUri = source["shareUri"];
	        this.sort = source["sort"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.ports = source["ports"];
	        this.uuid = source["uuid"];
	        this.alterId = source["alterId"];
	        this.security = source["security"];
	        this.flow = source["flow"];
	        this.network = source["network"];
	        this.headerType = source["headerType"];
	        this.host = source["host"];
	        this.path = source["path"];
	        this.streamSecurity = source["streamSecurity"];
	        this.allowInsecure = source["allowInsecure"];
	        this.sni = source["sni"];
	        this.alpn = source["alpn"];
	        this.fingerprint = source["fingerprint"];
	        this.publicKey = source["publicKey"];
	        this.shortId = source["shortId"];
	        this.spiderX = source["spiderX"];
	        this.coreType = source["coreType"];
	        this.extra = source["extra"];
	        this.muxEnabled = source["muxEnabled"];
	    }
	}
	export class RuleItem {
	    id: string;
	    outboundTag: string;
	    enabled: boolean;
	    remarks?: string;
	    domain?: string[];
	    domainSuffix?: string[];
	    domainKeyword?: string[];
	    domainRegex?: string[];
	    geosite?: string[];
	    ip?: string[];
	    ipCidr?: string[];
	    geoip?: string[];
	    port?: string;
	    protocol?: string[];
	    processName?: string[];
	    network?: string;
	    inbound?: string[];
	    ruleSet?: string[];
	
	    static createFrom(source: any = {}) {
	        return new RuleItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.outboundTag = source["outboundTag"];
	        this.enabled = source["enabled"];
	        this.remarks = source["remarks"];
	        this.domain = source["domain"];
	        this.domainSuffix = source["domainSuffix"];
	        this.domainKeyword = source["domainKeyword"];
	        this.domainRegex = source["domainRegex"];
	        this.geosite = source["geosite"];
	        this.ip = source["ip"];
	        this.ipCidr = source["ipCidr"];
	        this.geoip = source["geoip"];
	        this.port = source["port"];
	        this.protocol = source["protocol"];
	        this.processName = source["processName"];
	        this.network = source["network"];
	        this.inbound = source["inbound"];
	        this.ruleSet = source["ruleSet"];
	    }
	}
	export class RoutingItem {
	    id: string;
	    remarks: string;
	    domainStrategy: string;
	    rules: RuleItem[];
	    enabled: boolean;
	    locked: boolean;
	    sort: number;
	
	    static createFrom(source: any = {}) {
	        return new RoutingItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.remarks = source["remarks"];
	        this.domainStrategy = source["domainStrategy"];
	        this.rules = this.convertValues(source["rules"], RuleItem);
	        this.enabled = source["enabled"];
	        this.locked = source["locked"];
	        this.sort = source["sort"];
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
	
	export class ServerStatItem {
	    profileId: string;
	    totalUp: number;
	    totalDown: number;
	    todayUp: number;
	    todayDown: number;
	    dateNow: string;
	    lastUpdate: number;
	
	    static createFrom(source: any = {}) {
	        return new ServerStatItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.profileId = source["profileId"];
	        this.totalUp = source["totalUp"];
	        this.totalDown = source["totalDown"];
	        this.todayUp = source["todayUp"];
	        this.todayDown = source["todayDown"];
	        this.dateNow = source["dateNow"];
	        this.lastUpdate = source["lastUpdate"];
	    }
	}
	
	export class SpeedTestResult {
	    profileId: string;
	    latency: number;
	    speed: number;
	
	    static createFrom(source: any = {}) {
	        return new SpeedTestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.profileId = source["profileId"];
	        this.latency = source["latency"];
	        this.speed = source["speed"];
	    }
	}
	export class SubItem {
	    id: string;
	    remarks: string;
	    url: string;
	    enabled: boolean;
	    sort: number;
	    filter?: string;
	    autoUpdateInterval: number;
	    updateTime: number;
	    userAgent?: string;
	
	    static createFrom(source: any = {}) {
	        return new SubItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.remarks = source["remarks"];
	        this.url = source["url"];
	        this.enabled = source["enabled"];
	        this.sort = source["sort"];
	        this.filter = source["filter"];
	        this.autoUpdateInterval = source["autoUpdateInterval"];
	        this.updateTime = source["updateTime"];
	        this.userAgent = source["userAgent"];
	    }
	}
	
	

}

export namespace service {
	
	export class Metadata {
	    network: string;
	    type: string;
	    sourceIP: string;
	    destinationIP: string;
	    sourcePort: string;
	    destinationPort: string;
	    host: string;
	    dnsMode: string;
	    processPath: string;
	
	    static createFrom(source: any = {}) {
	        return new Metadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.network = source["network"];
	        this.type = source["type"];
	        this.sourceIP = source["sourceIP"];
	        this.destinationIP = source["destinationIP"];
	        this.sourcePort = source["sourcePort"];
	        this.destinationPort = source["destinationPort"];
	        this.host = source["host"];
	        this.dnsMode = source["dnsMode"];
	        this.processPath = source["processPath"];
	    }
	}
	export class Connection {
	    id: string;
	    metadata: Metadata;
	    upload: number;
	    download: number;
	    start: string;
	    chains: string[];
	    rule: string;
	    rulePayload: string;
	
	    static createFrom(source: any = {}) {
	        return new Connection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.metadata = this.convertValues(source["metadata"], Metadata);
	        this.upload = source["upload"];
	        this.download = source["download"];
	        this.start = source["start"];
	        this.chains = source["chains"];
	        this.rule = source["rule"];
	        this.rulePayload = source["rulePayload"];
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
	export class ConnectionsResponse {
	    downloadTotal: number;
	    uploadTotal: number;
	    connections: Connection[];
	
	    static createFrom(source: any = {}) {
	        return new ConnectionsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.downloadTotal = source["downloadTotal"];
	        this.uploadTotal = source["uploadTotal"];
	        this.connections = this.convertValues(source["connections"], Connection);
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
	export class GeoDataInfo {
	    geoipVersion: string;
	    geositeVersion: string;
	    geoipPath: string;
	    geositePath: string;
	    lastUpdated: number;
	
	    static createFrom(source: any = {}) {
	        return new GeoDataInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.geoipVersion = source["geoipVersion"];
	        this.geositeVersion = source["geositeVersion"];
	        this.geoipPath = source["geoipPath"];
	        this.geositePath = source["geositePath"];
	        this.lastUpdated = source["lastUpdated"];
	    }
	}
	
	export class UpdateInfo {
	    coreType: number;
	    currentVersion: string;
	    latestVersion: string;
	    hasUpdate: boolean;
	    downloadUrl: string;
	    assetName: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.coreType = source["coreType"];
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.hasUpdate = source["hasUpdate"];
	        this.downloadUrl = source["downloadUrl"];
	        this.assetName = source["assetName"];
	    }
	}

}

