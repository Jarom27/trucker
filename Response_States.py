from enum import Enum

class ResponseStates(Enum):
    CLOSE_CONNECTION = -1
    SENT_RESPONSE = 2
    AUTHENTICATION = 3