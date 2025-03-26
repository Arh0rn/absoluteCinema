package models

// User TODO: Relocate all tags to dao for every model
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpInput struct {
	Username string `json:"username" validate:"required,gte=3,lte=20" example:"user123"`
	Email    string `json:"email" validate:"required,email" example:"words@mail.com"`
	Password string `json:"password" validate:"required,gte=8,lte=16" example:"mypassword123"`
}

func (i SignUpInput) Validate() error {
	return validate.Struct(i)
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email" example:"words@mail.com"`
	Password string `json:"password" validate:"required,gte=8,lte=16" example:"mypassword123"`
}

func (i SignInInput) Validate() error {
	return validate.Struct(i)
}
