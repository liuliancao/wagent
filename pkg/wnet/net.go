package wnet

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func GetPrivateIPs() (private_ips []string) {
	private_ips = make([]string, 0)

	ifaces, e := net.Interfaces()
	if e != nil {
		log.Println("get net.Interfaces error: ", e.Error())
		return private_ips
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		fmt.Println("iface name:", iface.Name, "iface mac", iface.HardwareAddr, "iface flags", iface.Flags)
		ipaddrs, e := iface.Addrs()

		if e != nil {
			log.Println("get iface addrs error: ", e.Error())
			return private_ips
		}
		for _, addr := range ipaddrs {
			ipstr := strings.Split(addr.String(), "/")[0]

			ip := net.ParseIP(ipstr)
			// added in go1.17
			if ip.IsPrivate() {
				if ip.To4() != nil {
					private_ips = append(private_ips, ip.String())
				}

			}
		}
	}
	return private_ips
}
func GetPublicIPs() (public_ips []string) {
	public_ips = make([]string, 0)

	ifaces, e := net.Interfaces()
	if e != nil {
		log.Println("get net.Interfaces error: ", e.Error())
		return public_ips
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		ipaddrs, e := iface.Addrs()

		if e != nil {
			log.Println("get iface addrs error: ", e.Error())
			return public_ips
		}
		for _, addr := range ipaddrs {
			ipstr := strings.Split(addr.String(), "/")[0]

			ip := net.ParseIP(ipstr)
			// added in go1.17
			if !ip.IsPrivate() {
				if ip.To4() != nil {
					public_ips = append(public_ips, ip.String())
				}
			}
		}
	}
	return public_ips
}
func GetAllIPs() (ips []string) {
	ips = make([]string, 0)

	ifaces, e := net.Interfaces()
	if e != nil {
		log.Println("get net.Interfaces error: ", e.Error())
		return ips
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		ipaddrs, e := iface.Addrs()

		if e != nil {
			log.Println("get iface addrs error: ", e.Error())
			return ips
		}
		for _, addr := range ipaddrs {
			ipstr := strings.Split(addr.String(), "/")[0]

			ip := net.ParseIP(ipstr)
			if ip.To4() != nil {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips
}
