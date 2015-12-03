package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type vcsCmd struct {
	checkout     []string
	update       []string
	revision     []string
	revisionMask string
}

var (
	hg = &vcsCmd{
		[]string{"hg", "update"},
		[]string{"hg", "pull"},
		[]string{"hg", "id", "-i"},
		"^(.+)$",
	}
	git = &vcsCmd{
		[]string{"git", "checkout", "-q"},
		[]string{"git", "fetch"},
		[]string{"git", "rev-parse", "HEAD"},
		"^(.+)$",
	}
	bzr = &vcsCmd{
		[]string{"bzr", "revert", "-r"},
		[]string{"bzr", "pull"},
		[]string{"bzr", "log", "-r-1", "--line"},
		"^([0-9]+)",
	}
)

func (vcs *vcsCmd) Checkout(p, destination string) error {
	args := append(vcs.checkout, destination)
	return vcsExec(p, args...)
}

func (vcs *vcsCmd) Update(p string) error {
	return vcsExec(p, vcs.update...)
}

func (vcs *vcsCmd) Revision(dir string) (string, error) {
	args := append(vcs.revision)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	b, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return "", err
	}
	rev := strings.TrimSpace(string(b))
	if vcs.revisionMask != "" {
		return regexp.MustCompile(vcs.revisionMask).FindString(rev), nil
	}
	return rev, nil
}

func (vcs *vcsCmd) Sync(p, destination string) error {
	err := vcs.Checkout(p, destination)
	if err != nil {
		err = vcs.Update(p)
		if err != nil {
			return err
		}
		err = vcs.Checkout(p, destination)
	}
	return err
}

func vcsExec(dir string, args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func list(dir string) ([]string, error) {
	cmd := exec.Command("go", "list", "./...")
	cmd.Dir = dir
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return strings.Split(stdout.String(), "\n"), nil
}

func has(c interface{}, key string) bool {
	if m, ok := c.(map[string]interface{}); ok {
		_, ok := m[key]
		return ok
	} else if a, ok := c.([]string); ok {
		for _, s := range a {
			if ok && s == key {
				return true
			}
		}
	}
	return false
}

func (gom *Gom) Update() error {
	cmdArgs := []string{"go", "get", "-u"}
	if insecure, ok := gom.options["insecure"].(string); ok {
		if insecure == "true" {
			cmdArgs = append(cmdArgs, "-insecure")
		}
	}
	cmdArgs = append(cmdArgs, gom.name+"/...")

	fmt.Printf("updating %s\n", gom.name)
	return run(cmdArgs, Green)
}

func (gom *Gom) Clone(args []string) error {
	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}
	if command, ok := gom.options["command"].(string); ok {
		target, ok := gom.options["target"].(string)
		if !ok {
			target = gom.name
		}

		srcdir := filepath.Join(vendor, "src", target)
		if err := os.MkdirAll(srcdir, 0755); err != nil {
			return err
		}

		customCmd := strings.Split(command, " ")
		customCmd = append(customCmd, srcdir)

		fmt.Printf("fetching %s (%v)\n", gom.name, customCmd)
		err = run(customCmd, Blue)
		if err != nil {
			return err
		}
	} else if private, ok := gom.options["private"].(string); ok {
		if private == "true" {
			target, ok := gom.options["target"].(string)
			if !ok {
				target = gom.name
			}
			srcdir := filepath.Join(vendor, "src", target)
			if _, err := os.Stat(srcdir); err != nil {
				if err := os.MkdirAll(srcdir, 0755); err != nil {
					return err
				}
				if err := gom.clonePrivate(srcdir); err != nil {
					return err
				}
			} else {
				if err := gom.pullPrivate(srcdir); err != nil {
					return err
				}
			}
		}
	}

	if skipdep, ok := gom.options["skipdep"].(string); ok {
		if skipdep == "true" {
			return nil
		}
	}
	cmdArgs := []string{"go", "get", "-d"}
	if insecure, ok := gom.options["insecure"].(string); ok {
		if insecure == "true" {
			cmdArgs = append(cmdArgs, "-insecure")
		}
	}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, gom.name+"/...")

	fmt.Printf("downloading %s\n", gom.name)
	return run(cmdArgs, Blue)
}

func (gom *Gom) pullPrivate(srcdir string) (err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(srcdir); err != nil {
		return err
	}
	defer os.Chdir(cwd)

	fmt.Printf("fetching private repo %s\n", gom.name)
	pullCmd := "git pull origin master"
	pullArgs := strings.Split(pullCmd, " ")
	err = run(pullArgs, Blue)
	if err != nil {
		return
	}

	return
}

func (gom *Gom) clonePrivate(srcdir string) (err error) {
	name := strings.Split(gom.name, "/")
	privateUrl := fmt.Sprintf("git@%s:%s/%s", name[0], name[1], name[2])

	fmt.Printf("fetching private repo %s\n", gom.name)
	cloneCmd := []string{"git", "clone", privateUrl, srcdir}
	err = run(cloneCmd, Blue)
	if err != nil {
		return
	}

	return
}

func (gom *Gom) Checkout() error {
	commit_or_branch_or_tag := ""
	if has(gom.options, "branch") {
		commit_or_branch_or_tag, _ = gom.options["branch"].(string)
	}
	if has(gom.options, "tag") {
		commit_or_branch_or_tag, _ = gom.options["tag"].(string)
	}
	if has(gom.options, "commit") {
		commit_or_branch_or_tag, _ = gom.options["commit"].(string)
	}
	if commit_or_branch_or_tag == "" {
		return nil
	}
	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}
	p := filepath.Join(vendor, "src")
	target, ok := gom.options["target"].(string)
	if !ok {
		target = gom.name
	}
	for _, elem := range strings.Split(target, "/") {
		var vcs *vcsCmd
		p = filepath.Join(p, elem)
		if isDir(filepath.Join(p, ".git")) {
			vcs = git
		} else if isDir(filepath.Join(p, ".hg")) {
			vcs = hg
		} else if isDir(filepath.Join(p, ".bzr")) {
			vcs = bzr
		}
		if vcs != nil {
			p = filepath.Join(vendor, "src", target)
			return vcs.Sync(p, commit_or_branch_or_tag)
		}
	}
	fmt.Printf("Warning: don't know how to checkout for %v\n", gom.name)
	return errors.New("gom currently support git/hg/bzr for specifying tag/branch/commit")
}

func (gom *Gom) Build(args []string) error {
	installCmd := []string{"go", "get"}
	hasPkg := false
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			arg = path.Join(arg, "...")
			hasPkg = true
		}
		installCmd = append(installCmd, arg)
	}

	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}
	target, ok := gom.options["target"].(string)
	if !ok {
		target = gom.name
	}
	p := filepath.Join(vendor, "src", target)

	if hasPkg {
		return vcsExec(p, installCmd...)
	}

	pkgs, err := list(p)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		if isIgnorePackage(pkg) {
			continue
		}
		p = filepath.Join(vendor, "src", pkg)
		err := vcsExec(p, installCmd...)
		if err != nil {
			return err
		}
	}
	return nil
}

func isFile(p string) bool {
	if fi, err := os.Stat(filepath.Join(p)); err == nil && !fi.IsDir() {
		return true
	}
	return false
}

func isDir(p string) bool {
	if fi, err := os.Stat(filepath.Join(p)); err == nil && fi.IsDir() {
		return true
	}
	return false
}

func isIgnorePackage(pkg string) bool {
	if pkg == "" {
		return true
	}
	paths := strings.Split(pkg, "/")
	for _, path := range paths {
		if path == "examples" {
			return true
		}
		if strings.HasPrefix(path, "_") {
			return true
		}
	}
	return false
}

func moveSrcToVendorSrc(vendor string) error {
	vendorSrc := filepath.Join(vendor, "src")
	dirs, err := readdirnames(vendor)
	if err != nil {
		return err
	}
	err = os.MkdirAll(vendorSrc, 0755)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if dir == "bin" || dir == "pkg" || dir == "src" {
			continue
		}
		err = os.Rename(filepath.Join(vendor, dir), filepath.Join(vendorSrc, dir))
		if err != nil {
			return err
		}
	}
	return nil
}

func moveSrcToVendor(vendor string) error {
	vendorSrc := filepath.Join(vendor, "src")
	dirs, err := readdirnames(vendorSrc)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		err = os.Rename(filepath.Join(vendorSrc, dir), filepath.Join(vendor, dir))
		if err != nil {
			return err
		}
	}
	err = os.Remove(vendorSrc)
	if err != nil {
		return err
	}
	return nil
}

func readdirnames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func parseInstallFlags(args []string) (opts map[string]string, retargs []string) {
	opts = make(map[string]string)
	re := regexp.MustCompile(`^--([a-z][a-z_]*)(=\S*)?`)
	for _, arg := range args {
		ss := re.FindAllStringSubmatch(arg, -1)
		if len(ss) > 0 {
			opts[ss[0][1]] = opts[ss[0][2]]
		} else {
			retargs = append(retargs, arg)
		}
	}
	return
}

func hasSaveOpts(opts map[string]string) bool {
	if _, ok := opts["save"]; ok {
		return true
	}
	if _, ok := opts["save-dev"]; ok {
		return true
	}
	return false
}

func install(args []string) error {
	var opts map[string]string
	opts, args = parseInstallFlags(args)
	allGoms, err := parseGomfile("Gomfile")
	if err != nil {
		return err
	}
	if hasSaveOpts(opts) {
		found := false
		for _, arg := range args {
			for _, gom := range allGoms {
				if gom.name == arg {
					found = true
					break
				}
			}
			if !found {
				options := map[string]interface{}{}
				if _, ok := opts["save-dev"]; ok {
					options["envs"] = []string{"development"}
				}
				allGoms = append(allGoms, Gom{name: arg, options: options})
			}
		}
		err = writeGomfile("Gomfile", allGoms)
		if err != nil {
			return err
		}
	}
	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}
	_, err = os.Stat(vendor)
	if err != nil {
		err = os.MkdirAll(vendor, 0755)
		if err != nil {
			return err
		}
	}
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		return err
	}
	err = os.Setenv("GOBIN", filepath.Join(vendor, "bin"))
	if err != nil {
		return err
	}

	// 1. Filter goms to install
	goms := make([]Gom, 0)
	for _, gom := range allGoms {
		if group, ok := gom.options["group"]; ok {
			if !matchEnv(group) {
				continue
			}
		}
		if goos, ok := gom.options["goos"]; ok {
			if !matchOS(goos) {
				continue
			}
		}
		goms = append(goms, gom)
	}

	if go15VendorExperimentEnv {
		err = moveSrcToVendorSrc(vendor)
		if err != nil {
			return err
		}
	}

	// 2. Clone the repositories
	for _, gom := range goms {
		err = gom.Clone(args)
		if err != nil {
			return err
		}
	}

	// 3. Checkout the commit/branch/tag if needed
	for _, gom := range goms {
		err = gom.Checkout()
		if err != nil {
			return err
		}
	}

	// 4. Build and install
	for _, gom := range goms {
		if skipdep, ok := gom.options["skipdep"].(string); ok {
			if skipdep == "true" {
				continue
			}
		}
		err = gom.Build(args)
		if err != nil {
			return err
		}
	}

	if go15VendorExperimentEnv {
		err = moveSrcToVendor(vendor)
		if err != nil {
			return err
		}
	}

	return nil
}

func update() error {
	goms, err := parseGomfile("Gomfile")
	if err != nil {
		return err
	}
	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		return err
	}
	err = os.Setenv("GOBIN", filepath.Join(vendor, "bin"))
	if err != nil {
		return err
	}

	if go15VendorExperimentEnv {
		err = moveSrcToVendorSrc(vendor)
		if err != nil {
			return err
		}
	}

	for _, gom := range goms {
		err = gom.Update()
		if err != nil {
			return err
		}
		vcs, _, p := vcsScan(vendorSrc(vendor), gom.name)
		if vcs != nil {
			rev, err := vcs.Revision(p)
			if err == nil && rev != "" {
				gom.options["commit"] = rev
			}
		}
	}

	if go15VendorExperimentEnv {
		err = moveSrcToVendor(vendor)
		if err != nil {
			return err
		}
	}

	return writeGomfile("Gomfile", goms)
}
