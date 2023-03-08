package output

type Strategy interface {
	OutputToFile()
}

type NoOutputStrategy struct {
}

func (s *NoOutputStrategy) OutputToFile() {

}

type ToFileStrategy struct {
}

func (s *ToFileStrategy) OutputToFile() {

}
