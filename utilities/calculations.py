def calculate_xor_checksum(message: bytes):
    checksum = message[0]
    for byte in message[1:]:
        checksum = checksum ^ byte
    return bytes([checksum])