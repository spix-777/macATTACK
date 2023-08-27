# macATTACK

-- Description ---

macATTACK is a tool built for macOS that enables users to spoof their MAC address. It has capabilities to randomize, set specific MAC addresses, or reset to the original MAC address. This tool uses system-level commands and, as such, requires sudo permissions to execute most of its operations.

-- Features ---

* Customizable terminal banner display.
* Spoof your MAC address randomly.
* Set a specific MAC address from a list.
* Reset your MAC address to its original state.
* Custom logging to provide information on ongoing processes.
* Version and author display.

-- Dependencies ---

* networksetup command: Native to macOS for network configurations.
* spoof-mac: External command for MAC address spoofing.
* arp-scan: Tool to scan the local network and retrieve MAC addresses.
* ping: System command used to check internet connectivity.

Usage

Run the program with desired flags:

bash
Copy code
$ go run <filename>.go [FLAGS]

Available Flags:
-v: Display the version and author information.
-s: Spoof the MAC address either randomly or try MAC addresses from a list.
-r: Reset the MAC address to its original state.

-- Examples ---
To display version and author info:

$ go run <filename>.go -v
To spoof your MAC address:

$ go run <filename>.go -s
To reset your MAC address:

$ go run <filename>.go -r

--- Note  ---

Please use this tool responsibly and ethically. Spoofing MAC addresses may violate terms of service of some networks and could be illegal in some jurisdictions.

-- Contributing ---

Feel free to submit pull requests to enhance the capabilities of this tool. Ensure you test new features before submitting.

-- Credits ---

Author: @Spix-777

-- Disclaimer ---

The author of this tool, contributors, or the platform do not bear any responsibility for any misuse or damages caused by this tool. Always seek permission before conducting any penetration testing.

Happy Hacking! 🚀🛠