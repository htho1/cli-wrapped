package main

import (
	"os"
	"fmt"
	"flag"
	"time"
	"strings"
	"os/user"
	"encoding/json"

	"github.com/fatih/color"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getHomeDir() string {
	user, err := user.Current()
	check(err)

	return user.HomeDir
}

func getNextWrappedYear() int {

	now := time.Now()

	if now.Month() == 12 {
		return now.Year() + 1
	} else {
		return now.Year()
	}
}

func checkForDataFile() {
	_, err := os.Stat(getHomeDir() + "/.cli-wrapped")

	if os.IsNotExist(err) {
		color.Cyan("cli-wrapped data file does not exist, creating it now")
		check(os.Mkdir(getHomeDir() + "/.cli-wrapped", 0755))
		check(os.WriteFile(getHomeDir() + "/.cli-wrapped/data.json", []byte("{}"), 0644))
	}
}

func readDataFile() map[string]interface{} {

	checkForDataFile()
	var decoded map[string]interface{}

	fileContent, err := os.ReadFile(getHomeDir() + "/.cli-wrapped/data.json")
	check(err)

	check(json.Unmarshal(fileContent, &decoded))

	return decoded
}

func writeDataFile(data map[string]interface{}) {

	checkForDataFile()
	encoded, err := json.Marshal(data)
	check(err)

	check(os.WriteFile(getHomeDir() + "/.cli-wrapped/data.json", encoded, 0644))
}

func main() {
	
	aboutFlagPtr := flag.Bool("about", false, "Show information about cli-wrapped")
	trackFlagPtr := flag.String("track", "", "Track a command. Don't manually use unless you're cheating!")

	flag.Parse()

	if *aboutFlagPtr {
		color.Cyan("cli-wrapped\n\n")
		fmt.Println("A utility to track the commands you use and present them to you at\nthe end of the year, like Spotify Wrapped for your CLI.")
		fmt.Println("Your wrapped is revealed at the end of the year, on December 1st.\n")
		fmt.Println("Currently tracking for year:", getNextWrappedYear())
		fmt.Println("Location of data.json file:", getHomeDir() + "/.cli-wrapped")
		return
	}

	// Track command
	if *trackFlagPtr != "" {
		splitCommand := strings.Split(*trackFlagPtr, " ")
		trackCommand := splitCommand[0]

		if (trackCommand == "sudo" || trackCommand == "doas")  && len(splitCommand) > 1 {
			trackCommand = splitCommand[1]
		}

		data := readDataFile()
		var newValue float64

		if data[trackCommand] == nil {
			newValue = 0
		} else {
			newValue = data[trackCommand].(float64)
		}

		newValue += 1
		data[trackCommand] = newValue

		writeDataFile(data)
	}
}