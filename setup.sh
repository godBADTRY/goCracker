#!/usr/bin/env bash


# Checking dependencies



if [ "$EUID" -ne 0 ]; then
	echo "[!] You must run this script as root"
	exit 1
fi


if command -v git &>/dev/null; then
	echo "[+] Git installed"
else
	echo "[!] Git not installed"
	echo "Installing git..."
	apt install git
fi

if command -v go &>/dev/null; then
	echo "[+] Go installed"
else
	echo "[!] Go not installed"
	echo "Installing Go"
	Go_ver=$(curl -s https://go.dev/VERSION\?m\=text | head -n1)
	wget https://go.dev/dl/$Go_ver.linux-amd64.tar.gz
	tar -C /usr/local -xzf go*.tar.gz
fi
