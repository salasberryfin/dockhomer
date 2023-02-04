package network

import (
	"log"
	"net"

	"github.com/vishvananda/netlink"
)

// NewBridge creates a new bridge device in the host machine
func NewBridge(name string) {
	linkAttributes := netlink.NewLinkAttrs()
	linkAttributes.Name = name
	link := &netlink.Bridge{
		LinkAttrs: linkAttributes,
	}
	log.Printf("Creating network bridge %v\n", link)
}

// NewVethPair creates a veth interface pair for establishing connection between
// the host and the container through the bridge device:
//   - veth1 is attached to the bridge
//   - veth2 is attached to the container's namespace
func NewVethPair(name string) {
	if ifaceExists(name) {
		log.Printf("%s already exists\n", name)
		return
	}
	log.Printf("%s will be created\n", name)
}

// attach links a network interface to a bridge device or network namespace
func attach() {}

// generateIfaceName creates a network interface that does not exist yet
func generateIfaceName() {}

// ifaceExists checks if given name matches an existing network interface
func ifaceExists(iface string) bool {
	_, err := net.InterfaceByName(iface)
	if err != nil {
		return false
	}

	return true
}
