package processor

type DefaultProcessor struct {}

func (p DefaultProcessor) Process(input string) (res string, err error) {
	return runBlogCore(splitArgs(input))
}
