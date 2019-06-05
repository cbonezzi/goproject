package mysql

import (
	"cesarbon.net/goproject/pkg/models"
	"database/sql"
)

// SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database
func (m *SnippetModel) Insert(model models.Snip) (int, error) {
	
	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"
	//alternatively, we can use prepared statements... as seen below use with causing!
	//this comes with trade of performs and complexity.
	//for example prepared statements are attach to db connections, therefore, if conn A has stmt 1, and conn A is busy processing resource 2
	//stmt 1 will look for another conn, say C. At this point, conn C needs to re-prepare the statement from stmt 1 resulting in overhead.
	//one must balance between performance and complexity when deciding which strategy is feasible for which problem.
	//insert, err := m.DB.Prepare("INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))")

	result, err := m.DB.Exec(stmt, model.Title, model.Content, model.Expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"

	row := m.DB.QueryRow(stmt, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

// Latest will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10"

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	//critical to defer since if it not closed, we can quickly used up connections on our pool.
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		
		s := &models.Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}