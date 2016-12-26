package testController

import "syntinel_executor/structures"

type TestController struct {
	tests map[int]Test
	queue structures.Queue
}

func NewTestController() TestController {
	return TestController{make(map[int]Test), structures.NewQueue()}
}
