package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/zclconf/go-cty/cty"

	giturl "github.com/chainguard-dev/git-urls"
)

type LocalTerraformModule struct {
	Name      string
	Dir       string
	ModuleFQN string
}

const (
	moduleBlockType    = "module"
	sourceAttrib       = "source"
	terraformExtension = "*.tf"
	restoreMarker      = "[restore-marker]"
	linebreak          = "\n"
)

var (
	localModules = []LocalTerraformModule{}
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
	u, err := giturl.Parse(remote)
	if err != nil {
		return "", "", err
	}
	if u.Host != "github.com" {
		return "", "", fmt.Errorf("expected GitHub remote, got: %s", remote)
	}
	orgRepo := u.Path
	orgRepo = trimAnySuffixes(orgRepo, []string{"/", ".git"})
	orgRepo = strings.TrimPrefix(orgRepo, "/")

	split := strings.Split(orgRepo, "/")
	if len(split) != 2 {
		return "", "", fmt.Errorf("expected GitHub remote of form https://github.com/ModuleRegistry/ModuleRepo, got: %s", remote)
	}
	org, repoName := split[0], split[1]

	// module repos are prefixed with terraform-google-
	if !strings.HasPrefix(repoName, "terraform-google-") {
		return "", "", fmt.Errorf("expected to find repo name prefixed with terraform-google-. Got: %s", repoName)
	}
	moduleName := strings.ReplaceAll(repoName, "terraform-google-", "")
	log.Printf("Module name set from remote to %s", moduleName)
	return moduleName, org, nil
}

// findSubModules generates slice of LocalTerraformModule for submodules
func findSubModules(path, rootModuleFQN string) []LocalTerraformModule {
	var subModules = make([]LocalTerraformModule, 0)
	// if no modules dir, return empty slice
	if _, err := os.Stat(path); err != nil {
		log.Print("No submodules found")
		return subModules
	}
	files, err := os.ReadDir(path)
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
		return f, nil
	}
	lines := strings.Split(strFile, linebreak)
	for i, line := range lines {
		if strings.Contains(line, restoreMarker) {
			lines[i] = strings.Split(line, restoreMarker)[1]
		}
	}
	return []byte(strings.Join(lines, linebreak)), nil
}

// matchedModule returns matching local TF module based on local path.
func matchedModule(localPath string) *LocalTerraformModule {
	for _, l := range localModules {
		if localPath == l.Dir {
			return &l
		}
	}
	return nil
}

// localToRemote converts all local references in f to remote references.
func localToRemote(f []byte, p string) ([]byte, error) {
	if _, err := os.Stat(p); err != nil {
		return nil, err
	}
	absPath, err := filepath.Abs(filepath.Dir(p))
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}
	f, err = restoreModules(f, p)
	if err != nil {
		return nil, err
	}

	currentReferences, err := moduleSourceRefs(f, p)
	if err != nil {
		return nil, fmt.Errorf("failed to write find module sources: %v", err)
	}
	newReferences := map[string]string{}
	for label, source := range currentReferences {
		localModule := matchedModule(filepath.Clean(filepath.Join(absPath, source)))
		if localModule == nil {
			log.Printf("no matches for %s", source)
			continue
		}
		newReferences[label] = localModule.ModuleFQN
	}
	if len(currentReferences) == 0 {
		return f, nil
	}
	updated, err := writeModuleRefs(f, p, newReferences)
	if err != nil {
		return nil, fmt.Errorf("failed to write updated module sources: %v", err)
	}
	// print diff info
	log.Printf("Modifications made to file %s", p)
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(f)),
		B:        difflib.SplitLines(string(updated)),
		FromFile: "Original",
		ToFile:   "Modified",
		Context:  3,
	}
	diffInfo, _ := difflib.GetUnifiedDiffString(diff)
	log.Println(diffInfo)
	return updated, nil
}

// remoteToLocal converts all remote references in f to local references.
func remoteToLocal(f []byte, p string) ([]byte, error) {
	if _, err := os.Stat(p); err != nil {
		return nil, err
	}
	f = commentVersions(f)
	absPath, err := filepath.Abs(filepath.Dir(p))
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}
	fqnMap := make(map[string]LocalTerraformModule, len(localModules))
	for _, l := range localModules {
		fqnMap[l.ModuleFQN] = l
	}
	currentReferences, err := moduleSourceRefs(f, p)
	if err != nil {
		return nil, fmt.Errorf("failed to write find module sources: %v", err)
	}
	newReferences := map[string]string{}
	for label, source := range currentReferences {
		localModule, exists := fqnMap[source]
		if !exists {
			continue
		}
		newModulePath, err := filepath.Rel(absPath, localModule.Dir)
		if err != nil {
			return nil, fmt.Errorf("failed to find relative path: %v", err)
		}
		newReferences[label] = newModulePath
	}
	if len(currentReferences) == 0 {
		return f, nil
	}
	updated, err := writeModuleRefs(f, p, newReferences)
	if err != nil {
		return nil, fmt.Errorf("failed to write updated module sources: %v", err)
	}
	// print diff info
	log.Printf("Modifications made to file %s", p)
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(f)),
		B:        difflib.SplitLines(string(updated)),
		FromFile: "Original",
		ToFile:   "Modified",
		Context:  3,
	}
	diffInfo, _ := difflib.GetUnifiedDiffString(diff)
	log.Println(diffInfo)
	return updated, nil
}

// commentVersions comments version attributes for local modules.
func commentVersions(f []byte) []byte {
	strFile := string(f)
	lines := strings.Split(strFile, linebreak)
	for _, localModule := range localModules {
		// check if current file has module/submodules references that should be swapped
		if !strings.Contains(strFile, localModule.ModuleFQN) {
			continue
		}
		for i, line := range lines {
			if !strings.Contains(line, localModule.ModuleFQN) {
				continue
			}
			if i < len(lines)-1 && strings.Contains(lines[i+1], "version") && !strings.Contains(lines[i+1], restoreMarker) {
				leadingWhiteSpace := lines[i+1][:strings.Index(lines[i+1], "version")]
				lines[i+1] = fmt.Sprintf("%s# %s %s", leadingWhiteSpace, restoreMarker, lines[i+1])
			}
		}
	}
	newExample := strings.Join(lines, linebreak)
	return []byte(newExample)
}

// getTFFiles returns a slice of valid TF file paths
func getTFFiles(path string) []string {
	// validate path
	if _, err := os.Stat(path); err != nil {
		log.Fatal(fmt.Errorf("Unable to find %s : %v", path, err))
	}
	var files = make([]string, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil && info.IsDir() {
			return nil
		}
		isTFFile, _ := filepath.Match(terraformExtension, filepath.Base(path))
		if isTFFile {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking files: %v", err)
	}
	return files

}

var (
	// Partial schema of examples.
	exampleSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       moduleBlockType,
				LabelNames: []string{"name"},
			},
		},
	}
	// Partial schema of each module.
	moduleSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: sourceAttrib,
			},
		},
	}
)

// moduleSourceRefs returns a map of module label to corresponding source references.
func moduleSourceRefs(f []byte, TFFilePath string) (map[string]string, error) {
	refs := map[string]string{}
	p, err := hclparse.NewParser().ParseHCL(f, TFFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hcl: %v", err)
	}
	c, _, diags := p.Body.PartialContent(exampleSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse example content: %v", diags.Error())
	}

	for _, b := range c.Blocks {
		if b.Type != moduleBlockType {
			continue
		}
		if len(b.Labels) != 1 {
			log.Printf("got multiple labels %v, module should only have one", b.Labels)
			continue
		}

		content, _, diags := b.Body.PartialContent(moduleSchema)
		if diags.HasErrors() {
			log.Printf("skipping %s module, failed to parse module content: %v", b.Labels[0], diags.Error())
			continue
		}

		sourcrAttr, exists := content.Attributes[sourceAttrib]
		if !exists {
			log.Printf("skipping %s module, no source attribute", b.Labels[0])
			continue
		}
		var sourceName string
		diags = gohcl.DecodeExpression(sourcrAttr.Expr, nil, &sourceName)
		if diags.HasErrors() {
			log.Printf("skipping %s module, failed to decode source value: %v", b.Labels[0], diags.Error())
			continue
		}
		refs[b.Labels[0]] = sourceName
	}
	return refs, nil
}

// writeModuleRefs appends or overwrites provided moduleRefs to file f.
func writeModuleRefs(f []byte, p string, moduleRefs map[string]string) ([]byte, error) {
	wf, diags := hclwrite.ParseConfig(f, p, hcl.Pos{})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse hcl: %v", diags.Error())
	}
	for _, b := range wf.Body().Blocks() {
		if b.Type() != moduleBlockType {
			continue
		}
		if len(b.Labels()) != 1 {
			log.Printf("got multiple labels %v, module should only have one", b.Labels())
			continue
		}
		newSource, exists := moduleRefs[b.Labels()[0]]
		if !exists {
			continue
		}
		b.Body().SetAttributeValue(sourceAttrib, cty.StringVal(newSource))
	}

	var testS strings.Builder
	_, err := wf.WriteTo(&testS)
	if err != nil {
		return nil, fmt.Errorf("failed to write hcl: %v", diags.Error())
	}
	return []byte(testS.String()), nil
}

func SwapModules(rootPath, moduleRegistrySuffix, moduleRegistryPrefix, subModulesDir, examplesDir string, restore bool) {
	rootPath = filepath.Clean(rootPath)
	moduleName, foundRegistryPrefix, err := getModuleNameRegistry(rootPath)
	if err != nil && moduleRegistryPrefix == "" {
		log.Printf("failed to get module name and registry: %v", err)
		return
	}

	if moduleRegistryPrefix != "" {
		foundRegistryPrefix = moduleRegistryPrefix
	}

	// add root module to slice of localModules
	localModules = append(localModules, LocalTerraformModule{moduleName, rootPath, fmt.Sprintf("%s/%s/%s", foundRegistryPrefix, moduleName, moduleRegistrySuffix)})
	examplesPath := fmt.Sprintf("%s/%s", rootPath, examplesDir)
	subModulesPath := fmt.Sprintf("%s/%s", rootPath, subModulesDir)

	// add submodules, if any to localModules
	submods := findSubModules(subModulesPath, localModules[0].ModuleFQN)
	localModules = append(localModules, submods...)

	// find all TF files in examples dir to process
	exampleTFFiles := getTFFiles(examplesPath)
	for _, TFFilePath := range exampleTFFiles {
		file, err := os.ReadFile(TFFilePath)
		if err != nil {
			log.Printf("Error reading file: %v", err)
		}

		var newFile []byte
		if restore {
			newFile, err = localToRemote(file, TFFilePath)
		} else {
			newFile, err = remoteToLocal(file, TFFilePath)
		}
		if err != nil {
			log.Printf("Error processing file: %v", err)
		}

		if newFile != nil {
			err = os.WriteFile(TFFilePath, newFile, 0644)
			if err != nil {
				log.Printf("Error writing file: %v", err)
			}
		}
	}
}
