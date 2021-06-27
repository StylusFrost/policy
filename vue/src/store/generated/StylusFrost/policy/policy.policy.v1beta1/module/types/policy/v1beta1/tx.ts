/* eslint-disable */
import { Reader, util, configure, Writer } from 'protobufjs/minimal'
import * as Long from 'long'

export const protobufPackage = 'policy.policy.v1beta1'

/** MsgStoreRego submit Rego code to the system */
export interface MsgStoreRego {
  /** Sender is the that actor that signed the messages */
  sender: string
  /** REGOByteCode can be raw or gzip compressed */
  regoByteCode: Uint8Array
  /** Valid entry points json encoded */
  entryPoints: Uint8Array
  /**
   * Source is a valid absolute HTTPS URI to the policy's source code,
   * optional
   */
  source: string
}

/** MsgStoreCodeResponse returns store result data. */
export interface MsgStoreRegoResponse {
  /** RegoID is the reference to the stored REGO code */
  regoId: number
}

const baseMsgStoreRego: object = { sender: '', source: '' }

export const MsgStoreRego = {
  encode(message: MsgStoreRego, writer: Writer = Writer.create()): Writer {
    if (message.sender !== '') {
      writer.uint32(10).string(message.sender)
    }
    if (message.regoByteCode.length !== 0) {
      writer.uint32(18).bytes(message.regoByteCode)
    }
    if (message.entryPoints.length !== 0) {
      writer.uint32(26).bytes(message.entryPoints)
    }
    if (message.source !== '') {
      writer.uint32(34).string(message.source)
    }
    return writer
  },

  decode(input: Reader | Uint8Array, length?: number): MsgStoreRego {
    const reader = input instanceof Uint8Array ? new Reader(input) : input
    let end = length === undefined ? reader.len : reader.pos + length
    const message = { ...baseMsgStoreRego } as MsgStoreRego
    while (reader.pos < end) {
      const tag = reader.uint32()
      switch (tag >>> 3) {
        case 1:
          message.sender = reader.string()
          break
        case 2:
          message.regoByteCode = reader.bytes()
          break
        case 3:
          message.entryPoints = reader.bytes()
          break
        case 4:
          message.source = reader.string()
          break
        default:
          reader.skipType(tag & 7)
          break
      }
    }
    return message
  },

  fromJSON(object: any): MsgStoreRego {
    const message = { ...baseMsgStoreRego } as MsgStoreRego
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = String(object.sender)
    } else {
      message.sender = ''
    }
    if (object.regoByteCode !== undefined && object.regoByteCode !== null) {
      message.regoByteCode = bytesFromBase64(object.regoByteCode)
    }
    if (object.entryPoints !== undefined && object.entryPoints !== null) {
      message.entryPoints = bytesFromBase64(object.entryPoints)
    }
    if (object.source !== undefined && object.source !== null) {
      message.source = String(object.source)
    } else {
      message.source = ''
    }
    return message
  },

  toJSON(message: MsgStoreRego): unknown {
    const obj: any = {}
    message.sender !== undefined && (obj.sender = message.sender)
    message.regoByteCode !== undefined && (obj.regoByteCode = base64FromBytes(message.regoByteCode !== undefined ? message.regoByteCode : new Uint8Array()))
    message.entryPoints !== undefined && (obj.entryPoints = base64FromBytes(message.entryPoints !== undefined ? message.entryPoints : new Uint8Array()))
    message.source !== undefined && (obj.source = message.source)
    return obj
  },

  fromPartial(object: DeepPartial<MsgStoreRego>): MsgStoreRego {
    const message = { ...baseMsgStoreRego } as MsgStoreRego
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = object.sender
    } else {
      message.sender = ''
    }
    if (object.regoByteCode !== undefined && object.regoByteCode !== null) {
      message.regoByteCode = object.regoByteCode
    } else {
      message.regoByteCode = new Uint8Array()
    }
    if (object.entryPoints !== undefined && object.entryPoints !== null) {
      message.entryPoints = object.entryPoints
    } else {
      message.entryPoints = new Uint8Array()
    }
    if (object.source !== undefined && object.source !== null) {
      message.source = object.source
    } else {
      message.source = ''
    }
    return message
  }
}

const baseMsgStoreRegoResponse: object = { regoId: 0 }

export const MsgStoreRegoResponse = {
  encode(message: MsgStoreRegoResponse, writer: Writer = Writer.create()): Writer {
    if (message.regoId !== 0) {
      writer.uint32(8).uint64(message.regoId)
    }
    return writer
  },

  decode(input: Reader | Uint8Array, length?: number): MsgStoreRegoResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input
    let end = length === undefined ? reader.len : reader.pos + length
    const message = { ...baseMsgStoreRegoResponse } as MsgStoreRegoResponse
    while (reader.pos < end) {
      const tag = reader.uint32()
      switch (tag >>> 3) {
        case 1:
          message.regoId = longToNumber(reader.uint64() as Long)
          break
        default:
          reader.skipType(tag & 7)
          break
      }
    }
    return message
  },

  fromJSON(object: any): MsgStoreRegoResponse {
    const message = { ...baseMsgStoreRegoResponse } as MsgStoreRegoResponse
    if (object.regoId !== undefined && object.regoId !== null) {
      message.regoId = Number(object.regoId)
    } else {
      message.regoId = 0
    }
    return message
  },

  toJSON(message: MsgStoreRegoResponse): unknown {
    const obj: any = {}
    message.regoId !== undefined && (obj.regoId = message.regoId)
    return obj
  },

  fromPartial(object: DeepPartial<MsgStoreRegoResponse>): MsgStoreRegoResponse {
    const message = { ...baseMsgStoreRegoResponse } as MsgStoreRegoResponse
    if (object.regoId !== undefined && object.regoId !== null) {
      message.regoId = object.regoId
    } else {
      message.regoId = 0
    }
    return message
  }
}

/** Msg defines the Msg service. */
export interface Msg {
  /** StoreRego to submit Rego code to the system */
  StoreRego(request: MsgStoreRego): Promise<MsgStoreRegoResponse>
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc
  constructor(rpc: Rpc) {
    this.rpc = rpc
  }
  StoreRego(request: MsgStoreRego): Promise<MsgStoreRegoResponse> {
    const data = MsgStoreRego.encode(request).finish()
    const promise = this.rpc.request('policy.policy.v1beta1.Msg', 'StoreRego', data)
    return promise.then((data) => MsgStoreRegoResponse.decode(new Reader(data)))
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>
}

declare var self: any | undefined
declare var window: any | undefined
var globalThis: any = (() => {
  if (typeof globalThis !== 'undefined') return globalThis
  if (typeof self !== 'undefined') return self
  if (typeof window !== 'undefined') return window
  if (typeof global !== 'undefined') return global
  throw 'Unable to locate global object'
})()

const atob: (b64: string) => string = globalThis.atob || ((b64) => globalThis.Buffer.from(b64, 'base64').toString('binary'))
function bytesFromBase64(b64: string): Uint8Array {
  const bin = atob(b64)
  const arr = new Uint8Array(bin.length)
  for (let i = 0; i < bin.length; ++i) {
    arr[i] = bin.charCodeAt(i)
  }
  return arr
}

const btoa: (bin: string) => string = globalThis.btoa || ((bin) => globalThis.Buffer.from(bin, 'binary').toString('base64'))
function base64FromBytes(arr: Uint8Array): string {
  const bin: string[] = []
  for (let i = 0; i < arr.byteLength; ++i) {
    bin.push(String.fromCharCode(arr[i]))
  }
  return btoa(bin.join(''))
}

type Builtin = Date | Function | Uint8Array | string | number | undefined
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error('Value is larger than Number.MAX_SAFE_INTEGER')
  }
  return long.toNumber()
}

if (util.Long !== Long) {
  util.Long = Long as any
  configure()
}
