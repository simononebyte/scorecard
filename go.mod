module github.com/simonbuckner/scorecard

go 1.12

require (
	github.com/simononebyte/restup v0.0.0-20190911135854-81d615f2e651
	github.com/simononebyte/scorecard/psa v0.0.0
)

replace (
	github.com/simononebyte/restup => ../restup
	github.com/simononebyte/scorecard/psa => ./psa
)
