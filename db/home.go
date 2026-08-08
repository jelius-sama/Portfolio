package db

func createHomeTables() error {
    var schema = `
        CREATE TABLE IF NOT EXISTS home_page (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            cv_path TEXT NOT NULL,
            image_path TEXT NOT NULL,
            image_alt_text TEXT,
            first_name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            education_level INTEGER NOT NULL,
            study_domain TEXT NOT NULL,
            who_am_i TEXT NOT NULL,
            email TEXT NOT NULL,
            qr_image_path TEXT NOT NULL,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS home_specialized_domain (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            home_page_id INTEGER NOT NULL,
            domain TEXT NOT NULL,
            position INTEGER NOT NULL,
            FOREIGN KEY (home_page_id) REFERENCES home_page(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS home_about_message (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            home_page_id INTEGER NOT NULL,
            message TEXT NOT NULL,
            position INTEGER NOT NULL,
            FOREIGN KEY (home_page_id) REFERENCES home_page(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS home_skill (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            home_page_id INTEGER NOT NULL,
            category TEXT NOT NULL,
            skill TEXT NOT NULL,
            position INTEGER NOT NULL,
            FOREIGN KEY (home_page_id) REFERENCES home_page(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS home_project (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            home_page_id INTEGER NOT NULL,
            title TEXT NOT NULL,
            description TEXT NOT NULL,
            thumbnail TEXT NOT NULL,
            position INTEGER NOT NULL,
            FOREIGN KEY (home_page_id) REFERENCES home_page(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS home_project_tag (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            project_id INTEGER NOT NULL,
            tag TEXT NOT NULL,
            position INTEGER NOT NULL,
            FOREIGN KEY (project_id) REFERENCES home_project(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS home_project_link (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            project_id INTEGER NOT NULL,
            label TEXT NOT NULL,
            href TEXT NOT NULL,
            kind TEXT NOT NULL,
            position INTEGER NOT NULL,
            FOREIGN KEY (project_id) REFERENCES home_project(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS home_experience (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            home_page_id INTEGER NOT NULL,
            name TEXT NOT NULL,
            description TEXT NOT NULL,
            date_range TEXT NOT NULL,
            is_active BOOLEAN NOT NULL DEFAULT 0,
            position INTEGER NOT NULL,
            FOREIGN KEY (home_page_id) REFERENCES home_page(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS home_contact_quick_link (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            home_page_id INTEGER NOT NULL,
            platform TEXT NOT NULL,
            url TEXT NOT NULL,
            icon TEXT NOT NULL,
            position INTEGER NOT NULL,
            FOREIGN KEY (home_page_id) REFERENCES home_page(id) ON DELETE CASCADE
        );

        CREATE INDEX IF NOT EXISTS idx_home_specialized_domain_page_id ON home_specialized_domain(home_page_id);
        CREATE INDEX IF NOT EXISTS idx_home_about_message_page_id ON home_about_message(home_page_id);
        CREATE INDEX IF NOT EXISTS idx_home_skill_page_id ON home_skill(home_page_id);
        CREATE INDEX IF NOT EXISTS idx_home_project_page_id ON home_project(home_page_id);
        CREATE INDEX IF NOT EXISTS idx_home_project_tag_project_id ON home_project_tag(project_id);
        CREATE INDEX IF NOT EXISTS idx_home_project_link_project_id ON home_project_link(project_id);
        CREATE INDEX IF NOT EXISTS idx_home_experience_page_id ON home_experience(home_page_id);
        CREATE INDEX IF NOT EXISTS idx_home_contact_quick_link_page_id ON home_contact_quick_link(home_page_id);
    `

    if _, err := DB.Exec(schema); err != nil {
        return err
    }

    return nil
}

