# Go TCP客户端

这是一个用Go编写的TCP客户端代码。代码定义了一个名为`TcpClient`的结构体，以及相关的方法。

## 结构体

`TcpClient`结构体包含以下字段：

- `address`：连接的地址
- `conn`：实际的网络连接对象
- `sendChan`：发送数据的通道
- `recvChan`：接收数据的通道
- `isShutdown`：是否关闭
- `index`：客户端索引
- `app`：应用程序实例

## 方法

以下是`TcpClient`结构体的方法：

### NewTcpClient

用于创建新的TcpClient实例并建立连接。

### startSending

用于从发送通道读取数据并将其写入连接。

### startReceiving

用于从连接读取数据并将其写入接收通道。

### Send

向TcpClient发送数据。

### Shutdown

关闭TcpClient实例并释放资源。

### reconnect

尝试重新连接TcpClient。

## 功能

这个TCP客户端可以用来与TCP服务器通信，包括发送和接收数据。在发送数据时，它会将数据添加到发送通道，然后由`startSending`方法将其发送到服务器。接收数据时，`startReceiving`方法会将数据从服务器读取并将其添加到接收通道。如果发生错误，客户端会尝试重新连接到服务器。
