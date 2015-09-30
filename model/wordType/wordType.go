package wordType

type WordType int

const (
	New WordType = 1 + iota
	Success
	Failed
	Delimiter
	Ignore
)
