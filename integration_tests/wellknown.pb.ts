/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import { Observable } from 'rxjs';
import * as GoogleRpcStatus from "./google/rpc/status.pb"


export enum AEnum {
  VALUE_0 = "VALUE_0",
  VALUE_1 = "VALUE_1",
}

export type WellknownTypes = {
  timestamp?: string;
  duration?: string;
  enumValue?: AEnum;
  mapValue?: {[key: string]: string};
  struct?: unknown;
  listValue?: unknown[];
  nullValue?: null;
  fieldMask?: string[];
  any?: unknown;
  empty?: {};
  status?: GoogleRpcStatus.Status;
}