-- +goose Up
INSERT INTO comments (post_id, user_id, content)
VALUES
    (1, 2, 'Gratulacje! React potrafi dać w kość 😄'),
    (1, 3, 'Super robota! 🚀'),
    (1, 4, 'Też ostatnio kończyłem projekt w React – znam ten ból 😅'),
    (1, 5, 'Jakiego stacka użyłeś?'),

    (2, 1, 'Wygląda świetnie, minimalistyczny vibe ✨'),
    (2, 3, 'Fajny system! Zrobiłeś komponenty od zera?'),
    (2, 5, 'Mega mi się podoba!'),

    (3, 2, 'Klasyka debugowania 😅'),
    (3, 4, 'Najlepsze uczucie, gdy znajdziesz ten jeden błąd!'),
    (3, 5, 'Bywa tak... ważne, że działa!'),

    (4, 1, 'Powodzenia na nowym stanowisku! 💪'),
    (4, 3, 'Gratulacje!'),
    (4, 5, 'Kibicuję!'),

    (5, 2, 'Oho, update o 2:00 AM – życzę powodzenia 😄'),
    (5, 4, 'Znam ten ból nocnych aktualizacji...'),

    (6, 1, 'Też czasem tak mam, freelancing jest ciężki 😅'),
    (6, 3, 'Wolność + niepewność – klasyka!'),
    (6, 5, 'Ważne żeby znaleźć balans.'),

    (7, 1, 'Chętnie zobaczę prototyp!'),
    (7, 4, 'Brzmi super!'),

    (8, 2, 'Code review zawsze daje dużo wiedzy 👍'),
    (8, 3, 'Zgadzam się, można nauczyć się masy nowych rzeczy.')
;

-- +goose Down
DELETE FROM comments
WHERE post_id IN (1,2,3,4,5,6,7,8);
