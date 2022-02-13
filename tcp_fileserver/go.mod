module tcp_fileserver

go 1.17

// To import package created with: go mod init <name> 
require (
    tcp_fileserver/utils v1.0.0
)

replace (
    tcp_fileserver/utils v1.0.0 => ./utils
)