package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"reflect"
	"sigs.k8s.io/prow/prow/flagutil"
	jira2 "sigs.k8s.io/prow/prow/jira"
	"strings"
)

const (
	bugLink = `[Jira Issue %s](%s/browse/%s)`
)

type options struct {
	jira flagutil.JiraOptions
}

func main() {
	original := flag.CommandLine
	klog.InitFlags(original)
	original.Set("alsologtostderr", "true")
	original.Set("v", "2")

	opt := &options{}
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, arguments []string) {
			if err := opt.Run(); err != nil {
				klog.Exitf("error: %v", err)
			}
		},
	}
	flagset := cmd.Flags()

	goFlagSet := flag.NewFlagSet("prowflags", flag.ContinueOnError)
	opt.jira.AddFlags(goFlagSet)
	flagset.AddGoFlagSet(goFlagSet)

	flagset.AddGoFlag(original.Lookup("v"))

	if err := cmd.Execute(); err != nil {
		klog.Exitf("error: %v", err)
	}
}

func (o *options) Run() error {
	c, err := o.jira.Client()
	if err != nil {
		klog.Fatalf("Unable to create jira client: %w", err)
	}

	// issue, err := c.GetIssue("TRT-1716")
	issue, err := c.GetIssue("OCPBUGS-35865")
	if err != nil {
		klog.Errorf("Unable to get jira issue: %w", err)
		return err
	}

	b, err := json.MarshalIndent(issue, "", "    ")
	if err != nil {
		klog.Errorf("unable to marshal Jira Issue: %v", err)
		return nil
	}
	klog.V(2).Infof("Retrieved issue:\n%s", string(b))

	fields := getCustomFields(*issue)

	b, err = json.MarshalIndent(fields, "", "    ")
	if err != nil {
		klog.Errorf("failed to marshall the custon fields struct. Error: %v", err)
		return nil
	}
	klog.V(2).Infof("Custom Fields:\n%s", string(b))

	b, err = json.Marshal(fields)
	if err != nil {
		klog.Errorf("failed to marshall the custon fields struct. Error: %v", err)
		return nil
	}
	klog.V(2).Infof("Custom Fields:\n%s", string(b))

	//targetVersion, err := helpers.GetIssueTargetVersion(issue)
	//
	//klog.V(2).Infof("Retrieved issue target version:\n %v", targetVersion)

	// create deep copy of parent "Fields" field
	//data, err := json.Marshal(issue.Fields)
	//if err != nil {
	//	klog.Errorf("unable to marshal Jira Issue Fields: %v", err)
	//	return nil
	//}
	//childIssueFields := &jira.IssueFields{}
	//err = json.Unmarshal(data, childIssueFields)
	//if err != nil {
	//	klog.Errorf("unable to unmarshal Jira Issue Fields: %v", err)
	//	return nil
	//}
	//childIssue := &jira.Issue{
	//	Fields: childIssueFields,
	//}
	//// update description
	//childIssue.Fields.Description = fmt.Sprintf("This is a clone of issue %s. The following is the description of the original issue: \n---\n%s", issue.Key, issue.Fields.Description)
	//
	//// attempt to create the new issue
	//createdIssue, err := c.CreateIssue(childIssue)
	//if err != nil {
	//	// some fields cannot be set on creation; unset them
	//	if JiraErrorStatusCode(err) != 400 {
	//		klog.Errorf("Houstan, we have a problem: %v", err)
	//		return nil
	//	}
	//	var newErr error
	//	childIssue, newErr = unsetProblematicFields(childIssue, JiraErrorBody(err))
	//	if newErr != nil {
	//		// in this situation, it makes more sense to just return the original error; any error from unsetProblematicFields will be
	//		// a json marshalling error, indicating an error different from the standard non-settable fields error. The error from
	//		// unsetProblematicFields is not useful in these cases
	//		klog.Errorf("Houstan, we have another problem: %v", err)
	//		return nil
	//	}
	//
	//	b, err = json.MarshalIndent(childIssue, "", "    ")
	//	if err != nil {
	//		klog.Errorf("unable to marshal child Jira Issue: %v", err)
	//		return nil
	//	}
	//	klog.V(2).Infof("Child issue:\n %s", string(b))
	//
	//	createdIssue, err = c.CreateIssue(childIssue)
	//	if err != nil {
	//		klog.Errorf("Houstan, we have yet another problem: %v", err)
	//		return nil
	//	}
	//}
	//
	//b, err = json.MarshalIndent(createdIssue, "", "    ")
	//if err != nil {
	//	klog.Errorf("unable to marshal Created Jira Issue: %v", err)
	//	return nil
	//}
	//klog.V(2).Infof("Cloned issue:\n %s", string(b))

	//clone, err := c.CloneIssue(issue)
	//if err != nil {
	//	klog.Errorf("Unable to clone jira issue: %w", err)
	//	return err
	//}
	//
	//b, err = json.MarshalIndent(clone, "", "    ")
	//if err != nil {
	//	klog.Errorf("unable to marshal Jira Issue: %v", err)
	//	return nil
	//}

	//klog.V(2).Infof("%q, %q, %q", issue.Key, c.JiraURL(), issue.Key)
	//
	//klog.V(2).Infof(`This pull request references `+bugLink+`, which is valid.`, issue.Key, c.JiraURL(), issue.Key)
	//klog.V(2).Infof(`This pull request references %s, which is valid.`, generateMarkdownLink(c.JiraURL(), issue.Key))

	//klog.V(2).Infof("Retrieved issue:\n %s", string(b))

	//projects, err := c.ListProjects()
	//if err != nil {
	//	klog.Errorf("Unable to get jira projects: %w", err)
	//	return err
	//}
	//
	//for _, project := range *projects {
	//	klog.V(2).Infof("Found project: [%s] %s", project.Key, project.Name)
	//}
	return nil
}

func JiraErrorStatusCode(err error) int {
	if jiraErr := (&jira2.JiraError{}); errors.As(err, &jiraErr) {
		return jiraErr.StatusCode
	}
	jiraErr, ok := err.(*jira2.JiraError)
	if !ok {
		return -1
	}
	return jiraErr.StatusCode
}

func JiraErrorBody(err error) string {
	if jiraErr := (&jira2.JiraError{}); errors.As(err, &jiraErr) {
		return jiraErr.Body
	}
	jiraErr, ok := err.(*jira2.JiraError)
	if !ok {
		return ""
	}
	return jiraErr.Body
}

type createIssueError struct {
	ErrorMessages []string          `json:"errorMessages"`
	Errors        map[string]string `json:"errors"`
}

func unsetProblematicFields(issue *jira.Issue, responseBody string) (*jira.Issue, error) {
	// handle unsettable "unknown" fields
	processedResponse := createIssueError{}
	if newErr := json.Unmarshal([]byte(responseBody), &processedResponse); newErr != nil {
		return nil, fmt.Errorf("Error processing jira error: %w", newErr)
	}
	// turn issue into map to simplify unsetting process
	marshalledIssue, err := json.Marshal(issue)
	if err != nil {
		return nil, err
	}
	issueMap := make(map[string]interface{})
	if err := json.Unmarshal(marshalledIssue, &issueMap); err != nil {
		return nil, err
	}
	fieldsMap := issueMap["fields"].(map[string]interface{})
	for field := range processedResponse.Errors {
		delete(fieldsMap, field)
	}
	// Remove null value customfields because they are causing a: 500 Internal Server Error
	for field, value := range fieldsMap {
		if strings.HasPrefix(field, "customfield_") && value == nil {
			delete(fieldsMap, field)
		}
	}
	issueMap["fields"] = fieldsMap
	// turn back into jira.Issue type
	marshalledFixedIssue, err := json.Marshal(issueMap)
	if err != nil {
		return nil, err
	}
	newIssue := jira.Issue{}
	if err := json.Unmarshal(marshalledFixedIssue, &newIssue); err != nil {
		return nil, err
	}
	return &newIssue, nil
}

func generateMarkdownLink(url, id string) string {
	link := `[Jira Issue %s](%s/browse/%s)`
	if strings.HasSuffix(url, "/") {
		link = `[Jira Issue %s](%sbrowse/%s)`
	}
	return fmt.Sprintf(link, id, url, id)
}

type CustomField struct {
	FieldName   string  `bigquery:"fieldName" json:"fieldName,omitempty"`
	ID          string  `bigquery:"id" json:"id,omitempty"`
	Name        string  `bigquery:"name" json:"name,omitempty"`
	Key         string  `bigquery:"key" json:"key,omitempty"`
	DisplayName string  `bigquery:"displayName" json:"displayName,omitempty"`
	Description string  `bigquery:"description" json:"description,omitempty"`
	Value       string  `bigquery:"value" json:"value,omitempty"`
	Votes       float64 `bigquery:"votes" json:"votes,omitempty"`
}

func getCustomFields(i jira.Issue) []CustomField {
	var customFields []CustomField
	klog.Infof("Issue: %s", i.Key)
	for k, v := range i.Fields.Unknowns {
		if v == nil {
			continue
		}
		klog.Infof("Field: %s -> %s", k, v)
		field := processCustomFieldValue(k, v)
		if field != nil {
			customFields = append(customFields, *field)
		}
	}
	return customFields
}

func processCustomFieldValue(name string, value interface{}) *CustomField {
	var field *CustomField
	var fields []CustomField
	var valueStr string

	switch t := value.(type) {
	case int:
		valueStr = fmt.Sprintf("%d", value)
	case float64:
		valueStr = fmt.Sprintf("%f", value)
	case string:
		valueStr = fmt.Sprintf("%s", value)
	case bool:
		valueStr = fmt.Sprintf("%t", value)
	case []interface{}:
		for _, n := range t {
			cf := getCustomField(name, n)
			if cf != nil {
				fields = append(fields, *cf)
			}
		}
	case map[string]interface{}:
		field = getCustomField(name, value)
	default:
		var r = reflect.TypeOf(t)
		klog.Warningf("Unknown CustomField type: %v", r)
		return nil
	}

	switch {
	case field != nil:
		field.FieldName = name
		return field
	case fields != nil && len(fields) > 0:
		b, err := json.Marshal(fields)
		if err != nil {
			klog.Errorf("failed to marshall the fields struct for %s. Error: %v", name, err)
			return nil
		}
		return &CustomField{
			FieldName: name,
			Value:     fmt.Sprint(string(b)),
		}
	case len(valueStr) > 0:
		return &CustomField{
			FieldName: name,
			Value:     valueStr,
		}
	default:
		return nil
	}
}

func getCustomField(name string, value interface{}) *CustomField {
	field := &CustomField{}
	switch v := value.(type) {
	case string:
		field.Value = v
	default:
		bytes, err := json.Marshal(value)
		if err != nil {
			klog.Errorf("failed to process the custom field %s. Error : %v", name, err)
			return nil
		}
		if err = json.Unmarshal(bytes, field); err != nil {
			klog.Errorf("failed to unmarshall the json to struct for %s. Error: %v", name, err)
			return nil
		}
	}
	return field
}
