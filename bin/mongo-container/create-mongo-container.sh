#!/bin/bash

set -e
. "../functions.sh"

# Print the usage message
function printHelp() {
  echo "Usage: "
  echo "  create-mongo-container.sh [Flags]"
  echo "    Flags:"
  echo "    -n <container name> - defines the containers name (default: mongodb)"
  echo "    -p <ports range> - defines ports range (default: 27017-27019)"
  echo "    -v <mongo version> - defines mongo version (default: latest)"
  echo "    -h (print this message)"
  echo
  echo "Taking all defaults:"
  echo "  create-mongo-container.sh"
}

# default values
NAME=mongodb
PORTS=27017-27019
VERSION=latest

# parse flags

while [[ $# -ge 1 ]] ; do
  key="$1"
  case $key in
  -h )
    printHelp
    exit 0
    ;;
  -n )
    NAME="$2"
    shift
    ;;
  -p )
    PORTS="$2"
    shift
    ;;
  -v )
    VERSION="$2"
    shift
    ;;
  * )
    echo
    echo "Unknown flag: $key"
    echo
    printHelp
    exit 1
    ;;
  esac
  shift
done

pp_info "Deploy Container" "Define Mongo container definitions"
docker run -d --name $NAME -p $PORTS:27017-27019 mongo:$VERSION
pp_success "Deploy Container" "You're good to go."
