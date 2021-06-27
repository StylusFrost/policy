/* eslint-disable */
import { Writer, Reader } from 'protobufjs/minimal';
export const protobufPackage = 'policy.policy.v1beta1';
const baseRegoInfo = { creator: '', source: '', entryPoints: '' };
export const RegoInfo = {
    encode(message, writer = Writer.create()) {
        if (message.regoHash.length !== 0) {
            writer.uint32(10).bytes(message.regoHash);
        }
        if (message.creator !== '') {
            writer.uint32(18).string(message.creator);
        }
        if (message.source !== '') {
            writer.uint32(26).string(message.source);
        }
        for (const v of message.entryPoints) {
            writer.uint32(34).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseRegoInfo };
        message.entryPoints = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.regoHash = reader.bytes();
                    break;
                case 2:
                    message.creator = reader.string();
                    break;
                case 3:
                    message.source = reader.string();
                    break;
                case 4:
                    message.entryPoints.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseRegoInfo };
        message.entryPoints = [];
        if (object.regoHash !== undefined && object.regoHash !== null) {
            message.regoHash = bytesFromBase64(object.regoHash);
        }
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = String(object.creator);
        }
        else {
            message.creator = '';
        }
        if (object.source !== undefined && object.source !== null) {
            message.source = String(object.source);
        }
        else {
            message.source = '';
        }
        if (object.entryPoints !== undefined && object.entryPoints !== null) {
            for (const e of object.entryPoints) {
                message.entryPoints.push(String(e));
            }
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.regoHash !== undefined && (obj.regoHash = base64FromBytes(message.regoHash !== undefined ? message.regoHash : new Uint8Array()));
        message.creator !== undefined && (obj.creator = message.creator);
        message.source !== undefined && (obj.source = message.source);
        if (message.entryPoints) {
            obj.entryPoints = message.entryPoints.map((e) => e);
        }
        else {
            obj.entryPoints = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseRegoInfo };
        message.entryPoints = [];
        if (object.regoHash !== undefined && object.regoHash !== null) {
            message.regoHash = object.regoHash;
        }
        else {
            message.regoHash = new Uint8Array();
        }
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = object.creator;
        }
        else {
            message.creator = '';
        }
        if (object.source !== undefined && object.source !== null) {
            message.source = object.source;
        }
        else {
            message.source = '';
        }
        if (object.entryPoints !== undefined && object.entryPoints !== null) {
            for (const e of object.entryPoints) {
                message.entryPoints.push(e);
            }
        }
        return message;
    }
};
var globalThis = (() => {
    if (typeof globalThis !== 'undefined')
        return globalThis;
    if (typeof self !== 'undefined')
        return self;
    if (typeof window !== 'undefined')
        return window;
    if (typeof global !== 'undefined')
        return global;
    throw 'Unable to locate global object';
})();
const atob = globalThis.atob || ((b64) => globalThis.Buffer.from(b64, 'base64').toString('binary'));
function bytesFromBase64(b64) {
    const bin = atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
        arr[i] = bin.charCodeAt(i);
    }
    return arr;
}
const btoa = globalThis.btoa || ((bin) => globalThis.Buffer.from(bin, 'binary').toString('base64'));
function base64FromBytes(arr) {
    const bin = [];
    for (let i = 0; i < arr.byteLength; ++i) {
        bin.push(String.fromCharCode(arr[i]));
    }
    return btoa(bin.join(''));
}
