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
<link rel="icon" href="/faviconm.png" type="image/png">
<style>
    body { font-family: Arial, sans-serif; background: #fafafa; display: flex; justify-content: center; padding: 60px 20px; }
    .card { background: white; padding: 32px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.08); width: 100%; max-width: 420px; }
    h1 { color: #333; margin-bottom: 24px; font-size: 22px; }
    label { display: block; margin-bottom: 6px; color: #555; font-size: 14px; }
    input[type="text"], input[type="number"] { width: 100%; padding: 9px 12px; margin-bottom: 14px; border: 1px solid #ddd; border-radius: 6px; font-size: 14px; }
    input[type="submit"] { background: #4CAF50; color: white; padding: 10px; border: none; border-radius: 6px; cursor: pointer; font-size: 15px; width: 100%; }
    input[type="submit"]:hover { background: #45a049; }
    a { display: inline-block; margin-top: 14px; color: #777; font-size: 13px; text-decoration: none; }
    a:hover { text-decoration: underline; }
</style>
</head>
<body>
<div class="card">
    <h1> Agregar Serie</h1>
    <form method="POST" action="/create">
        <label>Nombre de la serie:</label>
        <input type="text" name="series_name" required>

        <label>Episodio actual:</label>
        <input type="number" name="current_episode" min="1" value="1" required>

        <label>Total de episodios:</label>
        <input type="number" name="total_episodes" min="1" required>

        <input type="submit" value="Agregar Serie">
    </form>
    <a href="/"> Volver al Track de Series</a>
</div>
</body>
</html>`

	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n%s", html) //se devuelve la respuesta HTTP con el código de estado 200 OK, el encabezado Content-Type indicando que el contenido es HTML, y luego el cuerpo de la respuesta que contiene el formulario HTML para agregar una nueva serie.

    

}
