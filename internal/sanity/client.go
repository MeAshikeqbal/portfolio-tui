package sanity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/MeAshikeqbal/portfolio-tui/internal/config"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type QueryResponse struct {
	Result []map[string]interface{} `json:"result"`
}

// Project represents a project from Sanity
type Project struct {
	ID           string      `json:"_id"`
	Title        string      `json:"title"`
	Slug         Slug        `json:"slug"`
	Description  string      `json:"description"`
	MainImage    MainImage   `json:"mainImage"`
	PublishedAt  string      `json:"publishedAt"`
	GitHub       string      `json:"github"`
	URL          string      `json:"url"`
	GitHubData   *GitHubData `json:"githubData"`
	Technologies []string    `json:"technologies"`
}

// Slug represents a Sanity slug object.
type Slug struct {
	Current string `json:"current"`
}

// MainImage represents a Sanity image object with custom fields.
type MainImage struct {
	Alt string `json:"alt"`
}

// GitHubData stores optional GitHub metadata for a project.
type GitHubData struct {
	Stars   int    `json:"stars"`
	Commits int    `json:"commits"`
	License string `json:"license"`
}

// Skill represents a skill from Sanity
type Skill struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

// About represents about content from Sanity
type About struct {
	ID         string `json:"_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Background string `json:"background"`
}

// Contact represents contact information from Sanity
type Contact struct {
	ID       string `json:"_id"`
	Platform string `json:"platform"`
	Value    string `json:"value"`
	URL      string `json:"url"`
}

// ContactFormSubmission represents a contact form submission
type ContactFormSubmission struct {
	ID      string `json:"_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// Post represents a blog post from Sanity
type Post struct {
	ID          string        `json:"_id"`
	Title       string        `json:"title"`
	Slug        interface{}   `json:"slug"`
	PublishedAt string        `json:"publishedAt"`
	Author      interface{}   `json:"author"`
	Categories  []interface{} `json:"categories"`
	MainImage   interface{}   `json:"mainImage"`
	Body        []interface{} `json:"body"`
}

// Category represents a blog category
type Category struct {
	ID          string `json:"_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Education represents an education entry from Sanity
type Education struct {
	ID     string `json:"_id"`
	Year   string `json:"year"`
	Degree string `json:"degree"`
	School string `json:"school"`
	Desc   string `json:"desc"`
}

// WorkExperience represents a single work entry within an experience year
type WorkExperience struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Desc    string `json:"desc"`
}

// Experience represents a year group of work experiences from Sanity
type Experience struct {
	ID    string           `json:"_id"`
	Year  string           `json:"year"`
	Works []WorkExperience `json:"works"`
}

func NewClient() *Client {
	sanityCfg := config.GetSanityConfig()

	// Sanity API endpoint format: https://{projectId}.api.sanity.io/v{apiVersion}/data/query/{dataset}
	// The apiVersion should be like "2024-12-21" and we add "v" prefix in the URL
	baseURL := fmt.Sprintf(
		"https://%s.api.sanity.io/v%s/data/query/%s",
		sanityCfg.ProjectID,
		sanityCfg.APIVersion,
		sanityCfg.Dataset,
	)
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Query(query string) ([]map[string]interface{}, error) {
	// URL encode the query
	encodedQuery := url.QueryEscape(query)
	resp, err := c.httpClient.Get(c.baseURL + "?query=" + encodedQuery)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sanity API error: %s", string(body))
	}

	var result QueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	return result.Result, nil
}

func (c *Client) GetProjects() ([]Project, error) {
	query := `*[_type == "project"] | order(order asc)`
	results, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	var projects []Project
	for _, r := range results {
		data, _ := json.Marshal(r)
		var p Project
		if err := json.Unmarshal(data, &p); err == nil {
			projects = append(projects, p)
		}
	}
	return projects, nil
}

func (c *Client) GetSkills() ([]Skill, error) {
	query := `*[_type == "skills"] | order(name asc)`
	results, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	var skills []Skill
	for _, r := range results {
		data, _ := json.Marshal(r)
		var s Skill
		if err := json.Unmarshal(data, &s); err == nil {
			skills = append(skills, s)
		}
	}
	return skills, nil
}

func (c *Client) GetAbout() (*About, error) {
	query := `*[_type == "about"][0]`
	results, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no about content found")
	}

	data, _ := json.Marshal(results[0])
	var about About
	if err := json.Unmarshal(data, &about); err != nil {
		return nil, err
	}
	return &about, nil
}

func (c *Client) GetContacts() ([]Contact, error) {
	query := `*[_type == "socialLink"]`
	results, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	for _, r := range results {
		data, _ := json.Marshal(r)
		var c Contact
		if err := json.Unmarshal(data, &c); err == nil {
			contacts = append(contacts, c)
		}
	}

	// If no contacts from Sanity, use config values
	if len(contacts) == 0 {
		contacts = getContactsFromConfig()
	}

	return contacts, nil
}

// getContactsFromConfig builds contact list from config values.
func getContactsFromConfig() []Contact {
	cfg := config.Get()
	contacts := []Contact{}

	// Email
	if email := cfg.Contact.Email; email != "" {
		contacts = append(contacts, Contact{
			Platform: "Email",
			Value:    email,
			URL:      "mailto:" + email,
		})
	}

	// Phone
	if phone := cfg.Contact.Phone; phone != "" {
		contacts = append(contacts, Contact{
			Platform: "Phone",
			Value:    phone,
			URL:      "tel:" + phone,
		})
	}

	// GitHub
	if github := cfg.Social.GitHub; github != "" {
		contacts = append(contacts, Contact{
			Platform: "GitHub",
			Value:    github,
			URL:      github,
		})
	}

	// LinkedIn
	if linkedin := cfg.Social.LinkedIn; linkedin != "" {
		contacts = append(contacts, Contact{
			Platform: "LinkedIn",
			Value:    linkedin,
			URL:      linkedin,
		})
	}

	// Twitter
	if twitter := cfg.Social.Twitter; twitter != "" {
		contacts = append(contacts, Contact{
			Platform: "Twitter",
			Value:    twitter,
			URL:      twitter,
		})
	}

	// Website
	if website := cfg.Social.Website; website != "" {
		contacts = append(contacts, Contact{
			Platform: "Website",
			Value:    website,
			URL:      website,
		})
	}

	// YouTube
	if youtube := cfg.Social.YouTube; youtube != "" {
		contacts = append(contacts, Contact{
			Platform: "YouTube",
			Value:    youtube,
			URL:      youtube,
		})
	}

	// Instagram
	if instagram := cfg.Social.Instagram; instagram != "" {
		contacts = append(contacts, Contact{
			Platform: "Instagram",
			Value:    instagram,
			URL:      instagram,
		})
	}

	// Fallback if no config values are set
	if len(contacts) == 0 {
		contacts = []Contact{
			{Platform: "Email", Value: "ashik@example.com", URL: "mailto:ashik@example.com"},
			{Platform: "GitHub", Value: "github.com/MeAshikeqbal", URL: "https://github.com/MeAshikeqbal"},
		}
	}

	return contacts
}

func (c *Client) GetPosts() ([]Post, error) {
	query := `*[_type == "post"] | order(publishedAt desc) {
		_id,
		title,
		slug,
		publishedAt,
		"author": author->name,
		"categories": categories[]->title,
		mainImage,
		body
	}`
	results, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, r := range results {
		data, _ := json.Marshal(r)
		var p Post
		if err := json.Unmarshal(data, &p); err == nil {
			posts = append(posts, p)
		}
	}
	return posts, nil
}

func (c *Client) GetCategories() ([]Category, error) {
	query := `*[_type == "category"]`
	results, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	var categories []Category
	for _, r := range results {
		data, _ := json.Marshal(r)
		var cat Category
		if err := json.Unmarshal(data, &cat); err == nil {
			categories = append(categories, cat)
		}
	}
	return categories, nil
}

func (c *Client) GetEducation() ([]Education, error) {
	query := `*[_type == "education"] | order(year desc)`
	results, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	var educations []Education
	for _, r := range results {
		data, _ := json.Marshal(r)
		var e Education
		if err := json.Unmarshal(data, &e); err == nil {
			educations = append(educations, e)
		}
	}
	return educations, nil
}

func (c *Client) GetExperiences() ([]Experience, error) {
	query := `*[_type == "experiences"] | order(year desc) {
		_id,
		year,
		"works": works[] {
			name,
			company,
			desc
		}
	}`
	results, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	var experiences []Experience
	for _, r := range results {
		data, _ := json.Marshal(r)
		var e Experience
		if err := json.Unmarshal(data, &e); err == nil {
			experiences = append(experiences, e)
		}
	}
	return experiences, nil
}
