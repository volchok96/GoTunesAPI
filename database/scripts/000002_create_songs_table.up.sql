-- Создаем таблицу songs
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,          -- Уникальный идентификатор песни
    group_name VARCHAR(255) NOT NULL,  -- Название группы
    song_name VARCHAR(255) NOT NULL,   -- Название песни
    release_date DATE,                 -- Дата релиза
    text TEXT,                         -- Текст песни
    link VARCHAR(2083)                 -- Ссылка на песню (например, YouTube)
);

-- Добавляем индекс для быстрого поиска по группе и названию песни
CREATE INDEX idx_group_song ON songs (group_name, song_name);
