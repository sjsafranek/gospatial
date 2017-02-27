#!/usr/bin/env python

import json
import time
import socket

TCP_IP = '127.0.0.1'
TCP_PORT = 3333
BUFFER_SIZE = 1024

def connect():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect((TCP_IP, TCP_PORT))
    return s

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
                data = data.decode()
                total_data.append(data)
                begin=time.time()
            else:
                time.sleep(0.1)
        except:
            pass
    return ''.join(total_data)


def newApikey():
    s = connect()
    MESSAGE = bytes( '{"method":"create_apikey"}\n', 'UTF-8' )
    s.send(MESSAGE)
    resp = recv_timeout(s,.15)
    data = json.loads(resp)
    if 'ok' is not data['status']:
        s.close()
        ValueError('Error creating apikey: ' + resp)
    return data['data']['apikey']
    s.close()



















