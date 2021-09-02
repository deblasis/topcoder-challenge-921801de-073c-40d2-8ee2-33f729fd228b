package model

type ServiceStatusRequest struct{}
type ServiceStatusReply struct {
	Code int64  `json:"code"`
	Err  string `json:"err,omitempty"`
}

type GetUserByUsernameRequest struct {
	Username string `json:"username"`
}
type GetUserByUsernameReply struct {
	User User   `json:"user"`
	Err  string `json:"err,omitempty"`
}

type CreateUserRequest struct {
	User User `json:"user"`
}
type CreateUserReply struct {
	Id  int64  `json:"id"`
	Err string `json:"err,omitempty"`
}
