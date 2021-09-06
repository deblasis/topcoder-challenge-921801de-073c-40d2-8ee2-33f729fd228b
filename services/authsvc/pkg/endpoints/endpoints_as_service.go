package endpoints

// ServiceStatus(ctx context.Context) (int64, error)
// func (s EndpointSet) ServiceStatus(ctx context.Context) (int64, error) {
// 	resp, err := s.StatusEndpoint(ctx, model.ServiceStatusRequest{})
// 	if err != nil {
// 		return 0, err
// 	}
// 	response := resp.(model.ServiceStatusReply)
// 	return response.Code, errors.Str2err(response.Err)
// }

// // GetUserByUsername(ctx context.Context, username string) (model.User, error)
// func (s EndpointSet) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
// 	var user model.User
// 	resp, err := s.GetUserByUsernameEndpoint(ctx, model.GetUserByUsernameRequest{Username: username})
// 	if err != nil {
// 		return user, err
// 	}
// 	response := resp.(model.GetUserByUsernameReply)
// 	return response.User, errors.Str2err(response.Err)
// }

// // CreateUser(ctx context.Context, user *model.User) (int64, error)
// func (s EndpointSet) CreateUser(ctx context.Context, user *model.User) (int64, error) {
// 	resp, err := s.CreateUserEndpoint(ctx, model.CreateUserRequest{
// 		Username: user.Username,
// 		Password: user.Password,
// 		Role:     user.Role,
// 	})
// 	if err != nil {
// 		return -1, err
// 	}
// 	response := resp.(model.CreateUserReply)
// 	return response.Id, errors.Str2err(response.Err)
// }

// var (
// 	_ endpoint.Failer = model.ServiceStatusReply{}
// 	_ endpoint.Failer = model.GetUserByUsernameReply{}
// 	_ endpoint.Failer = model.CreateUserReply{}
// )

// type ServiceStatusRequest struct{}
// type ServiceStatusReply struct {
// 	Code int64 `json:"code"`
// 	Err  error `json:"-"`
// }

// type GetUserByUsernameRequest struct {
// 	Username string `json:"username" validate:"required,notblank"`
// }
// type GetUserByUsernameReply struct {
// 	User model.User `json:"user"`
// 	Err  error      `json:"-"`
// }

// type CreateUserRequest model.User
// type CreateUserReply struct {
// 	Id  int64 `json:"id"`
// 	Err error `json:"-"`
// }

//func (r ServiceStatusReply) Failed() error { return r.Err }
//func (r GetUserByUsernameReply) Failed() error { return r.Err }
//func (r CreateUserReply) Failed() error { return r.Err }
