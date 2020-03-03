#!/bin/bash

set -e
. "functions.sh"

env=${1:-"dev"}

pp_info "setup-2" "Install Samples, Binaries, and Docker Images"

cd ../../
if [ -d "fabric-samples" ]; then
  pp_success "setup-2" "Directory fabric-samples exists." 
else
  sudo curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.0.0 1.4.4 0.4.18
fi
cd bin/fabric-setup/

pp_success "setup-2" "You're good to go!"
