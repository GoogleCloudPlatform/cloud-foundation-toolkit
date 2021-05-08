package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/pmezard/go-difflib/difflib"
)

type LocalTerraformModule struct {
	Name      string
	Dir       string
	ModuleFQN string
}

var (
	terraformExtension = "*.tf"
	restoreMarker      = "[restore-marker]"
	linebreak          = "\n"
	localModules       = []LocalTerraformModule{}
)

// getRemoteURL gets the URL of a given remote from git repo at dir
func getRemoteURL(dir, remoteName string) (string, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return "", err
	}
	rm, err := r.Remote(remoteName)
	if err != nil {
		return "", err
	}
	return rm.Config().URLs[0], nil
}

// trimAnySuffixes trims first matching suffix from slice of suffixes
func trimAnySuffixes(s string, suffixes []string) string {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			s = s[:len(s)-len(suffix)]
			return s
		}
	}
	return s
}

// getModuleNameRegistry returns module name and registry by parsing git remote
func getModuleNameRegistry(dir string) (string, string, error) {
	remote, err := getRemoteURL(dir, "origin")
	if err != nil {
		return "", "", err
	}

	// GH remote will be of form https://github.com/ModuleRegistry/ModuleName
	if !strings.Contains(remote, "https://github.com/") {
		return "", "", fmt.Errorf("Expected GitHub remote of form https://github.com/ModuleRegistry/ModuleRepo. Got: %s", remote)
	}

	// remotes maybe suffixed with a trailing / or .git
	remote = trimAnySuffixes(remote, []string{"/", ".git"})
	namePrefix := strings.ReplaceAll(remote, "https://github.com/", "")
	if !strings.Contains(namePrefix, "/") {
		return "", "", fmt.Errorf("Expected GitHub org/owner of form ModuleRegistry/ModuleRepo. Got: %s", namePrefix)
	}
	moduleRegistry := namePrefix[:strings.LastIndex(namePrefix, "/")]
	repoName := namePrefix[strings.LastIndex(namePrefix, "/")+1:]

	// module repos are prefixed with terraform-google-
	if !strings.HasPrefix(repoName, "terraform-google-") {
		return "", "", fmt.Errorf("Expected to find repo name prefixed with terraform-google-. Got: %s", repoName)
	}
	moduleName := strings.ReplaceAll(repoName, "terraform-google-", "")
	log.Printf("Module name set from remote to %s", moduleName)
	return moduleName, moduleRegistry, nil
}

// findSubModules generates slice of LocalTerraformModule for submodules
func findSubModules(path, rootModuleFQN string) []LocalTerraformModule {
	var subModules = make([]LocalTerraformModule, 0)
	// if no modules dir, return empty slice
	if _, err := os.Stat(path); err != nil {
		log.Print("No submodules found")
		return subModules
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Error finding submodules: %v", err)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Error finding submodule absolute path: %v", err)
	}
	for _, f := range files {
		if f.IsDir() {
			subModules = append(subModules, LocalTerraformModule{f.Name(), filepath.Join(absPath, f.Name()), fmt.Sprintf("%s//modules/%s", rootModuleFQN, f.Name())})
		}
	}
	return subModules
}

// restoreModules restores old config as marked by restoreMarker
func restoreModules(f []byte, p string) ([]byte, error) {
	if _, err := os.Stat(p); err != nil {
		return nil, err
	}
	strFile := string(f)
	if !strings.Contains(strFile, restoreMarker) {
		return nil, nil
	}
	lines := strings.Split(strFile, linebreak)
	for i, line := range lines {
		if strings.Contains(line, restoreMarker) {
			lines[i] = strings.Split(line, restoreMarker)[1]
		}
	}
	return []byte(strings.Join(lines, linebreak)), nil
}

// replaceLocalModules swaps current local module registry references with local path
func replaceLocalModules(f []byte, p string) ([]byte, error) {
	if _, err := os.Stat(p); err != nil {
		return nil, err
	}
	absPath, err := filepath.Abs(filepath.Dir(p))
	if err != nil {
		return nil, fmt.Errorf("Error finding example absolute path: %v", err)
	}
	strFile := string(f)
	lines := strings.Split(strFile, linebreak)
	for _, localModule := range localModules {
		// check if current file has module/submodules references that should be swapped
		if !strings.Contains(strFile, localModule.ModuleFQN) {
			continue
		}
		// get relative path from example to local module
		newModulePath, err := filepath.Rel(absPath, localModule.Dir)
		if err != nil {
			return nil, fmt.Errorf("Error finding relative path: %v", err)
		}
		for i, line := range lines {
			if strings.Contains(line, fmt.Sprintf("\"%s\"", localModule.ModuleFQN)) && !strings.Contains(line, restoreMarker) {
				// swap with local module and add restore point
				leadingWhiteSpace := line[:strings.Index(line, "source")]
				newSource := fmt.Sprintf("source = \"%s\"", newModulePath)
				lines[i] = leadingWhiteSpace + newSource + fmt.Sprintf(" # %s %s", restoreMarker, line)
				// if next line is a version declaration, disable that as well
				if i < len(lines)-1 && strings.Contains(lines[i+1], "version") {
					leadingWhiteSpace = lines[i+1][:strings.Index(lines[i+1], "version")]
					lines[i+1] = fmt.Sprintf("%s# %s %s", leadingWhiteSpace, restoreMarker, lines[i+1])
				}
			}
		}
	}
	newExample := strings.Join(lines, linebreak)
	// check if any swaps have been made
	if newExample == strFile {
		return nil, nil
	}
	// print diff info
	log.Printf("Modifications made to file %s", p)
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(strFile),
		B:        difflib.SplitLines(newExample),
		FromFile: "Original",
		ToFile:   "Modified",
		Context:  3,
	}
	diffInfo, _ := difflib.GetUnifiedDiffString(diff)
	log.Println(diffInfo)
	return []byte(newExample), nil

}

// getTFFiles returns a slice of valid TF file paths
func getTFFiles(path string) []string {
	// validate path
	if _, err := os.Stat(path); err != nil {
		log.Fatal(fmt.Errorf("Unable to find %s : %v", path, err))
	}
	var files = make([]string, 0)
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil && info.IsDir() {
			return nil
		}
		isTFFile, _ := filepath.Match(terraformExtension, filepath.Base(path))
		if isTFFile {
			files = append(files, path)
		}
		return nil
	})
	return files

}

func SwapModules(rootPath, moduleRegistrySuffix, subModulesDir, examplesDir string, restore bool) {
	moduleName, moduleRegistryPrefix, err := getModuleNameRegistry(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	// add root module to slice of localModules
	localModules = append(localModules, LocalTerraformModule{moduleName, rootPath, fmt.Sprintf("%s/%s/%s", moduleRegistryPrefix, moduleName, moduleRegistrySuffix)})
	examplesPath := fmt.Sprintf("%s/%s", rootPath, examplesDir)
	subModulesPath := fmt.Sprintf("%s/%s", rootPath, subModulesDir)

	// add submodules, if any to localModules
	submods := findSubModules(subModulesPath, localModules[0].ModuleFQN)
	localModules = append(localModules, submods...)

	// find all TF files in examples dir to process
	exampleTFFiles := getTFFiles(examplesPath)
	for _, TFFilePath := range exampleTFFiles {
		file, err := ioutil.ReadFile(TFFilePath)
		if err != nil {
			log.Printf("Error reading file: %v", err)
		}

		var newFile []byte
		if restore {
			newFile, err = restoreModules(file, TFFilePath)
		} else {
			newFile, err = replaceLocalModules(file, TFFilePath)
		}
		if err != nil {
			log.Printf("Error processing file: %v", err)
		}

		if newFile != nil {
			err = ioutil.WriteFile(TFFilePath, newFile, 0644)
			if err != nil {
				log.Printf("Error writing file: %v", err)
			}
		}

	}

}
