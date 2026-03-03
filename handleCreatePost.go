package main

import ( // Este es un servidor HTTP básico en Go que escucha en el puerto 8080 y responde con "Hello World" a cualquier solicitud.

	"database/sql" // "database/sql" se utiliza para interactuar con una base de datos SQL, aunque en este código no se muestra su uso específico, es común en aplicaciones web para manejar datos persistentes.
	"fmt"          // "fmt" se utiliza para formatear la respuesta HTTP que se enviará al cliente.

	"net/url"
	"strconv"
	_ "modernc.org/sqlite" // Este es un import anónimo que se utiliza para registrar el controlador de SQLite con el paquete "database/sql". Esto permite que el programa utilice SQLite como base de datos sin necesidad de importar explícitamente el paquete en el código. El guion bajo (_) indica que el paquete se importa solo por sus efectos secundarios, es decir, para registrar el controlador de la base de datos, sin utilizar directamente ninguna función o tipo del paquete en el código.
)




func handleCreatePost(body string, db *sql.DB) string {

	values, err := url.ParseQuery(body)
    if err != nil {
        fmt.Println("Error parseando body:", err)
        return "HTTP/1.1 400 Bad Request\r\n\r\n"
    }

	name := values.Get("series_name")
	currentEp, _ := strconv.Atoi(values.Get("current_episode"))
	totalEps, _ := strconv.Atoi(values.Get("total_episodes"))

    fmt.Println("Nombre:", name)
    fmt.Println("Episodio actual:", currentEp)
    fmt.Println("Total episodios:", totalEps)



	_, err = db.Exec("INSERT INTO series (name, current_episode, total_episodes) VALUES (?, ?, ?)", name, currentEp, totalEps)
    if err != nil {
        fmt.Println("Error insertando serie:", err)
        return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
    }

    // Redirigir al index
    return "HTTP/1.1 303 See Other\r\nLocation: /\r\n\r\n"

}