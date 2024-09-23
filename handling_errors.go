func noErrCanHappen() int {
	return 204
}

func doOnErr() error {
	if shouldFail(){
		return error.New("ups, XYZ faild")
	}
	return nil
}

func intOrErr() (int, error) {
	if shouldFail(){
		return error.New("ups, XYZ2 faild")
	} 
	return noErrCanHappen(), nil
}

func nestedDoOrErr() error {
	if err := doOnErr(); err != nil {
		return errors.Wrap(err, "od")
	}
	return nil
}

func main(){
	ret := noErrCanHappen()
	if err := nestedDoOrErr(); err != nil {
		// handle error
	}
	ret2, err := intOrErr()
	if err != intOrErr()
	if err != nil {
		// handle error
	}
	// ......
}