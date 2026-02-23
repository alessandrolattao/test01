-- Seed data for development

-- Admin user: admin@test.com / password
-- bcrypt hash of "password" with cost 10
INSERT INTO admins (id, email, password_hash) VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'admin@test.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'
);

-- Active questionnaire
INSERT INTO questionnaires (id, version, is_active) VALUES (
    'b0000000-0000-0000-0000-000000000001', 1, TRUE
);

-- Question 1: Problem solving
INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
    'c0000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000001',
    'You discover a critical bug in production on a Friday afternoon. What do you do?',
    1
);
INSERT INTO answers (id, question_id, text, score, sort_order) VALUES
    ('d0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000001', 'Fix it immediately and deploy the patch', 10, 1),
    ('d0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000001', 'Document it and schedule it for Monday', 3, 2),
    ('d0000000-0000-0000-0000-000000000003', 'c0000000-0000-0000-0000-000000000001', 'Pretend you did not see it', 0, 3),
    ('d0000000-0000-0000-0000-000000000004', 'c0000000-0000-0000-0000-000000000001', 'Notify the team and assess the impact before acting', 8, 4);

-- Question 2: Teamwork
INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
    'c0000000-0000-0000-0000-000000000002',
    'b0000000-0000-0000-0000-000000000001',
    'A colleague disagrees with your technical approach. How do you handle it?',
    2
);
INSERT INTO answers (id, question_id, text, score, sort_order) VALUES
    ('d0000000-0000-0000-0000-000000000005', 'c0000000-0000-0000-0000-000000000002', 'Listen to their perspective and find a compromise', 10, 1),
    ('d0000000-0000-0000-0000-000000000006', 'c0000000-0000-0000-0000-000000000002', 'Insist on your approach since you know it is better', 0, 2),
    ('d0000000-0000-0000-0000-000000000007', 'c0000000-0000-0000-0000-000000000002', 'Let them do it their way to avoid conflict', 3, 3),
    ('d0000000-0000-0000-0000-000000000008', 'c0000000-0000-0000-0000-000000000002', 'Propose a proof of concept for both approaches', 8, 4);

-- Question 3: Learning
INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
    'c0000000-0000-0000-0000-000000000003',
    'b0000000-0000-0000-0000-000000000001',
    'How do you stay updated with new technologies?',
    3
);
INSERT INTO answers (id, question_id, text, score, sort_order) VALUES
    ('d0000000-0000-0000-0000-000000000009', 'c0000000-0000-0000-0000-000000000003', 'I regularly take courses and build side projects', 10, 1),
    ('d0000000-0000-0000-0000-000000000010', 'c0000000-0000-0000-0000-000000000003', 'I read tech blogs occasionally', 5, 2),
    ('d0000000-0000-0000-0000-000000000011', 'c0000000-0000-0000-0000-000000000003', 'I only learn what is needed for my current job', 2, 3),
    ('d0000000-0000-0000-0000-000000000012', 'c0000000-0000-0000-0000-000000000003', 'I attend conferences and contribute to open source', 8, 4);

-- Question 4: Time management
INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
    'c0000000-0000-0000-0000-000000000004',
    'b0000000-0000-0000-0000-000000000001',
    'You have multiple deadlines approaching. How do you prioritize?',
    4
);
INSERT INTO answers (id, question_id, text, score, sort_order) VALUES
    ('d0000000-0000-0000-0000-000000000013', 'c0000000-0000-0000-0000-000000000004', 'Assess urgency and impact, then create a priority list', 10, 1),
    ('d0000000-0000-0000-0000-000000000014', 'c0000000-0000-0000-0000-000000000004', 'Work on whatever is closest to its deadline', 5, 2),
    ('d0000000-0000-0000-0000-000000000015', 'c0000000-0000-0000-0000-000000000004', 'Ask my manager to decide for me', 3, 3),
    ('d0000000-0000-0000-0000-000000000016', 'c0000000-0000-0000-0000-000000000004', 'Try to work on everything simultaneously', 0, 4);

-- Question 5: Communication
INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
    'c0000000-0000-0000-0000-000000000005',
    'b0000000-0000-0000-0000-000000000001',
    'How would you explain a complex technical concept to a non-technical stakeholder?',
    5
);
INSERT INTO answers (id, question_id, text, score, sort_order) VALUES
    ('d0000000-0000-0000-0000-000000000017', 'c0000000-0000-0000-0000-000000000005', 'Use analogies and simple language, focus on business impact', 10, 1),
    ('d0000000-0000-0000-0000-000000000018', 'c0000000-0000-0000-0000-000000000005', 'Show them the code and walk through it step by step', 2, 2),
    ('d0000000-0000-0000-0000-000000000019', 'c0000000-0000-0000-0000-000000000005', 'Send them a link to the documentation', 0, 3),
    ('d0000000-0000-0000-0000-000000000020', 'c0000000-0000-0000-0000-000000000005', 'Create a visual diagram and present it briefly', 8, 4);

-- Example candidate 1 (completed, high score)
INSERT INTO candidates (id, first_name, last_name, email, questionnaire_id, total_score, completed) VALUES (
    'e0000000-0000-0000-0000-000000000001',
    'Marco', 'Rossi', 'marco.rossi@example.com',
    'b0000000-0000-0000-0000-000000000001',
    46, TRUE
);
INSERT INTO candidate_answers (candidate_id, question_id, answer_id, score) VALUES
    ('e0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000001', 10),
    ('e0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000005', 10),
    ('e0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000003', 'd0000000-0000-0000-0000-000000000009', 10),
    ('e0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'd0000000-0000-0000-0000-000000000014', 8),
    ('e0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000005', 'd0000000-0000-0000-0000-000000000020', 8);

-- Example candidate 2 (completed, medium score)
INSERT INTO candidates (id, first_name, last_name, email, questionnaire_id, total_score, completed) VALUES (
    'e0000000-0000-0000-0000-000000000002',
    'Giulia', 'Bianchi', 'giulia.bianchi@example.com',
    'b0000000-0000-0000-0000-000000000001',
    23, TRUE
);
INSERT INTO candidate_answers (candidate_id, question_id, answer_id, score) VALUES
    ('e0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000002', 3),
    ('e0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000008', 8),
    ('e0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'd0000000-0000-0000-0000-000000000010', 5),
    ('e0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'd0000000-0000-0000-0000-000000000015', 5),
    ('e0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000005', 'd0000000-0000-0000-0000-000000000018', 2);
