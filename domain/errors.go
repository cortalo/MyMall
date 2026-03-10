package domain

import "errors"

var ErrRegionAbandoned = errors.New("region abandoned")
var ErrRegionNotAllowedStatus = errors.New("region not allowed status")
