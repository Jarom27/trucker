from abc import ABC, abstractmethod
from asyncio import Transport 
class IGPSProtocol(ABC):
    @abstractmethod
    def process_message(self, request):
        """Abstract method for GPS protocols"""
        pass
    @abstractmethod
    def execute_message(self,request,transport: Transport):
        pass
    
