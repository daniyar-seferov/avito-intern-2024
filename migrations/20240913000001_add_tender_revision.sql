-- +goose Up
-- +goose StatementBegin

CREATE TABLE tender_revision (
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status tender_status DEFAULT 'CREATED',
    type service_type NOT NULL,
    version INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tender_id, version)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tender_revision;

-- +goose StatementEnd
