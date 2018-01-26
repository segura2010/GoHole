#!/bin/sh
redis-server &
# Generate a new encryption key every time the container starts
gohole -gkey
# Run GoHole
gohole -s -c /root/gohole_config.json
