-- Create tables for repository, files, and elements.

CREATE SCHEMA cindex;

CREATE TABLE cindex.code_repositories (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    repo TEXT NOT NULL,
    repo_type TEXT NOT NULL
);

CREATE TABLE cindex.code_files (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    repo_id INTEGER NOT NULL REFERENCES code_repositories(id),
    file_path TEXT NOT NULL,
    file_name TEXT NOT NULL,
    programming_language TEXT NOT NULL,
    contents TEXT NOT NULL,
    node_type TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create a table to store tokenized code elements (functions, methods, classes, etc.)
CREATE TABLE cindex.code_elements (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    file_id INTEGER NOT NULL REFERENCES code_files(id),
    element_type TEXT NOT NULL, -- e.g., 'function', 'method', 'class'
    element_name TEXT NOT NULL,
    start_line INTEGER NOT NULL,
    end_line INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes to support faster searching
CREATE INDEX ON cindex.code_files (file_name);
CREATE INDEX ON cindex.code_files (programming_language);
CREATE INDEX ON cindex.code_elements (element_type);
CREATE INDEX ON cindex.code_elements (element_name);
CREATE INDEX ON cindex.code_elements (file_id);
