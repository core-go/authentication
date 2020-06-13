package auth

type AuthStatus int

const (
	StatusSuccess               = AuthStatus(0)
	StatusSuccessAndReactivated = AuthStatus(1)
	StatusTwoFactorRequired     = AuthStatus(2)
	StatusFail                  = AuthStatus(3)
	StatusWrongPassword         = AuthStatus(4)
	StatusPasswordExpired       = AuthStatus(5)
	StatusAccessTimeLocked      = AuthStatus(6)
	StatusLocked                = AuthStatus(7)
	StatusSuspended             = AuthStatus(8)
	StatusDisabled              = AuthStatus(9)
	StatusSystemError           = AuthStatus(10)
)
