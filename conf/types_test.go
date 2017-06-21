package conf

import (
	"fmt"
)

func ExampleReportStatus_String() {
	ss := []ReportStatus{
		ReportPendingStatus,
		ReportRunningStatus,
		ReportSucceededStatus,
		ReportFailedStatus,
		ReportUnknowStatus,
		1<<32 - 1}
	for _, s := range ss {
		fmt.Println(s)
	}

	// Output:
	// pending
	// running
	// succeeded
	// failed
	// unknow
	// unknow
}

func ExampleNode_String_name() {
	n := Node{
		Name:        "test",
		Identifiers: []string{"foo", "bar"},
	}
	fmt.Println(n)
	// Output:
	// test
}

func ExampleNode_String_identifiers() {
	n := Node{
		Identifiers: []string{"foo", "bar"},
	}
	fmt.Println(n)
	// Output:
	// foo,bar
}
