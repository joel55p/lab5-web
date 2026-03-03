package main // El paquete "main" es el punto de entrada de la aplicación en Go. Es necesario para ejecutar el programa.

import ( // Este es un servidor HTTP básico en Go que escucha en el puerto 8080 y responde con "Hello World" a cualquier solicitud.
	"bufio"        // "bufio" se utiliza para leer la solicitud del cliente de manera eficiente.
	"database/sql" // "database/sql" se utiliza para interactuar con una base de datos SQL, aunque en este código no se muestra su uso específico, es común en aplicaciones web para manejar datos persistentes.
	"fmt"          // "fmt" se utiliza para formatear la respuesta HTTP que se enviará al cliente.
	"log"          // "log" se utiliza para registrar errores y mensajes informativos en la consola.
	"net"          // "net" se utiliza para crear un servidor TCP que escuche en el puerto 8080 y acepte conexiones entrantes. uno de os paquetes más importantes para la comunicación de red en Go.
	"strings"      // "strings" se utiliza para manipular la cadena de texto de la solicitud HTTP, como extraer el path solicitado.
	"net/url"
	"strconv"
	_ "modernc.org/sqlite" // Este es un import anónimo que se utiliza para registrar el controlador de SQLite con el paquete "database/sql". Esto permite que el programa utilice SQLite como base de datos sin necesidad de importar explícitamente el paquete en el código. El guion bajo (_) indica que el paquete se importa solo por sus efectos secundarios, es decir, para registrar el controlador de la base de datos, sin utilizar directamente ninguna función o tipo del paquete en el código.
)

func main() { // La función "main" es el punto de entrada del programa. Aquí se configura el servidor TCP y se maneja la lógica principal para aceptar conexiones.
	// se necesita hacer primero una conexion a la base de datos
	db, err := sql.Open("sqlite", "series.db") //se abre conexion una vez en el main, se hace con sql.open que recibe el tipo de base de datos y el nombre del archivo
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()                              //se cierra cuando se termine el main, es decir, cuando se cierre el servidor, se cierra la conexion a la base de datos para liberar recursos. es importante cerrar la conexion a la base de datos para evitar fugas de memoria y asegurar que los recursos se liberen adecuadamente cuando el programa termine su ejecución.

	listener, err := net.Listen("tcp", ":8080") // "net.Listen" crea un servidor TCP que escucha en el puerto 8080. Si hay un error al crear el servidor, se registra y se termina el programa. si funciona se guarda el listener para aceptar conexiones entrantes.
	if err != nil {                             // Si ocurre un error al intentar escuchar en el puerto, se registra el error y se termina el programa.
		log.Fatal(err) // "log.Fatal" registra el error y termina la ejecución del programa. Esto es útil para asegurarse de que el servidor no intente continuar si no puede escuchar en el puerto especificado.
	}
	defer listener.Close()              // "defer" asegura que el listener se cerrará cuando la función "main" termine, lo que es importante para liberar recursos. y main va a esperar a que se cierre el listener antes de finalizar la ejecución del programa.
	log.Print("Listening on port 8080") // Se registra un mensaje en la consola indicando que el servidor está escuchando en el puerto 8080.

	for { // Este es un bucle infinito que acepta conexiones entrantes. Cada vez que se acepta una conexión, se maneja en una goroutine separada para permitir que el servidor continúe aceptando otras conexiones mientras se procesa la actual.
		conn, err := listener.Accept() // "listener.Accept" espera a que llegue una conexión entrante y la acepta. Si hay un error al aceptar la conexión, se registra el error y se continúa con el siguiente ciclo del bucle para esperar otra conexión.
		if err != nil {                // Si ocurre un error al aceptar la conexión, se registra el error y se continúa con el siguiente ciclo del bucle para esperar otra conexión.
			log.Print("Error accepting:", err) // "log.Print" registra el error que ocurrió al aceptar la conexión, pero no termina el programa. Esto permite que el servidor siga funcionando y acepte otras conexiones a pesar de los errores ocasionales.
			continue                           // "continue" se utiliza para saltar el resto del código en el bucle actual y pasar a la siguiente iteración, lo que permite que el servidor siga aceptando conexiones incluso si ocurre un error al aceptar una conexión específica.
		}
		go handle(conn, db) // "go handle(conn)" inicia una nueva goroutine para manejar la conexión aceptada. Esto permite que el servidor procese múltiples conexiones simultáneamente sin bloquear el bucle principal que acepta conexiones. Cada conexión se maneja de forma independiente en su propia goroutine, lo que mejora la capacidad de respuesta del servidor.
	}
} //cabe mencionar que conn, err := listener.Accept() lo que hace es que si acepta se guarda en la variable conn, y si no acepta se guarda el error en la variable err, y se maneja el error con un if para evitar que el programa se caiga.

// osea que if err !=nil lo que significa es que si si hay un error, entonces se ejecuta el bloque de código dentro del if, que en este caso es log.Print("Error accepting:", err) y continue para seguir aceptando conexiones. Si no hay error, entonces se ejecuta go handle(conn) para manejar la conexión aceptada en una goroutine separada.
// una gorutine es una función que se ejecuta de manera concurrente con otras funciones. Es una forma ligera de manejar múltiples tareas al mismo tiempo sin bloquear el programa principal. En este caso, cada vez que se acepta una conexión, se inicia una nueva goroutine para manejar esa conexión específica, lo que permite que el servidor siga aceptando otras conexiones mientras se procesa la actual.
func handle(conn net.Conn, db *sql.DB) {
	defer conn.Close() //cierra conexion al finalizar fun para liberar recursos
	reader := bufio.NewReader(conn)

	// Leer request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		log.Print("Error reading request:", err)
		return
	}

	parts := strings.Fields(requestLine)
	path := "/"
	method := "GET"
	if len(parts) >= 2 { 
		path = parts[1] //se obtiene el path solicitado, que es la segunda parte de la request line (después del método HTTP).
		method = parts[0] //se obtiene el metodo http, que es la primera parte de la request line (antes del path).
		
	}
	


	contentLength := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
		if strings.HasPrefix(line, "Content-Length:") {
			lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			contentLength, _ = strconv.Atoi(lengthStr)
		}
	}

	// Leer body
	body := ""
	if contentLength > 0 {
		bodyBytes := make([]byte, contentLength)
		reader.Read(bodyBytes)
		body = string(bodyBytes)
	}







		// ---Routing---
	var response string

	switch{
	case method == "GET" && path == "/":
		response = handleIndex( db) //se llama a la función handleIndex para manejar la solicitud GET al path "/". Esta función se encargará de construir la respuesta HTML con la tabla de series desde la base de datos y devolverla al cliente.
	case method == "GET" && path == "/create":
		response = handleCreate() //se llama a la función handleCreate para manejar la solicitud GET al path "/create". Esta función aún no está implementada, pero se espera que maneje la lógica para agregar una nueva serie a la base de datos y devolver una respuesta adecuada al cliente.
	case method == "POST" && path == "/create":
		response = handleCreatePost(body, db) //se llama a la función handleCreatePost para manejar la solicitud POST al path "/create". Esta función aún no está implementada, pero se espera que maneje la lógica para procesar los datos enviados desde el formulario de creación de una nueva serie, agregar esa serie a la base de datos y devolver una respuesta adecuada al cliente.
	default:
		response = "HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\n\r\n<h1>404 Not Found</h1>" //si la solicitud no coincide con ninguna de las rutas definidas, se devuelve una respuesta HTTP con el código de estado 404 Not Found y un mensaje HTML indicando que la página no fue encontrada.
	}

	conn.Write([]byte(response))  // handle() escribe
}

	

	

	
		
	
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
	<tr><th>No.</th><th>Nombre</th><th>Episodio actual</th><th>Total de episodios</th></tr>
		
<a href="/create" target="_blank">Agregar Serie</a>



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

			// Agregar una fila a la tabla por cada serie
		html += fmt.Sprintf(
			"<tr><td>%d</td><td>%s</td><td>%d</td><td>%d</td></tr>",
			id, name, currentEpisode, totalEpisodes,
		)
	}

	html += `</table></body></html>` //cierre del html


		// Enviar respuesta HTTP
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n%s", html) //se devuelve la respuesta HTTP con el código de estado 200 OK, el encabezado Content-Type indicando que el contenido es HTML, y luego el cuerpo de la respuesta que contiene el HTML construido con la tabla de series.
	//se hace return de tipo string porque la función handleIndex está definida para devolver un string, que es la respuesta HTTP completa que se enviará al cliente. El formato de la respuesta incluye el código de estado, los encabezados y el cuerpo HTML.
}
	




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