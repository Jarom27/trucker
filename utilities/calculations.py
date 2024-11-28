def calculate_xor_checksum(message: bytes):
    checksum = message[0]
    for byte in message[1:]:
        checksum = checksum ^ byte
    return bytes([checksum])

def convert_decimal_to_grades(decimal_number):
    grades = int(decimal_number)
    minutes = int((decimal_number - grades) * 60)
    seconds = ((decimal_number - grades) * 60 - minutes) * 60
    return f"{grades},{minutes},{seconds}"