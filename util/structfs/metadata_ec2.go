package structfs

type EC2Metadata struct {
	Latest struct {
		Metadata struct {
			AMIID              int      `json:"ami-id"`
			AMILaunchIndex     int      `json:"ami-launch-index"`
			AMIManifestPath    string   `json:"ami-manifest-path"`
			AncestorAMIIDs     []int    `json:"ancestor-ami-ids"`
			BlockDeviceMapping []string `json:"block-device-mapping"`
			InstanceID         int      `json:"instance-id"`
			InstanceType       string   `json:"instance-type"`
			LocalHostname      string   `json:"local-hostname"`
			LocalIPv4          string   `json:"local-ipv4"`
			kernelID           int      `json:"kernel-id"`
			Placement          string   `json:"placement"`
			AvailabilityZone   string   `json:"availability-zone"`
			ProductCodes       string   `json:"product-codes"`
			PublicHostname     string   `json:"public-hostname"`
			PublicIPv4         string   `json:"public-ipv4"`
			PublicKeys         []struct {
				Key []string `json:"-"`
			} `json:"public-keys"`
			RamdiskID      int      `json:"ramdisk-id"`
			ReservationID  int      `json:"reservation-id"`
			SecurityGroups []string `json:"security-groups"`
		} `json:"meta-data"`
		Userdata string `json:"user-data"`
	} `json:"latest"`
}
