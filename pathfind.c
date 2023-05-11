#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <assert.h>
#include <unistd.h>

#include <sys/time.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <netinet/ip_icmp.h>

static const short id = 1; // non zero id for icmp header
static const int timeout = 1; // timeout in seconds
static const char* dstPort = "80";
static const size_t maxIpHeaderSize = 60;

static const int defaultNRequests = 3;

// custom icmp message format
struct icmp_msg {
    // header
    uint8_t type, code;
    uint16_t checksum;
    uint16_t id;
    uint16_t seq;
    // body to fill to 64 bytes
    char bytes[54];
};

static const size_t recieveBufferSize = maxIpHeaderSize + sizeof(struct icmp_msg);

uint16_t checksum(const void *bytes, size_t size) {
    uint32_t sum = 0;
    uint16_t* p = (uint16_t*)bytes;
    
    // add everything moving in windows of 2 bytes
    // size is even ==> don't care about -1, size is odd ==> will miss last byte
    for (; (char*)p < (char*)bytes + size - 1; ++p) sum += *p;
    // check for one byte left over
    if (size % 2 == 1) sum += *(uint8_t*)p;
    
    // kinda doing something ...
    sum = (sum & 0xffff) + (sum >> 16);
    sum += (sum >> 16);
    
    return ~sum;
}

const char* icmpErrorString(int, int);

int main(int argc, const char * argv[]) {
    (void)argc; // do not need it

    // extract hostname and number of requests
    argv++;
    assert(*argv);
    const char* name = *argv;
    long nRequests = defaultNRequests;
    if (*(++argv)) nRequests = strtol(*argv, NULL, 10);
    assert(errno != EINVAL && errno != ERANGE);
    
    // get address and other info for the given hostname
    struct addrinfo *info, hint = {.ai_family = AF_INET, .ai_protocol = IPPROTO_ICMP};
    int err = getaddrinfo(name, dstPort, &hint, &info);
    assert(!err);

    // create socket
    int s = socket(info->ai_family, SOCK_DGRAM, info->ai_protocol);
    assert(s >= 0);
    
    // set a timeout 1.5s
    struct timeval t = {.tv_sec = timeout, .tv_usec = 500000};
    err = setsockopt(s, SOL_SOCKET, SO_RCVTIMEO, &t, sizeof t);
    assert(!err);
    
    int quit = 0;
    struct timeval tstart, tend;
    
    char dstString[100];
    getnameinfo(info->ai_addr, info->ai_addrlen, dstString, sizeof dstString, NULL, 0, 0);
    const char* addrString = inet_ntoa( ((struct sockaddr_in*)info->ai_addr)->sin_addr );
    printf("Tracing %s: %s ...\n", dstString, addrString);
    
    for (uint16_t i = 0; !quit; ++i) {
        // make echo request icmp message with time mark
        struct icmp_msg data = {.type = ICMP_ECHO, .id = id, .seq = htons(i)};
        data.checksum = checksum(&data, sizeof data);
        
        // set increasingly time to live starting from 1
        int ttl = i + 1;
        err = setsockopt(s, IPPROTO_IP, IP_TTL, &ttl, sizeof ttl);
        assert(!err);
        
        struct sockaddr_in respAddr = {};
        socklen_t len = sizeof respAddr;
        printf("%2d. ", ttl);
        int addrPrinted = 0;
        // send the same packet `nRequests` times
        for (int j = 0; j < nRequests; ++j) {
            // send to the address stored in `info`
            gettimeofday(&tstart, 0);
            long n = sendto(s, &data, sizeof data, 0, info->ai_addr, info->ai_addrlen);
            assert(n > 0);
            
            // recieve a respone, measure time
            char buffer[recieveBufferSize];
            n = recvfrom(s, &buffer, sizeof buffer, 0, (struct sockaddr*)&respAddr, &len);
            gettimeofday(&tend, 0);
            if (n < 0 && errno == EWOULDBLOCK) {
                printf("* ");
                continue;
            }
            
            // discard ip header, check the type, handle icmp errors
            const size_t ipHeaderSize = ((struct ip*)buffer)->ip_hl * 4;
            struct icmp_msg* reply = (struct icmp_msg*)(buffer + ipHeaderSize);
            quit = (reply->type == ICMP_ECHOREPLY);
            if (!quit && reply->type != ICMP_TIMXCEED) {
                assert(ICMP_ERRORTYPE(reply->type));
                printf("\nIcmp error: %s\n", icmpErrorString(reply->type, reply->code));
                goto cleanup;
            }
            
            // calculate rtt
            double rtt = (tend.tv_sec - tstart.tv_sec) * 1e3 + (tend.tv_usec - tstart.tv_usec) * 1e-3;
            
            // convert ip to name
            char hostName[100];
            getnameinfo((struct sockaddr*)&respAddr, sizeof respAddr, hostName, sizeof hostName, NULL, 0, 0);
            
            // get responce ip address as string
            char addrString[INET_ADDRSTRLEN];
            inet_ntop(AF_INET, &respAddr.sin_addr, addrString, sizeof addrString);

            // print stuff
            if (!addrPrinted) {
                printf("%s (%s) -- ", hostName, addrString);
                addrPrinted = 1;
            }
            printf("%g ms ", rtt);
        }
        printf("\n");
    }
cleanup:
    freeaddrinfo(info);
    close(s);
    
    return 0;
}

const char* icmpErrorString(int type, int code) {
    switch (type) {
    case ICMP_UNREACH:
        switch (code) {
        case ICMP_UNREACH_NET: return "Destination network unreachable";
        case ICMP_UNREACH_HOST: return "Destination host unreachable";
        case ICMP_UNREACH_PROTOCOL: return "Destination protocol unreachable";
        case ICMP_UNREACH_PORT: return "Destination port unreachable";
        case ICMP_UNREACH_NEEDFRAG: return "Fragmentation required, and DF flag set";
        case ICMP_UNREACH_SRCFAIL: return "Source route failed";
        case ICMP_UNREACH_NET_UNKNOWN: return "Destination network unknown";
        case ICMP_UNREACH_HOST_UNKNOWN: return "Destination host unknown";
        case ICMP_UNREACH_ISOLATED: return "Source host isolated";
        case ICMP_UNREACH_NET_PROHIB: return "Network administratively prohibited";
        case ICMP_UNREACH_HOST_PROHIB: return "Host administratively prohibited";
        case ICMP_UNREACH_TOSNET: return "Network unreachable for ToS";
        case ICMP_UNREACH_TOSHOST: return "Host unreachable for ToS";
        case ICMP_UNREACH_FILTER_PROHIB: return "Communication administratively prohibited";
        case ICMP_UNREACH_HOST_PRECEDENCE: return "Host Precedence Violation";
        case ICMP_UNREACH_PRECEDENCE_CUTOFF: return "Precedence cutoff in effect";
        }
    case ICMP_REDIRECT:
        switch (code) {
        case ICMP_REDIRECT_NET: return "Redirect Datagram for the Network";
        case ICMP_REDIRECT_HOST: return "Redirect Datagram for the Host";
        case ICMP_REDIRECT_TOSNET: return "Redirect Datagram for the ToS & network";
        case ICMP_REDIRECT_TOSHOST: return "Redirect Datagram for the ToS & host";
        }
    case ICMP_TIMXCEED:
        switch (code) {
        case ICMP_TIMXCEED_INTRANS: return "TTL expired in transit";
        case ICMP_TIMXCEED_REASS: return "Fragment reassembly time exceeded";
        }
    case ICMP_PARAMPROB:
        switch (code) {
        case ICMP_PARAMPROB_ERRATPTR: return "Pointer indicates the error";
        case ICMP_PARAMPROB_OPTABSENT: return "Missing a required option";
        case ICMP_PARAMPROB_LENGTH: return "Bad length";
        }
    }
    return "Unknown error";
}
