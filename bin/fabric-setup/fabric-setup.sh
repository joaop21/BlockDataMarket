#!/usr/bin/env sh

set -e
. "functions"

env=${1:-"dev"}

pp_info "setup" "Installing Prerequisites"

if not_installed "git"; then
  pp_error "setup" "
  We are using git to download the needed repositories, since it was not found on your  system we cannot ensure that you are using the correct versions of all the tools. 
  Please install it and run this script again, or proceed at your own peril.
  "

  ensure_confirmation
else
  pp_success "setup" "git is already installed"
fi

if not_installed "curl"; then
  pp_error "setup" "
  We are using curl, since it was not found on your  system we cannot ensure that you   are using the correct versions of all the tools. Please install it and run this 
  script again, or proceed at your own peril.
  "

  ensure_confirmation
else
  pp_success "setup" "curl is already installed"
fi


if not_installed "docker"; then
  pp_error "setup" "We are using docker for our fabric network which isn't installed."
  pp_info "setup" "We are installing docker"
  sudo apt-get update
  sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  sudo apt-key fingerprint 0EBFCD88
  sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
  sudo apt-get update
  sudo apt-get install docker-ce docker-ce-cli containerd.io
  pp_info "setup" "Enabling Docker to start on boot."
  sudo systemctl start docker
  sudo systemctl enable docker
  sudo usermod -a -G docker $(user)
  pp_success "setup" "docker is up-and-running"
else
  docker_state=$(sudo docker info >/dev/null 2>&1)
  if [[ $? -ne 0 ]]; then
    pp_warn "setup" "docker does not seem to be running, run it first and retry"
    exit 1
  else
    pp_success "setup" "docker is up-and-running"
  fi
fi


if not_installed "docker-compose"; then
  pp_error "setup" "We are using docker-compose for our fabric network which isn't installed."
  pp_info "setup" "We are installing docker-compose"
  sudo curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
  sudo chmod +x /usr/local/bin/docker-compose
  pp_success "setup" "docker-compose is up-and-running"
else
  pp_success "setup" "docker-compose is already installed"
fi


pp_info "setup" "Installing Required Languages"

if not_installed "go"; then
  pp_error "setup" "We are using Go as development language which isn't installed."
  pp_info "setup" "We are installing Go."
  wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.14.linux-amd64.tar.gz
  echo "export GOROOT=\"/usr/local/go\"" >> ~/.profile
  echo "export GOBIN=\"\$HOME/go/bin\"" >> ~/.profile
  echo "mkdir -p \$GOBIN" >> ~/.profile
  echo "export PATH=\$PATH:\$GOROOT/bin:\$GOBIN" >> ~/.profile
  touch ~/.profile
  rm go1.14.linux-amd64.tar.gz

  pp_warn "setup" "You just need to exit your account and enter again."
else
  pp_success "setup" "Go is already installed"
fi

pp_info "setup" "Install Samples, Binaries, and Docker Images"

cd ..
if [ -d "fabric-samples" ]; then
  pp_success "setup" "Directory fabric-samples exists." 
else
  curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.0.0 1.4.4 0.4.18
fi
cd bin/

pp_success "setup" "You're good to go!"
