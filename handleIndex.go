package main

import ( // Este es un servidor HTTP básico en Go que escucha en el puerto 8080 y responde con "Hello World" a cualquier solicitud.

	"database/sql" // "database/sql" se utiliza para interactuar con una base de datos SQL, aunque en este código no se muestra su uso específico, es común en aplicaciones web para manejar datos persistentes.
	"fmt"          // "fmt" se utiliza para formatear la respuesta HTTP que se enviará al cliente.
	"log"          // "log" se utiliza para registrar errores y mensajes informativos en la consola.

	_ "modernc.org/sqlite" // Este es un import anónimo que se utiliza para registrar el controlador de SQLite con el paquete "database/sql". Esto permite que el programa utilice SQLite como base de datos sin necesidad de importar explícitamente el paquete en el código. El guion bajo (_) indica que el paquete se importa solo por sus efectos secundarios, es decir, para registrar el controlador de la base de datos, sin utilizar directamente ninguna función o tipo del paquete en el código.
)

func handleIndex( db *sql.DB) string {
		// Construir la tabla HTML desde la base de datos
	var html string //se define una variable donde se colocara el html

	rows, err := db.Query("SELECT * FROM series") //se selecciona eso de mi base de datos
	if err != nil {
		log.Print("Error querying database:", err) //error

	}
	defer rows.Close() //para cerrar la consulta a  la base de datos

	// Inicio del HTML con estilos
	html = `<html> 
<head>
<title>Track de mis series actuales</title>
<link rel="icon" href="/faviconm.png" type="image/png">
<style>
    body { font-family: Arial, sans-serif; padding: 30px; background: #fafafa; }
    h1 { color: #333; margin-bottom: 16px; }
    a { color: #555; text-decoration: none; display: inline-block; margin-bottom: 16px; }
    a:hover { text-decoration: underline; } 
    table { width: 100%; border-collapse: collapse; background: white; box-shadow: 0 2px 8px rgba(0,0,0,0.08); border-radius: 8px; overflow: hidden; }
    th { background: #444; color: white; padding: 12px; text-align: left; }
    td { padding: 12px; border-bottom: 1px solid #eee; }
    tr:last-child td { border-bottom: none; }
    tr:hover td { background: #f5f5f5; }
    progress { width: 100px; height: 10px; accent-color: #4CAF50; }
    button { padding: 5px 12px; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; }
    .btn-add { background: #4CAF50; color: white; }
    .btn-sub { background: #e74c3c; color: white; }
    .btn-del { background: #aaa; color: white; }
</style>
</head>
<body>
<h1> Track de mis series actuales</h1>
<table>
    <tr>
        <th>No.</th>
        <th>Nombre</th>
        <th>Episodio  actual</th>
        <th>Total Episodios</th>
        <th>Progreso</th>
        <th>agregar episodio</th>
        <th>Eliminar episodio</th>
        <th>Eliminar serie</th>
    </tr>
<a href="/create" target="_blank">+ Agregar Serie</a>

<script src="/script.js"></script>



`


		// Iterar sobre cada fila de la base de datos para que se construya la tabla en el html
	for rows.Next() { //next justamente hace eso como un puntero que avanza y se detiene cuando ya no hay más filas, entonces el for se ejecuta mientras haya filas para procesar. cada vez que se llama a rows.Next(), se mueve al siguiente registro en el conjunto de resultados de la consulta. Si hay un registro disponible, devuelve true y permite que el bloque de código dentro del for se ejecute para procesar ese registro. Si no hay más registros disponibles, devuelve false y el bucle termina.
		var id, currentEpisode, totalEpisodes int
		var name string

		// Scan guarda los valores de la fila actual en las variables
		err := rows.Scan(&id, &name, &currentEpisode, &totalEpisodes)
		if err != nil {
			log.Print("Error scanning row:", err)
			continue
		}

		// Marcar serie como completada si ya se vieron todos los episodios
		completada := ""
		if currentEpisode == totalEpisodes {
			completada = "<span style='color: gold; font-weight: bold;'> <strong>COMPLETADA!!!!</strong></span>"
		}

		html += fmt.Sprintf(
			"<tr><td>%d</td><td>%s %s</td><td>%d</td><td>%d</td><td><progress value='%d' max='%d'></progress> %d/%d</td><td><button class='btn-add' onclick='nextEpisode(%d)'>+1</button></td><td><button class='btn-sub' onclick='prevEpisode(%d)'>-1</button></td><td><button class='btn-del' onclick='deleteSerie(%d)'>Eliminar</button></td></tr>",
			id, name, completada, currentEpisode, totalEpisodes, currentEpisode, totalEpisodes, currentEpisode, totalEpisodes, id, id, id,
		)

	}

	html += `</table></body></html>` //cierre del html


		// Enviar respuesta HTTP
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n%s", html) //se devuelve la respuesta HTTP con el código de estado 200 OK, el encabezado Content-Type indicando que el contenido es HTML, y luego el cuerpo de la respuesta que contiene el HTML construido con la tabla de series.
	//se hace return de tipo string porque la función handleIndex esta definida para devolver un string, que es la respuesta HTTP completa que se enviará al cliente. El formato de la respuesta incluye el código de estado, los encabezados y el cuerpo HTML.
}