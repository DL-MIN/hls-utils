package types

import (
	"net/http"

	problems "github.com/spacecafe/gobox/gin-problems"
)

var (
	ProblemNoSuchStream = problems.NewProblem(
		"",
		http.StatusText(http.StatusNotFound),
		http.StatusNotFound,
		"The specified stream does not exist.",
	)
)
