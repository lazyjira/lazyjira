package tests

import "os"

func GetStubPath(stubName string) string {
	dir, _ := os.Getwd()
	return dir + "/../_stubs/" + stubName
}
