package server

type User struct {
	Name     string `json:"name"`
	RollNo   int    `json:"rollno"`
	Password string `json:"password"`
	Batch    int    `json:"batch"`
}

type AuthUser struct {
	RollNo   int    `json:"rollno"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}
