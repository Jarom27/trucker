from protocols.IGPSProtocol import IGPSProtocol
from protocols.MicodusProtocol import MicodusProtocol

class ProtocolSelector:
    def __init__(self):
        self.protocols = {
            "micodus": MicodusProtocol(None),
        }

    def get_protocol(self, message:bytes) -> IGPSProtocol:
        if message[0] == 0x7e and message[-1] == 0x7e:
            return self.protocols["micodus"]
        else:
            raise ValueError("Protocol is unknown for the system")