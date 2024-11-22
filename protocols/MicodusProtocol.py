from .IGPSProtocol import IGPSProtocol
from asyncio import Transport

class MicodusProtocol(IGPSProtocol):
    def __init__(self):
        self.message_type = ""
        self.device_id = ""
        self.body_length = 0
        self.serial_number = 0
    
    def identify_message_type(self,data:bytes):
        message = ''.join([hex(byte) for byte in data[1:3]])
        message_type = message.replace("0x","")
        return message_type
        
    def identify_device_id(self,data:bytes):
        message = ''.join([hex(byte) for byte in data[5:11]])
        device_id = message.replace("0x","")
        return device_id

    def get_body_length(self,data:bytes):
        return int.from_bytes(data[3:5],byteorder="big")

    def process_message(self,request):
        print("Processing data using Micodus Protocol")
        data = request["data"]
        self.message_type = self.identify_message_type(data)
        self.device_id = self.identify_device_id(data)
        self.body_length = self.get_body_length(data)
        self.serial_number = (self.serial_number + 1) % 0xFFFF
        request["message_type"] = self.message_type
        request["device_id"] = self.device_id
        request["body_length"] = self.body_length
        return {"status": "200"}
    
    def execute_message(self,request,transport:Transport):
        print(request["message_type"])
        if request["message_type"] == "03":
            print("Closing connection")
            transport.close()
        elif request["message_type"] == "10":
            print(request["device_id"])
            reply_message = b"\x7e\x81\x00\x00\x0a"+request["data"][5:11]+b""+self.serial_number.to_bytes(2,byteorder="big")+b"\x00\x00"
            checksum = self.calculate_checksum(reply_message)
            reply_message += checksum
            reply_message += b'\x7e'
            print("Repply: "+reply_message.hex())
            transport.write(reply_message)
        elif request["message_type"] == "12":
            print(request["device_id"])
            reply_message = b"\x7e\x82\x01\x00\x0a"+request["data"][5:11]+b""+self.serial_number.to_bytes(2,byteorder="big")
            checksum = self.calculate_checksum(reply_message)
            reply_message += checksum
            reply_message += b'\x7e'
            print("Repply: "+reply_message.hex())
            transport.write(reply_message)
        elif request["message_type"] == "0200":
            print(request["device_id"])
            
    def calculate_checksum(self,message: bytes):
        checksum = message[1]
        for byte in message[2:]:
            checksum = checksum ^ byte
        return bytes([checksum])
            
        
        
