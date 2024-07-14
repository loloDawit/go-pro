package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLDatabase(t *testing.T) {
	tests := []struct {
		name           string
		sqlOpenError   error
		queryError     error
		rowsScanError  error
		expectedError  error
		expectedOutput string
	}{
		{
			name:           "Successful Connection",
			sqlOpenError:   nil,
			queryError:     nil,
			rowsScanError:  nil,
			expectedError:  nil,
			expectedOutput: "PostgreSQL 12.3",
		},
		{
			name:          "sql.Open Error",
			sqlOpenError:  errors.New("sql.Open error"),
			expectedError: errors.New("sql.Open error"),
		},
		{
			name:          "db.Query Error",
			sqlOpenError:  nil,
			queryError:    errors.New("db.Query error"),
			expectedError: errors.New("db.Query error"),
		},
		{
			name:          "rows.Scan Error",
			sqlOpenError:  nil,
			queryError:    nil,
			rowsScanError: errors.New("rows.Scan error"),
			expectedError: errors.New("rows.Scan error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock database connection
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Set up the sqlOpen mock
			sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
				if tt.sqlOpenError != nil {
					return nil, tt.sqlOpenError
				}
				return db, nil
			}

			// Set up expectations for the query
			if tt.queryError == nil && tt.sqlOpenError == nil {
				rows := sqlmock.NewRows([]string{"version"})
				if tt.rowsScanError == nil {
					rows.AddRow(tt.expectedOutput)
				} else {
					rows.AddRow("dummy_version")
					// We cannot directly inject a scan error; instead, we simulate a query error that results in no rows being scanned.
					mock.ExpectQuery("SELECT version()").WillReturnRows(rows).WillReturnError(tt.rowsScanError)
				}
				mock.ExpectQuery("SELECT version()").WillReturnRows(rows)
			} else if tt.sqlOpenError == nil {
				mock.ExpectQuery("SELECT version()").WillReturnError(tt.queryError)
			}

			// Call the function you want to test
			connStr := "mock_connection_string"
			actualDB, err := NewSQLDatabase(connStr)

			// Check results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, actualDB)
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}
