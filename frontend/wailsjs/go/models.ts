export namespace main {
	
	export class ClientConfig {
	    ServerIp: string;
	    ServerPort: string;
	
	    static createFrom(source: any = {}) {
	        return new ClientConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ServerIp = source["ServerIp"];
	        this.ServerPort = source["ServerPort"];
	    }
	}
	export class TransferConfig {
	    srcAddr: string;
	    srcPort: string;
	    dstAddr: string;
	    dstPort: string;
	
	    static createFrom(source: any = {}) {
	        return new TransferConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.srcAddr = source["srcAddr"];
	        this.srcPort = source["srcPort"];
	        this.dstAddr = source["dstAddr"];
	        this.dstPort = source["dstPort"];
	    }
	}
	export class ServerConfig {
	    tcpAddr: string;
	    tcpPort: string;
	    udpAddr: string;
	    udpPort: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tcpAddr = source["tcpAddr"];
	        this.tcpPort = source["tcpPort"];
	        this.udpAddr = source["udpAddr"];
	        this.udpPort = source["udpPort"];
	    }
	}
	export class Config {
	    Server: ServerConfig;
	    Transfer: TransferConfig;
	    Client: ClientConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Server = this.convertValues(source["Server"], ServerConfig);
	        this.Transfer = this.convertValues(source["Transfer"], TransferConfig);
	        this.Client = this.convertValues(source["Client"], ClientConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

