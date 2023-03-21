# MirCat

MirCat is a network game protocol analysis tool. Its front-end is developed using the React framework, while its back-end is written in Golang. Currently, the Golang part provides the following API functions:

1. GetConfig
2. SetConfig
3. ClientTcpOpen
4. ClientTcpSend
5. ClientTcpClose
6. ClientTcpCloseAll
7. ServerTcpStart
8. ServerTcpStop
9. ServerSendMessage
10. ServerBroadcastMessage
11. TransferTcpStart
12. TransferTcpStop
13. TransferSendToClient
14. TransferSendToServer
15. TransferBroadcastToClient
16. TransferBroadcastToServer

The events that have already been implemented are:

- client-tcp-error
- client-tcp-info
- client-tcp-data
- config-saved (deprecated)
- server-tcp-error
- server-tcp-info
- server-tcp-data
- transfer-tcp-error
- transfer-tcp-info
- transfer-src-data
- transfer-dst-data

For detailed back-end documentation, godoc can be started on the local machine and accessed through the following link:

http://localhost:6060/pkg/mir-cat/pkg/mircat

New features can be added by providing PRs to merge the code.
