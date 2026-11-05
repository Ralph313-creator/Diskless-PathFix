export namespace gamefix {
	
	export class CCBoot {
	    line: string;
	    script: string;
	    added: boolean;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new CCBoot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.line = source["line"];
	        this.script = source["script"];
	        this.added = source["added"];
	        this.error = source["error"];
	    }
	}
	export class Folder {
	    name: string;
	    source: string;
	    link: string;
	    copied: boolean;
	    watch: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Folder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.source = source["source"];
	        this.link = source["link"];
	        this.copied = source["copied"];
	        this.watch = source["watch"];
	    }
	}
	export class Game {
	    name: string;
	    folder: string;
	    folders: Folder[];
	
	    static createFrom(source: any = {}) {
	        return new Game(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.folder = source["folder"];
	        this.folders = this.convertValues(source["folders"], Folder);
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

export namespace main {
	
	export class PackageView {
	    Name: string;
	    Dir: string;
	    Launcher: string;
	    Version: string;
	    Folders: core.FolderStatus[];
	    // Go type: core
	    Update?: any;
	    LastLog: string;
	    Error: string;
	    ScheduleMinutes: number;
	    CCBoot: string;
	    Registry: string[];
	    Session: string[][];
	    Links: string[];
	    Icon: string;
	    OwnIcon: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PackageView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Dir = source["Dir"];
	        this.Launcher = source["Launcher"];
	        this.Version = source["Version"];
	        this.Folders = this.convertValues(source["Folders"], core.FolderStatus);
	        this.Update = this.convertValues(source["Update"], null);
	        this.LastLog = source["LastLog"];
	        this.Error = source["Error"];
	        this.ScheduleMinutes = source["ScheduleMinutes"];
	        this.CCBoot = source["CCBoot"];
	        this.Registry = source["Registry"];
	        this.Session = source["Session"];
	        this.Links = source["Links"];
	        this.Icon = source["Icon"];
	        this.OwnIcon = source["OwnIcon"];
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

export namespace manager {
	
	export class Program {
	    Name: string;
	    PID: number;
	    Path: string;
	
	    static createFrom(source: any = {}) {
	        return new Program(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.PID = source["PID"];
	        this.Path = source["Path"];
	    }
	}
	export class Use {
	    Programs: Program[];
	    Services: string[];
	
	    static createFrom(source: any = {}) {
	        return new Use(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Programs = this.convertValues(source["Programs"], Program);
	        this.Services = source["Services"];
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

