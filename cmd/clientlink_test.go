package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCLICommandEventsMapDirectCommand(t *testing.T) {
	root := &cobra.Command{Use: "pmock-cli"}
	phone := &cobra.Command{Use: "phone"}
	start := &cobra.Command{Use: "start", Run: func(*cobra.Command, []string) {}}
	start.Flags().Int("id", 0, "")
	phone.AddCommand(start)
	root.AddCommand(phone)
	root.SetArgs([]string{"phone", "start", "--id", "123"})

	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	events := cliCommandEvents(start)
	if len(events) != 2 {
		t.Fatalf("events = %+v, want activity + applied", events)
	}
	if events[0].Kind != "activity" || events[0].Action != "powerOn" || events[0].PhoneID != "123" {
		t.Fatalf("activity = %+v", events[0])
	}
	if events[1].Kind != "applied" || len(events[1].Scopes) != 1 || events[1].Scopes[0] != "cloudPhones" {
		t.Fatalf("applied = %+v", events[1])
	}
}

func TestCLICommandEventsReadOnlyCommandIsSilent(t *testing.T) {
	root := &cobra.Command{Use: "pmock-cli"}
	phone := &cobra.Command{Use: "phone"}
	list := &cobra.Command{Use: "list", Run: func(*cobra.Command, []string) {}}
	phone.AddCommand(list)
	root.AddCommand(phone)
	root.SetArgs([]string{"phone", "list"})

	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if events := cliCommandEvents(list); len(events) != 0 {
		t.Fatalf("events = %+v, want none", events)
	}
}
