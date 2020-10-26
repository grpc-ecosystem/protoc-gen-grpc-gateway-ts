#protoc-gen-grpc-gateway-ts

`protoc-gen-grpc-gateway-ts` is a Typescript client generator for the [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway/) project. It generates idiomatic Typescript clients that connect the web frontend and golang backend fronted by grpc-gateway. 


Features:
1. idiomatic Typescript clients and messages
2. Supports both One way and server side streaming gRPC calls
3. POJO request construction guarded by message type definitions, which is way easier compare to `grpc-web`
4. No need to use swagger/open api to generate client code for the web.

License
protoc-gen-grpc-gateway-ts is licensed under the MIT License. See LICENSE.txt for more details.