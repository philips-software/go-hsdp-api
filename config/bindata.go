// Code generated for package config by go-bindata DO NOT EDIT. (@generated)
// sources:
// hsdp.json
package config

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

var _hsdpJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x58\xd1\x6e\xab\x30\x0c\x7d\xef\x57\x44\x79\x1e\x4c\xbd\x7d\xdb\xcf\x5c\x79\x24\x2b\x51\x69\x12\x91\xd0\x6d\xba\xda\xbf\x5f\xb1\xb2\xa9\xd0\xd0\x38\xc4\x74\x7b\x98\x56\x68\x72\x8e\xe3\xda\x3e\x07\xfe\x6d\x18\xe3\xad\xdc\x2b\xa3\xf9\x13\xeb\xaf\x18\xe3\x60\xa1\xda\x7d\x5f\x32\xc6\xa5\x3e\x5d\x5c\x32\xc6\x6d\x6b\xc4\xe8\x0e\x63\xdc\xc9\xf6\xa4\x2a\x39\xb9\xcd\x18\x57\x70\xbc\xba\x79\xbe\xfd\xb7\x6b\x1b\xfe\xc4\x78\xed\xbd\x75\x4f\x8f\x8f\x0a\x8e\xc5\x00\x53\x82\xdd\x95\xb5\x13\xb6\x54\x86\x3f\x5c\xed\x15\xd7\x7b\x05\x72\x2f\x92\x73\xb4\xed\xe3\x61\x33\xe5\x0f\x1d\x09\x19\xd2\x4d\xe8\x83\x7a\x06\x0d\x08\xf4\xf3\x42\x3c\x70\x63\xf6\x7b\xa5\xf7\x08\xe4\xc6\xf4\x0b\xa5\xf3\xa6\xfd\x73\x03\x7f\x13\xfa\xfc\xf5\xe9\x9b\x3d\x58\x17\xbc\x7a\x99\x96\x8f\x30\x47\x50\x7d\x11\xf2\xb9\x5f\x6f\x1a\x26\x58\x55\xba\x77\x17\x8e\xf0\xe2\xf4\xbc\x32\xda\x99\x66\x5a\x99\x53\xb8\x61\x55\x1c\x4e\x98\xea\x20\xdb\xa2\xef\x1a\xe7\xdb\xf7\x29\x6c\x6d\x9c\xef\x71\xcf\xcb\xe2\x70\x7b\xf0\xf2\x15\xae\x60\x2e\xf3\x01\xd5\xae\xb4\x9f\x20\x95\x39\x8e\x73\xf2\xc5\xb6\x7f\x2d\xa6\xeb\x82\x6c\x1d\x4c\x6b\x6b\x9a\x87\x0e\xe0\x46\x5a\x37\x97\xff\x07\x60\x2e\xbb\xe2\x55\x7e\xc6\x31\x3f\x32\xaa\x46\x49\xed\x0b\x3f\x5e\x37\x5b\x21\xe7\xc3\x81\x43\x94\x6b\x0d\xae\xb8\x40\x2f\x87\x68\x4a\x5b\xab\x46\x59\x57\xd4\x12\x1a\x5f\xbb\x4e\x79\x39\xce\x0b\x0b\x75\xf6\xa2\x61\x15\xa3\xc4\x0e\xb0\x94\x63\x60\x06\x1b\x5d\x5a\xd0\x03\x8f\x8c\x32\x71\x10\xd2\x11\x2f\x1b\x94\x74\xfc\xf6\xa0\x10\xdc\xf6\xa0\x0a\xdb\x9a\xb7\x77\x3a\x62\xb7\xab\x5a\x29\x30\x2d\x37\xac\xcc\xa0\x0e\x8b\xc7\xc3\x5d\x2d\x06\x55\xd7\x2e\xc6\xcb\x88\x8d\xd8\x9e\xdc\xb3\x53\x7f\xa2\x3b\xef\xdd\x18\xcb\xf2\x9a\xe5\xa8\xa0\xf5\xb2\x99\x33\x23\xe7\x6f\x0b\xd9\x6d\xcb\xaa\x31\x9d\x88\xd9\x84\x5b\xfe\x2c\xa5\xcc\x43\x9e\xad\x0f\x22\x46\x9f\xe4\xd8\x7a\x40\x42\xc7\x16\x85\x8b\x3a\x36\x79\x94\xb0\x45\x38\x36\x44\x26\x70\x7e\x6d\x0e\x28\xe8\xd6\x1c\x6c\xc9\x9d\x5a\xda\x00\xbe\x94\x0d\x07\xdb\xe4\xe7\x3c\xf4\xfe\x04\xee\x15\xec\x0f\x1a\x3e\xc7\xea\xa0\x49\x08\x6c\xcd\x3c\xd7\x6f\x90\xf3\x25\x95\x84\xda\x8b\xe4\x24\x96\xe4\xb5\xaa\x67\xed\x8a\x49\xae\x12\x42\x7d\xeb\xa9\xb3\xf5\x6d\xae\x16\xe6\xde\x3f\x04\xcf\xbb\x5c\xcd\xa2\x70\x69\x6a\x16\x85\x8b\xaa\x99\x43\x69\xd9\x78\x55\x86\x96\xcd\xa7\x34\xa8\x66\x9d\x2b\x24\xac\xf0\xee\x61\xb9\xa2\x0d\x11\x91\xbe\x0c\x48\xc6\xcc\x8c\x71\x05\x35\xcc\xa2\xcc\x51\xc8\x2c\x62\x02\xd5\xcc\xe2\xcf\x7b\x26\x4f\xa3\x8e\x8a\xb8\x90\xa7\x95\xfa\x47\xc8\x93\x6c\x8c\x3d\x4a\x4d\xd7\x3f\x59\x98\x99\x31\x12\xf5\xcf\x72\xca\xdf\x60\xc8\xa8\x7e\xc8\xc5\x78\x19\xb1\x11\x9b\xb9\x7b\x0e\xbf\x9f\x18\x78\xf7\x1e\x72\xcb\xf2\x1a\x6d\x0a\x07\x5a\x3c\x9b\x37\x74\x5f\x54\xa2\x2d\x9c\xef\x76\x88\xc8\xbf\x96\x16\x03\x07\x49\xe8\x84\xd6\x59\x53\x58\xe7\x94\x0e\x0d\xd9\xe9\x21\x80\xea\xe5\x53\xbd\x88\x5e\x0f\x69\x5a\x43\x1d\x85\x8b\x1a\x6a\x8d\x32\xd4\x9a\xca\x50\xdf\x4a\xea\xd8\x52\x6f\xfa\xbf\x8f\xcd\xff\x00\x00\x00\xff\xff\x9c\xb3\x24\x13\x15\x20\x00\x00")

func hsdpJsonBytes() ([]byte, error) {
	return bindataRead(
		_hsdpJson,
		"hsdp.json",
	)
}

func hsdpJson() (*asset, error) {
	bytes, err := hsdpJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "hsdp.json", size: 8213, mode: os.FileMode(420), modTime: time.Unix(1606253327, 0)}
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
	"hsdp.json": hsdpJson,
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
	"hsdp.json": &bintree{hsdpJson, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
