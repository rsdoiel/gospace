//
// Package aspace is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
// copyright (c) 2015
// Caltech Library
//
package aspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// The library version
var Version = "0.0.0"

// ArchivesSpaceAPI is a struct holding the essentials for communicating
// with the ArchicesSpace REST API
type ArchivesSpaceAPI struct {
	URL       *url.URL `json:"api_url"`
	Username  string   `json:"username,omitempty"`
	Password  string   `json:"password,omitempty"`
	AuthToken string   `json:"token,omitempty"`
}

// ResponseMsg is a structure to hold the JSON portion of a response from the ArchivesSpaceAPI
type ResponseMsg struct {
	Status      string      `json:"status,omitempty"`
	ID          int         `json:"id,omitempty"`
	LockVersion int         `json:"lock_version,omitempty"`
	Stale       interface{} `json:"stale,omitempty"`
	URI         string      `json:"uri,omitempty"`
	Warnings    []string    `json:"warnings,omitempty"`
	Error       interface{} `json:"error,omitempty"`
}

// Repository represents an ArchivesSpace repository from the client point of view
type Repository struct {
	JSONModelType         string                 `json:"json_model_type,omitempty"`
	ID                    int                    `json:"id,omitempty"`
	RepoCode              string                 `json:"repo_code"`
	Name                  string                 `json:"name"`
	URI                   string                 `json:"uri,omitempty"`
	URL                   string                 `json:"url,omitempty"`
	AgentRepresentation   map[string]interface{} `json:"agent_representation,omitempty"`
	Country               string                 `json:"country,omitempty"`
	ImageURL              string                 `json:"image_url,omitempty"`
	OrgCode               string                 `json:"org_code,omitempty"`
	ParentInstitutionName string                 `json:"parent_institution_name,omitempty"`
	LockVersion           int                    `json:"lock_version"`
	CreatedBy             string                 `json:"created_by,omitempty"`
	CreateTime            string                 `json:"create_time,omitempty"`
	SystemMTime           string                 `json:"system_mtime,omitempty"`
	UserMTime             string                 `json:"user_mtime,omitempty"`
}

// Date an ArchivesSpace Date structure
type Date struct {
	JSONModelType  string `json:"jsonmodel_type,omitempty"`
	LockVersion    int    `json:"lock_version"`
	Begin          string `json:"begin,omitempty"`
	End            string `json:"end,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	CreateTime     string `json:"create_time,omitempty"`
	SystemMTime    string `json:"system_mtime,omitempty"`
	UserMTime      string `json:"user_mtime,omitempty"`
	LastModifiedBy string `json:"last_modified_by,omitempty"`
	DateType       string `json:"date_type,omitempty"`
	Label          string `json:"label,omitempty"`
}

// NoteText is the content type of subnotes
type NoteText struct {
	JSONModelType string `json:"jsonmodel_type,omitempty"`
	Content       string `json:"content,omitempty"`
	Publish       bool   `json:"publish,omitempty"`
}

// NoteBiogHist - Notes Biographical Historical
type NoteBiogHist struct {
	JSONModelType string     `json:"jsonmodel_type,omitempty"`
	Label         string     `json:"label,omitempty"`
	PersistentID  string     `json:"persistent_id,omitempty"`
	SubNotes      []NoteText `json:"subnotes,omitempty"`
	Publish       bool       `json:"publish,omitempty"`
}

// NamePerson a single agent name structure
type NamePerson struct {
	JSONModelType        string `json:"json_model_type,omitempty"`
	LockVersion          int    `json:"lock_version"`
	PrimaryName          string `json:"primary_name,omitempty"`
	RestOfName           string `json:"rest_of_name,omitempty"`
	SortName             string `json:"sort_name,omitempty"`
	SortNameAutoGenerate bool   `json:"sort_name_auto_generate,omitempty"`
	CreatedBy            string `json:"created_by,omitempty"`
	CreateTime           string `json:"create_time,omitempty"`
	SystemMTime          string `json:"system_mtime,omitempty"`
	UserMTime            string `json:"user_mtime,omitempty"`
	LastModifiedBy       string `json:"last_modified_by,omitempty"`
	Authorized           bool   `json:"authorized,omitempty"`
	IsDisplayName        bool   `json:"is_display_name,omitempty"`
	Source               string `json:"source,omitempty"`
	Rules                string `json:"rules,omitempty"`
	NameOrder            string `json:"name_order,omitempty"`
	UseDates             []Date `json:"use_dates,omitempty"`
}

// AgentContact is a JSONModel for the AgentContacts array/map
type AgentContact struct {
	JSONModelType  string   `json:"json_model_type,omitempty"`
	LockVersion    int      `json:"lock_version"`
	Name           string   `json:"name,omitempty"`
	CreatedBy      string   `json:"created_by,omitempty"`
	CreateTime     string   `json:"create_time,omitempty"`
	SystemMTime    string   `json:"system_mtime,omitempty"`
	UserMTime      string   `json:"user_mtime,omitempty"`
	LastModifiedBy string   `json:"last_modified_by,omitempty"`
	Telephones     []string `json:"telephones,omitempty"`
}

// Agent represents an ArchivesSpace complete agent record from the client point of view
type Agent struct {
	JSONModelType             string          `json:"json_model_type,omitempty"`
	LockVersion               int             `json:"lock_version"`
	ID                        int             `json:"id,omitempty"`
	Published                 bool            `json:"publish,omitempty"`
	CreatedBy                 string          `json:"created_by,omitempty"`
	CreateTime                string          `json:"create_time,omitempty"`
	SystemMTime               string          `json:"system_mtime,omitempty"`
	UserMTime                 string          `json:"user_mtime,omitempty"`
	LastModifiedBy            string          `json:"last_modified_by,omitempty"`
	AgentType                 string          `json:"agent_type,omitempty"`
	URI                       string          `json:"uri,omitempty"`
	Title                     string          `json:"title,omitempty"`
	IsLinkedToPublishedRecord bool            `json:"is_linked_to_published_record,omitempty"`
	Names                     []*NamePerson   `json:"names,omitempty"`
	DisplayName               *NamePerson     `json:"display_name,omitempty"`
	RelatedAgents             []interface{}   `json:"related_agents,omitempty"`
	DatesOfExistance          []Date          `json:"dates_of_existence,omitempty"`
	AgentContacts             []*AgentContact `json:"agent_contacts,omitempty"`
	LinkedAgentRoles          []interface{}   `json:"linked_agent_roles,omitempty"`
	ExternalDocuments         []interface{}   `json:"external_documents,omitempty"`
	RightsStatements          []interface{}   `json:"rights_statements,omitempty"`
	Notes                     []*NoteBiogHist `json:"notes,omitempty"`
}

// ExternalID represents an external ID as found in Accession records
type ExternalID struct {
	JSONModelType  string `json:"json_model_type,omitempty"`
	ID             string `json:"external_id,omitempty"`
	Source         string `json:"source,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	CreateTime     string `json:"create_time,omitempty"`
	SystemMTime    string `json:"system_mtime,omitempty"`
	UserMTime      string `json:"user_mtime,omitempty"`
	LastModifiedBy string `json:"last_modified_by,omitempty"`
}

// Extent represents an extends json model found in Accession records
type Extent struct {
	JSONModelType   string `json:"json_model_type,omitempty"`
	LockVersion     int    `json:"lock_version"`
	CreatedBy       string `json:"created_by,omitempty"`
	CreateTime      string `json:"create_time,omitempty"`
	SystemMTime     string `json:"system_mtime,omitempty"`
	UserMTime       string `json:"user_mtime,omitempty"`
	LastModifiedBy  string `json:"last_modified_by,omitempty"`
	Number          string `json:"number,omitempty"`
	PhysicalDetails string `json:"physical_details"`
	Portion         string `json:"portion,omitempty"`
	ExtendType      string `json:"extent_type,omitempty"`
}

// UserDefined struct used in accession records for holding user defined data.
type UserDefined struct {
	JSONModelType  string            `json:"json_model_type,omitempty"`
	LockVersion    int               `json:"lock_version"`
	CreatedBy      string            `json:"created_by,omitempty"`
	CreateTime     string            `json:"create_time,omitempty"`
	SystemMTime    string            `json:"system_mtime,omitempty"`
	UserMTime      string            `json:"user_mtime,omitempty"`
	LastModifiedBy string            `json:"last_modified_by,omitempty"`
	Boolean1       bool              `json:"boolean_1,omitempty"`
	Boolean2       bool              `json:"boolean_2,omitempty"`
	Boolean3       bool              `json:"boolean_3,omitempty"`
	Boolean4       bool              `json:"boolean_4,omitempty"`
	Boolean5       bool              `json:"boolean_5,omitempty"`
	Text1          string            `json:"text_1,omitempty"`
	Text2          string            `json:"test_2,omitempty"`
	Text3          string            `json:"text_3,omitempty"`
	Text4          string            `json:"text_4,omitempty"`
	Text5          string            `json:"text_5,omitempty"`
	Repository     map[string]string `json:"repository"`
}

// Accession represents an accession record in ArchivesSpace from the client point of view
type Accession struct {
	JSONModelType       string                   `json:"json_model_type,omitempty"`
	LockVersion         int                      `json:"lock_version"`
	ID                  int                      `json:"id,omitempty"`
	Suppressed          bool                     `json:"suppressed,omitempty"`
	Title               string                   `json:"title,omitempty"`
	DisplayString       string                   `json:"display_string,omitempty"`
	Publish             bool                     `json:"publish,omitempty"`
	ContentDescription  string                   `json:"content_description,omitempty"`
	Provenance          string                   `json:"provenance,omitempty"`
	AccessionDate       string                   `json:"accession_date,omitempty"`
	RestrictionsApply   bool                     `json:"restrictions_apply,omitempty"`
	UseRestrictions     bool                     `json:"use_restrictions,omitempty"`
	CreatedBy           string                   `json:"created_by,omitempty,omitempty"`
	CreateTime          string                   `json:"create_time,omitempty,omitempty"`
	SystemMTime         string                   `json:"system_mtime,omitempty,omitempty"`
	UserMTime           string                   `json:"user_mtime,omitempty,omitempty"`
	LastModifiedBy      string                   `json:"last_modified_by,omitempty"`
	ID0                 string                   `json:"id_0,omitempty"`
	ID1                 string                   `json:"id_1,omitempty"`
	ExternalIDs         []ExternalID             `json:"external_ids,omitempty"`
	RelelatedAccessions []map[string]interface{} `json:"related_accessions,omitempty"`
	Classifications     []map[string]interface{} `json:"classifications,omitempty"`
	Subjects            []map[string]interface{} `json:"subjects,omitempty"`
	LinkedEvents        []map[string]interface{} `json:"linked_events,omitempty"`
	Extents             []Extent                 `json:"extents,omitempty"`
	Dates               []string                 `json:"dates,omitempty"`
	ExternalDocuments   []map[string]interface{} `json:"external_documents,omitempty"`
	RightsStatements    []map[string]interface{} `json:"rights_statements,omitempty"`
	Deaccessions        []map[string]interface{} `json:"deaccessions,omitempty"`
	RelelatedResources  []map[string]interface{} `json:"related_resources,omitempty"`
	LinkedAgents        []Agent                  `json:"linked_agents,omitempty"`
	Instances           []map[string]interface{} `json:"instances,omitempty"`
	URI                 string                   `json:"uri,omitempty"`
	Repository          map[string]string        `json:"repository,omitempty"`
	UserDefined         map[string]interface{}   `json:"user_defined,omitempty"`
}

func checkEnv(apiURL, username, password string) bool {
	if strings.TrimSpace(apiURL) == "" {
		return false
	}
	if strings.TrimSpace(username) == "" {
		return false
	}
	if strings.TrimSpace(password) == "" {
		return false
	}
	return true
}

// New creates a new ArchivesSpaceAPI object for use with most of the functions
// in the gas package.
func New(apiURL, username, password string) *ArchivesSpaceAPI {
	aspace := new(ArchivesSpaceAPI)
	if checkEnv(apiURL, username, password) == false {
		log.Fatal("Cannot create a new ArchivesSpace API connection")
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		log.Fatalf("ArchivesSpace connection not configured, %s", err)
	}
	aspace.URL = u
	aspace.Username = username
	aspace.Password = password
	aspace.AuthToken = ""
	return aspace
}

// IsAuth returns true if the auth token has been set, false otherwise
func (aspace *ArchivesSpaceAPI) IsAuth() bool {
	if aspace.AuthToken == "" {
		return false
	}
	return true
}

// Login authenticates against the ArchivesSpace REST API setting the AuthToken
// value in the ArchivesSpaceAPI struct.
func (aspace *ArchivesSpaceAPI) Login() error {
	// See https://golang.org/pkg/net/url/#pkg-examples for example building a URL from parts.
	// Command line example: curl -F "password=admin" "http://localhost:8089/users/admin/login"
	var data map[string]interface{}

	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = fmt.Sprintf("/users/%s/login", aspace.Username)
	form := url.Values{}
	form.Add("password", aspace.Password)

	res, err := http.PostForm(aspace.URL.String(), form)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.Status != "200 OK" {
		return fmt.Errorf("ArchivesSpace returned HTTP status %s", res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("ArchivesSpace return unreadable body: %s", err)
	}

	if err = json.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("Can't process JSON response %s\n\t%s", content, err)
	}
	aspace.AuthToken = data["session"].(string)
	return nil
}

// Logout clear the authentication token for the session with the API
func (aspace *ArchivesSpaceAPI) Logout() error {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = `/logout`
	client := &http.Client{}
	req, err := http.NewRequest("GET", aspace.URL.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-ArchivesSpace-Session", aspace.AuthToken)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// API the common HTTP request processing for interacting with ArchivesSpaceAPI
func (aspace *ArchivesSpaceAPI) API(method string, url string, payload io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("Can't create request: %s", err)
	}
	req.Header.Add("X-ArchivesSpace-Session", aspace.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	if method == "POST" {
		res, err := client.Do(req)
		defer res.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("Request error: %s", err)
		}
		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Read body error: %s", err)
		}
		return content, nil
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Request error: %s", err)
	}
	if res.Status != "200 OK" {
		return nil, fmt.Errorf("ArchiveSpace API error %s", res.Status)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body error: %s", err)
	}
	return content, nil
}

// CreateRepository will create a respository via the REST API for the
// ArchivesSpace instance defined in the ArchivesSpaceAPI struct.
// It will return the created record.
func (aspace *ArchivesSpaceAPI) CreateRepository(repo *Repository) (*ResponseMsg, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = "/repositories"
	payload, err := json.Marshal(repo)
	if err != nil {
		return nil, err
	}
	content, err := aspace.API("POST", aspace.URL.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	// content should look something like
	// {"status":"Created","id":3,"lock_version":0,"stale":null,"uri":"/repositories/3","warnings":[]}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetRepository returns the repository details based on Id
func (aspace *ArchivesSpaceAPI) GetRepository(id int) (*Repository, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = fmt.Sprintf(`/repositories/%d`, id)

	content, err := aspace.API("GET", aspace.URL.String(), nil)
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	repo.ID = id
	err = json.Unmarshal(content, repo)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// UpdateRepository takes a respository structure and sends it to the ArchivesSpace REST API
func (aspace *ArchivesSpaceAPI) UpdateRepository(repo *Repository) (*ResponseMsg, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = repo.URI
	payload, err := json.Marshal(repo)
	if err != nil {
		return nil, fmt.Errorf("Can't JSON encode update %v %s", repo, err)
	}
	content, err := aspace.API("POST", aspace.URL.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	// content should look something like
	// {"status":"Created","id":3,"lock_version":0,"stale":null,"uri":"/repositories/3","warnings":[]}
	// OR
	// {"error":"Some error message here"}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("Could not unpack UpdateRepository() response [%s] %s", content, err)
	}
	return data, nil
}

// DeleteRepository takes a respository structure and sends it to the ArchivesSpace REST API
func (aspace *ArchivesSpaceAPI) DeleteRepository(repo *Repository) (*ResponseMsg, error) {
	/*
		Example Listing the repo with curl:
			curl -H "X-ArchivesSpace-Session: $TOKEN" --request GET "http://localhost:8089/repositories" | python -m json.tool

		Example Delete the repo with curl:
			curl -H "X-ArchivesSpace-Session: $TOKEN" -d '{"repo_code": "1448043078"}' --request DELETE "http://localhost:8089/repositories/8" | python -m json.tool
	*/
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = fmt.Sprintf("/repositories/%d", repo.ID)
	payload, err := json.Marshal(repo)
	if err != nil {
		return nil, fmt.Errorf("Can't JSON encode update %v %s", repo, err)
	}
	content, err := aspace.API("DELETE", aspace.URL.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("Cannnot decode DeleteRepository() response %s", err)
	}
	return data, nil
}

// ListRepositories returns a list of repositories available via the ArchivesSpace REST API
func (aspace *ArchivesSpaceAPI) ListRepositories() ([]Repository, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = `/repositories`

	content, err := aspace.API("GET", aspace.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	// content should look something like
	// {"lock_version":0,"repo_code":"1447893780","name":"This is a test generated from go_test","created_by":"admin","last_modified_by":"admin","create_time":"2015-11-19T00:43:00Z","system_mtime":"2015-11-19T00:43:00Z","user_mtime":"2015-11-19T00:43:00Z","jsonmodel_type":"repository","uri":"/repositories/16","agent_representation":{"ref":"/agents/corporate_entities/15"}}
	var repos []Repository
	err = json.Unmarshal(content, &repos)
	if err != nil {
		return nil, err
	}
	// Now I need to populate the repos[?].ID fields
	for i := range repos {
		if id, err := strconv.Atoi(strings.TrimPrefix(repos[i].URI, "/repositories/")); err == nil {
			repos[i].ID = id
		}
	}
	return repos, nil
}

// CreateAgent creates a Agent recod via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) CreateAgent(aType string, agent *Agent) (*ResponseMsg, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = fmt.Sprintf("/agents/%s", aType)
	payload, err := json.Marshal(agent)
	if err != nil {
		return nil, err
	}
	content, err := aspace.API("POST", aspace.URL.String(), bytes.NewReader(payload))

	// content should look something like
	// {"status":"Created","id":5,"lock_version":0,"stale":true,"uri":"/agents/people/5","warnings":[]}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetAgent return an Agent via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) GetAgent(agentType string, agentID int) (*Agent, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = fmt.Sprintf(`/agents/%s/%d`, agentType, agentID)

	content, err := aspace.API("GET", aspace.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	agent := new(Agent)
	err = json.Unmarshal(bytes.TrimSpace(content), &agent)
	if err != nil {
		return nil, err
	}
	// Make sure the ID comes from agent.URI
	p := strings.Split(agent.URI, "/")
	id, err := strconv.Atoi(p[len(p)-1])
	if err != nil {
		return agent, err
	}
	agent.ID = id
	return agent, nil
}

// UpdateAgent creates a Agent recod via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) UpdateAgent(agentRequest *Agent) (*ResponseMsg, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = agentRequest.URI
	payload, err := json.Marshal(agentRequest)
	if err != nil {
		return nil, err
	}
	content, err := aspace.API("POST", aspace.URL.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// content should look something like
	// {"status":"Updated","id":13,"lock_version":1,"stale":true,"uri":"/agents/people/13"}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteAgent creates a Agent record via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) DeleteAgent(agentRequest *Agent) (*ResponseMsg, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = agentRequest.URI
	payload, err := json.Marshal(agentRequest)
	if err != nil {
		return nil, err
	}
	content, err := aspace.API("DELETE", aspace.URL.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// content should look something like
	// {"status":"Deleted","id":13}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode DeleteAgent() %s", err)
	}
	return data, nil
}

// ListAgents return an array of Agents via the ArchivesSpace API
func (aspace *ArchivesSpaceAPI) ListAgents(agentType string) ([]int, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = fmt.Sprintf(`/agents/%s`, agentType)
	q := aspace.URL.Query()
	q.Set("all_ids", "true")
	aspace.URL.RawQuery = q.Encode()

	content, err := aspace.API("GET", aspace.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	// content should look something like
	// [1,2,3,4]
	var agentIDs []int
	err = json.Unmarshal(content, &agentIDs)
	if err != nil {
		return nil, err
	}
	return agentIDs, nil
}

// CreateAccession creates a new Accession record in a Repository
func (aspace *ArchivesSpaceAPI) CreateAccession(repoID int, accession *Accession) (*ResponseMsg, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = fmt.Sprintf("/repositories/%d/accessions", repoID)
	payload, err := json.Marshal(accession)
	if err != nil {
		return nil, err
	}
	content, err := aspace.API("POST", aspace.URL.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// content should look something like
	// {"status":"Created","id":5,"lock_version":0,"stale":true,"uri":"/repositories/28/accessions/5","warnings":[]}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetAccession retrieves an Accession record from a Repository
func (aspace *ArchivesSpaceAPI) GetAccession(repoID, accessionID int) (*Accession, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = fmt.Sprintf("/repositories/%d/accessions/%d", repoID, accessionID)

	content, err := aspace.API("GET", aspace.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	// content should look something like
	// {"status":"Created","id":5,"lock_version":0,"stale":true,"uri":"/repositories/2/accessions/5","warnings":[]}
	accession := new(Accession)
	err = json.Unmarshal(content, accession)
	if err != nil {
		return nil, err
	}
	p := strings.Split(accession.URI, "/")
	accession.ID, err = strconv.Atoi(p[len(p)-1])
	if err != nil {
		return accession, fmt.Errorf("Accession ID parse error %d %s", accession.ID, err)
	}
	return accession, nil
}

// UpdateAccession updates an existing Accession record in a Repository
func (aspace *ArchivesSpaceAPI) UpdateAccession(accession *Accession) (*ResponseMsg, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = accession.URI
	payload, err := json.Marshal(accession)
	if err != nil {
		return nil, fmt.Errorf("Can't JSON encode update %v %s", accession, err)
	}
	content, err := aspace.API("POST", aspace.URL.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// content should look something like
	// {"status":"Created","id":3,"lock_version":0,"stale":null,"uri":"/repositories/3","warnings":[]}
	// OR
	// {"error":"Some error message here"}
	data := new(ResponseMsg)
	err = json.Unmarshal(content, data)
	if err != nil {
		return nil, fmt.Errorf("Could not unpack UpdateRepository() response [%s] %s", content, err)
	}
	return data, nil
}

// DeleteAccession deleted an Accession record from a Repository
func (aspace *ArchivesSpaceAPI) DeleteAccession(accession *Accession) (*ResponseMsg, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = accession.URI
	payload, err := json.Marshal(accession)
	if err != nil {
		return nil, err
	}
	content, err := aspace.API("DELETE", aspace.URL.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// content should look something like
	// {"status":"Deleted","id":13}
	data := new(ResponseMsg)
	err = json.Unmarshal(bytes.TrimSpace(content), data)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode DeleteAccession() %s", err)
	}
	return data, nil
}

// ListAccessions return a list of Accession IDs from a Repository
func (aspace *ArchivesSpaceAPI) ListAccessions(repositoryID int) ([]int, error) {
	aspace.URL.RawPath = ""
	aspace.URL.RawQuery = ""
	aspace.URL.Path = fmt.Sprintf(`/repositories/%d/accessions`, repositoryID)
	q := aspace.URL.Query()
	q.Set("all_ids", "true")
	aspace.URL.RawQuery = q.Encode()

	content, err := aspace.API("GET", aspace.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	// content should look something like
	// [1,2,3,4]
	var accessionIDs []int
	err = json.Unmarshal(content, &accessionIDs)
	if err != nil {
		return nil, err
	}
	return accessionIDs, nil
}

//
// String functions for aspace public structures
//

// String convert an ArchicesSpaceAPI struct as a JSON formatted string
func (aspace *ArchivesSpaceAPI) String() string {
	src, _ := json.Marshal(aspace)
	return string(src)
}

// String return a Repository as a JSON formatted string
func (repository *Repository) String() string {
	src, _ := json.Marshal(repository)
	return string(src)
}

// String return an Agent as a JSON formatted string
func (agent *Agent) String() string {
	src, _ := json.Marshal(agent)
	return string(src)
}

// String return a ResponseMsg
func (responseMsg *ResponseMsg) String() string {
	src, _ := json.Marshal(responseMsg)
	return string(src)
}

// String return a UserDefined
func (userDefined *UserDefined) String() string {
	src, _ := json.Marshal(userDefined)
	return string(src)
}

// String return a ExternalID
func (externalID *ExternalID) String() string {
	src, _ := json.Marshal(externalID)
	return string(src)
}

// String return an Extent
func (extent *Extent) String() string {
	src, _ := json.Marshal(extent)
	return string(src)
}

// String return an Accession
func (accession *Accession) String() string {
	src, _ := json.Marshal(accession)
	return string(src)
}
