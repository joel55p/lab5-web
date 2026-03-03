package main

import ( // Este es un servidor HTTP básico en Go que escucha en el puerto 8080 y responde con "Hello World" a cualquier solicitud.

	"fmt"          // "fmt" se utiliza para formatear la respuesta HTTP que se enviará al cliente.

	_ "modernc.org/sqlite" // Este es un import anónimo que se utiliza para registrar el controlador de SQLite con el paquete "database/sql". Esto permite que el programa utilice SQLite como base de datos sin necesidad de importar explícitamente el paquete en el código. El guion bajo (_) indica que el paquete se importa solo por sus efectos secundarios, es decir, para registrar el controlador de la base de datos, sin utilizar directamente ninguna función o tipo del paquete en el código.
)



func handleCreate() string {

	var html string

	html = `<html>
<head>
<title>Agregar Serie</title>
<style>
	body { font-family: Arial, sans-serif; padding: 20px; }
	h1   { color: #333; }
	form { margin-top: 20px; }
	label { display: block; margin-bottom: 5px; }
	input[type="text"], input[type="number"] { width: 100%; padding: 8px; margin-bottom: 10px; border: 1px solid #ccc; border-radius: 4px; }
	input[type="submit"] { background-color: #4CAF50; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
	input[type="submit"]:hover { background-color: #45a049; }
</style>
</head>
<body>
<h1>Agregar Serie</h1>
<form method="POST" action="/create">
	<label for="name">Nombre de la serie:</label>
	<input type="text" id="name" name="series_name" required>

	<label for="currentEpisode">Episodio actual:</label>
	<input type="number" id="currentEpisode" name="current_episode" min="1" value="1" required>

	<label for="totalEpisodes">Total de episodios:</label>
	<input type="number" id="totalEpisodes" name="total_episodes" min="1"  required>

	<input type="submit" value="Agregar Serie">
</form>

<a href="/" target="_blank">Volver al Track de Series</a>
</body>
</html>`

	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n%s", html) //se devuelve la respuesta HTTP con el código de estado 200 OK, el encabezado Content-Type indicando que el contenido es HTML, y luego el cuerpo de la respuesta que contiene el formulario HTML para agregar una nueva serie.

    

}
