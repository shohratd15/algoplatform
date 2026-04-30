-- ======================================
-- Migration: Insert sample programming problems
-- ======================================

-- === A + B Problem ===
INSERT INTO problems (slug, difficulty)
VALUES ('a-plus-b', 'easy');

INSERT INTO problem_statements (problem_id, language, title, statement)
VALUES
(
    (SELECT id FROM problems WHERE slug = 'a-plus-b'),
    'eng',
    'A + B Problem',
    'Given two integers a and b, output their sum.'
);

INSERT INTO problem_statements (problem_id, language, title, statement)
VALUES
(
    (SELECT id FROM problems WHERE slug = 'a-plus-b'),
    'rus',
    'A + B Problem',
    'Даны два целых числа a и b, выведите их сумму.'
);

INSERT INTO problem_statements (problem_id, language, title, statement)
VALUES
(
    (SELECT id FROM problems WHERE slug = 'a-plus-b'),
    'tkm',
    'A + B Problem',
    'A we B bitin sanlar berlen. Olaryň jemini tapmaly.'
);

INSERT INTO problem_tests (problem_id, input_data, expected_output, is_sample)
VALUES
((SELECT id FROM problems WHERE slug = 'a-plus-b'), '2 3', '5', TRUE),
((SELECT id FROM problems WHERE slug = 'a-plus-b'), '10 20', '30', FALSE),
((SELECT id FROM problems WHERE slug = 'a-plus-b'), '-5 5', '0', FALSE);


-- === Factorial Problem ===
INSERT INTO problems (slug, difficulty)
VALUES ('factorial', 'easy');

INSERT INTO problem_statements (problem_id, language, title, statement)
VALUES
(
    (SELECT id FROM problems WHERE slug = 'factorial'),
    'eng',
    'Factorial',
    'Given an integer n (0 ≤ n ≤ 10), output n factorial (n!).'
);

INSERT INTO problem_tests (problem_id, input_data, expected_output, is_sample)
VALUES
((SELECT id FROM problems WHERE slug = 'factorial'), '3', '6', TRUE),
((SELECT id FROM problems WHERE slug = 'factorial'), '5', '120', FALSE),
((SELECT id FROM problems WHERE slug = 'factorial'), '0', '1', FALSE);


-- === Palindrome Problem ===
INSERT INTO problems (slug, difficulty)
VALUES ('palindrome', 'medium');

INSERT INTO problem_statements (problem_id, language, title, statement)
VALUES
(
    (SELECT id FROM problems WHERE slug = 'palindrome'),
    'eng',
    'Palindrome',
    'Given a string, determine if it is a palindrome. Output YES or NO.'
);

INSERT INTO problem_tests (problem_id, input_data, expected_output, is_sample)
VALUES
((SELECT id FROM problems WHERE slug = 'palindrome'), 'abba', 'YES', TRUE),
((SELECT id FROM problems WHERE slug = 'palindrome'), 'abc', 'NO', FALSE),
((SELECT id FROM problems WHERE slug = 'palindrome'), 'a', 'YES', FALSE);
