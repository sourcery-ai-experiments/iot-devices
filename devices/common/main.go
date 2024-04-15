package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kloudlite/iot-devices/constants"
	"github.com/kloudlite/iot-devices/pkg/conf"
	"github.com/kloudlite/iot-devices/pkg/k3s"
	"github.com/kloudlite/iot-devices/pkg/networkmanager"
	"github.com/kloudlite/iot-devices/types"
)

func getConfig(ip, token string) string {
	temp := `
runAs: primaryMaster
primaryMaster:
  publicIP: {{ip}}
  token: {{token}}
  nodeName: master-1
  labels: {"kloudlite.io/node.has-role":"primary-master","kloudlite.io/provider.name":"raspberry","kloudlite.io/release":"v1.0.5-nightly"}
  SANs: ["{{ip}}"]
  taints: ["node-role.kubernetes.io/master=:NoSchedule"]
  extraServerArgs: ["--disable-helm-controller","--disable","traefik","--disable","servicelb","--node-external-ip",{{ip}},"--flannel-external-ip","--cluster-domain","cluster.local","--kubelet-arg","--system-reserved=cpu=100m,memory=200Mi,ephemeral-storage=1Gi,pid=1000","--disable-agent"]
    `

	s := strings.ReplaceAll(temp, "{{ip}}", ip)
	s = strings.ReplaceAll(s, "{{token}}", token)

	return s
}

func StartPing(ctx types.MainCtx) {
	for {
		if err := ping(ctx); err != nil {
			ctx.GetLogger().Errorf(err, "sending ping to server")
		}

		time.Sleep(constants.PingInterval * time.Second)
	}
}

func ping(ctx types.MainCtx) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	c, err := conf.GetConf()
	if err != nil {
		return err
	}

	var data = struct {
		PublicKey string `json:"publicKey"`
	}{
		// PublicKey: c.PublicKey,
		PublicKey: "10.2.2.2",
	}

	dataBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	resp, err := client.Post(constants.GetPingUrl(), "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response types.Response

	// read all the response body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)

	ctx.GetLogger().Infof("Ping response: %s", buf.String())
	if resp.StatusCode == http.StatusOK {

		if err := response.FromJson(buf.Bytes()); err != nil {
			return err
		}

		if response.Reset {
			ctx.GetLogger().Infof("Resetting device")
			return k3s.New(ctx).Reset()
		}

		ctx.UpdateDevice(&response)

		ip, err := networkmanager.GetIfIp()
		if err != nil {
			return err
		}

		conf := getConfig(ip, string(c.PrivateKey))
		if err := k3s.New(ctx).UpsertConfig(conf); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("status code: %d", resp.StatusCode)
}
