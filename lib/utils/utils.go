package utils

func FatalError(err error) {
	if err != nil {
		Fatal(err)
	}
}
