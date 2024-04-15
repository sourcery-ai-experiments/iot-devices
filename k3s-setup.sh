#! /usr/bin/env bash

[ $EUID -ne 0 ] && echo "this script must be run as root. current EUID is $EUID"

arch=$(uname -m)
if [ "$arch" != "x86_64" ] && [ "$arch" != "aarch64" ]; then
  echo "CPU Architecture is not x86_64 or aarch64. Exiting ..." && exit 1
fi

if [ "$arch" == "x86_64" ]; then
  RUNNER_NAME="runner-amd64"
  K3S_BINARY="k3s"
elif [ "$arch" == "aarch64" ]; then
  RUNNER_NAME="runner-arm64"
  K3S_BINARY="k3s-arm64"
fi

#sleep 30s

kloudlite_release=v1.0.3

curl -L0 "https://github.com/kloudlite/infrastructure-as-code/releases/download/$kloudlite_release/$K3S_BINARY" >/usr/local/bin/k3s
chmod +x /usr/local/bin/k3s

curl -L0 "https://github.com/kloudlite/infrastructure-as-code/releases/download/$kloudlite_release/$RUNNER_NAME" >./runner
chmod +x ./runner

cat >/etc/systemd/system/kloudlite-k3s.service <<EOF
[Unit]
Description=This script will start kloudlite k3s runner. It is maintained by kloudlite.io, and is used to run k3s with a custom set of args.

[Service]
ExecStart=$PWD/runner --config $PWD/runner-config.yml

[Install]
WantedBy=multi-user.target
EOF

systemctl enable --now kloudlite-k3s.service
