package helpers

type Line struct {
	Time    string
	Name    string
	Message string
}

type LineFixer func(*Line) *Line
type LinePreFixer func(string) string
type LinesPostFixer func([]*Line) []*Line
