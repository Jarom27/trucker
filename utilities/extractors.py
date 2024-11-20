def identify_message_type(message: bytes):
    message_type = ""
    if message[1:3] == b"\x01\x00":
        message_type = "0100"
    elif message[1:3] == b"\x00\x03":
        message_type = "0003"
    else:
        message_type = "0000"
    return message_type

def identify_device_id(message: bytes):
    device_id = message[5:11]
    return device_id