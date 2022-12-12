import * as jspb from 'google-protobuf'

import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';


export class InfoReq extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InfoReq.AsObject;
  static toObject(includeInstance: boolean, msg: InfoReq): InfoReq.AsObject;
  static serializeBinaryToWriter(message: InfoReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InfoReq;
  static deserializeBinaryFromReader(message: InfoReq, reader: jspb.BinaryReader): InfoReq;
}

export namespace InfoReq {
  export type AsObject = {
  }
}

export class InfoRes extends jspb.Message {
  getPublicKey(): string;
  setPublicKey(value: string): InfoRes;

  getHost(): google_protobuf_wrappers_pb.StringValue | undefined;
  setHost(value?: google_protobuf_wrappers_pb.StringValue): InfoRes;
  hasHost(): boolean;
  clearHost(): InfoRes;

  getPort(): number;
  setPort(value: number): InfoRes;

  getHostVpnIp(): string;
  setHostVpnIp(value: string): InfoRes;

  getIsAdmin(): boolean;
  setIsAdmin(value: boolean): InfoRes;

  getAllowedIps(): string;
  setAllowedIps(value: string): InfoRes;

  getDnsEnabled(): boolean;
  setDnsEnabled(value: boolean): InfoRes;

  getDnsAddress(): string;
  setDnsAddress(value: string): InfoRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InfoRes.AsObject;
  static toObject(includeInstance: boolean, msg: InfoRes): InfoRes.AsObject;
  static serializeBinaryToWriter(message: InfoRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InfoRes;
  static deserializeBinaryFromReader(message: InfoRes, reader: jspb.BinaryReader): InfoRes;
}

export namespace InfoRes {
  export type AsObject = {
    publicKey: string,
    host?: google_protobuf_wrappers_pb.StringValue.AsObject,
    port: number,
    hostVpnIp: string,
    isAdmin: boolean,
    allowedIps: string,
    dnsEnabled: boolean,
    dnsAddress: string,
  }
}

