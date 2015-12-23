package surefire

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="TestSuite" time="0.529" tests="3" errors="0" skipped="0" failures="0">
  <properties>
    <property name="java.runtime.name" value="Java(TM) SE Runtime Environment"/>
    <property name="sun.boot.library.path" value="/Library/Java/JavaVirtualMachines/jdk1.8.0_66.jdk/Contents/Home/jre/lib"/>
  </properties>
  <testcase name="testApply" classname="org.btrplace.plan.ActionTest" time="0.045"/>
  <testcase name="testBasics" classname="org.btrplace.plan.ActionTest" time="2"/>
  <testcase name="testEvents" classname="org.btrplace.plan.ActionTest" time="0"/>
</testsuite>
`

func TestDecode(t *testing.T) {
	in := strings.NewReader(data)
	suite, err := Decode(in)
	assert.Nil(t, err)
	log.Println(suite)
}
