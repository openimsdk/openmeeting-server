package definition

type User struct {
	UserID   string `json:"user_id" column:"user_id"`
	Account  string `json:"account" column:"account"`
	Nickname string `json:"nickname" column:"nickname"`
	Password string `json:"password" column:"password"`
	Birth    string `json:"birth" column:"birth"`
	Gender   string `json:"gender" column:"gender"`
	Email    string `json:"email" column:"email"`
}

func (User) SheetName() string {
	return "user"
}
