package models

var (
	InternalServerError = Response{Success: false, Error: "Internal server error"}
	NotFoundError       = Response{Success: false, Error: "Not found"}
	DatabaseError       = Response{Success: false, Error: "Database error"}
	InvalidRequestError = Response{Success: false, Error: "Invalid request"}
	InvalidTokenError   = Response{Success: false, Error: "Invalid token"}
	ForbiddenError      = Response{Success: false, Error: "Forbidden"}
	AlreadyExistsError  = Response{Success: false, Error: "Already exists"}

	Success = Response{Success: true}
)
