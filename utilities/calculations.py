def calculate_checksum(message: bytes):
    checksum = message[1]
    for byte in message[2:]:
        checksum = checksum ^ byte
    return bytes([checksum])