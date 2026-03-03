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
<style>
	body { font-family: Arial, sans-serif; padding: 20px; }
	h1   { color: #333; }
	table { width: 100%; border-collapse: collapse; margin-top: 20px; }
	th, td { border: 1px solid #ccc; padding: 10px; text-align: left; }
	th { background-color: #f4f4f4; }
</style>
</head>
<body>
<h1>Track de mis series actuales</h1>
<table>
	<tr><th>No.</th><th>Nombre</th><th>Episodio actual</th><th>Total de episodios</th><th>Agregar a episodio actual</th></tr>
<a href="/create" target="_blank">Agregar Serie</a>

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

			// Agregar una fila a la tabla por cada serie y se agrega ahora boton
		html += fmt.Sprintf( //aqui se coloca el boton ya que se necesita el id de cada serie para que el boton funcione, entonces se hace dentro del for para que se agregue un boton por cada serie en la tabla, y cada boton tenga el id de su respectiva serie para que al hacer click en el boton se pueda identificar a qué serie se le debe actualizar el episodio actual.
			"<tr><td>%d</td><td>%s</td><td>%d</td><td>%d</td><td><button onclick='nextEpisode(%d)'>+1</button></td></tr>",
		id, name, currentEpisode, totalEpisodes, id,
		)
	}

	html += `</table></body></html>` //cierre del html


		// Enviar respuesta HTTP
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n%s", html) //se devuelve la respuesta HTTP con el código de estado 200 OK, el encabezado Content-Type indicando que el contenido es HTML, y luego el cuerpo de la respuesta que contiene el HTML construido con la tabla de series.
	//se hace return de tipo string porque la función handleIndex está definida para devolver un string, que es la respuesta HTTP completa que se enviará al cliente. El formato de la respuesta incluye el código de estado, los encabezados y el cuerpo HTML.
}