-- Textes de l'interface, par langue : le site est multilingue sans modifier le code.
-- Ajouter une langue = une ligne dans langues + ses lignes dans traductions.

CREATE TABLE traductions (
    langue_id INT NOT NULL REFERENCES langues(id) ON DELETE CASCADE,
    cle VARCHAR(100) NOT NULL,
    texte TEXT NOT NULL,
    PRIMARY KEY (langue_id, cle)
);
