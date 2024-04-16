1. Execute the following command to enable ip forwarding on each node

```bash
echo 1 | sudo tee /proc/sys/net/ipv4/ip_forward
```

2. Add this service to `/etc/systemd/system/broad.service`
```bash
Description=My Go Program
After=network.target

[Service]
Type=simple
User=root
ExecStart=/bin/bash /home/raspberry/startup.sh
Restart=always
RestartSec=4

[Install]
WantedBy=multi-user.target
```

3. Add this script to `/home/raspberry/startup.sh`
```bash
#!/bin/bash
cd /home/raspberry && ./server
```

4. Add `./server` app to /home/raspberry

5. Add the following configuration `/home/raspberry/kloudlite-conf.yaml`
```yaml
privateKey: "your_private_key"
publicKey: "your_public_key"
```
