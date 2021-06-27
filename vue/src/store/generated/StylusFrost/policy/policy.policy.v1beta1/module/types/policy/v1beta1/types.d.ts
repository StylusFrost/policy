import { Writer, Reader } from 'protobufjs/minimal';
export declare const protobufPackage = "policy.policy.v1beta1";
/** RegoInfo is data for the uploaded policy REGO code */
export interface RegoInfo {
    /** RegoHash is the unique identifier */
    regoHash: Uint8Array;
    /** Creator address who initially stored the code */
    creator: string;
    /**
     * Source is a valid absolute HTTPS URI to the policy's source code,
     * optional
     */
    source: string;
    /** Valid entry points when using the Rego code */
    entryPoints: string[];
}
export declare const RegoInfo: {
    encode(message: RegoInfo, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): RegoInfo;
    fromJSON(object: any): RegoInfo;
    toJSON(message: RegoInfo): unknown;
    fromPartial(object: DeepPartial<RegoInfo>): RegoInfo;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
