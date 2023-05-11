#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <assert.h>
#include <unistd.h>
#include <string.h>

#include <sys/socket.h>
#include <arpa/inet.h>
#include <netdb.h>

void printUsage(void) {
    printf("Usage: ./client <host> <port> <message>\n");
    exit(1);
}

int main(int argc, const char * argv[]) {
    (void)argc; // do not need it

    // extract hostname and a port
    argv++;
    assert(*argv);
    const char* name = *argv;
    const char* port;
    const char* msg;
    if (*(++argv)) port = *argv;
    else printUsage();
    if (*(++argv)) msg = *argv;
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
    
    char dstString[100];
    getnameinfo(info->ai_addr, info->ai_addrlen, dstString, sizeof dstString, NULL, 0, 0);
    char addrString[INET6_ADDRSTRLEN];
    inet_ntop(info->ai_family, &((struct sockaddr_in6*)info->ai_addr)->sin6_addr, addrString, sizeof addrString);
    printf("Sending %s to %s: %s ...\n", msg, dstString, addrString);
    
    err = connect(s, info->ai_addr, info->ai_addrlen);
    assert(!err);
    
    long n = send(s, msg, strlen(msg), 0);
    assert(n > 0);
    
    // recieve a respone
    char buffer[200];
    n = recv(s, buffer, sizeof buffer, 0);
    assert(n > 0);
    
    buffer[n] = '\0';
    printf("got back: %s\n", buffer);
    
    freeaddrinfo(info);
    close(s);
    
    return 0;
}
