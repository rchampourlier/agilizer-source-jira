package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB() *DB {
	return &DB{openDB()}
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// MaxOpenConns defines the maximum number of open connections
// to the DB.
const MaxOpenConns = 5 // for Heroku Postgres

type IssueEvent struct {
	EventTime        time.Time
	EventKind        string
	EventAuthor      string
	IssueKey         string
	CommentBody      *string
	StatusChangeFrom *string
	StatusChangeTo   *string
}

type IssueState struct {
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Key               string
	Project           *string
	Status            *string
	ResolvedAt        *time.Time
	Priority          *string
	Summary           *string
	Description       *string
	Type              *string
	Labels            *string
	Reporter          *string
	Assignee          *string
	DeveloperBackend  *string
	DeveloperFrontend *string
	Reviewer          *string
	ProductOwner      *string
	BugCause          *string
	Epic              *string
	Tribe             *string
	Components        *string
	FixVersions       *string
}

func (db *DB) InsertIssueEvent(e IssueEvent, s IssueState) {
	query := `
	INSERT INTO jira_issues_events (
		event_time,
		event_kind,
		event_author,
		comment_body,
		status_change_from,
		status_change_to,
		issue_key,
		issue_created_at,
		issue_updated_at,
		issue_project,
		issue_status,
		issue_resolved_at,
		issue_priority,
		issue_summary,
		issue_description,
		issue_type,
		issue_labels,
		issue_assignee,
		issue_developer_backend,
		issue_developer_frontend,
		issue_reviewer,
		issue_product_owner,
		issue_bug_cause,
		issue_epic,
		issue_tribe,
		issue_components,
		issue_fix_versions
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27);
	`

	rows, err := db.Query(
		query,
		e.EventTime,
		e.EventKind,
		e.EventAuthor,
		e.CommentBody,
		e.StatusChangeFrom,
		e.StatusChangeTo,
		e.IssueKey,
		s.CreatedAt,
		s.UpdatedAt,
		s.Project,
		s.Status,
		s.ResolvedAt,
		s.Priority,
		s.Summary,
		s.Description,
		s.Type,
		s.Labels,
		s.Assignee,
		s.DeveloperBackend,
		s.DeveloperFrontend,
		s.Reviewer,
		s.ProductOwner,
		s.BugCause,
		s.Epic,
		s.Tribe,
		s.Components,
		s.FixVersions,
	)
	if err != nil {
		log.Fatalln(fmt.Errorf("error in `InsertIssueEvent`: %s", err))
	}
	rows.Close()
}

func (db *DB) InsertIssueState(s IssueState) {
	query := `
	INSERT INTO jira_issues_states (
		issue_created_at,
		issue_updated_at,
		issue_key,
		issue_project,
		issue_status,
		issue_resolved_at,
		issue_priority,
		issue_summary,
		issue_description,
		issue_type,
		issue_labels,
		issue_assignee,
		issue_developer_backend,
		issue_developer_frontend,
		issue_reviewer,
		issue_product_owner,
		issue_bug_cause,
		issue_epic,
		issue_tribe,
		issue_components,
		issue_fix_versions
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21);
	`
	rows, err := db.Query(
		query,
		s.CreatedAt,
		s.UpdatedAt,
		s.Key,
		s.Project,
		s.Status,
		s.ResolvedAt,
		s.Priority,
		s.Summary,
		s.Description,
		s.Type,
		s.Labels,
		s.Assignee,
		s.DeveloperBackend,
		s.DeveloperFrontend,
		s.Reviewer,
		s.ProductOwner,
		s.BugCause,
		s.Epic,
		s.Tribe,
		s.Components,
		s.FixVersions,
	)
	if err != nil {
		log.Fatalln(fmt.Errorf("error in `InsertIssueState`: %s", err))
	}
	rows.Close()
}

func openDB() *sql.DB {
	connStr := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", connStr)
	db.SetMaxOpenConns(MaxOpenConns)
	if err != nil {
		log.Fatalln(fmt.Errorf("error in `openDB`: %s", err))
	}
	return db
}

func (db *DB) Reset() {
	queries := make([]string, 0)
	queries = append(queries, queriesResetTableJiraIssuesEvents()...)
	queries = append(queries, queriesResetTableJiraIssuesStates()...)
	err := db.doQueries(queries)
	if err != nil {
		log.Fatalln(fmt.Errorf("error in `Reset`: %s", err))
	}
}

func (db *DB) DropDBTables() {
	queries := []string{
		`DROP TABLE IF EXISTS "jira_issues_events";`,
		`DROP TABLE IF EXISTS "jira_issues_states";`,
	}
	err := db.doQueries(queries)
	if err != nil {
		log.Fatalln(fmt.Errorf("error in `dropDBTables`: %s", err))
	}
}

// doQueries performs the specified queries on the passed db.
// If an error occurs, it returns the error. This function can't
// be used for queries where you need the result rows.
func (db *DB) doQueries(queries []string) error {
	for _, q := range queries {
		rows, err := db.Query(q)
		if err != nil {
			return err
		}
		defer rows.Close()
		if err = rows.Err(); err != nil {
			return err
		}
	}
	return nil
}

func queriesResetTableJiraIssuesEvents() []string {
	return []string{
		`DROP TABLE IF EXISTS "jira_issues_events";`,
		`CREATE TABLE "jira_issues_events" (
		  "id" serial primary key not null,
		  "inserted_at" TIMESTAMP(6) NOT NULL DEFAULT statement_timestamp(),
		  "event_time" TIMESTAMP NOT NULL,
		  "event_kind" TEXT NOT NULL,
		  "event_author" TEXT NOT NULL,
	          "issue_created_at" TIMESTAMP NOT NULL,
		  "issue_updated_at" TIMESTAMP NOT NULL,
		  "issue_key" TEXT NOT NULL,
		  "issue_project" TEXT NOT NULL,
		  "issue_status" TEXT NOT NULL,
		  "issue_resolved_at" TIMESTAMP,
		  "issue_priority" TEXT NOT NULL,
		  "issue_summary" TEXT NOT NULL,
		  "issue_description" TEXT,
		  "issue_type" TEXT NOT NULL,
		  "issue_labels" TEXT,
		  "issue_assignee" TEXT,
		  "issue_developer_backend" TEXT,
		  "issue_developer_frontend" TEXT,
		  "issue_reviewer" TEXT,
		  "issue_product_owner" TEXT,
		  "issue_bug_cause" TEXT,
		  "issue_epic" TEXT,
		  "issue_tribe" TEXT,
		  "issue_components" TEXT,
		  "issue_fix_versions" TEXT,
		  "comment_body" TEXT,
		  "status_change_from" TEXT,
		  "status_change_to" TEXT
		);`,
	}
}

func queriesResetTableJiraIssuesStates() []string {
	return []string{
		`DROP TABLE IF EXISTS "jira_issues_states";`,
		`CREATE TABLE "jira_issues_states" (
			"id" SERIAL PRIMARY KEY NOT NULL,
			"inserted_at" TIMESTAMP(6) NOT NULL DEFAULT statement_timestamp(),
			"issue_created_at" TIMESTAMP NOT NULL,
			"issue_updated_at" TIMESTAMP NOT NULL,
			"issue_key" TEXT NOT NULL,
			"issue_project" TEXT NOT NULL,
			"issue_status" TEXT NOT NULL,
			"issue_resolved_at" TIMESTAMP,
			"issue_priority" TEXT NOT NULL,
			"issue_summary" TEXT NOT NULL,
			"issue_description" TEXT,
			"issue_type" TEXT NOT NULL,
			"issue_labels" TEXT,
			"issue_assignee" TEXT,
			"issue_developer_backend" TEXT,
			"issue_developer_frontend" TEXT,
			"issue_reviewer" TEXT,
			"issue_product_owner" TEXT,
			"issue_bug_cause" TEXT,
			"issue_epic" TEXT,
			"issue_tribe" TEXT,
			"issue_components" TEXT,
			"issue_fix_versions" TEXT
		);`,
	}
}
