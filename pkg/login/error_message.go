package login

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ErrorMessage) Error() string {
	return e.Message
}
