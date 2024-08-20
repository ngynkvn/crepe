-- Create the main 'code_files' table
CREATE TABLE code_files (
    id SERIAL PRIMARY KEY,
    file_path TEXT NOT NULL,
    file_name TEXT NOT NULL,
    programming_language TEXT NOT NULL,
    contents TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create a table to store tokenized code elements (functions, methods, classes, etc.)
CREATE TABLE code_elements (
    id SERIAL PRIMARY KEY,
    file_id INTEGER NOT NULL REFERENCES code_files(id),
    element_type TEXT NOT NULL, -- e.g., 'function', 'method', 'class'
    element_name TEXT NOT NULL,
    start_line INTEGER NOT NULL,
    end_line INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes to support faster searching
CREATE INDEX ON code_files (file_name);
CREATE INDEX ON code_files (programming_language);
CREATE INDEX ON code_elements (element_type);
CREATE INDEX ON code_elements (element_name);
CREATE INDEX ON code_elements (file_id);
