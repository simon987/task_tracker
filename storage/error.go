package storage

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
