package auth

type AuthStatus int

const (
	Success               = AuthStatus(0)
	SuccessAndReactivated = AuthStatus(1)
	Fail                  = AuthStatus(2)
	WrongPassword         = AuthStatus(3)
	PasswordExpired       = AuthStatus(4)
	AccessTimeLocked      = AuthStatus(5)
	Locked                = AuthStatus(6)
	Suspended             = AuthStatus(7)
	Disabled              = AuthStatus(8)
	SystemError           = AuthStatus(9)
)
