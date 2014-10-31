package main

const MaximumOutNodes = 4
const MaximumInNodes = 4

const MinimumOutNodes = 2
const MinimumInNodes = 1

const NumRetinaSlices = 12
const NumActuators = 24
const NumHiddenNodes = 6

//                    two layers øf retina in prey
const NumTotalNodes = NumRetinaSlices*2 + NumActuators + NumHiddenNodes
