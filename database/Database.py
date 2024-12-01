import aiomysql
import os
import dotenv


async def create_pool():
    
    return await aiomysql.create_pool(
        host=os.environ.get("MYSQL_HOST"),
        port=int(os.environ.get("MYSQL_PORT",3306)),
        user=os.environ.get("MYSQL_USER"),
        password=os.environ.get("MYSQL_PASS"),
        db=os.environ.get("MYSQL_NAME"),
        autocommit=True,
    )

# Consultar la base de datos usando el pool
async def insert_database(pool, query, params):
    async with pool.acquire() as connection:  # Adquiere una conexión del pool
            async with connection.cursor() as cursor:
                try:
                    await cursor.execute(query, params)
                    return await connection.commit()
                except Exception as e:
                    print(f"An error happen in Database.py {e}")