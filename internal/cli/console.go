package cli

import "github.com/fatih/color"

var printfBold = color.New(color.Bold).Printf
var printfRed = color.New(color.FgRed).Printf
var printRed = color.New(color.FgRed).Print
var sprintRed = color.New(color.FgRed).Sprint
var sprintGreen = color.New(color.FgGreen).Sprint
var sprintBold = color.New(color.Bold).Sprint
