# GoHole

GoHole is a DNS server written in Golang with the same idea than the [PiHole](https://pi-hole.net), blocking advertisements's and tracking's domains.

It uses a Redis DB as cache.

### Installation

1. Clone this repository and rename the folder to GoHole if it is not the name.
2. Run the install script `install.sh` to install all the dependencies.
3. Compile using Makefile (`make`). Or run `make install` to install (you should have included your $GOPATH/bin to your $PATH).

Finally, run using the executable for your platform.

### Usage

To start the DNS server you have to run the following command:

`gohole -s`

You can specify a config file with the command line argument `-c`. See the `config_example.json` file to see the structure.

You can also provide the `-p` argument to specify the port in which the DNS server will listen.

To block ads domains, you must add them to the cache DB. In order to do that, you must pass a blocklist file using the following command:

`gohole -ab path/to/blacklist_file`

If the list is published in a web server, you can provide the URL: 

`gohole -ab http://domain/path/to/blacklist_file`

If you does not know any blacklist, you can see the file `blacklists/list.txt`. It contains the blacklists used by the PiHole. You can use a file with a list of blacklist like the `blacklists/list.txt` file to automatically add all the lists:

`gohole -abl blacklists/list.txt`

You can also block domains by using the following command:

`gohole -ad google.com -ip4 0.0.0.0 -ip6 "::1"`

#### Flush cache and logs

You can flush cache and logs DBs.

**Flush domains cache**

`gohole -fcache`

**Flush logs**

`gohole -flog`


#### Statistics and Logs

You can see the stats and logs by using the following command line arguments:

**See all the clients that have made a request**

`gohole -lc`

**See all request made by a client**

`gohole -lip <clientip>`

**See all clients that queried a domain**

`gohole -ld <clientip>`

**See all clients that queried a domain**

`gohole -ld <clientip>`

### Docker

You can use GoHole in a Docker container. To do that, you can use the Docker image: https://hub.docker.com/r/segura2010/gohole/ running:

`docker pull segura2010/gohole`

Once you pull the Docker image, you can run a container using the command: 

`docker run -d --name gohole -p 53:53/udp --restart=unless-stopped gohole/master`

Then, you will have a Docker container running GoHole. But it does not install any blacklist (domains will not be blocked). In order to do that, you must open a shell in the container with:

`docker exec -it gohole /bin/bash`

After that, you can run `gohole -abl blacklists/list.txt` to set up the blocked domains.


**Tested on Go 1.8.3**
