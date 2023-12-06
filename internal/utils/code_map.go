package utils

var CodeMap map[int64]string

func init() {
	ResponseStatusCode()
}

func ResponseStatusCode() {
	CodeMap = map[int64]string{
		200:  "Successfully.",
		1001: "Wrong with the data format trasnfer from POST.",
		1002: "Failed to execute SQL syntax.",
		1003: "No search this account.",
		1004: "User password is not correct.",
		1005: "Failed to create session value JSON data.",
		1006: "User ID not found.",
		1007: "Totp code verify failed.",
		1008: "Generate TOTP key failed.",
		1009: "Failed to save QRcode.",
		1010: "Encode PNG failed.",
		1011: "Check whether the user_id exists or not failed.",
		1012: "Update totp table failed.",
		1013: "Insert totp table failed.",
		1014: "Generate the JWT string failed.",
		1015: "JWT signature is invalid.",
		1016: "Parse JWT claim data failed.",
		1017: "Token is invalid.",
		1018: "Generate new JWT failed.",
		1019: "Cookie key is not found.",
		1020: "Get cookie key failed.",
		1021: "The sign up data is not satisfied of rule.",
		1022: "Fail to send verification mail.",
		1023: "Email address verification failed.",
		1024: "Email address already verified.",
		1025: "Fail to send password reset mail.",
		1026: "Reset code invalid or expired.",
		1027: "Reset password failed.",
		1028: "Get RBAC members error",
		1029: "Get RBAC roles error",
		1030: "Add RBAC policy error",
		1031: "Add RBAC grouping policy error",
		1032: "Delete RBAC policy error",
		1033: "Delete RBAC grouping policy error",
		1034: "Delete RBAC role error",
		1035: "Delete RBAC member error",
		1036: "Parse date failed.",
		1037: "Update user_info table failed.",
		1038: "Query string should be include username key.",
		1039: "Checking whether key exist or not happen something wrong.",
		1040: "Delete Redis key failed.",
		1041: "Failed to create session value JSON data.",
		1042: "Redis SET data failed.",
		1043: "Redis key does not exist.",
		1044: "Redis GET data failed.",
		1045: "Check whether the username exists or not failed.",
		1046: "Check whether the mail exists or not failed.",
		1047: "OTP code is not correct.",
		1048: "Username not found.",
		1049: "Get rowsAffected failed.",
		1050: "Token has expired.",
		1051: "OTP verify result value is not bool.",
		1052: "username and mail_otp_verify are not exists.",
		1053: "2FA feature doesn't not enable.",
		1054: "The result of the email regular expression doesn't match the expected outcome.",
		1055: "Login failed because of session or jwt problem.",
		1056: "This account have already existed",
		1057: "No search this mail.",
		1058: "The new password is not satisfied of rule.",
		1059: "New password is the same old one.",
		1060: "The user password have already existed.",
		1061: "The password can't empty.",
		1062: "The mail format is not correct.",
		1063: "The json data unmarshal failed.",
		1064: "You don't have the required permission.",
		1065: "Check user permission failed.",
		1101: "Fail to parse POST form data.",
		1102: "Fail to bind POST form data.",
		1103: "Fail to parse path parameters.",
		1104: "Invalid data to request.",
	}
}
