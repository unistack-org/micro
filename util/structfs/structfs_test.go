package structfs

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

var doOrig = []byte(`{
  "droplet_id":2756294,
  "hostname":"sample-droplet",
  "vendor_data":"#cloud-config\ndisable_root: false\nmanage_etc_hosts: true\n\ncloud_config_modules:\n - ssh\n - set_hostname\n - [ update_etc_hosts, once-per-instance ]\n\ncloud_final_modules:\n - scripts-vendor\n - scripts-per-once\n - scripts-per-boot\n - scripts-per-instance\n - scripts-user\n",
  "public_keys":["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi6cygCUmuNlB0KqzBpHXf7CFYb3VE4pDOf/RLJ8OFDjOM+fjF83a24QktSVIpQnHYpJJT2pQMBxD+ZmnhTbKv+OjwHSHwAfkBullAojgZKzz+oN35P4Ea4J78AvMrHw0zp5MknS+WKEDCA2c6iDRCq6/hZ13Mn64f6c372JK99X29lj/B4VQpKCQyG8PUSTFkb5DXTETGbzuiVft+vM6SF+0XZH9J6dQ7b4yD3sOder+M0Q7I7CJD4VpdVD/JFa2ycOS4A4dZhjKXzabLQXdkWHvYGgNPGA5lI73TcLUAueUYqdq3RrDRfaQ5Z0PEw0mDllCzhk5dQpkmmqNi0F sammy@digitalocean.com"],
  "region":"nyc3",
  "interfaces":{
    "private":[
      {
        "ipv4":{
          "ip_address":"10.132.255.113",
          "netmask":"255.255.0.0",
          "gateway":"10.132.0.1"
        },
        "mac":"04:01:2a:0f:2a:02",
        "type":"private"
      }
    ],
    "public":[
      {
        "ipv4":{
          "ip_address":"104.131.20.105",
          "netmask":"255.255.192.0",
          "gateway":"104.131.0.1"
        },
        "ipv6":{
          "ip_address":"2604:A880:0800:0010:0000:0000:017D:2001",
          "cidr":64,
          "gateway":"2604:A880:0800:0010:0000:0000:0000:0001"
        },
        "mac":"04:01:2a:0f:2a:01",
        "type":"public"}
    ]
  },
  "floating_ip": {
    "ipv4": {
      "active": false
    }
  },
  "dns":{
    "nameservers":[
      "2001:4860:4860::8844",
      "2001:4860:4860::8888",
      "8.8.8.8"
    ]
  },
  "features":{
    "dhcp_enabled": true
  }
}
`)

func server(t *testing.T, ch chan error) {
	stfs := DigitalOceanMetadata{}
	err := json.Unmarshal(doOrig, &stfs.Metadata.V1)
	if err != nil {
		t.Fatal(err)
	}

	http.Handle("/metadata/v1/", FileServer(&stfs, "json", time.Now()))
	http.Handle("/metadata/v1.json", &stfs)
	go func() {
		ch <- http.ListenAndServe("127.0.0.1:8080", nil)
	}()
	time.Sleep(2 * time.Second)
}

func get(path string) ([]byte, error) {
	res, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	return io.ReadAll(res.Body)
}

func TestAll(t *testing.T) {
	ch := make(chan error)
	server(t, ch)

	tests := []struct {
		in  string
		out string
	}{
		{"http://127.0.0.1:8080/metadata/v1/", "features\nhostname\nvendor_data\nregion\npublic_keys\ndns\ninterfaces\ndroplet_id\nfloating_ip"},
		{"http://127.0.0.1:8080/metadata/v1/droplet_id", "2756294"},
		{"http://127.0.0.1:8080/metadata/v1/dns/", "nameservers"},
		{"http://127.0.0.1:8080/metadata/v1/dns/nameservers", "2001:4860:4860::8844\n2001:4860:4860::8888\n8.8.8.8"},
		{"http://127.0.0.1:8080/metadata/v1/features/dhcp_enabled", "true"},
	}

	for _, tt := range tests {
		select {
		case err := <-ch:
			t.Fatal(err)
		default:
			buf, err := get(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if string(buf) != tt.out {
				t.Errorf("req %s output %s not match requested %s", tt.in, string(buf), tt.out)
			}
		}
	}

	select {
	case err := <-ch:
		t.Fatal(err)
	default:
		doTest, err := get("http://127.0.0.1:8080/metadata/v1.json")
		if err != nil {
			t.Fatal(err)
		}

		oSt := DigitalOceanMetadata{}
		err = json.Unmarshal(doOrig, &oSt.Metadata.V1)
		if err != nil {
			t.Fatal(err)
		}

		nSt := DigitalOceanMetadata{}

		err = json.Unmarshal(doTest, &nSt.Metadata.V1)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(oSt, nSt) {
			t.Fatalf("%v not match %v", oSt, nSt)
		}
	}
}

func TestFileInfo_Methods(t *testing.T) {
	now := time.Now()
	fi := &fileInfo{name: "test.txt", size: 42, modtime: now}

	if fi.Sys() != nil {
		t.Error("expected Sys() == nil")
	}
	if fi.Size() != 42 {
		t.Errorf("want size 42, got %d", fi.Size())
	}
	if fi.Name() != "test.txt" {
		t.Errorf("want name 'test.txt', got %q", fi.Name())
	}
	if fi.Mode() != os.FileMode(0o644) {
		t.Errorf("want mode 0644, got %v", fi.Mode())
	}
	if fi.IsDir() {
		t.Error("expected IsDir() == false")
	}
	if !fi.ModTime().Equal(now) {
		t.Errorf("want modtime %v, got %v", now, fi.ModTime())
	}
}

func TestFileInfo_DirMode(t *testing.T) {
	fi := &fileInfo{name: "dir/", size: 0, modtime: time.Now()}
	if fi.Mode()&os.ModeDir == 0 {
		t.Error("expected directory mode for name ending in '/'")
	}
}

func TestFile_CloseReadSeekStatReaddir(t *testing.T) {
	now := time.Now()
	f := &file{name: "f.txt", data: []byte("hello world"), modtime: now}

	// Close
	if err := f.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Read
	buf := make([]byte, 5)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read: %v", err)
	}
	if n == 0 {
		t.Error("expected bytes read > 0")
	}

	// Seek to start and read again to exhaust
	_, _ = f.Seek(0, io.SeekStart)
	full := make([]byte, 100)
	n, err = f.Read(full)
	if err != nil && err != io.EOF {
		t.Fatalf("Read full: %v", err)
	}
	if string(full[:n]) != "hello world" {
		t.Errorf("want 'hello world', got %q", string(full[:n]))
	}

	// Read past EOF
	n, err = f.Read(buf)
	if err != io.EOF {
		t.Errorf("expected EOF after data exhausted, got err=%v n=%d", err, n)
	}

	// Seek SeekCurrent
	_, _ = f.Seek(0, io.SeekStart)
	_, _ = f.Seek(2, io.SeekCurrent)

	// Seek SeekEnd
	off, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		t.Fatalf("Seek SeekEnd: %v", err)
	}
	if off != int64(len(f.data)) {
		t.Errorf("want offset %d, got %d", len(f.data), off)
	}

	// Stat
	info, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Name() != "f.txt" {
		t.Errorf("Stat name: want 'f.txt', got %q", info.Name())
	}

	// Readdir
	entries, err := f.Readdir(0)
	if err != nil {
		t.Fatalf("Readdir: %v", err)
	}
	if entries != nil {
		t.Error("expected nil from Readdir")
	}
}

func TestFileServer_ServeHTTP_BadPath(t *testing.T) {
	stfs := DigitalOceanMetadata{}
	_ = json.Unmarshal(doOrig, &stfs.Metadata.V1)

	h := FileServer(&stfs, "json", time.Now())

	// request a path that doesn't exist in the struct — should 500
	req := httptest.NewRequest(http.MethodGet, "/nonexistent_field", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", w.Code)
	}
}

func TestFileServer_ZeroModtime(t *testing.T) {
	stfs := DigitalOceanMetadata{}
	_ = json.Unmarshal(doOrig, &stfs.Metadata.V1)
	// zero modtime should default to time.Now() inside FileServer
	h := FileServer(&stfs, "json", time.Time{})
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}
