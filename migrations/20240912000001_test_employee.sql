-- +goose Up
-- +goose StatementBegin

INSERT INTO employee (username, first_name, last_name)
VALUES
    ('alice_walker', 'Alice', 'Walker'),
    ('robert_brown', 'Robert', 'Brown'),
    ('linda_clark', 'Linda', 'Clark'),
    ('daniel_johnson', 'Daniel', 'Johnson'),
    ('emma_miller', 'Emma', 'Miller'),
    ('olivia_wilson', 'Olivia', 'Wilson'),
    ('noah_moore', 'Noah', 'Moore'),
    ('liam_taylor', 'Liam', 'Taylor');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM employee
WHERE username IN ('alice_walker', 'robert_brown', 'linda_clark', 'daniel_johnson', 'emma_miller', 'olivia_wilson', 'noah_moore', 'liam_taylor');

-- +goose StatementEnd
