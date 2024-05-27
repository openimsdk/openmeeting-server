package definition

type User struct {
	UserID   string `column:"user_id"`
	Account  string `column:"account"`
	Nickname string `column:"nickname"`
	Password string `column:"password"`
	Birth    string `column:"birth"`
	Gender   string `column:"gender"`
	Email    string `column:"email"`
}

func (User) SheetName() string {
	return "user"
}
