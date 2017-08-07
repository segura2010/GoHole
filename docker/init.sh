#!/bin/sh
redis-server &
gohole -s -c /root/gohole_config.json
