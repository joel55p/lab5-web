package main

import (
    "fmt"
    "log"
    "os" // "os" se utiliza para leer el archivo "script.js" desde el sistema de archivos. Esto es necesario para servir el contenido del archivo JavaScript al cliente cuando se solicita la ruta "/script.js".
)

func handleScript() string { // sirve el archivo script.js al browser cuando lo solicita
    content, err := os.ReadFile("script.js") // Lee el contenido del archivo "script.js" y lo almacena en la variable "content". Si ocurre un error al leer el archivo, se maneja el error y se devuelve una respuesta HTTP con el código de estado 404 Not Found.
    if err != nil {
        log.Print("Error leyendo script.js:", err)
        return "HTTP/1.1 404 Not Found\r\n\r\n"
    }
    return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/javascript\r\n\r\n%s", content)
}