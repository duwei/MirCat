// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

export function ClientTcpClose(arg1:number):Promise<void>;

export function ClientTcpCloseAll():Promise<void>;

export function ClientTcpOpen():Promise<number>;

export function ClientTcpSend(arg1:number,arg2:string):Promise<void>;

export function ServerBroadcastMessage(arg1:string):Promise<void>;

export function ServerSendMessage(arg1:string,arg2:string):Promise<void>;

export function ServerTcpStart():Promise<boolean>;

export function ServerTcpStop():Promise<boolean>;

export function TransferBroadcastToClient(arg1:string):Promise<void>;

export function TransferBroadcastToServer(arg1:string):Promise<void>;

export function TransferSendToClient(arg1:string,arg2:string):Promise<void>;

export function TransferSendToServer(arg1:string,arg2:string):Promise<void>;

export function TransferTcpStart():Promise<boolean>;

export function TransferTcpStop():Promise<boolean>;
