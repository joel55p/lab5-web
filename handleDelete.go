package main

import (
    "database/sql"
    "log"
)

func handleDelete(db *sql.DB, id string) string { // elimina una serie específica de la base de datos por su id
    _, err := db.Exec("DELETE FROM series WHERE id = ?", id)
    if err != nil {
        log.Print("Error sl querer elimar serie:", err)
        return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
    }
    return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nok"
}