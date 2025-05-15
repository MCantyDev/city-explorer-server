package models

type UserCreationError struct {
	Error        string                        `json:"error"`
	Requirements UserCreationErrorRequirements `json:"requirements"`
}

type UserCreationErrorRequirements struct {
	MinFirstNameLength      int  `json:"min_firstname_length"`
	MaxFirstNameLength      int  `json:"max_firstname_length"`
	MinLastNameLength       int  `json:"min_lastname_length"`
	MaxLastNameLength       int  `json:"max_lastname_length"`
	MinPasswordLength       int  `json:"min_password_length"`
	MaxPasswordLength       int  `json:"max_password_length"`
	PasswordRequiresSymbols bool `json:"password_requires_symbols"`
	PasswordRequiresNumbers bool `json:"password_requires_numbers"`
}

type UserLoginError struct {
	Error string `json:"error"`
}

// Needs to have a default one for when encountering errors (As will only display 1 error at a time regardless)
func NewUserCreationError() UserCreationError {
	return UserCreationError{
		Error:        "",
		Requirements: NewUserCreationErrorRequirements(),
	}
}

// MIGHT BE CHANGED TO USE ENV Variables...but not rn
func NewUserCreationErrorRequirements() UserCreationErrorRequirements {
	return UserCreationErrorRequirements{
		MinFirstNameLength:      2,
		MaxFirstNameLength:      50,
		MinLastNameLength:       2,
		MaxLastNameLength:       50,
		MinPasswordLength:       8,
		MaxPasswordLength:       20,
		PasswordRequiresSymbols: true,
		PasswordRequiresNumbers: true,
	}
}
