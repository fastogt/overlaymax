package store

import "errors"

var ErrNilCollection = errors.New("pogreb collection is nil")
var ErrIDNotFind = errors.New("overlay with ID can't find")
