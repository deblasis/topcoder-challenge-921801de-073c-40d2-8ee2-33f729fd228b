package transport

// func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.UserManager {
// 	// We construct a single ratelimiter middleware, to limit the total outgoing
// 	// QPS from this client to all methods on the remote instance. We also
// 	// construct per-endpoint circuitbreaker middlewares to demonstrate how
// 	// that's done, although they could easily be combined into a single breaker
// 	// for the entire remote instance, too.
// 	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

// 	// global client middlewares
// 	var options []grpctransport.ClientOption

// 	if zipkinTracer != nil {
// 		// Zipkin GRPC Client Trace can either be instantiated per gRPC method with a
// 		// provided operation name or a global tracing client can be instantiated
// 		// without an operation name and fed to each Go kit client as ClientOption.
// 		// In the latter case, the operation name will be the endpoint's grpc method
// 		// path.
// 		//
// 		// In this example, we demonstrace a global tracing client.
// 		options = append(options, zipkin.GRPCClientTrace(zipkinTracer))

// 	}
// 	var statusEndpoint endpoint.Endpoint
// 	{
// 		statusEndpoint = grpctransport.NewClient(
// 			conn,
// 			service.ServiceName,
// 			"ServiceStatus",
// 			encodeGRPCServiceStatusRequest,
// 			decodeGRPCServiceStatusReply,
// 			pb.ServiceStatusReply{},
// 			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
// 		).Endpoint()

// 		statusEndpoint = opentracing.TraceClient(otTracer, "ServiceStatus")(statusEndpoint)
// 		statusEndpoint = limiter(statusEndpoint)
// 		statusEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
// 			Name:    "ServiceStatus",
// 			Timeout: 30 * time.Second,
// 		}))(statusEndpoint)
// 	}

// 	var getUserByUsernameEndpoint endpoint.Endpoint
// 	{
// 		getUserByUsernameEndpoint = grpctransport.NewClient(
// 			conn,
// 			service.ServiceName,
// 			"GetUserByUsername",
// 			encodeGRPCGetUserByUsernameRequest,
// 			decodeGRPCGetUserByUsernameReply,
// 			pb.UserReply{},
// 			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
// 		).Endpoint()

// 		getUserByUsernameEndpoint = opentracing.TraceClient(otTracer, "GetUserByUsername")(getUserByUsernameEndpoint)
// 		getUserByUsernameEndpoint = limiter(getUserByUsernameEndpoint)
// 		getUserByUsernameEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
// 			Name:    "GetUserByUsername",
// 			Timeout: 30 * time.Second,
// 		}))(getUserByUsernameEndpoint)
// 	}
// 	var createUserEndpoint endpoint.Endpoint
// 	{
// 		createUserEndpoint = grpctransport.NewClient(
// 			conn,
// 			service.ServiceName,
// 			"CreateUser",
// 			encodeGRPCCreateUserRequest,
// 			decodeGRPCCreateUserReply,
// 			pb.UserReply{},
// 			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
// 		).Endpoint()

// 		createUserEndpoint = opentracing.TraceClient(otTracer, "CreateUser")(createUserEndpoint)
// 		createUserEndpoint = limiter(createUserEndpoint)
// 		createUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
// 			Name:    "CreateUser",
// 			Timeout: 30 * time.Second,
// 		}))(createUserEndpoint)
// 	}

// 	return endpoints.EndpointSet{
// 		StatusEndpoint:            statusEndpoint,
// 		GetUserByUsernameEndpoint: getUserByUsernameEndpoint,
// 		CreateUserEndpoint:        createUserEndpoint,
// 	}

// }

// func encodeGRPCServiceStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
// 	return &pb.ServiceStatusRequest{}, nil
// }

// func decodeGRPCServiceStatusReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
// 	reply := grpcReply.(*pb.ServiceStatusReply)
// 	return model.ServiceStatusReply{Code: reply.Code, Err: reply.Err}, nil
// }

// func encodeGRPCGetUserByUsernameRequest(_ context.Context, request interface{}) (interface{}, error) {
// 	req := request.(model.GetUserByUsernameRequest)
// 	return &pb.GetUserByUsernameRequest{
// 		Username: req.Username,
// 	}, nil
// }

// func decodeGRPCGetUserByUsernameReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
// 	reply := grpcReply.(*pb.UserReply)
// 	return model.GetUserByUsernameReply{User: model.User{
// 		ID:       reply.User.Id,
// 		Username: reply.User.Username,
// 		Password: reply.User.Password,
// 		Role:     strings.Title(strings.ToLower(strings.TrimLeft(reply.User.Role.String(),"ROLE_")),
// 	}, nil
// }

// func encodeGRPCCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
// 	req := request.(model.CreateUserRequest)

// 	//TODO: centralise
// 	roleId := pb.User_Role_value[strings.ToUpper("ROLE_" + req.Role)]
// 	if roleId <= 0 {
// 		return nil, errors.New("cannot unmarshal role")
// 	}
// 	return &pb.CreateUserRequest{
// 		User: &pb.User{
// 			Id:       req.ID,
// 			Username: req.Username,
// 			Password: req.Password,
// 			Role:     pb.User_Role(roleId),
// 		},
// 	}, nil
// }
// func decodeGRPCCreateUserReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
// 	reply := grpcReply.(*pb.UserReply)
// 	return model.CreateUserReply{
// 		Id:  reply.User.Id,
// 		Err: "",
// 	}, nil
// }
