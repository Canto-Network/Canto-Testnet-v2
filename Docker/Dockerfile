FROM ubuntu:latest 

# add your keyname
ARG KEY_NAME
ARG MONIKER
ARG NETWORK_ID

# arbitrary validator metadata
ARG DETAILS

ARG TOKEN_DELEGATION

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update &&\ 
	apt-get install -y apt-utils git git-extras software-properties-common vim \
	wget ca-certificates curl build-essential libssl-dev make openssh-server snap

RUN snap install go --classic

# builds out cantod
RUN git clone https://github.com/Canto-Network/Canto-Testnet
WORKDIR /Canto-Testnet/cmd/cantod
RUN go install -tags ledger ./... &&\
    mv $HOME/go/bin/cantod /usr/bin/

# generate keys
RUN cantod keys add $KEY_NAME &&\
    cantod keys add $KEY_NAME --recover &&\
    cantod keys add $KEY_NAME --ledger 

# VALIDATOR SETUP 

# initialize node
RUN cantod init $MONIKER --chain-id canto_7777-1

# get genesis file
RUN wget https://github.com/Canto-Network/Canto-Testnet/raw/main/Networks/Testnet/genesis.json -P $HOME/.cantod/config/

RUN sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.0001acanto"/g' $HOME/.cantod/config/app.toml

WORKDIR /etc/systemd/system/
RUN echo "[Unit]
Description=Canto Node
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/
ExecStart=/root/go/bin/cantod start --trace --log_level info --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable
Restart=on-failure
StartLimitInterval=0
RestartSec=3
LimitNOFILE=65535
LimitMEMLOCK=209715200

[Install]
WantedBy=multi-user.target" >> cantod.service

WORKDIR /Canto-Testnet/cmd/cantod
RUN sudo systemctl daemon-reload &&\
    sudo systemctl enable cantod.service &&\
    systemctl start cantod && journalctl -u cantod -f

RUN cantod tx staking create-validator \
--from {{KEY_NAME}} \
--chain-id canto_7777-1 \
--moniker="$MONIKER" \
--commission-max-change-rate=0.01 \
--commission-max-rate=1.0 \
--commission-rate=0.05 \
--details="$DETAILS" \
--security-contact="$CONTACT" \
--website="$WEBSITE" \
--pubkey $(cantod tendermint show-validator) \
--min-self-delegation="1" \
--amount 100000acanto \
--node https://canto.rpc.chandrastation.com:443


ENTRYPOINT ["/bin/bash"]
