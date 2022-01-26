declare namespace report {

    export class EventReport {

        constructor(serverUrl: string, appid: string, appkey: string, debug: number);

        public identify(uuid :string): void;

        public login(uid :string): void;

        public logout(): void;

        public getDistinctId():string;

        public track(eventName :string,properties :any): void;

        public userSet(properties :any):EventReport;

        public userUnset(key :string):void;
        
        public userAdd(properties :any):void;

        public userSetOnce(properties :any):void;

        public getSuperProperties():any;
        
        public getUserProperties():any;

        public setSuperProperties(properties :any):void;

        public unsetSuperProperties(key :string):void;

        public clearSuperProperties():void;

        public getUserProperties():any;

        public trackUserData():void;

    }
}