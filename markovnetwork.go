package main

import ()

// PLGMN represents the whole MN of gates
type MarkovNetwork struct {
}

func (s *MarkovNetwork) Run(sensorValues []bool) []bool {
	return []bool{}
}

//TODO
func (s *MarkovNetwork) ToString() string {
	return ""
}

func New(retinaSlices int) *MarkovNetwork {
	mn := &MarkovNetwork{}
	return mn
}
