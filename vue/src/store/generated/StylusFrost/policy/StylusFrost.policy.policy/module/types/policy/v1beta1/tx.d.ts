import { Reader, Writer } from 'protobufjs/minimal';
export declare const protobufPackage = "policy.policy.v1beta1";
/** MsgStoreRego submit Rego code to the system */
export interface MsgStoreRego {
    /** Sender is the that actor that signed the messages */
    sender: string;
    /** REGOByteCode can be raw or gzip compressed */
    regoByteCode: Uint8Array;
    /** Valid entry points json encoded */
    entryPoints: Uint8Array;
    /**
     * Source is a valid absolute HTTPS URI to the policy's source code,
     * optional
     */
    source: string;
}
/** MsgStoreCodeResponse returns store result data. */
export interface MsgStoreRegoResponse {
    /** RegoID is the reference to the stored REGO code */
    regoId: number;
}
export declare const MsgStoreRego: {
    encode(message: MsgStoreRego, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgStoreRego;
    fromJSON(object: any): MsgStoreRego;
    toJSON(message: MsgStoreRego): unknown;
    fromPartial(object: DeepPartial<MsgStoreRego>): MsgStoreRego;
};
export declare const MsgStoreRegoResponse: {
    encode(message: MsgStoreRegoResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgStoreRegoResponse;
    fromJSON(object: any): MsgStoreRegoResponse;
    toJSON(message: MsgStoreRegoResponse): unknown;
    fromPartial(object: DeepPartial<MsgStoreRegoResponse>): MsgStoreRegoResponse;
};
/** Msg defines the Msg service. */
export interface Msg {
    /** StoreRego to submit Rego code to the system */
    StoreRego(request: MsgStoreRego): Promise<MsgStoreRegoResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    StoreRego(request: MsgStoreRego): Promise<MsgStoreRegoResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
