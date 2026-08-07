INSERT INTO links_page (
    handle, tag_line, image_path, image_alt_text, who_am_i, qr_image_path, updated_at
) VALUES (
    'jelius_sama',
    'Undergraduate Student • Open Source Enthusiast',
    '/assets/compressed/jelius.webp',
    'Portrait of Jelius Basumatary',
    'I''m an undergraduate CS student who loves building things that work well and are easy to use.',
    '/assets/compressed/links.webp',
    datetime('now')
);

INSERT INTO link_entries (links_page_id, icon, title, subtitle, href, position, updated_at) VALUES
(1, 'forgejo', 'Git', 'Check out my code and projects', 'https://git.jelius.dev/jelius-sama', 1, datetime('now')),
(1, 'forgejo', 'Git for Uni. Projects', 'University Project Repositories', 'https://git.jelius.dev/Jelius-ADBU', 2, datetime('now')),
(1,  'briefcase', 'Freelance', 'Hire me for freelance projects and collaborations', 'https://freelance.jelius.dev', 3, datetime('now')),
(1, 'linkedin', 'LinkedIn', 'Connect with me professionally', 'https://www.linkedin.com/in/jelius-basumatary-485044339/', 4, datetime('now')),
(1, 'code', 'Portfolio', 'View my portfolio site', '/', 5, datetime('now')),
(1, 'x.com', 'X', 'Follow me for tech updates', 'https://x.com/jelius_sama', 6, datetime('now')),
(1, 'youtube', 'YouTube', 'Watch my videos when I start uploading content.', 'https://youtube.com/@jelius-sama', 7, datetime('now')),
(1, 'instagram', 'Instagram', 'Behind the scenes content', 'https://instagram.com/_jelius_sama', 8, datetime('now')),
(1, 'mail', 'Email Me', 'Get in touch directly', 'mailto:contact@jelius.dev', 9, datetime('now')),
(1, 'server', 'Status', 'Uptime monitoring of my services', 'https://status.jelius.dev', 10, datetime('now'));


INSERT INTO home_page (
    cv_path, image_path, image_alt_text, first_name, last_name,
    education_level, study_domain, who_am_i, email, qr_image_path, updated_at
) VALUES (
    '/assets/cv.pdf',
    '/assets/compressed/jelius.webp',
    'Portrait of Jelius Basumatary',
    'Jelius',
    'Basumatary',
    0,
    'Computer Science',
    'An undergraduate CS student who loves building things that work well and are easy to use.',
    'contact@jelius.dev',
    '/assets/jelius-dev-dark.png',
    datetime('now')
);

INSERT INTO home_specialized_domain (home_page_id, domain, position) VALUES
(1, 'Backend Developer', 1),
(1, 'Systems Programmer', 2);

INSERT INTO home_about_message (home_page_id, message, position) VALUES
(1, 'I''m an undergraduate CS student who loves building things that work well and are easy to use. I enjoy solving problems, making code simple, and learning new things every day.', 1),
(1, 'Right now, I''m working toward my bachelor''s degree. I enjoy developing software that makes a real difference in everyday life, and I''m always curious how things work under the hood.', 2),
(1, 'When I''m not coding, you''ll usually find me exploring new technologies or sharing what I''ve learned with the developer community.', 3);

INSERT INTO home_skill (home_page_id, category, skill, position) VALUES
(1, 'language', 'Rust', 1),
(1, 'language', 'Swift', 2),
(1, 'language', 'C', 3),
(1, 'language', 'TypeScript', 4),
(1, 'language', 'Go', 5),
(1, 'framework', 'Node.js', 1),
(1, 'framework', 'SolidJS', 2),
(1, 'framework', 'React', 3),
(1, 'framework', 'HTMX', 4),
(1, 'framework', 'Templ', 5),
(1, 'other', 'Redis', 1),
(1, 'other', 'PostgreSQL', 2),
(1, 'other', 'AWS', 3),
(1, 'other', 'Cloudflare', 4),
(1, 'other', 'Git', 5);

INSERT INTO home_project (home_page_id, title, description, thumbnail, position) VALUES
(1, 'Pixelle', 'An anime image gallery application with collections of illustrations and photography.', '/assets/compressed/project-pixelle.webp', 1),
(1, 'VPS Watch Dog', 'A watch dog program that monitors your VPS''s system usage and alerts via mail if usage is running high.', '/assets/compressed/VPSWatchDog.webp', 2),
(1, 'Storage Watch Dog', 'A program that watches a specific directory and if the storage space is running low alerts via mails.', '/assets/compressed/StorageWatchDog.webp', 3),
(1, 'AWS Mail Parser', 'Polls your AWS SQS and when a new mail event is detected fetches that mail from S3 bucket and saves it to your maildir after parsing it.', '/assets/compressed/AWSMailParser.webp', 4),
(1, 'Convert CBZ', 'A high-performance, concurrent tool for converting folders containing images into CBZ (Comic Book Archive) files. Built in Go for speed and reliability.', '/assets/compressed/convert_cbz.webp', 5);

INSERT INTO home_project_tag (project_id, tag, position) VALUES
(1, 'Next.js', 1), (1, 'TypeScript', 2), (1, 'PostgreSQL', 3),
(2, 'Go', 1), (2, 'SMTP', 2), (2, 'POSIX API', 3),
(3, 'Go', 1), (3, 'SMTP', 2),
(4, 'Go', 1), (4, 'AWS', 2), (4, 'File Parsing', 3), (4, 'IMAP', 4),
(5, 'Go', 1), (5, 'Archive', 2), (5, 'Concurrency', 3);

INSERT INTO home_project_link (project_id, label, href, kind, position) VALUES
(1, 'Code', 'https://github.com/jelius-sama/pixelle-demo', 'code', 1),
(1, 'Live Demo', 'https://pixelle.jelius.dev/', 'demo', 2),
(2, 'Code', 'https://github.com/jelius-sama/VPSWatchDog', 'code', 1),
(3, 'Code', 'https://github.com/jelius-sama/StorageWatchDog', 'code', 1),
(4, 'Code', 'https://github.com/jelius-sama/AWSMailParser', 'code', 1),
(5, 'Code', 'https://github.com/jelius-sama/convert_cbz', 'code', 1),
(5, 'Blog Post', 'https://jelius.dev/blog/1a42b1d', 'blog', 2);

INSERT INTO home_experience (home_page_id, name, description, date_range, is_active, position) VALUES
(1, 'Birth', 'ST. Augustine Hospital • New born', 'November 2007', 0, 1),
(1, 'Student', 'Holy Child • Primary School', 'Jan 2009 – Dec 2013', 0, 2),
(1, 'Student', 'ST. John''s • High School', 'Jan 2014 – April 2025', 0, 3),
(1, 'Undergraduate', 'Assam Don Bosco • University', 'April 2025 – Present', 1, 4);

INSERT INTO home_contact_quick_link (home_page_id, platform, url, icon, position) VALUES
(1, 'X', 'https://x.com/jelius_sama', 'x.com', 1),
(1, 'LinkedIn', 'https://www.linkedin.com/in/jelius-basumatary-485044339', 'linkedin', 2),
(1, 'Git', 'https://git.jelius.dev/jelius-sama', 'forgejo', 3),
(1, 'Linktree', '/links', 'linktree', 4);
