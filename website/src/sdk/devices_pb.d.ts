import * as jspb from 'google-protobuf'

import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';


export class Device extends jspb.Message {
  getName(): string;
  setName(value: string): Device;

  getOwner(): string;
  setOwner(value: string): Device;

  getPublicKey(): string;
  setPublicKey(value: string): Device;

  getAddress(): string;
  setAddress(value: string): Device;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Device;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Device;

  getConnected(): boolean;
  setConnected(value: boolean): Device;

  getLastHandshakeTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastHandshakeTime(value?: google_protobuf_timestamp_pb.Timestamp): Device;
  hasLastHandshakeTime(): boolean;
  clearLastHandshakeTime(): Device;

  getReceiveBytes(): number;
  setReceiveBytes(value: number): Device;

  getTransmitBytes(): number;
  setTransmitBytes(value: number): Device;

  getEndpoint(): string;
  setEndpoint(value: string): Device;

  getOwnerName(): string;
  setOwnerName(value: string): Device;

  getOwnerEmail(): string;
  setOwnerEmail(value: string): Device;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Device.AsObject;
  static toObject(includeInstance: boolean, msg: Device): Device.AsObject;
  static serializeBinaryToWriter(message: Device, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Device;
  static deserializeBinaryFromReader(message: Device, reader: jspb.BinaryReader): Device;
}

export namespace Device {
  export type AsObject = {
    name: string,
    owner: string,
    publicKey: string,
    address: string,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    connected: boolean,
    lastHandshakeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    receiveBytes: number,
    transmitBytes: number,
    endpoint: string,
    ownerName: string,
    ownerEmail: string,
  }
}

export class AddDeviceReq extends jspb.Message {
  getName(): string;
  setName(value: string): AddDeviceReq;

  getPublicKey(): string;
  setPublicKey(value: string): AddDeviceReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDeviceReq.AsObject;
  static toObject(includeInstance: boolean, msg: AddDeviceReq): AddDeviceReq.AsObject;
  static serializeBinaryToWriter(message: AddDeviceReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddDeviceReq;
  static deserializeBinaryFromReader(message: AddDeviceReq, reader: jspb.BinaryReader): AddDeviceReq;
}

export namespace AddDeviceReq {
  export type AsObject = {
    name: string,
    publicKey: string,
  }
}

export class ListDevicesReq extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDevicesReq.AsObject;
  static toObject(includeInstance: boolean, msg: ListDevicesReq): ListDevicesReq.AsObject;
  static serializeBinaryToWriter(message: ListDevicesReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDevicesReq;
  static deserializeBinaryFromReader(message: ListDevicesReq, reader: jspb.BinaryReader): ListDevicesReq;
}

export namespace ListDevicesReq {
  export type AsObject = {
  }
}

export class ListDevicesRes extends jspb.Message {
  getItemsList(): Array<Device>;
  setItemsList(value: Array<Device>): ListDevicesRes;
  clearItemsList(): ListDevicesRes;
  addItems(value?: Device, index?: number): Device;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDevicesRes.AsObject;
  static toObject(includeInstance: boolean, msg: ListDevicesRes): ListDevicesRes.AsObject;
  static serializeBinaryToWriter(message: ListDevicesRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDevicesRes;
  static deserializeBinaryFromReader(message: ListDevicesRes, reader: jspb.BinaryReader): ListDevicesRes;
}

export namespace ListDevicesRes {
  export type AsObject = {
    itemsList: Array<Device.AsObject>,
  }
}

export class DeleteDeviceReq extends jspb.Message {
  getName(): string;
  setName(value: string): DeleteDeviceReq;

  getOwner(): google_protobuf_wrappers_pb.StringValue | undefined;
  setOwner(value?: google_protobuf_wrappers_pb.StringValue): DeleteDeviceReq;
  hasOwner(): boolean;
  clearOwner(): DeleteDeviceReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteDeviceReq.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteDeviceReq): DeleteDeviceReq.AsObject;
  static serializeBinaryToWriter(message: DeleteDeviceReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteDeviceReq;
  static deserializeBinaryFromReader(message: DeleteDeviceReq, reader: jspb.BinaryReader): DeleteDeviceReq;
}

export namespace DeleteDeviceReq {
  export type AsObject = {
    name: string,
    owner?: google_protobuf_wrappers_pb.StringValue.AsObject,
  }
}

export class ListAllDevicesReq extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAllDevicesReq.AsObject;
  static toObject(includeInstance: boolean, msg: ListAllDevicesReq): ListAllDevicesReq.AsObject;
  static serializeBinaryToWriter(message: ListAllDevicesReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAllDevicesReq;
  static deserializeBinaryFromReader(message: ListAllDevicesReq, reader: jspb.BinaryReader): ListAllDevicesReq;
}

export namespace ListAllDevicesReq {
  export type AsObject = {
  }
}

export class ListAllDevicesRes extends jspb.Message {
  getItemsList(): Array<Device>;
  setItemsList(value: Array<Device>): ListAllDevicesRes;
  clearItemsList(): ListAllDevicesRes;
  addItems(value?: Device, index?: number): Device;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAllDevicesRes.AsObject;
  static toObject(includeInstance: boolean, msg: ListAllDevicesRes): ListAllDevicesRes.AsObject;
  static serializeBinaryToWriter(message: ListAllDevicesRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAllDevicesRes;
  static deserializeBinaryFromReader(message: ListAllDevicesRes, reader: jspb.BinaryReader): ListAllDevicesRes;
}

export namespace ListAllDevicesRes {
  export type AsObject = {
    itemsList: Array<Device.AsObject>,
  }
}

