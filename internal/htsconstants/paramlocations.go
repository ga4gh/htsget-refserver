// Package htsconstants contains program constants
//
// Module paramlocations contains constants relating to parameter locations within
// an HTTP request
package htsconstants

// ParamLoc enum for parameter location
type ParamLoc int

// enum values for ParamLoc
const (
	ParamLocPath    ParamLoc = 0
	ParamLocQuery   ParamLoc = 1
	ParamLocHeader  ParamLoc = 2
	ParamLocReqBody ParamLoc = 3
)
