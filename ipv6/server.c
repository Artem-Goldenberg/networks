#include <stdio.h>
#include <ctype.h>
#include <stdlib.h>
#include <errno.h>
#include <assert.h>
#include <unistd.h>
#include <string.h>

#include <sys/socket.h>
#include <arpa/inet.h>
#include <netdb.h>

void printUsage(void) {
    printf("Usage: ./server <host> <port>\n");
    exit(1);
}

void toUpper(char *str) {
    for (int i = 0; i < strlen(str); ++i) {
        str[i] = toupper(str[i]);
    }
}

int main(int argc, const char * argv[]) {
    (void)argc; // do not need it

    // extract hostname and a port
    argv++;
    assert(*argv);
    const char* name = *argv;
    const char* port;
    if (*(++argv)) port = *argv;
    else printUsage();
    
    // get address and other info for the given hostname
    struct addrinfo *info, hint = {
        .ai_family = AF_INET6, // force IPV6
        .ai_socktype = SOCK_STREAM
    };
    int err = getaddrinfo(name, port, &hint, &info);
    assert(!err);

    // create socket
    int s = socket(info->ai_family, info->ai_socktype, 0);
    assert(s >= 0);
    
    int yes = 1;
    if (setsockopt(s, SOL_SOCKET, SO_REUSEADDR, &yes,
                    sizeof(int)) == -1) {
        perror("setsockopt");
        exit(1);
    }
            
    err = bind(s, info->ai_addr, info->ai_addrlen);
    assert(!err);
    
    err = listen(s, 1);
    assert(!err);
    
    char addrString[INET6_ADDRSTRLEN];
    inet_ntop(info->ai_family, &((struct sockaddr_in6*)info->ai_addr)->sin6_addr, addrString, sizeof addrString);
    printf("Listening on %s: %s ...\n", addrString, port);
    
    while(1) {
        struct sockaddr_in6 connAddr = {};
        socklen_t len = sizeof connAddr;
        
        int conn = accept(s, (struct sockaddr*)&connAddr, &len);
        assert(conn > 0);
        
        char connString[INET6_ADDRSTRLEN];
        inet_ntop(info->ai_family, &connAddr.sin6_addr, connString, sizeof connString);
        printf("Accepted connection from: %s\n", connString);
        
        char buffer[1000];
        long n = recv(conn, buffer, sizeof buffer, 0);
        assert(n > 0);
        
        buffer[n] = '\0';
        printf("Recieved %s\n", buffer);
        toUpper(buffer);
        printf("Sending back: %s\n", buffer);
        
        n = send(conn, buffer, strlen(buffer), 0);
        assert(n > 0);
        close(conn);
    }
    
    return 0;
}
