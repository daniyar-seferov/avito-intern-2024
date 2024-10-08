-- +goose Up
-- +goose StatementBegin

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tender_status') THEN
        CREATE TYPE tender_status AS ENUM (
            'CREATED',    
            'PUBLISHED',  
            'CLOSED'
        );
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'service_type') THEN
        CREATE TYPE service_type AS ENUM (
            'CONSTRUCTION',    
            'DELIVERY',  
            'MANUFACTURE'
        );
    END IF;
END $$;


CREATE TABLE tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status tender_status DEFAULT 'CREATED',
    type service_type NOT NULL,
    version INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS tender CASCADE;

DROP TYPE IF EXISTS tender_status;
DROP TYPE IF EXISTS service_type;

-- +goose StatementEnd
