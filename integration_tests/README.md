# Integration test

The integration test first runs `./scripts/gen-protos.sh` again to generate Typescript file for the proto `service.proto`.

Then it starts `main.go` server that loads up the protos and run tests via `Karma` to verify if the generated client works properly.

The JS integration test file is `integration_test.ts`.

Changes on the server side needs to run `./scripts/gen-server-proto.sh` to update the protos and the implementation is in `service.go`.

Changes on the test client side is in `integration_test.ts`.

CI test script starts with `test-ci.sh` will make sure the client typescript file to be regenerated before running the test.
