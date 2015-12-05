package main

import (
	"os"
	"path/filepath"
)

// Vendor TBD
type Vendor struct {
	path string
}

// NewVendor - vendor construction
func NewVendor(path string) (*Vendor, error) {
	vendorAbsPath, err := filepath.Abs(vendorFolder)
	if err != nil {
		return nil, err
	}
	return &Vendor{
		path: vendorAbsPath,
	}, nil
}

// SetEnvVariables -- set env according vendor path
func (v *Vendor) SetEnvVariables() error {
	if err := os.Setenv("GOPATH", v.path); err != nil {
		return err
	}
	if err := os.Setenv("GOBIN", filepath.Join(v.path, "bin")); err != nil {
		return err
	}

	return nil
}

// MoveSrcToVendorSrc TBD
func (v *Vendor) MoveSrcToVendorSrc() error {
	vendorSrc := filepath.Join(v.path, "src")
	dirs, err := v.readDirNames(v.path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(vendorSrc, 0755); err != nil {
		return err
	}
	for _, dir := range dirs {
		if dir == "bin" || dir == "pkg" || dir == "src" {
			continue
		}
		err = os.Rename(filepath.Join(v.path, dir), filepath.Join(vendorSrc, dir))
		if err != nil {
			return err
		}
	}
	return nil
}

// MoveSrcToVendor TBD
func (v *Vendor) MoveSrcToVendor() error {
	vendorSrc := filepath.Join(v.path, "src")
	dirs, err := v.readDirNames(vendorSrc)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		err = os.Rename(filepath.Join(vendorSrc, dir), filepath.Join(v.path, dir))
		if err != nil {
			return err
		}
	}
	if err := os.Remove(vendorSrc); err != nil {
		return err
	}
	return nil
}

func (v *Vendor) readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	return list, nil
}
