package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"
)

type tmplIssueTemplateData struct {
	Labels   []string
	Versions []string
	Proxies  []string
}

type tmplConfigurationKeysData struct {
	Timestamp time.Time
	Keys      []string
	Package   string
}

type tmplScriptsGEnData struct {
	Package          string
	VersionSwaggerUI string
}

// ConfigurationKey is the docs json model for the Authelia configuration keys.
type ConfigurationKey struct {
	Path   string `json:"path"`
	Secret bool   `json:"secret"`
	Env    string `json:"env"`
}

// Languages is the docs json model for the Authelia languages configuration.
type Languages struct {
	Defaults   DefaultsLanguages `json:"defaults"`
	Namespaces []string          `json:"namespaces"`
	Languages  []Language        `json:"languages"`
}

type DefaultsLanguages struct {
	Language  Language `json:"language"`
	Namespace string   `json:"namespace"`
}

// Language is the docs json model for a language.
type Language struct {
	Display    string   `json:"display"`
	Locale     string   `json:"locale"`
	Namespaces []string `json:"namespaces,omitempty"`
	Fallbacks  []string `json:"fallbacks,omitempty"`

	Tag language.Tag `json:"-"`
}

const (
	labelAreaPrefixPriority = "priority"
	labelAreaPrefixType     = "type"
	labelAreaPrefixStatus   = "status"
)

type labelPriority int

//nolint:deadcode // Kept for future use.
const (
	labelPriorityCritical labelPriority = iota
	labelPriorityHigh
	labelPriorityMedium
	labelPriorityNormal
	labelPriorityLow
)

var labelPriorityDescriptions = [...]string{
	"Critical",
	"High",
	"Medium",
	"Normal",
	"Low",
}

func (p labelPriority) String() string {
	return fmt.Sprintf("%s/%d/%s", labelAreaPrefixPriority, p+1, strings.ToLower(labelPriorityDescriptions[p]))
}

func (p labelPriority) Description() string {
	return labelPriorityDescriptions[p]
}

type labelStatus int

const (
	labelStatusNeedsDesign labelStatus = iota
	labelStatusNeedsTriage
)

var labelStatusDescriptions = [...]string{
	"needs-design",
	"needs-triage",
}

func (s labelStatus) String() string {
	return fmt.Sprintf("%s/%s", labelAreaPrefixStatus, labelStatusDescriptions[s])
}

type labelType int

//nolint:deadcode // Kept for future use.
const (
	labelTypeFeature labelType = iota
	labelTypeBugUnconfirmed
	labelTypeBug
)

var labelTypeDescriptions = [...]string{
	"feature",
	"bug/unconfirmed",
	"bug",
}

func (t labelType) String() string {
	return fmt.Sprintf("%s/%s", labelAreaPrefixType, labelTypeDescriptions[t])
}
