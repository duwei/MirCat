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

The events that have already been implemented are:

- client-tcp-error
- client-tcp-info
- client-tcp-data
- config-saved (deprecated)
- server-tcp-error
- server-tcp-info
- server-tcp-data

For detailed back-end documentation, godoc can be started on the local machine and accessed through the following link:

http://localhost:6060/pkg/mir-cat/pkg/mircat

New features can be added by providing PRs to merge the code.
