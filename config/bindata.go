// Code generated for package config by go-bindata DO NOT EDIT. (@generated)
// sources:
// hsdp.toml
package config

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _hsdpToml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x58\x4d\x6f\xdb\x38\x10\x3d\xaf\x7e\xc5\x6c\x7c\x5d\xbb\xe8\xe6\x16\xa0\x07\x37\x31\x52\xa3\xdd\xa4\xb0\x53\x14\x8b\x45\xb1\x18\x53\x23\x89\x88\x44\x12\x24\x65\xc7\xff\x7e\x41\x7d\xc4\xb2\x25\x59\x54\x62\x60\x6f\x36\xc5\x79\x33\xef\xf1\x71\x48\x69\x02\x4f\x09\x37\xc0\x0d\x20\xc4\xa9\xdc\x60\x0a\x9a\x62\x6e\xac\xde\x43\x24\x35\xa8\x7c\x93\x72\x06\x5f\xd6\x77\xdf\xc1\x90\xde\x72\x46\x26\x98\x00\x93\xc2\x22\x17\x5c\xc4\x90\x49\x63\xc1\x58\xb4\x9c\xb9\xe1\x88\xc7\xb9\x46\xcb\xa5\x98\xc1\xd2\x1a\xc8\x90\x0b\x50\xb9\x56\xd2\x90\x4b\x63\x25\x04\x13\xd0\x14\xe6\x8c\x8e\xe7\x83\x42\x8d\x19\x59\xd2\x06\x50\x84\x10\x6a\xbe\x25\xc0\xdc\xca\xac\x7c\x2e\x23\xb0\x09\xb9\x70\x95\xa2\x8d\xa4\xce\x20\xca\xb5\x4d\x48\xc3\x66\x0f\x09\x6e\x5d\x39\x08\x49\x9e\xa1\x28\x10\x32\x64\x09\x17\x04\x9a\x30\xc4\x4d\x4a\x60\x64\xae\x19\x81\x8c\x0a\x0a\x59\x26\x45\x6f\x09\xc1\xc4\x25\xfa\x4c\x91\xd4\x04\xb9\x71\xd0\xdc\x02\x1a\xd8\xcb\x5c\x1f\x80\xc0\xea\xdc\x26\xa0\x52\x42\x53\xf1\xd1\x99\x9b\x99\xa0\x83\x70\xf5\xd2\x8b\x22\x66\x29\x84\x2d\xa6\x39\x99\x59\x30\x09\x26\x30\x4f\x8d\x74\xd3\x62\x49\x06\x76\xdc\x26\x32\xb7\x60\x70\xcf\x45\xfc\x07\x6c\x72\x0b\x3b\x82\x1d\x4f\x53\x37\x56\x24\x16\xfb\x1d\xee\x6f\x82\x09\x2c\xad\x93\xb1\xa8\x42\x93\x51\x52\x18\xbe\xe1\xa9\xdd\x3b\x65\x49\x98\x5c\xbb\x32\xb4\x26\x66\x4f\xb8\xc9\xc8\x31\x2a\x02\x43\x52\xa9\xdc\x67\x24\x6c\x55\xce\x32\x72\x0f\xc0\x28\x69\x81\xb4\x96\xda\x80\xd4\xc0\x05\x93\x99\x4a\xc9\x12\x70\xe1\xf4\xae\x44\x2a\xc9\x6a\xe4\x86\x00\x1d\xe8\xf7\xd5\x0c\x9e\x12\x14\xcf\xe6\xf7\x02\xee\x8e\x22\x2e\xb8\x9b\x6c\x6e\x8a\x81\x15\xc5\x2e\x72\xea\x5c\x46\x32\xd6\xa8\x12\xce\x2a\xaf\x49\x01\xbb\x84\x34\xbd\xfa\x0b\x50\x53\x55\x22\x85\xc1\x04\x16\x62\xcb\xb5\x14\xae\xdc\x02\xc1\x58\x8c\xa9\x10\x8d\x8b\x1e\x40\x2b\x01\x79\x08\xce\x7a\x09\x1a\x0a\x1b\x8c\x5d\xbd\xeb\x32\x53\x81\xc6\xa4\x30\x79\x56\xd9\x23\xb2\x3b\x3c\x54\xe2\x84\x11\xc6\xa2\xa8\xa6\x1a\x45\x8c\x47\x9c\x35\xd1\x64\xe4\x1e\xd4\xf3\x83\x09\x7c\xe1\xa4\x51\xb3\x64\x0f\xf4\x82\x4e\xbd\x52\x81\x05\x6a\x9b\x04\x13\x80\x5a\x8b\xdc\x4c\x09\x8d\x2d\x86\xa0\xae\xa8\xfa\x07\xc0\xa2\xea\x67\x93\x3c\x4b\x39\x09\x3b\xb5\xf4\x1a\xd6\x0e\x04\xe0\x98\x05\x93\xdf\x00\x78\x98\x35\x46\x19\x6a\x4b\xe9\x31\x28\x65\xc2\x82\xd2\x32\xcc\x99\x5b\xab\x21\xcc\xc3\xbf\x7e\xe4\x62\x20\xd4\xaf\xff\x9a\x8c\x29\x9f\xee\x0e\xa5\xbf\x91\x58\x4f\x11\xe5\x76\xbb\x38\x37\x73\xcd\x74\xe1\xc1\x60\x02\x0f\x64\xac\x6b\x03\x66\x2f\x2c\xbe\xdc\xc0\x3f\xa5\xd7\x66\x74\x20\x32\xab\x8c\x30\xe3\x95\x6d\x7e\x15\x6b\xbf\xae\x8d\xcd\x50\xc0\x86\xc0\x72\x0a\x0b\x8b\xd6\x76\x95\x1a\x50\x40\x03\x28\x98\xc0\x4f\x02\x74\x4d\x02\x85\xe5\x8c\x2b\xb4\x54\x37\xe8\xc3\x3e\xb1\x60\x64\x46\xa0\x24\x2f\x42\x16\xa5\xe1\x4a\x4f\x56\xd0\x1b\x99\x8b\xb0\x0e\xb9\x81\xdb\x54\xe6\x21\x44\x6e\x50\xef\x4f\x43\x48\xf0\x6d\xbd\x24\x27\x71\xcb\xf9\x5f\x41\x63\xe3\x1c\xc3\xd4\x4a\x54\x96\x7e\x55\x81\x45\xbf\x82\x5c\xa7\xf0\x09\xae\x12\x6b\x95\xb9\xf9\xf0\x01\x15\x9f\x31\x17\x3c\x53\x2c\x72\xeb\x3d\x63\x32\xbb\x0a\x42\x59\x9c\x13\x9f\xe0\xaa\xc6\x50\x09\x4f\xb9\x32\xd3\x84\x30\xb5\x89\xc9\xb9\xa5\x72\xea\xab\xec\xa5\x9b\x86\x92\x51\xfe\x71\xa6\x12\x13\xaa\xd3\x44\x75\xfc\x60\x22\x83\x1f\x87\x92\x98\xbd\x29\xa6\x15\x79\xb8\x6c\xa6\x39\x1a\xae\x21\x51\x21\xbb\xf6\x01\x45\x75\xdd\x05\x7a\x34\xdc\x58\x96\xbb\xc7\xdb\xaf\x8b\xd5\x74\xb5\xb8\x5f\xae\x9f\x56\x7f\xf7\x2e\x4c\x28\xd9\x33\xe9\x69\x7d\xd2\xff\x0a\x12\x77\x88\x7f\x82\xab\xf2\xc1\x4c\x74\x15\x7d\x2a\xf8\x10\x88\x53\xbe\x05\xd2\x14\x73\x08\xc0\x43\xba\x21\x88\x3e\xa1\x6e\x1f\x1f\xd6\x8f\xdf\x16\xfd\xce\x95\xc2\xc8\x94\x5a\xeb\x52\x8d\xfb\x09\x34\x04\x32\x28\xd0\x10\x80\x8f\xb7\x06\x20\xfa\x04\xfa\x31\x9f\xf7\x8a\x93\x23\xb6\xf0\x72\xc4\xae\x7d\xdd\x27\x4d\x1f\xc4\xc9\x6e\xed\x12\xa5\x2f\xb4\xb5\x07\xbb\x05\x39\x17\xde\x27\x86\x6b\x7e\xa7\x62\x90\xd8\xce\x1a\xc7\xd6\xa1\xf3\x63\xd6\xc2\xe7\x98\x4d\x9b\x73\x07\x7b\x1c\xc7\xec\xdf\x77\x63\x84\x6d\x8c\x70\x24\x46\x17\x69\x77\xa6\x0e\xb2\x7d\x55\xfb\x8d\x4c\xfd\xe3\xbb\x59\x7a\xc7\x9f\x1a\xf4\x3d\xcb\x3a\x78\xa2\x0c\x90\x1d\x8e\x1f\x5e\x52\xef\x53\xad\x49\x78\xd4\x92\xfe\x5f\x2c\xbd\xe3\x9b\x4d\xe3\x3d\xcb\x79\xd4\x4c\x3c\x76\xe4\xf1\xfc\xe1\xa5\xea\x6c\x56\x75\xd1\xa3\x96\xc4\xa7\xd2\xee\xb9\xe7\xa5\x3e\xd3\x4e\x47\xd7\x78\xd4\x5b\x07\x6a\x3c\x9e\x7b\xbe\xc6\xa3\xb9\x5d\xfd\x2a\xa4\xed\x60\x91\x21\x6d\x29\x95\xaa\xb8\xc2\xbf\xb5\x65\x8d\xc3\xe8\x26\x35\x0a\xa3\x79\x40\xdd\x8d\x38\xa0\xc2\x0e\x0d\x2e\x7f\x30\xf4\x64\xb9\x7c\x63\xf6\xa0\x73\x81\xa6\x38\x40\xe7\x32\x4d\xc9\x83\xca\x88\xa6\x31\x50\xf2\xa8\xcd\x3d\x80\x35\x6e\x13\xf6\x80\x8d\x33\x7f\xf3\x32\x3f\x5f\x3d\x2d\xbe\xf5\xdf\xe5\x8b\xef\x13\x87\x57\x83\xf2\xff\xd4\x5d\xe1\xab\x1b\x6b\xfb\xc2\xd9\xba\xca\x77\x63\xb8\xeb\x6a\x2f\xc6\xd1\x4d\xbe\x3b\xde\x74\xd5\xd0\x60\x76\x3f\x7f\x5a\xfc\x9c\xf7\xbf\xc7\xc5\x68\x69\x87\x8d\xd7\x9e\x78\x57\xf0\xea\x7c\xe1\x3d\x79\xd0\x47\xb5\x0b\xb3\xff\x25\x3a\xa3\x6e\xd4\x26\xf9\x2e\x44\xd3\x57\xa5\xe9\xc6\x3b\xbe\xc4\x77\x21\x96\x33\x3a\x31\x5b\x8f\x1a\x12\xaf\xaf\x6f\x57\x8b\xbb\xb5\x77\xfb\x2c\xbf\x08\x99\x96\x83\xab\xf1\x0b\xb7\xd1\xa1\x6c\x17\x6f\xa7\x63\xe8\xbd\xbf\xad\xfa\xd2\x1b\xcc\xd4\x58\xd0\x2f\xf3\xb5\x37\xdb\x04\xdb\xb9\x13\x1c\xc9\xb2\xd9\x89\xee\x56\xd3\xf5\xd3\x8f\xeb\xce\x85\x35\x28\xc2\x8d\x7c\x39\xf4\x84\x50\x4f\x8d\xcd\xaf\xdb\xef\xe6\xd5\x83\x69\x1d\x31\xe6\x2a\xf0\xed\xf1\xfe\x7e\xf9\x70\xef\x2d\x41\x2a\xe3\x98\x8b\xb8\x55\x44\x2a\xdd\x30\x19\x2b\xf5\x9f\x17\x5e\x75\x9f\x94\xfe\x69\x86\x36\xec\x68\x82\xef\xdf\xb5\x5e\x04\xbd\xd3\x9c\xbb\x32\x8c\x26\xe7\x7f\x7f\xf0\x22\x31\xe6\x12\xe1\x05\xd8\xf7\x09\xe6\xeb\xf2\xf3\xfc\x61\xee\xed\xea\x67\xbe\x41\xd1\xfe\xd0\x53\x0e\x5f\xd8\xce\x67\x73\x5d\xce\xc7\xfe\x94\xde\x6f\xe0\xf3\x94\x2e\xe2\x5c\x7f\x3a\xfe\x96\x3d\x5f\xf6\x18\xaf\x9e\x47\xea\x33\xe9\xf7\xaf\x4b\x6f\x87\xaa\x67\xde\x82\x57\xcf\x7c\xaa\xb4\x7c\xd9\x8f\x72\xe8\x7f\x01\x00\x00\xff\xff\x29\x98\x7f\x62\xc3\x1f\x00\x00")

func hsdpTomlBytes() ([]byte, error) {
	return bindataRead(
		_hsdpToml,
		"hsdp.toml",
	)
}

func hsdpToml() (*asset, error) {
	bytes, err := hsdpTomlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "hsdp.toml", size: 8131, mode: os.FileMode(420), modTime: time.Unix(1604996572, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"hsdp.toml": hsdpToml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"hsdp.toml": &bintree{hsdpToml, map[string]*bintree{}},
}}
