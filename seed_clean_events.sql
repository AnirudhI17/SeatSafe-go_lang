-- Clean Database and Add 15 Professional Events
-- Run this script in your PostgreSQL database

-- Step 1: Delete all existing data (in correct order due to foreign keys)
DELETE FROM tickets;
DELETE FROM registrations;
DELETE FROM events;

-- Step 2: Insert 15 professional events
-- Note: Replace 'YOUR_ORGANIZER_USER_ID' with an actual organizer user ID from your users table

-- You can get an organizer user ID by running:
-- SELECT id FROM users WHERE role = 'organizer' LIMIT 1;

-- For now, I'll use a placeholder. You need to replace it with a real UUID.

INSERT INTO events (id, organizer_id, title, description, location, is_online, starts_at, ends_at, capacity, price_cents, currency, status, registered_count) VALUES

-- Tech & Business Events
('a1111111-1111-1111-1111-111111111111', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1), 
'Tech Innovation Summit 2024', 
'Join industry leaders and innovators for a day of cutting-edge technology discussions, networking, and hands-on workshops. Explore AI, blockchain, and the future of tech.',
'San Francisco Convention Center, CA',
false,
'2024-06-15 09:00:00',
'2024-06-15 18:00:00',
100,
15000,
'USD',
'published',
0),

('a2222222-2222-2222-2222-222222222222', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Digital Marketing Masterclass',
'Learn the latest strategies in SEO, social media marketing, and content creation from industry experts. Perfect for marketers and business owners.',
'New York Marriott, NY',
false,
'2024-06-20 10:00:00',
'2024-06-20 17:00:00',
100,
9900,
'USD',
'published',
0),

('a3333333-3333-3333-3333-333333333333', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Startup Pitch Competition',
'Watch innovative startups pitch their ideas to top investors. Network with entrepreneurs and VCs. Cash prizes for winners!',
'Austin Tech Hub, TX',
false,
'2024-06-25 14:00:00',
'2024-06-25 20:00:00',
100,
2500,
'USD',
'published',
0),

-- Arts & Culture Events
('a4444444-4444-4444-4444-444444444444', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Summer Music Festival',
'Three days of live music featuring top artists across multiple genres. Food trucks, art installations, and unforgettable performances.',
'Golden Gate Park, San Francisco, CA',
false,
'2024-07-10 12:00:00',
'2024-07-12 23:00:00',
100,
12500,
'USD',
'published',
0),

('a5555555-5555-5555-5555-555555555555', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Contemporary Art Exhibition',
'Explore groundbreaking works from emerging and established artists. Gallery tours, artist talks, and exclusive previews.',
'Museum of Modern Art, New York, NY',
false,
'2024-07-15 10:00:00',
'2024-07-15 19:00:00',
100,
3500,
'USD',
'published',
0),

('a6666666-6666-6666-6666-666666666666', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Broadway Night: Hamilton',
'Experience the award-winning musical Hamilton with premium orchestra seats. Pre-show dinner included.',
'Richard Rodgers Theatre, New York, NY',
false,
'2024-07-20 19:00:00',
'2024-07-20 22:30:00',
100,
25000,
'USD',
'published',
0),

-- Sports & Fitness Events
('a7777777-7777-7777-7777-777777777777', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'City Marathon 2024',
'Join thousands of runners in this annual marathon through scenic city routes. All skill levels welcome. Medals for all finishers!',
'Downtown Chicago, IL',
false,
'2024-08-05 07:00:00',
'2024-08-05 14:00:00',
100,
7500,
'USD',
'published',
0),

('a8888888-8888-8888-8888-888888888888', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Yoga & Wellness Retreat',
'Weekend retreat featuring yoga sessions, meditation workshops, healthy cuisine, and spa treatments in a peaceful mountain setting.',
'Sedona Wellness Center, AZ',
false,
'2024-08-10 15:00:00',
'2024-08-12 12:00:00',
100,
45000,
'USD',
'published',
0),

-- Food & Drink Events
('a9999999-9999-9999-9999-999999999999', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Wine Tasting Experience',
'Sample premium wines from renowned vineyards. Expert sommeliers guide you through tasting notes and food pairings.',
'Napa Valley Winery, CA',
false,
'2024-08-15 17:00:00',
'2024-08-15 21:00:00',
100,
8500,
'USD',
'published',
0),

('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'International Food Festival',
'Taste cuisines from around the world! Over 50 food vendors, live cooking demos, and cultural performances.',
'Millennium Park, Chicago, IL',
false,
'2024-08-20 11:00:00',
'2024-08-20 22:00:00',
100,
0,
'USD',
'published',
0),

-- Education & Workshops
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Photography Workshop: Landscape',
'Learn professional landscape photography techniques. Includes outdoor shooting session and post-processing tutorial.',
'Yosemite National Park, CA',
false,
'2024-09-01 08:00:00',
'2024-09-01 17:00:00',
100,
12000,
'USD',
'published',
0),

('cccccccc-cccc-cccc-cccc-cccccccccccc', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Public Speaking Bootcamp',
'Overcome stage fright and master the art of public speaking. Interactive exercises, video feedback, and expert coaching.',
'Boston Conference Center, MA',
false,
'2024-09-05 09:00:00',
'2024-09-05 18:00:00',
100,
19900,
'USD',
'published',
0),

-- Networking & Social Events
('dddddddd-dddd-dddd-dddd-dddddddddddd', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Professional Networking Mixer',
'Connect with professionals across industries. Speed networking, keynote speaker, and cocktail reception.',
'Seattle Convention Center, WA',
false,
'2024-09-10 18:00:00',
'2024-09-10 21:00:00',
100,
3500,
'USD',
'published',
0),

('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Women in Tech Conference',
'Empowering women in technology through inspiring talks, mentorship sessions, and career development workshops.',
'Silicon Valley Tech Center, CA',
false,
'2024-09-15 09:00:00',
'2024-09-15 17:00:00',
100,
7500,
'USD',
'published',
0),

-- Entertainment
('ffffffff-ffff-ffff-ffff-ffffffffffff', (SELECT id FROM users WHERE role = 'organizer' LIMIT 1),
'Comedy Night Live',
'An evening of laughter with top stand-up comedians. Two-drink minimum. 18+ only.',
'The Comedy Store, Los Angeles, CA',
false,
'2024-09-20 20:00:00',
'2024-09-20 23:00:00',
100,
4500,
'USD',
'published',
0);

-- Verify the insert
SELECT COUNT(*) as total_events FROM events WHERE status = 'published';
SELECT title, location, capacity FROM events ORDER BY starts_at;
