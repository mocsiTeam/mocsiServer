package errs

var CodeErr = map[uint]string{
	0:  "unknown error",
	1:  "access denied",
	11: "name already exists",
	12: "email already exists",
	21: "wrong username or password",
	22: "user not found",
	31: "name alredy exists",
	32: "user not found",
	33: "no rights",
	34: "group_not_found",
	35: "failed to delete room",
	36: "failed to delete rights",
	37: "failed to create room",
	38: "failed to create rights",
}
