package sign

type SignError struct {
	Code    int    `json:"error.code"`
	Message string `json:"error.message"`
	Status  string `json:"error.status"`
}

func (e SignError) Error() string {
	return e.Message
}
