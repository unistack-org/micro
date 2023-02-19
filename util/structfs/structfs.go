package structfs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

// FileServer creates new file server from the struct iface with specific tag and specific modtime
func FileServer(iface interface{}, tag string, modtime time.Time) http.Handler {
	if modtime.IsZero() {
		modtime = time.Now()
	}
	return &fs{iface: iface, tag: tag, modtime: modtime}
}

func (fs *fs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	f, err := fs.Open(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, r.URL.Path, fs.modtime, f)
}

type fs struct {
	iface   interface{}
	tag     string
	modtime time.Time
}

type file struct {
	name    string
	offset  int64
	data    []byte
	modtime time.Time
}

type fileInfo struct {
	name    string
	size    int64
	modtime time.Time
}

func (fi *fileInfo) Sys() interface{} {
	return nil
}

func (fi *fileInfo) Size() int64 {
	return fi.size
}

func (fi *fileInfo) Name() string {
	return fi.name
}

func (fi *fileInfo) Mode() os.FileMode {
	if strings.HasSuffix(fi.name, "/") {
		return os.FileMode(0755) | os.ModeDir
	}
	return os.FileMode(0644)
}

func (fi *fileInfo) IsDir() bool {
	// disables additional open /index.html
	return false
}

func (fi *fileInfo) ModTime() time.Time {
	return fi.modtime
}

func (f *file) Close() error {
	return nil
}

func (f *file) Read(b []byte) (int, error) {
	var err error
	var n int

	if f.offset >= int64(len(f.data)) {
		return 0, io.EOF
	}

	if len(f.data) > 0 {
		n = copy(b, f.data[f.offset:])
	}

	if n < len(b) {
		err = io.EOF
	}

	f.offset += int64(n)
	return n, err
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	//	log.Printf("seek %d %d %s\n", offset, whence, f.name)
	switch whence {
	case os.SEEK_SET:
		f.offset = offset
	case os.SEEK_CUR:
		f.offset += offset
	case os.SEEK_END:
		f.offset = int64(len(f.data)) + offset
	}
	return f.offset, nil

}

func (f *file) Stat() (os.FileInfo, error) {
	return &fileInfo{name: f.name, size: int64(len(f.data)), modtime: f.modtime}, nil
}

func (fs *fs) Open(path string) (http.File, error) {
	return newFile(path, fs.iface, fs.tag, fs.modtime)
}

func newFile(name string, iface interface{}, tag string, modtime time.Time) (*file, error) {
	var err error

	f := &file{name: name, modtime: modtime}
	f.data, err = structItem(name, iface, tag)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func structItem(path string, iface interface{}, tag string) ([]byte, error) {
	var buf []byte
	var err error
	var curiface interface{}

	if path == "/" {
		return getNames(iface, tag)
	}

	idx := strings.Index(path[1:], "/")
	switch {
	case idx > 0:
		curiface, err = getStruct(path[1:idx+1], iface, tag)
		if err != nil {
			return nil, err
		}
		buf, err = structItem(path[idx+1:], curiface, tag)
	case idx == 0:
		return getNames(iface, tag)
	case idx < 0:
		return getValue(path[1:], iface, tag)
	}

	return buf, err
}

func getNames(iface interface{}, tag string) ([]byte, error) {
	var lines []string
	s := reflectValue(iface)
	typeOf := s.Type()
	for i := 0; i < s.NumField(); i++ {
		value := typeOf.Field(i).Tag.Get(tag)
		if value != "" {
			lines = append(lines, value)
		}
	}
	if len(lines) > 0 {
		return []byte(strings.Join(lines, "\n")), nil
	}
	return nil, fmt.Errorf("failed to find names")
}

func getStruct(name string, iface interface{}, tag string) (interface{}, error) {
	s := reflectValue(iface)
	typeOf := s.Type()
	for i := 0; i < s.NumField(); i++ {
		if typeOf.Field(i).Tag.Get(tag) == name {
			return s.Field(i).Interface(), nil
		}
	}
	return nil, fmt.Errorf("failed to find iface %T with name %s", iface, name)
}

func getValue(name string, iface interface{}, tag string) ([]byte, error) {
	s := reflectValue(iface)
	typeOf := s.Type()
	switch typeOf.Kind() {
	case reflect.Map:
		return []byte(fmt.Sprintf("%v", s.MapIndex(reflect.ValueOf(name)).Interface())), nil
	default:
		for i := 0; i < s.NumField(); i++ {
			if typeOf.Field(i).Tag.Get(tag) != name {
				continue
			}
			ifs := s.Field(i).Interface()
			switch s.Field(i).Kind() {
			case reflect.Slice:
				var lines []string
				for k := 0; k < s.Field(i).Len(); k++ {
					lines = append(lines, fmt.Sprintf("%v", s.Field(i).Index(k)))
				}
				return []byte(strings.Join(lines, "\n")), nil
			default:
				return []byte(fmt.Sprintf("%v", ifs)), nil
			}
		}
	}
	return nil, fmt.Errorf("failed to find %s in interface %T", name, iface)
}

func hasValidType(obj interface{}, types []reflect.Kind) bool {
	for _, t := range types {
		if reflect.TypeOf(obj).Kind() == t {
			return true
		}
	}

	return false
}

func reflectValue(obj interface{}) reflect.Value {
	var val reflect.Value

	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		val = reflect.ValueOf(obj).Elem()
	} else {
		val = reflect.ValueOf(obj)
	}

	return val
}
