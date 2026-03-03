INSERT OR IGNORE INTO groups (id, name) VALUES 
(1, 'basic_greetings'),
(2, 'seasons');

INSERT OR IGNORE INTO words (id, bengali, english, parts) VALUES 
(1, 'স্বাগতম', 'welcome', 'অব্যয়'),
(2, 'নমস্কার', 'hello', 'অব্যয়'),
(3, 'ধন্যবাদ', 'thank you', 'অব্যয়'),
(4, 'আপনি', 'you', 'সর্বনাম'),
(5, 'আমি', 'I', 'সর্বনাম'),
(6, 'গ্রীষ্ম', 'summer', 'বিশেষ্য'),
(7, 'শীত', 'winter', 'বিশেষ্য'),
(8, 'বসন্ত', 'spring', 'বিশেষ্য'),
(9, 'শরৎ', 'autumn', 'বিশেষ্য'),
(10, 'বর্ষা', 'rainy season', 'বিশেষ্য');

INSERT OR IGNORE INTO words_groups (word_id, group_id) VALUES 
(1, 1), (2, 1), (3, 1), (4, 1), (5, 1),
(6, 2), (7, 2), (8, 2), (9, 2), (10, 2);

INSERT OR IGNORE INTO study_sessions (id, group_id, study_activities_id) VALUES 
(1, 1, 1),
(2, 2, 2);

INSERT OR IGNORE INTO word_review_items (word_id, study_session_id, correct) VALUES 
(1, 1, 1),
(2, 1, 0),
(6, 2, 1);
