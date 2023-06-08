package request

type UserRequest struct {
	Firstname string `json:"firstname" validate:"required,alpha"`
	Lastname  string `json:"lastname" validate:"required,alpha"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,numeric"`
	Address   string `json:"address" validate:"required"`
	Username  string `json:"username" validate:"required,alphanum,gte=3"`
	Password  string `json:"password" validate:"required,alphanum,gte=8"`
}

type UserRequestService struct {
	Id        int64  `json:"id" validate:"required"`
	Firstname string `json:"firstname" validate:"required,alpha"`
	Lastname  string `json:"lastname" validate:"required,alpha"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,numeric"`
	Address   string `json:"address" validate:"required"`
	CreatedAt int64  `json:"created_at" validate:"-"`
	Username  string `json:"username" validate:"required,alphanum,gte=3"`
	Password  string `json:"password" validate:"required,alphanum,gte=8"`
}

type UserUpdateRequest struct {
	Firstname string `json:"firstname" validate:"-"`
	Lastname  string `json:"lastname" validate:"-"`
	Email     string `json:"email" validate:"email"`
	Phone     string `json:"phone" validate:"-"`
	Address   string `json:"address" validate:"-"`
	Username  string `json:"username" validate:"-"`
	Password  string `json:"password" validate:"-"`
}

type UserUpdateService struct {
	Id        int64  `json:"id" validate:"-"`
	Firstname string `json:"firstname" validate:"-"`
	Lastname  string `json:"lastname" validate:"-"`
	Email     string `json:"email" validate:"email"`
	Phone     string `json:"phone" validate:"-"`
	Address   string `json:"address" validate:"-"`
	Username  string `json:"username" validate:"gte=3"`
	Password  string `json:"password" validate:"gte=8"`
}

type UserForgotPasswordRequest struct {
	Username string `json:"username_or_email" validate:"required,gte=3"`
	Password string `json:"password" validate:"required,gte=8"`
}

type UserLoginRequest struct {
	Username string `json:"username_or_email" validate:"required,gte=3"`
	Password string `json:"password" validate:"required,gte=8"`
}
