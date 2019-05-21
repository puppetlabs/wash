package primary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/puppetlabs/wash/cmd/internal/find/types"
	"github.com/stretchr/testify/suite"
)

type MetaPrimaryTestSuite struct {
	primaryTestSuite
	e types.Entry
}

func (s *MetaPrimaryTestSuite) TestMetaPrimaryErrors() {
	s.RETC("", "expected a key sequence")
	s.RETC("foo", "key sequences must begin with a '.'")
	s.RETC(".", "expected a key sequence after '.'")
	s.RETC(".[", "expected a key sequence after '.'")
	s.RETC(".key", "expected a predicate expression")
	s.RETC(".key +{", "expected.*closing.*}")
	s.RETC(".key]", `expected an opening '\['`)
	s.RETC(".key[", `expected a closing '\]'`)
	// Test some inner predicate expression parser errors
	s.RETC(".key1 .key2 (", `\(: missing closing '\)'`)
	s.RETC(".key1 .key2 ( -foo", "unknown predicate -foo")
	s.RETC(".key1 [?] (", `\(: missing closing '\)'`)
	s.RETC(".key1 [?] ( -foo", "unknown predicate -foo")
}

func (s *MetaPrimaryTestSuite) TestMetaPrimaryValidInputTruePredicates() {
	s.RTC(".architecture x86_64 -primary", "-primary", s.e)
	s.RTC(".blockDeviceMappings[?] .deviceName /dev/sda1 -primary", "-primary", s.e)
	s.RTC(".cpuOptions.coreCount 4 -primary", "-primary", s.e)
	s.RTC(".tags[?] .key termination_date -a .value +1h -primary", "-primary", s.e)
	s.RTC(".tags[?] .key foo -o .key department -primary", "-primary", s.e)
	s.RTC(".elasticGpuAssociations -null -primary", "-primary", s.e)
	s.RTC(".networkInterfaces[?] .association.ipOwnerID amazon -a .privateIpAddresses[?] .association.ipOwnerID amazon -primary", "-primary", s.e)
	// Test some inner predicate expressions
	s.RTC(".tags[?] .key ( foo -o department ) -primary", "-primary", s.e)
	s.RTC(".blockDeviceMappings[?] .ebs ( .attachTime +1h -a .status attached ) -primary", "-primary", s.e)
	s.RTC(".cpuOptions ( .coreCount ( ( -1 -a +5 ) -o 4 ) ) .threadsPerCore 1 -primary", "-primary", s.e)
}

func (s *MetaPrimaryTestSuite) TestMetaPrimaryValidInputFalsePredicates() {
	// Should fail b/c arch key does not exist
	s.RNTC(".arch x86_64 -primary", "-primary", s.e)
	// Should fail b/c architecture is a string, not a number
	s.RNTC(".architecture +10 -primary", "-primary", s.e)
	// Should fail b/c the termination_date tag's value is in the past while
	// +{1h} queries the future. 
	s.RNTC(".tags[?] .key termination_date -a .value +{1h} -primary", "-primary", s.e)
	// Should fail b/c the tags array has elements whose ".key" value is _not_ termination_date.
	// Informally, this means that this EC2 instance has more than one tag.
	s.RNTC(".tags[*] .key termination_date -primary", "-primary", s.e)
	// Should fail b/c architecture cannot be both a number and a string
	s.RNTC(".architecture +10 -a x86_64 -primary", "-primary", s.e)
}

func (s *MetaPrimaryTestSuite) TestMetaPrimaryValidInputNegation() {
	s.RNTC(".architecture ! x86_64 -primary", "-primary", s.e)
	s.RNTC(".blockDeviceMappings[?] ! .deviceName /dev/sda1 -primary", "-primary", s.e)
	s.RNTC(".cpuOptions.coreCount ! 4 -primary", "-primary", s.e)
	// De'Morgan's Law: !(p1(a) && p2(a)) == ! p1(a) || ! p2(a). Since there's more than one
	// tag (e.g. "department"), the negation of this predicate will evaluate to true.
	s.RTC(".tags[?] ! ( .key termination_date -a .value +1h ) -primary", "-primary", s.e)
	// De'Morgan's Law: !(p1(a) || p2(a)) == !p1(a) && !p2(a). Substituting, this translates to
	// "at least one tag that does _not_ have the "key" key set to 'foo' AND 'department'". Since
	// we have a tag with "key" set to "termination_date", and since "termination_date" is not "foo"
	// and "department", this predicate evaluates to true.
	s.RTC(".tags[?] ! ( .key foo -o .key department ) -primary", "-primary", s.e)
	s.RNTC(".elasticGpuAssociations ! -null -primary", "-primary", s.e)
	// There's only one network interface, so the negation here evaluates to false.
	s.RNTC(".networkInterfaces[?] ! ( .association.ipOwnerID amazon -a .privateIpAddresses[?] .association.ipOwnerID amazon ) -primary", "-primary", s.e)
	// Test negation for inner predicate expressions
	s.RNTC(".tags[0] .key ( ! ( foo -o department ) ) -primary", "-primary", s.e)
	s.RNTC(".cpuOptions ( .coreCount ( ( -1 -a +5 ) -o ! 4 ) ) .threadsPerCore 1 -primary", "-primary", s.e)
}

func (s *MetaPrimaryTestSuite) TestOptionsSetter() {
	opts := types.NewOptions()

	// Test that if maxdepth is not set, then the options setter
	// sets it to 1.
	metaPrimary.optionsSetter(&opts)
	s.Equal(int(1), opts.Maxdepth)

	// Test that if maxdepth is set, then the options setter
	// leaves it alone.
	opts.Maxdepth = 10
	opts.MarkAsSet(types.MaxdepthFlag)
	metaPrimary.optionsSetter(&opts)
	s.Equal(int(10), opts.Maxdepth)
}

func TestMetaPrimary(t *testing.T) {
	s := new(MetaPrimaryTestSuite)
	rawMeta, err := ioutil.ReadFile("testdata/metadata.json")
	if err != nil {
		t.Fatal(fmt.Sprintf("Failed to read testdata/metadata.json"))
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMeta, &m); err != nil {
		t.Fatal(fmt.Sprintf("Failed to unmarshal testdata/metadata.json: %v", err))
	}
	s.e.Attributes.SetMeta(m)
	s.Parser = metaPrimary
	s.ConstructEntry = func(v interface{}) types.Entry {
		return v.(types.Entry)
	}
	suite.Run(t, s)
}
