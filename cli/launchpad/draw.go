package launchpad

import (
	"fmt"
	"io"
	"strings"
	"os"
)

// diagram represents a GCP draw diagram
type diagram struct {
	paths []*diagramPath
	group *diagramGroup
}

// diagramGroup represents a group of cards (including a top-level diagram)
type diagramGroup struct {
	d      *diagram
	name   string
	cards  []*diagramCard
	groups []*diagramGroup
}

// diagramPath represents a path within a GCP diagram
type diagramPath struct {
	from      *diagramCard
	to        *diagramCard
	connector string
}

// diagramCard represents a card within a GCP diagram
type diagramCard struct {
	id   string
	name string
}

// addCard adds a card into a diagram
func (dg *diagramGroup) addCard(id string, displayName string) (*diagramCard, error) {
	card := &diagramCard{
		id:   id,
		name: displayName,
	}
	dg.cards = append(dg.cards, card)
	return card, nil
}

// addPath adds a path into a diagram
func (dg *diagramGroup) addPath(from *diagramCard, to *diagramCard, connector string) (*diagramPath, error) {
	path := &diagramPath{
		from:      from,
		to:        to,
		connector: connector,
	}
	dg.d.paths = append(dg.d.paths, path)
	return path, nil
}

// addGroup adds a group to a diagram
func (dg *diagramGroup) addGroup(name string) (*diagramGroup, error) {
	group := &diagramGroup{
		d:    dg.d,
		name: name,
	}
	dg.groups = append(dg.groups, group)
	return group, nil
}

func (dg *diagramGroup) dump(ind int, buff io.Writer) error {
	indent := strings.Repeat(" ", ind)

	for _, card := range dg.cards {
		fmt.Fprintf(buff, "%scard generic as %s {\n", indent, card.id)
		cardIndent := strings.Repeat(" ", ind+defaultIndentSize)
		fmt.Fprintf(buff, "%sname \"%s\"\n", cardIndent, card.name)
		fmt.Fprintf(buff, "%s}\n", indent)
	}

	for _, child := range dg.groups {
		fmt.Fprintf(buff, "%sgroup %s {\n", indent, child.name)
		err := child.dump(ind+defaultIndentSize, buff)
		if err != nil {
			return err
		}
		fmt.Fprintf(buff, "%s}\n", indent)
	}
	return nil
}

func (path *diagramPath) dump(ind int, buff io.Writer) error {
	indent := strings.Repeat(" ", ind)
	fmt.Fprintf(buff, "%s%s%s%s\n", indent, path.from.id, path.connector, path.to.id)
	return nil
}

func (d *diagram) dump(buff io.Writer) error {
	fmt.Fprintf(buff, "elements {\n")
	d.group.dump(defaultIndentSize, buff)
	fmt.Fprintf(buff, "}\n")

	fmt.Fprintf(buff, "paths {\n")
	for _, path := range d.paths {
		path.dump(defaultIndentSize, buff)
	}
	fmt.Fprintf(buff, "}\n")

	return nil
}

// String implements Stringer and generates a string representation.
func (d *diagram) String() string {
	buff := &strings.Builder{}

	d.dump(buff)

	return buff.String()
}

func (ao *assembledOrg) makeDiagram() (*diagram, error) {
	orgDiagram := &diagram{}
	orgDiagram.group = &diagramGroup{
		d: orgDiagram,
	}

	ao.org.draw(orgDiagram.group, nil)

	return orgDiagram, nil
}

// draw prints diagram(s) for a given org
func (ao *assembledOrg) draw(fp *os.File) error {
	orgDiagram, err := ao.makeDiagram()
	if err != nil {
		return err
	}

	print(orgDiagram.String())

	orgDiagram.dump(fp)

	return nil
}
