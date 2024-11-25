from abc import ABC, abstractmethod
class IGPSProtocol(ABC):
    @abstractmethod
    def process_message(self, request):
        """Abstract method for GPS protocols"""
        pass
    
