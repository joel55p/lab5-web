package main

import (
    "log"
    "net"
    "os"
)

func handleFavicon(conn net.Conn) { // sirve el archivo favicon.png escribiendo directo a la conexion porque es binario
    content, err := os.ReadFile("faviconm.png")
    if err != nil {
        log.Print("Error leyendo favicon.png:", err)
        conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
        return
    }

    header := "HTTP/1.1 200 OK\r\nContent-Type: image/png\r\n\r\n"
    conn.Write([]byte(header))
    conn.Write(content) // escribe los bytes binarios directamente sin pasar por Sprintf
}