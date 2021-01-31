package logparser

type Parser interface {
	Scan() bool
	ParsedLog() ParsedLog
}
