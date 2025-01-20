package utils

type JError struct {
	Error string `json:"error"`
}

/*
Funcao que melhora o erro à ser enviado ao client
*/
func NewJError(err error) JError {
	jerr := JError{"generic error"}
	if err != nil {
		jerr.Error = err.Error()
	}
	return jerr
}
