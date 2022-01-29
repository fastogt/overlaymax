package store

import "errors"

var ErrNilCollection = errors.New("mongoDB Collection is nil")
var ErrIDNotFind = errors.New("overlay with ID can't find")
