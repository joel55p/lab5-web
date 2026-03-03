package main
import ( // Este es un servidor HTTP básico en Go que escucha en el puerto 8080 y responde con "Hello World" a cualquier solicitud.

	"database/sql" // "database/sql" se utiliza para interactuar con una base de datos SQL, aunque en este código no se muestra su uso específico, es común en aplicaciones web para manejar datos persistentes.
	"fmt"          // "fmt" se utiliza para formatear la respuesta HTTP que se enviará al cliente.

	_ "modernc.org/sqlite" // Este es un import anónimo que se utiliza para registrar el controlador de SQLite con el paquete "database/sql". Esto permite que el programa utilice SQLite como base de datos sin necesidad de importar explícitamente el paquete en el código. El guion bajo (_) indica que el paquete se importa solo por sus efectos secundarios, es decir, para registrar el controlador de la base de datos, sin utilizar directamente ninguna función o tipo del paquete en el código.
)



func handleUpdate(db *sql.DB, id string) string { //para actualizar el episodio actual de una serie  espeecifica por eso es que se le pasa como parametro el id
    _, err := db.Exec(
        "UPDATE series SET current_episode = current_episode + 1 WHERE id = ? AND current_episode < total_episodes",
        id,
    )
    if err != nil {
        fmt.Println("Error actualizando:", err)
        return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
    }
    return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nok"
}