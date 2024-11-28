from .IGPSProtocol import IGPSProtocol
from utilities.calculations import calculate_xor_checksum,convert_decimal_to_grades
from asyncio import Transport
from Response_States import ResponseStates

class MicodusProtocol(IGPSProtocol):
    def __init__(self, send_response_callback):
        self.send_response = send_response_callback
        self.serial_number = 0
        self.device_id = None
        
    def identify_message_type(self,data:bytes):
        message = ''.join([hex(byte) for byte in data[1:3]])
        message_type = message.replace("0x","")
        return message_type

    def handle_identification(self):
        print("Indentification process")
        response_message = self.build_message(b"\x81\x00", message_length=b"\x02", message_result=b"\x00\x00")
        print(f"Message data: {response_message.hex()}")
        #self.send_response(response_message)
        return {"status" : ResponseStates.SENT_RESPONSE, "message" : response_message}

    def handle_authenticacion(self):
        print("Autentiffication process")
        return self.request_position()

    def request_position(self):
        print("Requesting location")
        response_message = self.build_message(message_type=b"\x82\x00")
        #self.send_response(response_message)
        return {"status" : ResponseStates.SENT_RESPONSE, "message" : response_message}
    
    def parse_data(self,message:bytes) -> dict:
        hexa_lat = int.from_bytes(message[23:27]) / 10**6
        latitude = convert_decimal_to_grades(hexa_lat)
        hexa_long = int.from_bytes(message[27:31])
        longitude = convert_decimal_to_grades(hexa_long)
        altitude = int.from_bytes(message[31:33],byteorder="big")

        location_report = {
                "latitude" : latitude, 
                "longitude" : longitude, 
                "altitude" : altitude
            }
        return location_report
        

    def build_message(self, message_type, message_body = None, message_length=b"\x00\x00", message_result = None) -> bytes:
        """Build a Micodus Message"""
        response_message_type = message_type
        response_message_length = message_length
        response_message_serial = bytes([self.serial_number])

        response = b"\x7e"
        
        #Header
        response += response_message_type + response_message_length + self.device_id + response_message_serial

        #Content
        if message_body and type(message_body) == bytes:
            response += message_body

        if message_result and type(message_result) == bytes:
            response += message_result

        response_checksum = calculate_xor_checksum(response[1:])
        response += response_checksum

        response += b"\x7e"
        return response

    def process_message(self, message: bytes):
        print("Processing data using Micodus Protocol")

        message_type = self.identify_message_type(message)
        self.device_id = message[5:11]
        response = 0

        if message_type == "03":
            return {"status" : ResponseStates.CLOSE_CONNECTION}
        elif message_type == "10":
            response = self.handle_identification()
        elif message_type == "12":
            response = self.handle_authenticacion()
        elif message_type == "20":
            print(f"Location: {message.hex()}")
            location_report = self.parse_data(message)
            response = {"status" : ResponseStates.AUTHENTICATION,"location" :  location_report}
        elif message_type == "21":
            pass

        return response
    
   