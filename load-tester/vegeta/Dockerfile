FROM ubuntu:22.04

# Required ENVs
# SERVER_API_IP
# RATE

# Update and install dependencies
RUN apt-get update && apt-get install -y curl 

# Install Vegeta
RUN curl -Lo vegeta.tar.gz "https://github.com/tsenart/vegeta/releases/latest/download/vegeta_$(curl -s "https://api.github.com/repos/tsenart/vegeta/releases/latest" | grep -Po '"tag_name": "v\K[0-9.]+')_linux_amd64.tar.gz" \
    && mkdir vegeta-temp \
    && tar xf vegeta.tar.gz -C vegeta-temp \
    && mv vegeta-temp/vegeta /usr/local/bin \
    && rm -rf vegeta.tar.gz vegeta-temp

WORKDIR /app
COPY *.json .
COPY *.sh .
RUN chmod +x /app/metrics.sh

CMD /app/metrics.sh $RATE
