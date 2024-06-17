#!/bin/sh
# On openwrt routers, this should print the value of the WAN interface's ip address
# which should be your internet-facing IP.
# Warning: YMMV with this script :)
ip -4 addr show eth0 | sed -Ene 's/^.*inet ([0-9.]+)\/.*$/\1/p'
