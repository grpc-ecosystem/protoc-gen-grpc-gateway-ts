/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as ComSquareupCashGapDatasourceDatasource from "./datasource/datasource"
import * as ComSquareupCashGapEnvironment from "./environment"
import * as gap from "gap/admin/lib/useGapFetch"

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };
type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (
    keyof T extends infer K ?
      (K extends string & keyof T ? { [k in K]: T[K] } & Absent<T, K>
        : never)
    : never);

export enum LogEntryLevel {
  DEBUG = "DEBUG",
  INFO = "INFO",
  WARN = "WARN",
  ERROR = "ERROR",
}

export type LogEntryStackTraceException = {
  type?: string
  message?: string
}

export type LogEntryStackTraceMethod = {
  identifier?: string
  file?: string
  line?: string
}

export type LogEntryStackTrace = {
  exception?: LogEntryStackTraceException
  lines?: LogEntryStackTraceMethod[]
}


type BaseLogEntry = {
  hostname?: string
  level?: LogEntryLevel
  elapsed?: number
  timestamp?: number
  env?: ComSquareupCashGapEnvironment.Environment
  hasStackTrace?: boolean
  message?: string
  tags?: string[]
  stackTraces?: LogEntryStackTrace[]
}

export type LogEntry = BaseLogEntry
  & OneOf<{ application: string; service: string }>


type BaseLogStream = {
}

export type LogStream = BaseLogStream
  & OneOf<{ dataCentre: DataCentreLogEntries; cloud: CloudLogEntries }>

export type DataCentreLogEntries = {
  logs?: LogEntry[]
}

export type CloudLogEntries = {
  logs?: LogEntry[]
}


type BaseFetchLogRequest = {
  source?: ComSquareupCashGapDatasourceDatasource.DataSource
}

export type FetchLogRequest = BaseFetchLogRequest
  & OneOf<{ application: string; service: string }>

export type FetchLogResponse = {
  result?: LogStream
}

export type PushLogRequest = {
  entry?: LogEntry
  source?: ComSquareupCashGapDatasourceDatasource.DataSource
}

export type PushLogResponse = {
  success?: boolean
}

export class LogService {
  static FetchLog(req: FetchLogRequest): Promise<gap.FetchState<FetchLogResponse>> {
    return gap.gapFetchGRPC<FetchLogRequest, FetchLogResponse>("/api/com.squareup.cash.gap.LogService/FetchLog", req)
  }
  static StreamLog(req: FetchLogRequest, entityNotifier?: gap.NotifyStreamEntityArrival<FetchLogResponse>): Promise<gap.FetchState<undefined>> {
    return gap.gapFetchGRPCStream<FetchLogRequest, FetchLogResponse>("/api/com.squareup.cash.gap.LogService/StreamLog", req, entityNotifier)
  }
  static PushLog(req: PushLogRequest): Promise<gap.FetchState<PushLogResponse>> {
    return gap.gapFetchGRPC<PushLogRequest, PushLogResponse>("/api/com.squareup.cash.gap.LogService/PushLog", req)
  }
}