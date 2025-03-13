package request

type ReqAuthSignIn struct {
	Email     string `json:"email" validate:"required,email"`    // Email must be unique
	Password  string `json:"password" validate:"required,min=7"` // Password minimal 7 characters
	IpV4      string `json:"-"`                                  // Ip Address v4
	IpV6      string `json:"-"`                                  // Ip Address v6
	UserAgent string `json:"-"`                                  // User Agent
	Device    string `json:"-"`                                  // Device
	OS        string `json:"-"`                                  // Operating System
} // @name ReqAuthSignIn

type ReqAuthSignUp struct {
	Email    string `json:"email" validate:"required,email"`     // Email is required, and must be a valid email address
	FullName string `json:"full_name" validate:"required,min=5"` // Full name is required and minimal 5 characters
	Password string `json:"password" validate:"required,min=7"`  // Password is required and minimal 7 characters
	Role     string `json:"role" validate:"required"`            // Role is required, must be admin_director/admin/admin_finance/admin_it
} // @name ReqAuthSignUp
