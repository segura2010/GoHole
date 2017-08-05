# GoHole

GoHole is a DNS server written in Golang with the same idea than the [PiHole](https://pi-hole.net), blocking advertisements's and tracking's domains.

It uses a Redis DB as cache.

### Installation

1. Clone this repository and rename the folder to GoHole if it is not the name.
2. Run the install script `install.sh` to install all the dependencies.
3. Compile using Makefile (`make`).

Finally, run using the executable for your platform.

**Tested on Go 1.7.3**
