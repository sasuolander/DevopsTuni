from http.server import BaseHTTPRequestHandler, HTTPServer
from socket import gethostname
from socket import gethostbyname

class handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text')
        self.end_headers()
        hostname = gethostname()
        local_ip = gethostbyname(hostname)
        message = "Hello from " + self.client_address[0] + ":" + str(self.client_address[1]) + \
                  " to " + str(local_ip) + ":" + str(server.server_port)
        self.wfile.write(bytes(message, "utf8"))

with HTTPServer(('', 8080), handler) as server:
    print("start")
    server.serve_forever()
