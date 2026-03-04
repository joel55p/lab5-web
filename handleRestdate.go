package main

import (
    "database/sql"

    "log"
)

func handleRestdate(db *sql.DB, id string) string { // decrementa el episodio actual de una serie x
    _, err := db.Exec(
        "UPDATE series SET current_episode = current_episode - 1 WHERE id = ? AND current_episode > 0",
        id,
    )
    if err != nil {
        log.Print("Error decrementando:", err)
        return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
    }
    return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nok"
}