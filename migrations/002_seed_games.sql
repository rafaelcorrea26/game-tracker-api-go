INSERT INTO games (name, status, year_completed, rating, notes)
VALUES
('Elden Ring','completed',2024,9,'jogo incrível')
ON CONFLICT DO NOTHING;