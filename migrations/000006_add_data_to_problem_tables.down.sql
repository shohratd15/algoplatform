-- ======================================
-- Down: Remove inserted sample problems
-- ======================================

DELETE FROM problem_tests
WHERE problem_id IN (SELECT id FROM problems WHERE slug IN ('a-plus-b','factorial','palindrome'));

DELETE FROM problem_statements
WHERE problem_id IN (SELECT id FROM problems WHERE slug IN ('a-plus-b','factorial','palindrome'));

DELETE FROM problems
WHERE slug IN ('a-plus-b','factorial','palindrome');
