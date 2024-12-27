import socket

if __name__ == '__main__':
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    server_socket.bind(('127.0.0.1', 8080))

    server_socket.listen(5)

    while True:
        client_socket, client_address = server_socket.accept()
        print(f"接受来自 {client_address} 的连接")

        try:
            data = client_socket.recv(1024)
            if data:
                print(f"接收到的数据: {data.decode('utf-8', errors='ignore')}")
            else:
                print("没有接收到数据")
        except Exception as e:
            print(f"发生错误: {e}")
        finally:
            client_socket.close()
