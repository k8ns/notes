package auth

type ResultCode string

const (
	ResultSuccess ResultCode = "success"
	ResultNoIdentity ResultCode = "no_identity"
	ResultWrongCredentials ResultCode = "wrong_credentials"
	ResultUncategorized ResultCode = "uncategorized"
)

type Result struct {
	Code ResultCode
	Message string
	Identity *User
}

func (r *Result) IsValid() bool {
	return r.Code == ResultSuccess
}