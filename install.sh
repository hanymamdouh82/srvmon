#!/bin/bash

BINARY_PATH="/usr/local/bin"
CONF_PATH="/etc/srvmon"
UNIT_PATH="/etc/systemd/system"

echo "Installing/Updating srvmon..."

# Stop the service if it's running
if systemctl is-active --quiet srvmon; then
    sudo systemctl stop srvmon
fi

# Copy binary
sudo install -m 755 srvmon "${BINARY_PATH}/srvmon"

# Ensure config directory exists (preserve old config)
sudo mkdir -p "$CONF_PATH"
if [ ! -f "${CONF_PATH}/conf.yaml" ]; then
    sudo cp conf.yaml "${CONF_PATH}/"
    echo "New config installed."
else
    echo "Config file exists, keeping the old one."
fi
#
# Generate systemd service file
sudo tee "${UNIT_PATH}/srvmon.service" > /dev/null <<EOF
[Unit]
Description=Hany Server Monitoring Daemon
After=network.target

[Service]
ExecStart=${BINARY_PATH}/srvmon -server -conf=${CONF_PATH}/conf.yaml
Restart=always
User=root
LimitNOFILE=65535
StandardOutput=journal
StandardError=journal
Environment="GODEBUG=madvdontneed=1"

[Install]
WantedBy=multi-user.target
EOF

echo "Systemd service file created."

# Install or update systemd service
sudo systemctl daemon-reload
sudo systemctl enable srvmon

# Start the service again
sudo systemctl restart srvmon

echo "Update complete! srvmon is now running."

