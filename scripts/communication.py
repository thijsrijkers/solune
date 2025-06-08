import socket

server_ip = '127.0.0.1'
server_port = 8743

client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

command = ''

try:
    client_socket.connect((server_ip, server_port))
    message = command + "\n"
    client_socket.sendall(message.encode())

    sock_file = client_socket.makefile('r')

    for line in sock_file:
        print(line.strip())

finally:
    client_socket.close()
