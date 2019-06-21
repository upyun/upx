package upyun

type Error struct {
	error
	statusCode int
}

func (e Error) Error() string {
	return e.error.Error()
}

func IsNotExist(err error) bool {
	if e, ok := err.(Error); ok {
		if e.statusCode == 404 {
			return true
		}
	}
	return false
}
