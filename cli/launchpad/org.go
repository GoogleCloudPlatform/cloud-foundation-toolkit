package launchpad

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type orgYAML struct {
	headerYAML `yaml:",inline"`
	Spec       orgSpecYAML `yaml:"spec"`
}

func (o *orgYAML) enroll(e *eval) error { return o.Spec.enroll(e) }
func (o *orgYAML) validate() error      { return o.Spec.validate() }

// orgSpecYAML defines Organization Kind's spec.
type orgSpecYAML struct {
	Id          string          `yaml:"id"`
	DisplayName string          `yaml:"displayName"`
	Folders     folderSpecYAMLs `yaml:"folders"`
}

func (o *orgSpecYAML) refId() string { return fmt.Sprintf("%s.%s", Organization, o.Id) }
func (o *orgSpecYAML) dump(ind int) []string {
	indent := strings.Repeat(" ", ind)
	rep := fmt.Sprintf("%s%s.%s (\"%s\")", indent, Organization, o.Id, o.DisplayName)
	buff := []string{rep}
	for _, sf := range o.Folders {
		buff = append(buff, sf.dump(ind+indentSize)...)
	}
	return buff
}

func (o *orgSpecYAML) validate() error {
	// TODO validate ORG spec
	return nil
}

// enroll adds organization into eval along with its sub-resources.
func (o *orgSpecYAML) enroll(e *eval) error {
	e.register(o, nil)
	for _, sf := range o.Folders {
		if err := sf.enroll(e); err != nil {
			return err
		}
	}
	return nil
}

func (o *orgSpecYAML) String() string {
	return strings.Join(o.dump(0), "\n")
}

// initializeByRef initializes an organization through other resource's reference.
func (o *orgSpecYAML) initializeByRef(ref *referenceYAML) error {
	if o.Id != "" && o.Id != ref.TargetId {
		return errors.New("unable to initialize organization to a different id")
	} else if o.Id == "" && ref.TargetId == "" {
		return errors.New("unset org id")
	}
	o.Id = ref.TargetId
	return nil
}

func (o *orgSpecYAML) initialize(oNew *orgSpecYAML) error {
	if o.Id != "" && o.Id != oNew.Id {
		return errors.New("unable to initialize organization to a different id")
	} else if o.Id == "" && oNew.Id == "" {
		return errors.New("unset org id")
	}
	// TODO validation?
	return nil
}

// resolveRefs adds references to this organization as its sub-resources.
func (o *orgSpecYAML) resolveRefs(refs []reference) error {
	for _, r := range refs {
		sf, ok := r.srcPtr.(*folderSpecYAML)
		if !ok {
			return errors.New("unable to add non folder as a sub-folder")
		}
		o.Folders.add(sf)
	}
	return nil
}
