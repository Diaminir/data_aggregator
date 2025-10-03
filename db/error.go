package db

// ErrorDB структура для передачи код статуса на уровне работы с базой данных
type ErrorDB struct {
	Err  error
	Code int
}

// NewErrorDB функция создающая экземпляр ErrorDB
func NewErrorDB(err error, code int) *ErrorDB {
	return &ErrorDB{
		Err:  err,
		Code: code,
	}
}
