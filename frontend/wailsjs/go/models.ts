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

