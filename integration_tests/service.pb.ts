/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import { Observable } from 'rxjs';
import * as fm from "./fetch.pb"
import * as Msg from "./msg.pb"

export type UnaryRequest = {
  counter?: number;
}

export type UnaryResponse = {
  result?: number;
}

export type BinaryRequest = {
  data?: Uint8Array;
}

export type BinaryResponse = {
  data?: Uint8Array;
}

export type StreamingRequest = {
  counter?: number;
}

export type StreamingResponse = {
  result?: number;
}

export type HttpGetRequest = {
  numToIncrease?: number;
}

export type HttpGetResponse = {
  result?: number;
}

export type HttpPostRequest = {
  a?: number;
  req?: PostRequest;
  c?: number;
}

export type PostRequest = {
  b?: number;
}

export type HttpPostRequest2 = {
  a?: number;
  reqCamel?: PostRequest;
  c?: number;
}

export type HttpPostResponse = {
  postResult?: number;
}

export type HttpPatchRequest = {
  a?: number;
  c?: number;
}

export type HttpPatchResponse = {
  patchResult?: number;
}

export type HttpDeleteRequest = {
  a?: number;
}

export type HTTPGetWithURLSearchParamsRequest = {
  a?: number;
  postReq?: PostRequest;
  c?: number[];
  extMsg?: Msg.ExternalMessage;
}

export type HTTPGetWithURLSearchParamsResponse = {
  urlSearchParamsResult?: number;
}

export type ZeroValueMsg = {
  c?: number;
  d?: number[];
  e?: boolean;
}

export type HTTPGetWithZeroValueURLSearchParamsRequest = {
  a?: string;
  b?: string;
  zeroValueMsg?: ZeroValueMsg;
}

export type HTTPGetWithZeroValueURLSearchParamsResponse = {
  a?: string;
  b?: string;
  zeroValueMsg?: ZeroValueMsg;
}

export type HttpGetRequest2 = {
  name?: string;
}

export type HttpGetRequest3 = {
  rCamel?: HttpGetRequest2;
}

export class CounterService {
  static Increment(req: UnaryRequest, initReq?: fm.InitReq): Promise<UnaryResponse> {
    return fm.fetchReq<UnaryRequest, UnaryResponse>(` + "`/main.CounterService/Increment`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static StreamingIncrements(req: StreamingRequest, entityNotifier?: fm.NotifyStreamEntityArrival<StreamingResponse>, initReq?: fm.InitReq): Promise<void> {
    return fm.fetchStreamingRequest<StreamingRequest, StreamingResponse>(` + "`/main.CounterService/StreamingIncrements`" + `, entityNotifier, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static FailingIncrement(req: UnaryRequest, initReq?: fm.InitReq): Promise<UnaryResponse> {
    return fm.fetchReq<UnaryRequest, UnaryResponse>(` + "`/main.CounterService/FailingIncrement`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static EchoBinary(req: BinaryRequest, initReq?: fm.InitReq): Promise<BinaryResponse> {
    return fm.fetchReq<BinaryRequest, BinaryResponse>(` + "`/main.CounterService/EchoBinary`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static HTTPGet(req: HttpGetRequest, initReq?: fm.InitReq): Promise<HttpGetResponse> {
    return fm.fetchReq<HttpGetRequest, HttpGetResponse>(` + "`/api/${req["numToIncrease"]}?${fm.renderURLSearchParams(req, ["numToIncrease"])}`" + `, {...initReq, method: "GET"});
  }
  static HTTPGet2(req: HttpGetRequest2, initReq?: fm.InitReq): Promise<HttpGetResponse> {
    return fm.fetchReq<HttpGetRequest2, HttpGetResponse>(` + "`/api/${req["name"]}:hello?${fm.renderURLSearchParams(req, ["name"])}`" + `, {...initReq, method: "GET"});
  }
  static HTTPGet3(req: HttpGetRequest3, initReq?: fm.InitReq): Promise<HttpGetResponse> {
    return fm.fetchReq<HttpGetRequest3, HttpGetResponse>(` + "`/api/${req["rCamel"]["nameCamel"]}:hello?${fm.renderURLSearchParams(req, ["rCamel.nameCamel"])}`" + `, {...initReq, method: "GET"});
  }
  static HTTPPostWithNestedBodyPath(req: HttpPostRequest, initReq?: fm.InitReq): Promise<HttpPostResponse> {
    return fm.fetchReq<HttpPostRequest, HttpPostResponse>(` + "`/post/${req["a"]}`" + `, {...initReq, method: "POST", body: JSON.stringify(req["req"], fm.replacer)});
  }
  static HTTPPostWithStarBodyPath(req: HttpPostRequest, initReq?: fm.InitReq): Promise<HttpPostResponse> {
    return fm.fetchReq<HttpPostRequest, HttpPostResponse>(` + "`/post/${req["a"]}/${req["c"]}`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static HttpPost2(req: HttpPostRequest2, initReq?: fm.InitReq): Promise<HttpPostResponse> {
    return fm.fetchReq<HttpPostRequest2, HttpPostResponse>(` + "`/post/${req["a"]}/${req["c"]}`" + `, {...initReq, method: "POST", body: JSON.stringify(req["reqCamel"], fm.replacer)});
  }
  static HttpPost2Nested(req: HttpPostRequest2, initReq?: fm.InitReq): Promise<HttpPostResponse> {
    return fm.fetchReq<HttpPostRequest2, HttpPostResponse>(` + "`/post/${req["a"]}/${req["c"]}`" + `, {...initReq, method: "POST", body: JSON.stringify(req["reqCamel"]["b"], fm.replacer)});
  }
  static HTTPPatch(req: HttpPatchRequest, initReq?: fm.InitReq): Promise<HttpPatchResponse> {
    return fm.fetchReq<HttpPatchRequest, HttpPatchResponse>(` + "`/patch`" + `, {...initReq, method: "PATCH", body: JSON.stringify(req, fm.replacer)});
  }
  static HTTPDelete(req: HttpDeleteRequest, initReq?: fm.InitReq): Promise<{}> {
    return fm.fetchReq<HttpDeleteRequest, {}>(` + "`/delete/${req["a"]}`" + `, {...initReq, method: "DELETE"});
  }
  static ExternalMessage(req: Msg.ExternalRequest, initReq?: fm.InitReq): Promise<Msg.ExternalResponse> {
    return fm.fetchReq<Msg.ExternalRequest, Msg.ExternalResponse>(` + "`/main.CounterService/ExternalMessage`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static HTTPGetWithURLSearchParams(req: HTTPGetWithURLSearchParamsRequest, initReq?: fm.InitReq): Promise<HTTPGetWithURLSearchParamsResponse> {
    return fm.fetchReq<HTTPGetWithURLSearchParamsRequest, HTTPGetWithURLSearchParamsResponse>(` + "`/api/query/${req["a"]}?${fm.renderURLSearchParams(req, ["a"])}`" + `, {...initReq, method: "GET"});
  }
  static HTTPGetWithZeroValueURLSearchParams(req: HTTPGetWithZeroValueURLSearchParamsRequest, initReq?: fm.InitReq): Promise<HTTPGetWithZeroValueURLSearchParamsResponse> {
    return fm.fetchReq<HTTPGetWithZeroValueURLSearchParamsRequest, HTTPGetWithZeroValueURLSearchParamsResponse>(` + "`/path/query?${fm.renderURLSearchParams(req, [])}`" + `, {...initReq, method: "GET"});
  }
}

export class ObservableCounterService {
  static Increment(req: UnaryRequest, initReq?: fm.InitReq): Observable<UnaryResponse> {
    return fm.fromFetchReq<UnaryRequest, UnaryResponse>(` + "`/main.CounterService/Increment`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static StreamingIncrements(req: StreamingRequest, initReq?: fm.InitReq): Observable<StreamingResponse> {
    return fm.fromFetchStreamingRequest<StreamingRequest, StreamingResponse>(` + "`/main.CounterService/StreamingIncrements`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static FailingIncrement(req: UnaryRequest, initReq?: fm.InitReq): Observable<UnaryResponse> {
    return fm.fromFetchReq<UnaryRequest, UnaryResponse>(` + "`/main.CounterService/FailingIncrement`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static EchoBinary(req: BinaryRequest, initReq?: fm.InitReq): Observable<BinaryResponse> {
    return fm.fromFetchReq<BinaryRequest, BinaryResponse>(` + "`/main.CounterService/EchoBinary`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static HTTPGet(req: HttpGetRequest, initReq?: fm.InitReq): Observable<HttpGetResponse> {
    return fm.fromFetchReq<HttpGetRequest, HttpGetResponse>(` + "`/api/${req["numToIncrease"]}?${fm.renderURLSearchParams(req, ["numToIncrease"])}`" + `, {...initReq, method: "GET"});
  }
  static HTTPGet2(req: HttpGetRequest2, initReq?: fm.InitReq): Observable<HttpGetResponse> {
    return fm.fromFetchReq<HttpGetRequest2, HttpGetResponse>(` + "`/api/${req["name"]}:hello?${fm.renderURLSearchParams(req, ["name"])}`" + `, {...initReq, method: "GET"});
  }
  static HTTPGet3(req: HttpGetRequest3, initReq?: fm.InitReq): Observable<HttpGetResponse> {
    return fm.fromFetchReq<HttpGetRequest3, HttpGetResponse>(` + "`/api/${req["rCamel"]["nameCamel"]}:hello?${fm.renderURLSearchParams(req, ["rCamel.nameCamel"])}`" + `, {...initReq, method: "GET"});
  }
  static HTTPPostWithNestedBodyPath(req: HttpPostRequest, initReq?: fm.InitReq): Observable<HttpPostResponse> {
    return fm.fromFetchReq<HttpPostRequest, HttpPostResponse>(` + "`/post/${req["a"]}`" + `, {...initReq, method: "POST", body: JSON.stringify(req["req"], fm.replacer)});
  }
  static HTTPPostWithStarBodyPath(req: HttpPostRequest, initReq?: fm.InitReq): Observable<HttpPostResponse> {
    return fm.fromFetchReq<HttpPostRequest, HttpPostResponse>(` + "`/post/${req["a"]}/${req["c"]}`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static HttpPost2(req: HttpPostRequest2, initReq?: fm.InitReq): Observable<HttpPostResponse> {
    return fm.fromFetchReq<HttpPostRequest2, HttpPostResponse>(` + "`/post/${req["a"]}/${req["c"]}`" + `, {...initReq, method: "POST", body: JSON.stringify(req["reqCamel"], fm.replacer)});
  }
  static HttpPost2Nested(req: HttpPostRequest2, initReq?: fm.InitReq): Observable<HttpPostResponse> {
    return fm.fromFetchReq<HttpPostRequest2, HttpPostResponse>(` + "`/post/${req["a"]}/${req["c"]}`" + `, {...initReq, method: "POST", body: JSON.stringify(req["reqCamel"]["b"], fm.replacer)});
  }
  static HTTPPatch(req: HttpPatchRequest, initReq?: fm.InitReq): Observable<HttpPatchResponse> {
    return fm.fromFetchReq<HttpPatchRequest, HttpPatchResponse>(` + "`/patch`" + `, {...initReq, method: "PATCH", body: JSON.stringify(req, fm.replacer)});
  }
  static HTTPDelete(req: HttpDeleteRequest, initReq?: fm.InitReq): Observable<{}> {
    return fm.fromFetchReq<HttpDeleteRequest, {}>(` + "`/delete/${req["a"]}`" + `, {...initReq, method: "DELETE"});
  }
  static ExternalMessage(req: Msg.ExternalRequest, initReq?: fm.InitReq): Observable<Msg.ExternalResponse> {
    return fm.fromFetchReq<Msg.ExternalRequest, Msg.ExternalResponse>(` + "`/main.CounterService/ExternalMessage`" + `, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static HTTPGetWithURLSearchParams(req: HTTPGetWithURLSearchParamsRequest, initReq?: fm.InitReq): Observable<HTTPGetWithURLSearchParamsResponse> {
    return fm.fromFetchReq<HTTPGetWithURLSearchParamsRequest, HTTPGetWithURLSearchParamsResponse>(` + "`/api/query/${req["a"]}?${fm.renderURLSearchParams(req, ["a"])}`" + `, {...initReq, method: "GET"});
  }
  static HTTPGetWithZeroValueURLSearchParams(req: HTTPGetWithZeroValueURLSearchParamsRequest, initReq?: fm.InitReq): Observable<HTTPGetWithZeroValueURLSearchParamsResponse> {
    return fm.fromFetchReq<HTTPGetWithZeroValueURLSearchParamsRequest, HTTPGetWithZeroValueURLSearchParamsResponse>(` + "`/path/query?${fm.renderURLSearchParams(req, [])}`" + `, {...initReq, method: "GET"});
  }
}