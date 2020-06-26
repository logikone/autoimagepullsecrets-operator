package main

func checkError(err error) {
	if err != nil {
		setupLog.Error(err, "error during setup")
	}
}
