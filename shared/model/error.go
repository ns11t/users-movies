package model

// Error codes for API clients
const (
	ErrorCodeBadToken                 = "bad_token"
	ErrorCodeBadUserModel             = "bad_user_model"
	ErrorCodeUserAlreadyExists        = "user_already_exists"
	ErrorCodeUserNotFound             = "user_not_found"
	ErrorCodeUserIncorrectPassword    = "user_incorrect_password"
	ErrorCodeUserMovieNotFound        = "user_movie_not_found"
	ErrorCodePersonMovieAlreadyExists = "person_movie_already_exists"
	ErrorCodeMovieNotFound            = "movie_not_found"
	ErrorCodeBadLimitOffsetFormat     = "bad_limit_offset_format"
	ErrorCodeBadMovieIdFormat         = "bad_movie_id_format"
	ErrorCodeBadMovieYearFormat       = "bad_movie_year_format"
	ErrorCodeInternal                 = "internal_server_error"
)
