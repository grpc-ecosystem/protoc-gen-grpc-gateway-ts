package data

// Service is the data representation of Service in proto
type Service struct {
	// Name is the name of the Service
	Name string
	// Methods is a list of methods data
	Methods []*Method
	// Comment is the comment of the Service
	Comment string
}

// Services is an alias of Service array
type Services []*Service

// HasServerStreamingMethod indicates whether there is server side streaming calls inside any of the services
func (s Services) HasServerStreamingMethod() bool {
	for _, service := range s {
		for _, method := range service.Methods {
			if method.ServerStreaming {
				return true
			}
		}
	}
	return false
}

// HasUnaryCallMethod indicates whether there is unary methods inside any of the services
func (s Services) HasUnaryCallMethod() bool {
	for _, service := range s {
		for _, method := range service.Methods {
			if !method.ServerStreaming && !method.ClientStreaming {
				return true
			}
		}
	}
	return false
}

// NeedsFetchModule returns whether the given services needs fetch module support
func (s Services) NeedsFetchModule() bool {
	hasServices := len(s) > 0
	return hasServices && (s.HasUnaryCallMethod() || s.HasServerStreamingMethod())
}

// NewService returns an initialised service
func NewService() *Service {
	return &Service{
		Methods: make([]*Method, 0),
	}
}

// Method represents the rpc calls in protobuf service
type Method struct {
	// Name is the name of the method
	Name string
	// URL is the method url path to invoke from client side
	URL string
	// Input is the input argument
	Input *MethodArgument
	// Output is the output argument
	Output *MethodArgument
	// ServerStreaming indicates the RPC call is a server streaming call
	ServerStreaming bool
	// ClientStreaming indicates the RPC call is a client streaming call, which will not be supported by GRPC Gateway
	ClientStreaming bool
	// HTTPMethod indicates the http method for this function
	HTTPMethod string
	// HTTPBody is the path for request body in the body's payload
	HTTPRequestBody *string
	// Comment is the comment of the Method
	Comment string
}

// MethodArgument stores the type information about method argument
type MethodArgument struct {
	// Type is the type of the argument
	Type string
	// IsExternal indicate if this type is an external dependency
	IsExternal bool
	// IsRepeated indicates whether the field is a repeated field
	IsRepeated bool
}

// GetType returns some information of the type to aid the rendering
func (m *MethodArgument) GetType() *TypeInfo {
	return &TypeInfo{
		Type:       m.Type,
		IsRepeated: m.IsRepeated,
		IsExternal: m.IsExternal,
	}
}

// SetExternal mutates the IsExternal attribute of the type
func (m *MethodArgument) SetExternal(external bool) {
	m.IsExternal = external
}
