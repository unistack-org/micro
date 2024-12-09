package structfs

type EC2Metadata struct {
	Latest struct {
		Userdata string `json:"user-data"`
		Metadata struct {
			AMIManifestPath  string `json:"ami-manifest-path"`
			InstanceType     string `json:"instance-type"`
			LocalHostname    string `json:"local-hostname"`
			LocalIPv4        string `json:"local-ipv4"`
			Placement        string `json:"placement"`
			AvailabilityZone string `json:"availability-zone"`
			ProductCodes     string `json:"product-codes"`
			PublicHostname   string `json:"public-hostname"`
			PublicIPv4       string `json:"public-ipv4"`
			PublicKeys       []struct {
				Key []string `json:"-"`
			} `json:"public-keys"`
			AncestorAMIIDs     []int    `json:"ancestor-ami-ids"`
			BlockDeviceMapping []string `json:"block-device-mapping"`
			SecurityGroups     []string `json:"security-groups"`
			RamdiskID          int      `json:"ramdisk-id"`
			ReservationID      int      `json:"reservation-id"`
			AMIID              int      `json:"ami-id"`
			AMILaunchIndex     int      `json:"ami-launch-index"`
			KernelID           int      `json:"kernel-id"`
			InstanceID         int      `json:"instance-id"`
		} `json:"meta-data"`
	} `json:"latest"`
}
