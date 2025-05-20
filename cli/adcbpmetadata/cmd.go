package adcbpmetadata

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


var flags struct {
	rootModulePath          string
}

type markdownHeading struct{
	headLevel int
	headOrder int
	headTitle string
	content bool
	description string
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVarP(&flags.rootModulePath, "path", "p", ".", "Path to the blueprint for generating metadata.")
}



var Cmd = &cobra.Command{
	Use:   "adc_validate",
	Short: "Validated terraform module for ingestion into ADC",
	Long:  `Validated terraform module for ingestion into ADC`,
	Args:  cobra.NoArgs,
	RunE: generate,
}




// The top-level command function that generates metadata based on the provided flags
func generate(cmd *cobra.Command, args []string) error {
	wdPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working dir: %w", err)
	}

	currBpPath := wdPath
	if !path.IsAbs(flags.rootModulePath) {
		currBpPath = path.Join(wdPath, flags.rootModulePath)
	}

	err=ValidateRootModuleForADC(currBpPath)
	if(err!=nil){
		return err
	}
	return nil
}



func ValidateRootModuleForADC(bpPath string) error{

	// files root module must have TODO: Add description describing need of each file
	// TODO: ADC doesn't require readme file but metadata generation code require it : we can add one readme
	requiredFilesForRoot:=[]string {"README.md", "main.tf", "versions.tf"}
	missingFiles:=checkFilePresence(bpPath, requiredFilesForRoot)

	if len(missingFiles) > 0 {
		return fmt.Errorf("top-level module must have following files also:%s", missingFiles)
	}

	goodToHaveFiles:=[]string{"test/setup/iam.tf",  "test/setup/main.tf"}
	missingGoodToHaveFiles:=checkFilePresence(bpPath, goodToHaveFiles)

	if len(missingGoodToHaveFiles) > 0 {
		Log.Warn("It is good to have these files also for generating metadata: ["+ strings.Join(missingGoodToHaveFiles, ", ") +"]\n")
	}

	otherFiles:=[]string{"assets/icon.png", "modules/", "examples",}
	missingOtherFiles:=checkFilePresence(bpPath, otherFiles)
	if len(missingOtherFiles) > 0 {
		Log.Info("Cft module can have these files also: [" + strings.Join(missingOtherFiles, ", ")+"]\n")
	}

 // Validate Readme.md file content
 err:= validateReadme(path.Join(bpPath, "README.md"))
 if err !=nil{
	 return err
 }

 err = validateVersionsFiles(path.Join(bpPath, "versions.tf"))
 if err !=nil {
	 return err
 }

 return nil
}


func validateVersionsFiles(versionsConfigPath string)error {
 Log.Info("**************** Validating  versions.tf****************")

 //set of allowed providers
 allowedTerraformProviders:=	[]string{
	 "hashicorp/google","hashicorp/google-beta",
 }

 p := hclparse.NewParser()
 versionsFile, diags := p.ParseHCLFile(versionsConfigPath)
 err := bpmetadata.HasHclErrors(diags)
 if err != nil {
	 return  err
 }

 // parse out the required providers from the config
 var hclModule tfconfig.Module
 hclModule.RequiredProviders = make(map[string]*tfconfig.ProviderRequirement)
 diags = tfconfig.LoadModuleFromFile(versionsFile, &hclModule)
 err = bpmetadata.HasHclErrors(diags)
 if err != nil {
	 return  err
 }


	var errors []string
 for _, providerData := range hclModule.RequiredProviders {
	 if providerData.Source == "" {
		 errors=append(errors, fmt.Sprintln("source not found in provider settings"))
	 }else if slices.Index(allowedTerraformProviders, providerData.Source)==-1{
		 errors=append(errors, fmt.Sprintln("Incorrect terraform provider: "+providerData.Source + ". Only following terraform providers are allowed as of now: ["+ strings.Join(allowedTerraformProviders, ", ")+"]"))
	 }

	 if len(providerData.VersionConstraints) == 0 {
		 errors=append(errors, fmt.Sprintln("version not found in provider settings"))
	 }
 }

if len(errors)>0{
 return fmt.Errorf(strings.Join(errors, "\n"))
}

Log.Info("****************Done Validating  versions.tf****************\n\n")
return nil
}


func validateReadme(readmeFilePath string ) error{
 Log.Info("\n\n**************** Validating  README.md****************")
 readmeContent, err := os.ReadFile(readmeFilePath)
 if err != nil {
	 return  fmt.Errorf("blueprint readme markdown is missing, create one using https://tinyurl.com/tf-mod-readme | error: %w", err)
 }
 /// Repo details validation

 // Must have headings
	mustHaveHeadings:= []markdownHeading{
	 {
		 headLevel:1,
		 headOrder:1,
		 headTitle:"",
		 content:false,
		 description:"Title of the module",

	 },
	}

 missingMustHaveHeadings:=missingHeadings(readmeContent,mustHaveHeadings)
 var errors []string
 if len(missingMustHaveHeadings)>0{
	 for _, heading :=range missingMustHaveHeadings{
		 e:= fmt.Sprintf("%s\t %s",heading.headTitle, heading.description)
		 errors = append(errors, e)
	 }
	 if len(errors) > 0 {
		 return fmt.Errorf("%s", strings.Join(errors, "\n"))
	 }
 }

 goodToHaveMarkdownHeadings:= []markdownHeading{
	 {
		 headLevel:1,
		 headOrder:1,
		 headTitle:"",
		 content:false,
		 description:"Title of the module",

	 },
	 {
		 headLevel:-1,
		 headOrder:-1,
		 headTitle:"Tagline",
		 content:true,
		 description:"corresponds to `BlueprintMetadataSpec.BlueprintInfo.description.Tagline` field in metadata.yaml",

	 },
 }
 missingGoodToHaveHeadings:=missingHeadings( readmeContent,goodToHaveMarkdownHeadings)


 if len(missingGoodToHaveHeadings)>0{
	 var warningMsg []string
	 for _, heading :=range missingGoodToHaveHeadings{
		 warningMsg=append(warningMsg,"# "+heading.headTitle+": \t"+ heading.description)
		 // Log.Warn("# "+heading.headTitle+": \t"+ heading.description)/
	 }
	 Log.Warn("\nGood to have headings- \n"+ strings.Join(warningMsg, "\n")+"\n")
 }

	otherHeadings:= []markdownHeading{
	 {
		 headLevel:-1,
		 headOrder:-1,
		 headTitle:"Tagline",
		 content:true,
		 description:"corresponds to `BlueprintMetadataSpec.BlueprintInfo.description.Tagline` field in metadata.yaml",

	 },
	 {
		 headLevel:-1,
		 headOrder:-1,
		 headTitle:"Detailed",
		 content:true,
		 description:"corresponds to `BlueprintMetadataSpec.BlueprintInfo.Description.Detailed` field in metadata.yaml",

	 },
	 {
		 headLevel:-1,
		 headOrder:-1,
		 headTitle:"PreDeploy",
		 content:true,
		 description:"corresponds to `BlueprintMetadataSpec.BlueprintInfo.Description.PreDeploy` field in metadata.yaml",

	 },

	 {
		 headLevel:-1,
		 headOrder:-1,
		 headTitle:"Architecture",
		 content:true,
		 description:"corresponds to `BlueprintMetadataSpec.BlueprintInfo.Description.Architecture` field in metadata.yaml",

	 },
	 {
		 headLevel:-1,
		 headOrder:-1,
		 headTitle:"Deployment Duration",
		 content:true,
		 description:"corresponds to `BlueprintMetadataSpec.BlueprintInfo.DeploymentDuration` field in metadata.yaml",

	 },
	 {
		 headLevel:-1,
		 headOrder:-1,
		 headTitle:"Cost",
		 content:true,
		 description:"corresponds to `BlueprintMetadataSpec.BlueprintInfo.CostEstimate` field in metadata.yaml",

	 },
	 {
		 headLevel:-1,
		 headOrder:-1,
		 headTitle:"Documentation",
		 content:true,
		 description:"corresponds to `BlueprintMetadataSpec.BlueprintContent.Documentation` field in metadata.yaml,"+
									 "if Documentation heading is present as a root chidren of markdown the very next paragraph is scanned,"+
									 "expecting to find an image or a link (as the diagram) followed immediately by a block of text"+
									 "(as the description). If this structure is present, it extracts the diagram URL and"+
									 "the description lines; otherwise, it does not include this field in the metadata.yaml.",

	 },
 }
 otherMissingHeadings:=missingHeadings( readmeContent,otherHeadings)

 if len(otherMissingHeadings)>0 {
	 var infoMsg []string
	 for _, heading :=range otherMissingHeadings{
		 infoMsg=append(infoMsg,"# "+heading.headTitle+ ":\t"+ heading.description )
	 }
	 Log.Info("\nOther missing headings-\n"+ strings.Join(infoMsg, "\n")+"\n")
 }
 Log.Info("\n****************Done Validating  README.md****************\n\n")
return nil
}

func missingHeadings(content []byte, markdownHeadings []markdownHeading)[]markdownHeading{
 missingHeadings:=[]markdownHeading{}
 for _, heading:= range markdownHeadings{
	 _, err := bpmetadata.GetMdContent(content, heading.headLevel, heading.headOrder, heading.headTitle, heading.content)
	 if err != nil {
		 missingHeadings=append(missingHeadings, heading)
	 }
 }
 return missingHeadings
}

func checkFilePresence(bpPath string, filePaths[] string) []string{

 missingFiles:=[]string{}
 for _, fileName := range filePaths{
	 filePath := path.Join(bpPath,fileName)
	 _, err := os.Stat(filePath)

	 if err != nil {
		 missingFiles = append(missingFiles, fileName)
	 }

 }
 return missingFiles
}
