package skywalking

import "fmt"

type Conf struct {
	ServerAddr       string
	ServiceName      string
	ServiceInstance  string
	ServiceNamespace string
}

func (c *Conf) String() string {
	return fmt.Sprintf("{ServerAddr: %s, ServiceName: %s, ServiceInstance: %s, ServiceNamespace: %s}",
		c.ServerAddr, c.ServiceName, c.ServiceInstance, c.ServiceNamespace)
}
