package nominal

import "time"

type MyStruct struct {
	Bool     bool
	Integers []int
	String   *string
	Times    []*time.Time
}
