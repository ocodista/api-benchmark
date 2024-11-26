#!/bin/bash
echo "SERVER_API_IP=${SERVER_API_IP}" >> /etc/environment

sudo apt-get update
sudo apt-get install -y curl 

# Installing vegeta
curl -Lo vegeta.tar.gz "https://github.com/tsenart/vegeta/releases/latest/download/vegeta_$(curl -s "https://api.github.com/repos/tsenart/vegeta/releases/latest" | grep -Po '"tag_name": "v\K[0-9.]+')_linux_amd64.tar.gz"

mkdir vegeta-temp
tar xf vegeta.tar.gz -C vegeta-temp

sudo mv vegeta-temp/vegeta /usr/local/bin

rm -rf vegeta.tar.gz
rm -rf vegeta-temp

# Remove existing Node.js (if installed)
sudo apt-get remove -y nodejs
sudo apt-get purge -y nodejs
sudo apt-get autoremove -y

# Install Node.js 21.X.X
sudo apt-get update && sudo apt-get install -y ca-certificates curl gnupg
curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --yes --dearmor -o /etc/apt/keyrings/nodesource.gpg
NODE_MAJOR=21
echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
sudo apt-get update && sudo apt-get install nodejs -y
echo "Installed Node!"
node -v

# Installing Go 1.21.4
wget https://golang.org./dl/go1.21.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go_path.sh
echo "Installed go!"

# Installing jagger
go install github.com/rs/jaggr@latest
go install github.com/rs/jplot@latest

touch /home/ubuntu/ended_startup
