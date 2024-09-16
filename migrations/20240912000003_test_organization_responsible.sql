-- +goose Up
-- +goose StatementBegin

INSERT INTO organization_responsible (organization_id, user_id)
VALUES
    (
        (SELECT id FROM organization WHERE name = 'Smart Solutions'),
        (SELECT id FROM employee WHERE username = 'alice_walker')
    ),
    (
        (SELECT id FROM organization WHERE name = 'Smart Solutions'),
        (SELECT id FROM employee WHERE username = 'robert_brown')
    ),
    (
        (SELECT id FROM organization WHERE name = 'Eco Future'),
        (SELECT id FROM employee WHERE username = 'linda_clark')
    ),
    (
        (SELECT id FROM organization WHERE name = 'Blue Ocean'),
        (SELECT id FROM employee WHERE username = 'daniel_johnson')
    ),
    (
        (SELECT id FROM organization WHERE name = 'Global Trends'),
        (SELECT id FROM employee WHERE username = 'emma_miller')
    ),
    (
        (SELECT id FROM organization WHERE name = 'Creative Minds'),
        (SELECT id FROM employee WHERE username = 'olivia_wilson')
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM organization_responsible
WHERE id IN (
    SELECT id
    FROM organization
    WHERE name IN ('Smart Solutions', 'Health for All', 'Eco Future', 'Blue Ocean', 'Global Trends', 'Creative Minds')
);

-- +goose StatementEnd
