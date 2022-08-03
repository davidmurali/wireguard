

## This document can be used to compile and run user-space wireguard or to set up native-space wireguard

### step-1: Install native space wireguard
    sudo apt update
    sudo apt install wireguard
    sudo apt install net-tools

### step-2: Create a private and public key
     wg genkey | sudo tee /etc/wireguard/private.key
     sudo chmod go= /etc/wireguard/private.key
	
### step-3: Setting up the go env variable in profile
    nano ~/.profile

### step-4: Add below to end of file (~/.profile)
    export PATH=$PATH:/usr/local/go/bin

### step-5: Save env variables
    source ~/.profile

### step-6: Create link for resolvconf, to use a peer as a DNS server
    ln -s /usr/bin/resolvectl /usr/local/bin/resolvconf
    apt-get install build-essential

## Above steps(1-6) needs to be done for both client and server

### step-7: Create wg0.conf file in server and paste the below content
	sudo nano /etc/wireguard/wg0.conf

	[Interface]
	Address = 10.0.0.1/24
	MTU=1420
	ListenPort = 51820	
	#privatekey of server
	PrivateKey = eNEAagtiA/CxwpTb8jvTkilwsYNbMojZXyH005lwVVQ= 
    	PostUp = iptables -t nat -I POSTROUTING 1 -s 10.0.0.0/24 -o eth0 -j MASQUERADE
	PostUp = iptables -I INPUT 1 -i wg0 -j ACCEPT
	PostUp = iptables -I FORWARD 1 -i eth0 -o wg0 -j ACCEPT
	PostUp = iptables -I FORWARD 1 -i wg0 -o eth0 -j ACCEPT
    	PostDown = iptables -t nat -D POSTROUTING -s 10.0.0.0/24 -o eth0 -j MASQUERADE
	PostDown = iptables -D INPUT -i wg0 -j ACCEPT
	PostDown = iptables -D FORWARD -i eth0 -o wg0 -j ACCEPT
	PostDown = iptables -D FORWARD -i wg0 -o eth0 -j ACCEPT

	[Peer]
    	#publickey of client
	PublicKey = J3wlDiIOaEiGImF8dZbCmCssqduBjcgtChcNXs9PjSo=
	AllowedIPs = 10.0.0.2/32

--------------------------------------

### step-8: Ip forwarding - other devices on the network will connect to server
    sysctl -w net.ipv4.ip_forward=1


### step-9: Create wg0.conf file in client and paste the below content
	sudo nano /etc/wireguard/wg0.conf

   	#This client is using the VPN for internet access.
	[Interface]
              #privatekey of client
	PrivateKey = wJ91tXL92kYM7bQFJhlFOY6hQtbGfEU07f+EPiexN3Q=
	Address = 10.0.0.2/24
	MTU=1384

	[Peer]
              #publickey of server
	PublicKey = 78eIdV4yEnLgB2ecXmViEc9/Y3PIhYxwIdUT9mIVry0=
              Endpoint = 35.172.27.241:51820
              AllowedIPs = 10.0.0.0/24,172.31.0.0/24
-----------------------------------

## Below steps(10-17) needs to be done for user-space in client and server

### step-10: Download and install go
    wget https://go.dev/dl/go1.18.3.linux-amd64.tar.gz
    tar -xf go1.18.3.linux-amd64.tar.gz
    mv ./go /usr/local
    
### step-11: Download wireguard-go source from below path and build
    git clone https://github.com/davidmurali/wireguard.git

### step-12: Download wg-quick-go source from below path
    https://github.com/davidmurali/wg-quick-go.git

wg-quick-go usually runs native space.
So, we modified the code to run user space wireguard

### step-13: Build wg-quick-go
    cd wg-quick-go/cmd/wg-quick
    go build -v -o wg-quick

### step-14: Copy the executable to /usr/bin
    cp wg-quick /usr/bin/

### step-15: Open the terminal and run 
    ./wireguard-go wg0

### step-16: Start the wireguard userspace tool
    wg-quick up wg0	

### step-17: Check the connection between server and client
    wg show
