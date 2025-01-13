package constants

const (
	APIVersion = "v1"

	BadRequestMessage = "request not fulfilled"

	// schedular constants
	HealthCheckRoute = "/health"
	MDBUri           = "localhost:27017"
	Database         = "ecommerce"
	Sender           = "puneetvishnoiias@gmail.com"

	// restaurnt
)

const (
	NormalUser = "user"
	AdminUser  = "admin"
)

const (
	// time slot for otp validation
	OtpValidation = 60
)

// collections
const (
	VerificationsCollection = "verifications"
	UserCollection          = "user"
	ProductCollection       = "products"
	AddressCollection       = "user_addresses"
	CartCollection          = "user_cart"
	RestaurantCollection    = "restaurant"
	MenuCollection          = "menu_items"
	OrderCollection         = "order"
)

// messages
const (
	AlreadyRegisterWithThisEmail = "already register with this email"
	EmailIsNotVerified           = "your email is not verified please verify your email"
	EmailValidationError         = "wrong email passed"
	OtpValidationError           = "wrong otp passed"
	OtpExpiredValidationError    = "otp expired"
	AlreadyVerifiedError         = "already verified"
	OptAlreadySentError          = "otp already sent to email"
	NotRegisteredUser            = "you are not register user"
	PasswordNotMatchedError      = "password doesn't match"
	NotAuthorizedUserError       = "you are not authorized to do this"
	NoProductAvaliable           = "no product avaliable"
	UserDoesNotExists            = "user not exists"
	AddressNotExists             = "address not exists. please add one address"
)
