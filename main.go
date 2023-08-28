package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// networkDeviceOnMac is the name of the network device on Mac
// Constants are defined outside of the main function
const networkDeviceOnMac = "en0"

// sudoPassword function initiates a sudo command execution with the provided loggers
func sudoPassword(uLogger, minLogger, plussLogger *log.Logger) {
	commandName := "echo"
	// Use an empty string as the first argument since it's not relevant for 'echo'
	execSudo(commandName, "", "", "", 2, uLogger, minLogger, plussLogger)
}

// wifi function controls the WiFi state using the networksetup command
// onOff should be either "on" or "off"
func wifi(onOff string, uLogger, minLogger, plussLogger *log.Logger) {
	// Construct the command
	cmd := exec.Command("networksetup", "-setairportpower", networkDeviceOnMac, onOff)

	// Run the command and capture any errors
	if err := cmd.Run(); err != nil {
		uLogger.Fatalf("Error toggling WiFi: %v\x1b[0m", err)
	}
}

// execSudo function executes commands with sudo and varying arguments
func execSudo(arg string, arg2 string, arg3 string, arg4 string, isMac int, uLogger, minLogger, plussLogger *log.Logger) string {
	var cmdArgs []string

	// Determine the command arguments based on the provided 'isMac' value
	switch isMac {
	case 1:
		cmdArgs = []string{arg2, networkDeviceOnMac}
	case 2:
		cmdArgs = []string{arg2}
	case 3:
		cmdArgs = []string{arg2, arg3, networkDeviceOnMac}
	case 4:
		cmdArgs = []string{arg2, arg3, arg4}
	case 5:
		cmdArgs = []string{arg2, arg3, arg4}
	default:
		uLogger.Fatalf("Invalid isMac value: %d\x1b[0m", isMac)
	}

	// Construct the command with 'sudo' prefix
	cmd := exec.Command("sudo", append([]string{arg}, cmdArgs...)...)

	// Execute the command and capture any errors
	output, err := cmd.CombinedOutput()
	if err != nil {
		//uLogger.Fatalf("Error executing command: %v", err)
	}

	return string(output)
}

// removeAllbutMac function extracts MAC addresses from captured output and returns a slice
func removeAllbutMac(output string, uLogger, minLogger, plussLogger *log.Logger) ([]string, []string) {
	var macAddresses []string
	var ipAddresses []string

	// Split the output into lines
	lines := strings.Split(output, "\n")

	// Iterate over relevant lines to extract MAC addresses
	for _, line := range lines[2 : len(lines)-3] {
		if len(line) == 0 {
			break // Exit loop if an empty line is encountered
		}

		// Split the line into words
		words := strings.Fields(line)
		if len(words) >= 2 {
			ipAddresses = append(ipAddresses, words[0])   // First word is the IP address
			macAddresses = append(macAddresses, words[1]) // Second word is the MAC address
		}

		// Look for error messages in the IP and MAC addresses
		// If an error is found, print the error message and exit
		// IP addresses should be numbers only
		// MAC addresses should be in the format xx:xx:xx:xx:xx:xx
		for _, line := range ipAddresses {
			buffer := strings.ReplaceAll(line, ".", "")
			_, err := strconv.Atoi(buffer)
			if err != nil {
				uLogger.Fatalln("Error:", err)
			}
		}

		for _, line := range macAddresses {
			buffer := strings.Split(line, ":")
			if len(buffer) != 6 {
				uLogger.Fatalln("Error: Invalid MAC address")
			}
		}

	}
	return ipAddresses, macAddresses
}

// reset function resets the MAC address
func reset(uLogger, minLogger, plussLogger *log.Logger) {
	//sudoPassword(uLogger, minLogger, plussLogger)
	// Turn off WiFi
	minLogger.Println("Turning off WiFi")
	wifi("off", uLogger, minLogger, plussLogger)

	// Reset MAC address
	plussLogger.Println("Resetting MAC address")
	cmd := exec.Command("sudo", "spoof-mac", "reset", networkDeviceOnMac)
	if err := cmd.Run(); err != nil {
		uLogger.Fatalf("Error resetting MAC address: %v\x1b[0m", err)
	}

	// Reset IP address
	plussLogger.Println("Resetting IP address")
	cmd = exec.Command("sudo", "networksetup", "-setdhcp", "Wi-Fi")
	if err := cmd.Run(); err != nil {
		uLogger.Fatalf("Error resetting MAC address: %v\x1b[0m", err)
	}

	// Turn on WiFi
	plussLogger.Println("Turning on WiFi\x1b[0m")
	wifi("on", uLogger, minLogger, plussLogger)
}

// banner function to print the banner
func banner() {
	banner := `
	███▄ ▄███▓ ▄▄▄       ▄████▄   ▄▄▄     ▄▄▄█████▓▄▄▄█████▓ ▄▄▄       ▄████▄   ██ ▄█▀
	▓██▒▀█▀ ██▒▒████▄    ▒██▀ ▀█  ▒████▄   ▓  ██▒ ▓▒▓  ██▒ ▓▒▒████▄    ▒██▀ ▀█   ██▄█▒ 
	▓██    ▓██░▒██  ▀█▄  ▒▓█    ▄ ▒██  ▀█▄ ▒ ▓██░ ▒░▒ ▓██░ ▒░▒██  ▀█▄  ▒▓█    ▄ ▓███▄░ 
	▒██    ▒██ ░██▄▄▄▄██ ▒▓▓▄ ▄██▒░██▄▄▄▄██░ ▓██▓ ░ ░ ▓██▓ ░ ░██▄▄▄▄██ ▒▓▓▄ ▄██▒▓██ █▄ 
	▒██▒   ░██▒ ▓█   ▓██▒▒ ▓███▀ ░ ▓█   ▓██▒ ▒██▒ ░   ▒██▒ ░  ▓█   ▓██▒▒ ▓███▀ ░▒██▒ █▄
	░ ▒░   ░  ░ ▒▒   ▓▒█░░ ░▒ ▒  ░ ▒▒   ▓▒█░ ▒ ░░     ▒ ░░    ▒▒   ▓▒█░░ ░▒ ▒  ░▒ ▒▒ ▓▒
	░  ░      ░  ▒   ▒▒ ░  ░  ▒     ▒   ▒▒ ░   ░        ░      ▒   ▒▒ ░  ░  ▒   ░ ░▒ ▒░
	░      ░     ░   ▒   ░          ░   ▒    ░        ░        ░   ▒   ░        ░ ░░ ░ 
		   ░         ░  ░░ ░            ░  ░                       ░  ░░ ░      ░  ░   
						 ░                                             ░               
	`
	fmt.Println(banner)
}

func main() {

	// Declare variables for storing output and MAC addresses
	var output string
	var ipAddress []string
	var macAddress []string

	// Create custom loggers with different prefixes for different types of messages
	uLogger := log.New(os.Stdout, "\x1b[31m [ ! ] ", log.LstdFlags)     // Logger for user-level messages
	verLogger := log.New(os.Stdout, " ", 0)                             // Logger for version information
	minLogger := log.New(os.Stdout, " [ - ] ", log.LstdFlags)           // Logger for minor status updates
	plussLogger := log.New(os.Stdout, "\x1b[32m [ + ] ", log.LstdFlags) // Logger for successful actions

	// Define command-line flags for different functionalities
	versionFlag := flag.Bool("v", false, "Version") // Display version information
	spoofFlag := flag.Bool("s", false, "Spoof Mac") // Spoof MAC addresses
	resetFlag := flag.Bool("r", false, "Reset Mac") // Reset MAC address
	flag.Parse()                                    // Parse the command-line flags

	sudoPassword(uLogger, minLogger, plussLogger) // Run a sudo command that requires a password
	banner()                                      // Display a custom banner

	if *versionFlag {
		// If version flag is set, display version information and exit
		verLogger.Println("                                   Version:     \x1b[32m1.0.1\x1b[0m")
		verLogger.Println("                                   Author :     \x1b[31m@Spix-777\x1b[0m")
		os.Exit(0)
	} else if *spoofFlag {
		// If spoof flag is set, perform MAC spoofing process
		minLogger.Println("Turning off wifi")
		wifi("off", uLogger, minLogger, plussLogger) // Turn off WiFi
		minLogger.Println("Randomizing MAC address")
		_ = execSudo("spoof-mac", "randomize", "", "", 1, uLogger, minLogger, plussLogger) // Randomize MAC address
		minLogger.Println("Turning on wifi")
		wifi("on", uLogger, minLogger, plussLogger) // Turn on WiFi
		time.Sleep(10 * time.Second)                // Pause for 10 seconds
		minLogger.Println("Get the MAC address of the network device")
		output = execSudo("arp-scan", "-l", "", "", 2, uLogger, minLogger, plussLogger)  // Get IP and MAC addresses using arp-scan
		ipAddress, macAddress = removeAllbutMac(output, uLogger, minLogger, plussLogger) // Extract IP and MAC addresses
		count := len(ipAddress)                                                          // Get the number of IP addresses
		for i := 0; i < count; i++ {
			minLogger.Println("Turning off wifi")
			wifi("off", uLogger, minLogger, plussLogger) // Turn off WiFi
			minLogger.Println("Try to spoof " + macAddress[i])
			_ = execSudo("spoof-mac", "set", macAddress[i], "", 3, uLogger, minLogger, plussLogger) // Set MAC address

			minLogger.Println("Setting up manual IP address " + ipAddress[i] + " with DHCP router")
			_ = execSudo("networksetup", "-setmanualwithdhcprouter", "Wi-Fi", ipAddress[i], 5, uLogger, minLogger, plussLogger) // Show MAC address

			minLogger.Println("Turning on wifi")
			wifi("on", uLogger, minLogger, plussLogger) // Turn on WiFi
			time.Sleep(10 * time.Second)                // Pause for 10 seconds
			minLogger.Println("Pinging google.com")
			output = execSudo("ping", "-c", "5", "8.8.8.8", 4, uLogger, minLogger, plussLogger) // Ping Google
			if strings.Contains(output, "5 packets received") {
				plussLogger.Println("Spoofed " + macAddress[i] + " successfully\x1b[0m")
				plussLogger.Println("Spoofed " + ipAddress[i] + " successfully\x1b[0m")
				plussLogger.Println("You have internet access!\x1b[0m")
				break // Exit the loop if spoofing is successful
			} else {
				uLogger.Println("Spoofed " + macAddress[i] + " failed\x1b[0m")
				uLogger.Println("Spoofed " + ipAddress[i] + " failed\x1b[0m")
			}
		}

	} else if *resetFlag {
		// If reset flag is set, reset the MAC address
		reset(uLogger, minLogger, plussLogger) // Call the reset function
	} else {
		// If no flag is provided, display an error message
		uLogger.Println("You must have a flag!\x1b[31m")
	}
}
