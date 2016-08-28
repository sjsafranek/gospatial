#!/usr/bin/env python
import os
import re
import time
import socket

# http://code.activestate.com/recipes/408859-socketrecv-three-ways-to-turn-it-into-recvall/

logFile = os.path.join("log", "db.log")

TCP_IP = '127.0.0.1'
TCP_PORT = 3333
BUFFER_SIZE = 1024
MESSAGE = '{"method":"authenticate", "authkey": "TLAJhDe925es"}\n'

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((TCP_IP, TCP_PORT))
s.send(MESSAGE)


#assume a socket disconnect (data returned is empty string) means  all data was #done being sent.
def recv_basic(the_socket):
    total_data=[]
    while True:
        data = the_socket.recv(8192)
        if not data: break
        total_data.append(data)
    return ''.join(total_data)
    
def recv_timeout(the_socket,timeout=2):
    the_socket.setblocking(0)
    total_data=[];data='';begin=time.time()
    while 1:
        #if you got some data, then break after wait sec
        if total_data and time.time()-begin>timeout:
            break
        #if you got no data at all, wait a little longer
        elif time.time()-begin>timeout*2:
            break
        try:
            data=the_socket.recv(8192)
            if data:
                total_data.append(data)
                begin=time.time()
            else:
                time.sleep(0.1)
        except:
            pass
    return ''.join(total_data)

def recv_end(the_socket):
    total_data=[];data=''
    while True:
            data=the_socket.recv(8192)
            if End in data:
                total_data.append(data[:data.find(End)])
                break
            total_data.append(data)
            if len(total_data)>1:
                #check if end_of_data was split
                last_pair=total_data[-2]+total_data[-1]
                if End in last_pair:
                    total_data[-2]=last_pair[:last_pair.find(End)]
                    total_data.pop()
                    break
    return ''.join(total_data)

def recv_size(the_socket):
    #data length is packed into 4 bytes
    total_len=0;total_data=[];size=sys.maxint
    size_data=sock_data='';recv_size=8192
    while total_len<size:
        sock_data=the_socket.recv(recv_size)
        if not total_data:
            if len(sock_data)>4:
                size_data+=sock_data
                size=struct.unpack('>i', size_data[:4])[0]
                recv_size=size
                if recv_size>524288:recv_size=524288
                total_data.append(size_data[4:])
            else:
                size_data+=sock_data
        else:
            total_data.append(sock_data)
        total_len=sum([len(i) for i in total_data ])
    return ''.join(total_data)



with open(logFile, "r") as fileHandler:
    queries = []
    print("collecting queries")
    for line in fileHandler:
        line = line.replace("\n", "")
        matches = re.findall(r"WRITE.*db.go:\d{1,3}: ", line)
        for match in matches:
            line = line.replace(match, "")
        queries.append(line)
        print(line)
    print("sending queries")
    for query in queries:
        s.send(query + "\n")
        print(recv_timeout(s,.15))

s.close()



















