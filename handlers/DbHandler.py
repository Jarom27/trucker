from handlers.Handler import Handler
from database.Database import insert_database
class DbHandler(Handler):
    def __init__(self, pool,next_handler=None):
        self.pool = pool
        super().__init__(next_handler)

    async def handle(self,request):
        Device_id = "".join([hex(b) for b in request["device_id"]])
        Device_id = Device_id.replace('0x','')
        location_report = request["location"]
        latitude, longitude, altitude = location_report.values()
        params = (latitude,longitude,altitude,Device_id)
        print(params)
        query = "INSERT INTO location_logs(latitude,longitude,altitude,Device_id) VALUES(%s,%s,%s,%s);"
        await insert_database(self.pool,query,params)

        return request
