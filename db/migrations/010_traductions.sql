CREATE TABLE traductions (
    langue_id INT NOT NULL REFERENCES langues(id) ON DELETE CASCADE,
    cle VARCHAR(100) NOT NULL,
    texte TEXT NOT NULL,
    PRIMARY KEY (langue_id, cle)
);
