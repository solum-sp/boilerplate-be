package model

import "proposal-template/pkg/utils"

var (
	ErrUnknown       = utils.NewCustomError("unknown")
	ErrMalformedJSON = utils.NewCustomError("malformed_json")
	ErrUnimplemented = utils.NewCustomError("unimplemented method")
)

var (
	ErrJWTSecretNotConfigured        = utils.NewCustomError("jwt_secret_not_configured")
	ErrJWTMissingAuthorizationHeader = utils.NewCustomError("jwt_missing_authorization_header")
	ErrJWTInvalidAuthorizationFormat = utils.NewCustomError("jwt_invalid_authorization_format")
	ErrJWTInvalidToken               = utils.NewCustomError("jwt_invalid_token")
	ErrJWTInvalidTokenClaims         = utils.NewCustomError("jwt_invalid_token_claims")
	ErrJWTTokenExpired               = utils.NewCustomError("jwt_token_expired")
	ErrJWTInvalidIssuer              = utils.NewCustomError("jwt_invalid_issuer")
	ErrJWTTokenNotYetValid           = utils.NewCustomError("jwt_token_not_yet_valid")
	ErrJWTUnexpectedSigningMethod    = utils.NewCustomError("jwt_unexpected_signing_method")
	ErrJWTFailToGenerateToken        = utils.NewCustomError("jwt_fail_to_generate_token")
)

var (
	ErrFailToChangePassword = utils.NewCustomError("fail_to_change_password")
	ErrSavingUser           = utils.NewCustomError("err_saving_user")
	ErrEmailNotAvailable    = utils.NewCustomError("email_not_available")
	ErrWrongPassword        = utils.NewCustomError("wrong_password")
	ErrUserNotFound         = utils.NewCustomError("user_not_found")
)
