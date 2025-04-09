import socket

server_ip = '127.0.0.1'
server_port = 9000

client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

command = 'instruction=set|store=user_data'

try:
    client_socket.connect((server_ip, server_port))
    
    message = command + "\n"
    
    client_socket.sendall(message.encode())

    sock_file = client_socket.makefile('r')
    response_line = sock_file.readline()

    print("Response:", response_line.strip())
    
finally:
    client_socket.close()
