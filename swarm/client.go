package swarm

import (
  //"net"
  "strings"

  "github.com/docker/docker/api/types"
  "github.com/docker/docker/api/types/swarm"
  dockerclient "github.com/docker/docker/client"
  "golang.org/x/net/context"

)

type Client interface {
  ListActiveNodeIPs() ([]string, error)
}

type swarmClient struct {
  api *dockerclient.Client
}

func NewClient() (Client, error) {
  cli, err := dockerclient.NewEnvClient()
  if err != nil {
    return nil, err
  }

  return swarmClient{api: cli}, nil
}

func (client swarmClient) ListActiveNodeIPs() ([]string, error) {
  var listOptions types.NodeListOptions
  nodes, err := client.api.NodeList(context.Background(), listOptions)
  if err != nil {
    return nil, err
  }

  var ips []string
  for _, node := range nodes {
    if node.Status.State == swarm.NodeStateReady {
      if node.Status.Addr == "0.0.0.0" {

        leaderIp := getIPFromAddr(node.ManagerStatus.Addr)
        if err != nil {
          return nil, err
        }
        ips = append(ips, leaderIp)
      } else {
        ips = append(ips, node.Status.Addr)
      }
    }
  }


  return ips, nil
}

// Get preferred outbound ip of this machine
// (see https://stackoverflow.com/a/37382208/3515146)
//func getOutboundIP() (string, error) {
    //conn, err := net.Dial("udp", "8.8.8.8:80")
    //if err != nil {
      //return "", err
    //}
    //defer conn.Close()

    //localAddr := conn.LocalAddr().String()
    //idx := strings.LastIndex(localAddr, ":")

    //return localAddr[0:idx], nil
//}

func getIPFromAddr(addr string) string {
  ipAndPort := strings.Split(addr, ":")
  return ipAndPort[0]
}
