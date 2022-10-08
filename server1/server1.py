import os
from http.server import BaseHTTPRequestHandler, HTTPServer
from socket import gethostname
from socket import gethostbyname
from urllib.request import urlopen

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        hostname = gethostname()
        local_ip = gethostbyname(hostname)
        message1 = "Hello from " + self.client_address[0] + ":" + str(self.client_address[1]) + \
                   " to " + str(local_ip) + ":" + str(server.server_port)
        url = os.environ['SERVER_2']
        message2 = urlopen(url).read()
        finalmessage = message1 + "\n" + message2.decode('UTF-8')
        self.wfile.write(bytes(finalmessage,'UTF-8'))


with HTTPServer(('', 8081), Handler) as server:
    print("start")
    server.serve_forever()
