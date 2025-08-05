package types

type Student struct {
	ID    int64 `json:"id"`
	Name string  `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Age int `json:"age"`
}