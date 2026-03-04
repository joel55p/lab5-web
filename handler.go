package main
import ( // Este es un servidor HTTP básico en Go que escucha en el puerto 8080 y responde con "Hello World" a cualquier solicitud.
	"bufio"        // "bufio" se utiliza para leer la solicitud del cliente de manera eficiente.
	"database/sql" // "database/sql" se utiliza para interactuar con una base de datos SQL, aunque en este código no se muestra su uso específico, es común en aplicaciones web para manejar datos persistentes.
          // "fmt" se utiliza para formatear la respuesta HTTP que se enviará al cliente.
	"log"          // "log" se utiliza para registrar errores y mensajes informativos en la consola.
	"net"          // "net" se utiliza para crear un servidor TCP que escuche en el puerto 8080 y acepte conexiones entrantes. uno de os paquetes más importantes para la comunicación de red en Go.
	"strings"      // "strings" se utiliza para manipular la cadena de texto de la solicitud HTTP, como extraer el path solicitado.
	"net/url"
	"strconv"
	_ "modernc.org/sqlite" // Este es un import anónimo que se utiliza para registrar el controlador de SQLite con el paquete "database/sql". Esto permite que el programa utilice SQLite como base de datos sin necesidad de importar explícitamente el paquete en el código. El guion bajo (_) indica que el paquete se importa solo por sus efectos secundarios, es decir, para registrar el controlador de la base de datos, sin utilizar directamente ninguna función o tipo del paquete en el código.
)
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


	routeParts := strings.SplitN(path, "?", 2) //El SplitN separa /update?id=3 en dos partes: route = "/update" y queryParams = "id=3". Por eso el switch necesita usar route para comparar con "/update" y no con "/update?id=3", porque el path completo incluye los parámetros de consulta, pero el route es solo la parte del path sin los parámetros. Esto permite que el switch maneje la ruta correctamente sin preocuparse por los parámetros de consulta que pueden variar.
	route := routeParts[0]
	queryParams := ""
	if len(routeParts) > 1 {
		queryParams = routeParts[1]
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
	case method == "GET" && route == "/":
		response = handleIndex( db) //se llama a la función handleIndex para manejar la solicitud GET al path "/". Esta función se encargará de construir la respuesta HTML con la tabla de series desde la base de datos y devolverla al cliente.
	case method == "GET" && route == "/create":
		response = handleCreate() //se llama a la función handleCreate para manejar la solicitud GET al path "/create". Esta función aún no está implementada, pero se espera que maneje la lógica para agregar una nueva serie a la base de datos y devolver una respuesta adecuada al cliente.
	case method == "POST" && route == "/create":
		response = handleCreatePost(body, db) //se llama a la función handleCreatePost para manejar la solicitud POST al path "/create". Esta función aún no está implementada, pero se espera que maneje la lógica para procesar los datos enviados desde el formulario de creación de una nueva serie, agregar esa serie a la base de datos y devolver una respuesta adecuada al cliente.
	
	case method == "POST" && route == "/update":
		params, _ := url.ParseQuery(queryParams) // convierte "id=x" en un mapa {id: "x"}
		id := params.Get("id")
		response = handleUpdate(db, id) //da el response

	case method == "GET" && route == "/script.js":
    	response = handleScript()	//para que lea el js 
		case method == "GET" && route == "/faviconm.png":
		handleFavicon(conn)
		return // salir de handle() directamente, no usar conn.Write al final

	case method == "POST" && route == "/downdate":
		params, _ := url.ParseQuery(queryParams) // convierte "id=x" en un mapa {id: "x"}
		id := params.Get("id")
		response = handleRestdate(db, id) //da el response
	default:
		response = "HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\n\r\n<h1>404 Not Found</h1>" //si la solicitud no coincide con ninguna de las rutas definidas, se devuelve una respuesta HTTP con el código de estado 404 Not Found y un mensaje HTML indicando que la página no fue encontrada。
	}

	conn.Write([]byte(response))  // handle() escribe
}
