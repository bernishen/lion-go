package result

type R struct {
	Flag    bool         `json:"flag"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data"`
}

func OK(data *interface{}) *R {
	r := R{
		Flag:    false,
		Message: "success",
		Data:    data,
	}
	return &r
}

func Faid(message string) *R {
	r := R{
		Flag:    false,
		Message: message,
	}
	return &r
}
