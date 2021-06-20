package db

type User struct {
	Name     string  `json:"name"`
	RollNo      int  `json:"rollno"`
	Password string  `json:"password"`
	Batch       int  `json:"batch"`
}

type AuthUser struct {
	RollNo      int  `json:"rollno"`
	Password string  `json:"password"`
}

type EntryParams struct {
	RollNo    int  `json:"id"`
	Amount    int  `json:"amount"`
}
  
type TransferParams struct {
	Receiver   int  `json:"receiver"`
	Sender     int  `json:"sender"`
	Amount     int  `json:"amount"`
	Remarks string  `json:"remarks"`
}

type WalletParams struct {
	RollNo   int 
	Balance  int
}

type RollNo struct {
	RollNo   int  `json:"rollno"`
}