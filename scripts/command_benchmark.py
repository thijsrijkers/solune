import socket
import time

server_ip = '127.0.0.1'
server_port = 8743
num_requests = 100
command = ''

latencies = []

def benchmark():
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect((server_ip, server_port))
    sock_file = client_socket.makefile('r')

    try:
        for i in range(num_requests):
            start_time = time.time()

            message = command + "\n"
            client_socket.sendall(message.encode())

            response = sock_file.readline().strip()

            end_time = time.time()
            latency = end_time - start_time
            latencies.append(latency)

            print(f"[{i+1}] Response: {response}, Latency: {latency:.6f} sec")

    finally:
        client_socket.close()

    if latencies:
        print("\n=== Benchmark Results ===")
        print(f"Total requests: {num_requests}")
        print(f"Average latency: {sum(latencies)/len(latencies):.6f} sec")
        print(f"Max latency: {max(latencies):.6f} sec")
        print(f"Min latency: {min(latencies):.6f} sec")
        print(f"Throughput: {num_requests / sum(latencies):.2f} requests/sec")

if __name__ == "__main__":
    benchmark()
