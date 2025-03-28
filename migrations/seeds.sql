-- Script 1: Inserting students with names Student1 ... Student10 and corresponding card credits
INSERT INTO students (name, age, sex, card_id) VALUES
('Student1', 20, TRUE, 1),
('Student2', 22, FALSE, 2),
('Student3', 25, TRUE, 3),
('Student4', 30, FALSE, 4),
('Student5', 19, TRUE, 5),
('Student6', 21, FALSE, 6),
('Student7', 23, TRUE, 7),
('Student8', 27, FALSE, 8),
('Student9', 29, TRUE, 9),
('Student10', 31, FALSE, 10);

INSERT INTO card_credits (student_id, card_number, expiration, cvv) VALUES
(1, 1234567812345678, 1226, 123),
(2, 2345678923456789, 1125, 456),
(3, 3456789034567890, 1027, 789),
(4, 4567890145678901, 0928, 987),
(5, 5678901256789012, 0729, 654),
(6, 6789012367890123, 0627, 321),
(7, 7890123478901234, 0326, 159),
(8, 8901234589012345, 0428, 753),
(9, 9012345690123456, 0927, 852),
(10, 1234567801234567, 1129, 951);

-- Script 2: Inserting students with different names and corresponding card credits
-- Note: To avoid unique constraint violations on card_id and adjust student_id references,
-- we use new card_id values (11-20) and assume these students will receive new auto-generated IDs.
-- For demonstration, we explicitly assign card_id values here.
INSERT INTO students (name, age, sex, card_id) VALUES
('Anna', 20, TRUE, 11),
('Boris', 22, FALSE, 12),
('Irina', 25, TRUE, 13),
('Dmitry', 30, FALSE, 14),
('Olga', 19, TRUE, 15),
('Sergey', 21, FALSE, 16),
('Natalia', 23, TRUE, 17),
('Ivan', 27, FALSE, 18),
('Elena', 29, TRUE, 19),
('Pavel', 31, FALSE, 20);

INSERT INTO card_credits (student_id, card_number, expiration, cvv) VALUES
(11, 19934567812345678, 19926, 123),
(12, 99345678923456789, 19925, 456),
(13, 34567890399567890, 19927, 789),
(14, 45678990145678901, 09928, 987),
(15, 56789012996789012, 09929, 654),
(16, 67890123678901239, 09927, 321),
(17, 78909923478901234, 09926, 159),
(18, 89012345999012345, 09928, 753),
(19, 90123456901234569, 08888, 852),
(20, 12345699801234567, 19929, 951);