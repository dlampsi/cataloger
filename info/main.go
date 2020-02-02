package info

import "fmt"

// Version number.
var Version = "0.0.0"

// Release release num
var Release = ""

// BuildTime label of build time.
var BuildTime = ""

// ForPrintFull Returns formated version and build time string for print.
func ForPrintFull() string {
	return fmt.Sprintf("cataloger %s\nBuild time %s\n", Version, BuildTime)
}

// ForPrint Returns formated version string for print.
func ForPrint() string {
	grafic := `                   
  _____       __         __                       
 / ___/___ _ / /_ ___ _ / /___  ___ _ ___  ____   
/ /__ / _ '// __// _ '// // _ \/ _ '// -_)/ __/   
\___/ \_,_/ \__/ \_,_//_/ \___/\_, / \__//_/      
                              /___/               `

	return fmt.Sprintf("%s \n\nVersion %s\nBuild time %s", grafic, Version, BuildTime)

}
