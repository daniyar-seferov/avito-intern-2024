-- +goose Up
-- +goose StatementBegin

INSERT INTO organization (name, description, type)
VALUES
    ('Smart Solutions', 'A consulting company providing IT services.', 'LLC'),
    ('Health for All', 'A healthcare organization.', 'IE'),
    ('Eco Future', 'A company focused on sustainable energy solutions.', 'JSC'),
    ('Blue Ocean', 'A logistics company.', 'LLC'),
    ('Global Trends', 'An international trade firm.', 'JSC'),
    ('Creative Minds', 'A marketing agency.', 'IE');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM organization
WHERE name IN ('Smart Solutions', 'Health for All', 'Eco Future', 'Blue Ocean', 'Global Trends', 'Creative Minds');

-- +goose StatementEnd
