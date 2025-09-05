# protoc-gen-connectrpc-permify

[connectrpc-permify](https://github.com/nrf110/connectrpc-permify) is a Go library providing a unary interceptor for [ConnectRPC](https://connectrpc.com) servers built. The interceptor provides authentication (currently via OAuth2) and Zanzibar-style fine-grained authorization support (via [Permify](https://permify.co)). It requires that all request messages conform to the [Checkable](https://github.com/nrf110/connectrpc-permify/blob/main/pkg/check.go#L81) interface - they must have a `GetChecks()` method implemented which returns the parameters needed for a call to Permify, to see if the current principal is allowed to make the RPC call.

This protobuf compiler plugin is a companion to connectrpc-permify. Using the protobuf custom options defined in connectrpc-permify, we can annotate our service methods and request messages to tell it which type of resource/entity we're acting on, the permission required, the specific id of the entity, and which tenant the entity belongs to (in a multi-tenant system). The plugin will process these annotations and generate our `GetChecks()` methods to satisfy the `Checkable` interface.

## Example

```protobuf
import "nrf110/permify/v1/permify.proto";

message GetUserRequest {
  string user_id = 1 [(nrf110.permify.v1.resource_id) = true];
  string organization_id = 2 [(nrf110.permify.v1.tenant_id) = true]
}

message User {
  option (nrf110.permify.v1.resource_type) = "user";
  string id = 1;
  string email = 2;
}

service UserService {
  rpc GetUser(GetUserRequest) returns (User) {
    option (nrf110.permify.v1.action) = "read";
  }

  rpc PublicHealthCheck(HealthCheckRequest) returns (HealthCheckResponse) {
    option (nrf110.permify.v1.public) = true;  // No auth required
  }
}
```

## Local development

### Dependencies

A devcontainer is provided for ready-to-go development environment. If you want to develop without it, just make sure you have installed:

- The currently used version of go
- The Buf CLI
- protoc-gen-go
- protoc-gen-connect-go

## Makefile

There are 2 Makefiles in this repo:

- The root level, which is responsible for building the plugin. Run `make help` for make target descriptions.
- Under `testdata`, another Makefile is responsible for running the protobuf compiler with your plugin build to generate output

### Testing

Testing protobuf compiler plugins is unfortunately tricky, as a lot of work would be required to mock out all of the AST nodes provided representing non-trivial protobuf files. Additionally, the compiler plugin must be built and tested from a shell command,

Given that, there are very few true unit tests. Instead, we validate a "golden", manually validated set of output files against freshly generated output files. New features or behavior changes should include new/updated .proto files under `testdata/input`. After validating the new behavior, copy the output files to `tesdata/golden` and commit them.
