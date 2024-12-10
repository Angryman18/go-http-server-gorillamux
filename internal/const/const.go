package constants

const (
	Regex string = `^(=.*[a-z])(=.*[A-Z])(=.*\d)(=.*[@$!%*&])[A-Za-z\d@$!%*&]{8,}$`
)

type ClaimType string

var Claim ClaimType = "Claim"
