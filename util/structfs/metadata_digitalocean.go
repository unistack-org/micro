package structfs

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type DigitalOceanMetadata struct {
	Metadata struct {
		V1 struct {
			Features map[string]interface{} `json:"features"`

			Hostname   string `json:"hostname"`
			VendorData string `json:"vendor_data"`
			Region     string `json:"region"`

			PublicKeys []string `json:"public_keys"`

			DNS struct {
				Nameservers []string `json:"nameservers"`
			} `json:"dns"`

			Interfaces struct {
				Private []struct {
					IPv4 struct {
						Address string `json:"ip_address"`
						Netmask string `json:"netmask"`
						Gateway string `json:"gateway"`
					}
					Mac  string `json:"mac"`
					Type string `json:"type"`
				} `json:"private"`
				Public []struct {
					IPv4 struct {
						Address string `json:"ip_address"`
						Netmask string `json:"netmask"`
						Gateway string `json:"gateway"`
					} `json:"ipv4"`
					Mac  string `json:"mac"`
					Type string `json:"type"`
					IPv6 struct {
						Address string `json:"ip_address"`
						Gateway string `json:"gateway"`
						CIDR    int    `json:"cidr"`
					} `json:"ipv6"`
				} `json:"public"`
			} `json:"interfaces"`

			DropletID int64 `json:"droplet_id"`

			FloatingIP struct {
				IPv4 struct {
					Active bool `json:"active"`
				} `json:"ipv4"`
			} `json:"floating_ip"`
		} `json:"v1"`
	} `json:"metadata"`
}

func (stfs *DigitalOceanMetadata) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/metadata/v1.json":
		_ = json.NewEncoder(w).Encode(stfs.Metadata.V1)
	default:
		fs := FileServer(stfs, "json", time.Now())
		idx := strings.Index(r.URL.Path[1:], "/")
		r.URL.Path = strings.Replace(r.URL.Path[idx+1:], "/metadata/v1/", "", 1)
		r.RequestURI = r.URL.Path
		fs.ServeHTTP(w, r)
	}
}
