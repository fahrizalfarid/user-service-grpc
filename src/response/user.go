package response

type Error struct {
	Message string `json:"error_message"`
}

type UserProfile struct {
	Id        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Username  string `json:"username"`
}

type UserProfileResponse struct {
	Id        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Username  string `json:"username"`
}

type UserFound struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

type UserLoginResponse struct {
	Id       int64  `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
