import struct,socket

# USAGE: byond_topic(data, server(Tuple, (addr,port)))

TYPE_NULL = 0x00
TYPE_FLOAT = 0x2a
TYPE_STRING = 0x06

def byond_topic(data,server):
    # BYOND packet forging data.
    data = bytes(data,'latin-1')
    pref = b'\x00\x83\x00'
    leng = struct.pack('b',len(data)+6)
    pad = b'\x00'*5
    suff = b'\x00'
    # b_msg will be our final message to send.
    b_msg = pref+leng+pad+data+suff

    s = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
    s.settimeout(3)
    s.connect(server)
    s.send(b_msg)
    # Grab the entire response at once, with regard for maximum BYOND response length.
    ret = s.recv(65541)
    s.close()
    return ret

def ping(addr):
    data = byond_topic("?ping",addr)
    if int(data[4]) != TYPE_FLOAT:
        print("Wrong data type returned")
        return -1
    return int(struct.unpack('<f',data[5:])[0])

print(str(ping(("play.cm-ss13.com",1400)))+" players online")
