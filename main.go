package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sensu/sensu-go/types"
	"github.com/spf13/cobra"
)

var (
	webhook     string
	annotations string
	debug       bool
	stdin       *os.File
)

// Properties struct includes the two fields (subject and message) which Zenduty webhook expects to be POSTed to it in JSON
type Properties struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
	Status string `json:"status"`
	Id string `json:"id"`
}

// Payload struct includes Properties struct for post to Zenduty webhook URL
type Payload struct {
	Properties `json:"properties"`
}

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sensu-go-zenduty-handler",
		Short: "The Sensu Go Zenduty handler for incident alerting",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&webhook,
		"webhook",
		"w",
		os.Getenv("ZENDUTY_WEBHOOK"),
		"The Webhook URL, use default from ZENDUTY_WEBHOOK env var")

	cmd.Flags().StringVarP(&annotations,
		"withAnnotations",
		"a",
		os.Getenv("ZENDUTY_ANNOTATIONS"),
		"The Zenduty handler will parse check and entity annotations with these values. Use ZENDUTY_ANNOTATIONS env var with commas, like: documentation,playbook")

	cmd.Flags().BoolVarP(&debug,
		"debug",
		"d",
		false,
		"Enable debug mode, which prints JSON object which would be POSTed to the Zenduty webhook instead of actually POSTing it")

	_ = cmd.MarkFlagRequired("webhook")

	return cmd
}

// formattedEventAction func
func formattedEventAction(event *types.Event) string {
	switch event.Check.Status {
	case 0:
		return "RESOLVED"
	case 1:
		return "WARNING"
	case 2:
		return "CRITICAL"
	default:
		return "ALERT"
	}
}

// stringInSlice checks if a slice contains a specific string
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// parseAnnotations func try to find a predeterminated keys
func parseAnnotations(event *types.Event) string {
	var output string
	// localannotations := make(map[string]string)
	tags := strings.Split(annotations, ",")
	if event.Check.Annotations != nil {
		for key, value := range event.Check.Annotations {
			if stringInSlice(key, tags) {
				output += fmt.Sprintf("  %s: %s ,\n", key, value)
			}
		}
	}
	if event.Entity.Annotations != nil {
		for key, value := range event.Check.Annotations {
			if stringInSlice(key, tags) {
				output += fmt.Sprintf("  %s: %s ,\n", key, value)
			}
		}
	}
	return output
}

// Gets ID of the event to set as Zenduty Entity Id, if ID doesn't exist(?) check name is used.
func getID(event *types.Event) string {
	var output string
	// if len(event.ID)>0{
	// 	output = string(event.ID)
	// } else{
		output = event.Check.Name
	// }
	return output
}

// eventSubject func returns a one-line short summary
func eventSubject(event *types.Event) string {
	return fmt.Sprintf("The %s check has changed to %s on host %s", event.Check.Name, formattedEventAction(event), event.Entity.Name)
}

// eventDescription func returns a formatted message
func eventDescription(event *types.Event) string {
	return fmt.Sprintf("Server: %s, \nCheck: %s, \nStatus: %s, \nCheck Output: %s, \nAnnotation Information: %s\n", event.Entity.Name, event.Check.Name, formattedEventAction(event), event.Check.Output, parseAnnotations(event))
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	if webhook == "" {
		_ = cmd.Help()
		return fmt.Errorf("webhook is empty")
	}

	if annotations == "" {
		annotations = "documentation"
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	eventJSON, err := ioutil.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("failed to read stdin: %s", err)
	}

	event := &types.Event{}
	err = json.Unmarshal(eventJSON, event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal stdin data: %s", err)
	}

	if err = event.Validate(); err != nil {
		return fmt.Errorf("failed to validate event: %s", err)
	}

	if !event.HasCheck() {
		return fmt.Errorf("event does not contain check")
	}

	formPost := &Payload{
		Properties: Properties{
			Subject: eventSubject(event),
			Message: eventDescription(event),
			Status: formattedEventAction(event),
			Id : getID(event),
		},
	}
	bodymarshal, err := json.Marshal(formPost)
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
	}

	if debug == true {
		fmt.Printf("[DEBUG] JSON output: %s\n", bodymarshal)
		fmt.Println("[DEBUG] Not posting JSON object to Zenduty webhook since we're in debug mode ")
		os.Exit(1)
	}

	Post(webhook, bodymarshal)
	return nil
}

//Post func to send the json to zenduty "generic webhook" integration
func Post(url string, body []byte) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
	}
	if resp.StatusCode != 200 {
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("[ERROR] %s", err)
		}
		s := string(bodyText)
		fmt.Printf("[LOG]: %s ; %s", resp.Status, s)
	}
	defer resp.Body.Close()
	return nil
}
