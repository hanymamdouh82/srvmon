# Server Monitoring Tool

A lightweight server monitoring tool that collects and sends system data to a monitoring client.

---

## 🚀 Server Side  

### 📁 File Locations  

| File               | Location                         | Description                                          |
|--------------------|---------------------------------|------------------------------------------------------|
| `srvmon`          | `/usr/local/bin/srvmon`         | Executable binary for the server-side application.  |
| `conf.yaml`       | `/etc/srvmon/conf.yaml`         | Configuration file defining collected data.         |
| `srvmon.service`  | `/etc/systemd/system/srvmon.service` | `systemd` unit file to manage the service. |

### ⚙️ Server Arguments  

- `-server` → Runs the application in **server mode**.
- `-conf=<path_to_conf_yaml>` → Specifies a custom configuration file (overrides default `/etc/srvmon/conf.yaml`).

---

## 🖥️ Client Side  

### 📁 File Locations  

| File          | Location                         | Description                                    |
|--------------|---------------------------------|------------------------------------------------|
| `srvmon`     | Any location in `PATH`         | Executable file for the client-side application. |
| `servers.yaml` | `$HOME/.config/srvmon/servers.yaml` | List of servers to be monitored. |

---

## 📌 Installation & Usage  

### 🔧 Server Setup  
1. **Place the binary in the correct location:**  
   ```sh
   sudo mv srvmon /usr/local/bin/
   sudo chmod +x /usr/local/bin/srvmon
   ```
2. **Create a configuration file:**  
   ```sh
   sudo mkdir -p /etc/srvmon
   sudo nano /etc/srvmon/conf.yaml
   ```
3. **Setup systemd service:**  
    
  `systemd` unit file:

  ```unit
  [Unit]
  Description=Hany Server Monitoring Daemon
  After=network.target

  [Service]
  ExecStart=/usr/local/bin/srvmon -server -conf=/etc/srvmon/conf.yaml
  Restart=always
  User=root
  WorkingDirectory=/usr/local/bin 
  LimitNOFILE=65535
  StandardOutput=journal
  StandardError=journal
  Environment="GODEBUG=madvdontneed=1"

  [Install]
  WantedBy=multi-user.target
  ```
    
   ```sh
   sudo mv srvmon.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable srvmon
   sudo systemctl start srvmon
   ```
4. **Check service status:**  
   ```sh
   sudo systemctl status srvmon
   ```

### 🔧 Client Setup  
1. **Move the binary to a location in `PATH`:**  
   ```sh
   sudo mv srvmon /usr/local/bin/
   sudo chmod +x /usr/local/bin/srvmon
   ```
2. **Create the servers list file:**  
   ```sh
   mkdir -p $HOME/.config/srvmon
   nano $HOME/.config/srvmon/servers.yaml
   ```
3. **Run the client:**  
   ```sh
   srvmon
   ```

---

## 📝 Notes  
- Ensure the server's configuration file is correctly formatted to avoid runtime errors.
- The client should have the correct list of servers in `servers.yaml` to function properly.
- Logs for the server can be accessed via:
  ```sh
  journalctl -u srvmon -f
  ```

---

## 🛠️ Troubleshooting  
### 🔹 Check Service Logs  
```sh
journalctl -u srvmon --no-pager | tail -n 50
```

### 🔹 Restart the Service  
```sh
sudo systemctl restart srvmon
```

### 🔹 Verify Binary is in `PATH`  
```sh
which srvmon
```

---

## 🐜 License  
This project is licensed under the [MIT License](LICENSE).

---

## 👨‍💻 Author  
Developed by **[Your Name]**.

