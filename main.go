package main // El paquete "main" es el punto de entrada de la aplicación en Go. Es necesario para ejecutar el programa.

import ( // Este es un servidor HTTP básico en Go que escucha en el puerto 8080 y responde con "Hello World" a cualquier solicitud.
	"bufio"        // "bufio" se utiliza para leer la solicitud del cliente de manera eficiente.
	"database/sql" // "database/sql" se utiliza para interactuar con una base de datos SQL, aunque en este código no se muestra su uso específico, es común en aplicaciones web para manejar datos persistentes.
	"fmt"          // "fmt" se utiliza para formatear la respuesta HTTP que se enviará al cliente.
	"log"          // "log" se utiliza para registrar errores y mensajes informativos en la consola.
	"net"          // "net" se utiliza para crear un servidor TCP que escuche en el puerto 8080 y acepte conexiones entrantes. uno de os paquetes más importantes para la comunicación de red en Go.
	"strings"      // "strings" se utiliza para manipular la cadena de texto de la solicitud HTTP, como extraer el path solicitado.

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
	if len(parts) >= 2 {
		path = parts[1]
	}

	// Consumir headers hasta línea vacía
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
	}

	// Construir la tabla HTML desde la base de datos
	var html string //se define una variable donde se colocara el html

	if path == "/" { //se empieza con /
		rows, err := db.Query("SELECT * FROM series") //se selecciona eso de mi base de datos
		if err != nil {
			log.Print("Error querying database:", err) //error
			return
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
		
	<script>

    alert("Cuidado estas a punto de ver las mejores pinches series de la historia");
  	</script>`
	

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

	} else {
		html = "<html><body><h1>404 - Página no encontrada</h1></body></html>"
	}

	// Construir y enviar la respuesta HTTP
	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+ //la primera línea de la respuesta HTTP, que indica el protocolo (HTTP/1.1), el código de estado (200) y el mensaje de estado (OK). Esto le dice al cliente que la solicitud fue exitosa y que se está enviando una respuesta con contenido.
			"Content-Type: text/html\r\n"+ //el encabezado "Content-Type" indica al cliente que el contenido de la respuesta es HTML, lo que le permite al navegador interpretar correctamente el contenido y mostrarlo como una página web.
			"Content-Length: %d\r\n"+ //el encabezado "Content-Length" especifica la longitud del cuerpo de la respuesta en bytes. Esto es importante para que el cliente sepa cuánto contenido esperar y pueda manejarlo adecuadamente.
			"Connection: close\r\n"+ //el encabezado "Connection: close" indica al cliente que el servidor cerrará la conexión después de enviar la respuesta. Esto es útil para liberar recursos en el servidor y evitar conexiones persistentes innecesarias, especialmente en un servidor simple como este.
			"\r\n"+//la línea en blanco después de los encabezados indica el final de los encabezados HTTP y el comienzo del cuerpo de la respuesta. Es un requisito en el formato HTTP para separar los encabezados del contenido.
			"%s", //el cuerpo(html)
		len(html), html, //valores para el formato, primero el largo del html para el header Content-Length, y luego el html mismo para el cuerpo de la respuesta.
	)

	conn.Write([]byte(response)) //se escribe la respuesta al cliente. "conn.Write" envía los bytes de la respuesta HTTP al cliente que hizo la solicitud. La respuesta incluye el código de estado, los encabezados y el cuerpo HTML que se construyó a partir de la base de datos.
}
