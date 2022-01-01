syntax = "proto3";

//option go_package = "cloudedgenetwork/proxyserver/pb";

package pb;

service ProxyHttpService {
    rpc ProxyCloud2Edge(stream Response) returns (stream Request) {}
    rpc ProxyEdge2Cloud(stream Request) returns (stream Response) {}
}

message Request {
    string id = 1;
    bool hello = 2;
    string target_addr = 3;
    string target_node = 4;
    HTTPHeader header = 5;
    string http_method = 6;
    string url = 7;
    bytes body = 8;

    // edge --> cloud
    string nodeid = 9;
}

message Response {
    string id = 1;
    bool hello = 2;
    string nodeid = 3;
    HTTPHeader header = 4;
    bytes body = 5;

    // edge --> cloud
    int64 statuscode = 10;
}

message HTTPHeader {
    repeated HTTPHeaderValue item = 1;
}

message HTTPHeaderValue {
    repeated string value = 1;
}

