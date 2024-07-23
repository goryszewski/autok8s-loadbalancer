package haproxy

func NewHaproxyrOperator(path string) Haproxy {
	return Haproxy{Path: path}
}
