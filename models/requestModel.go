package models

type SignUpRequest struct {
	Email            *string `json:"email" validate:"required,email"`
	Username         *string `json:"username" validate:"required,min=4,max=32"`
	Password         *string `json:"password" validate:"required,min=6,max=32"`
	Password_confirm *string `json:"password_confirm" validate:"eqfield=Password"`
}

type LoginRequest struct {
	UsernameOrEmail *string `json:"username_or_email" validate:"required"`
	Password        *string `json:"password" validate:"required"`
}

type LogoutRequest struct {
	Token   *string `json:"token" validate:"required"`
	User_id string  `json:"user_id" validate:"required"`
}

type DeleteAccountRequest struct {
	Token   *string `json:"token" validate:"required"`
	User_id string  `json:"user_id" validate:"required"`
}

type GetProfileInfoRequest struct {
	Token   *string `json:"token" validate:"required"`
	User_id string  `json:"user_id" validate:"required"`
}

type UpdateAccountInfoRequest struct {
	Token            *string `json:"token" validate:"required"`
	User_id          string  `json:"user_id" validate:"required"`
	Email            *string `json:"email" validate:"omitempty,email"`
	Username         *string `json:"username" validate:"omitempty,min=4,max=32"`
	Password         *string `json:"password" validate:"omitempty,min=6,max=32"`
	Password_confirm *string `json:"password_confirm" validate:"omitempty,eqfield=Password"`
}

type CreateProfileInfoRequest struct {
	Token       *string      `json:"token" validate:"required"`
	User_id     string       `json:"user_id" validate:"required"`
	ProfileInfo *ProfileInfo `json:"profile_info" validate:"required"`
}

type EditProfileInfoRequest struct {
	Token       *string      `json:"token" validate:"required"`
	User_id     string       `json:"user_id" validate:"required"`
	ProfileInfo *ProfileInfo `json:"profile_info"`
}

type CreateBunkieRequest struct {
	Token           *string `json:"token" validate:"required"`
	User_id         string  `json:"user_id" validate:"required"`
	City            *string `json:"city" validate:"required"`
	District        *string `json:"district" validate:"required"`
	Quarter         *string `json:"quarter" validate:"required"`
	Header          *string `json:"header" validate:"required"`
	Description     *string `json:"description" validate:"required"`
	NumberOfRooms   *string `json:"number_of_rooms" validate:"required"`
	Price           float64 `json:"price" validate:"required"`
	GenderPreferred *string `json:"gender_preferred" validate:"required"`
}

type GetBunkieAdRequest struct {
	Token   *string `json:"token" validate:"required"`
	Ad_id   *string `json:"ad_id" validate:"required"`
	User_id string  `json:"user_id" validate:"required"`
}

type UpdateBunkieRequest struct {
	Token           *string `json:"token" validate:"required"`
	Ad_id           *string `json:"ad_id" validate:"required"`
	User_id         string  `json:"user_id" validate:"required"`
	City            *string `json:"city"`
	District        *string `json:"district"`
	Quarter         *string `json:"quarter"`
	Header          *string `json:"header"`
	Description     *string `json:"description"`
	NumberOfRooms   *string `json:"number_of_rooms"`
	Price           float64 `json:"price"`
	GenderPreferred *string `json:"gender_preferred"`
}

type DeleteBunkieRequest struct {
	Token   *string `json:"token" validate:"required"`
	Ad_id   *string `json:"ad_id" validate:"required"`
	User_id string  `json:"user_id" validate:"required"`
}

type CreateRoomAdRequest struct {
	Token            *string `json:"token" validate:"required"`
	User_id          string  `json:"user_id" validate:"required"`
	Header_bytearray *string `json:"header_bytearray" validate:"required"`
	Other_bytearrays *string `json:"other_bytearrays"`
	Header           *string `json:"header" validate:"required"`
	Description      *string `json:"description" validate:"required"`
	City             *string `json:"city" validate:"required"`
	District         *string `json:"district" validate:"required"`
	Quarter          *string `json:"quarter" validate:"required"`
	Price            float64 `json:"price" validate:"required"`
	GenderPreferred  *string `json:"gender_preferred" validate:"required"`
	NumberOfRooms    *string `json:"number_of_rooms" validate:"required"`
}

type GetRoomAdRequest struct {
	Token   *string `json:"token" validate:"required"`
	Ad_id   *string `json:"ad_id" validate:"required"`
	User_id string  `json:"user_id" validate:"required"`
}

type UpdateRoomAdRequest struct {
	Token            *string `json:"token" validate:"required"`
	Ad_id            *string `json:"ad_id" validate:"required"`
	User_id          string  `json:"user_id" validate:"required"`
	Header_bytearray *string `json:"header_bytearray"`
	Other_bytearrays *string `json:"other_bytearrays"`
	Header           *string `json:"header"`
	Description      *string `json:"description"`
	City             *string `json:"city"`
	District         *string `json:"district"`
	Quarter          *string `json:"quarter"`
	Price            float64 `json:"price"`
	NumberOfRooms    *string `json:"number_of_rooms"`
	GenderPreferred  *string `json:"gender_preferred"`
}

type DeleteRoomAdRequest struct {
	Token   *string `json:"token" validate:"required"`
	Ad_id   *string `json:"ad_id" validate:"required"`
	User_id string  `json:"user_id" validate:"required"`
}