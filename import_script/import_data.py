import json
import logging
import csv

from pymongo import MongoClient
import psycopg2
import psycopg2.extras

logging.basicConfig(level="DEBUG")
logger = logging.getLogger(__name__)

class ConfigManager(object):

    def __init__(self, config_file="config.json"):
        with open(config_file, "rb") as fh:
            self.config = json.load(fh)
    
    def get_mongo_client(self):
        conn_str = self.config["db_connection"]["mongodb"]
        return MongoClient(conn_str)

    def get_pg_connection(self):
        conn_str = self.config["db_connection"]["postgres"]
        return psycopg2.connect(conn_str)
        

class Util(object):
    @staticmethod
    def parse_csv(input_file: str):
        header = None
        data = []
        with open(input_file, 'r') as fh:
            for row in csv.reader(fh):
                if header is None:
                    header = row
                    continue
                data.append(row)
        return header, data

    @staticmethod
    def convert_fields(header: list, data: list, converters: dict) -> list:
        for row in data:
            for key, converter in converters.items():
                idx = header.index(key)
                row[idx] = converter(row[idx])
        return data


class ImportMongo(object):
    def __init__(self, config_manager):
        mongo_client = config_manager.get_mongo_client()
        self.db = mongo_client.get_default_database()

    def import_data(self, collection_name, input_file, converters={}):
        # parse csv
        header, data = Util.parse_csv(input_file)
        # convert data type/format
        data = Util.convert_fields(header, data, converters)

        data = [dict(zip(header, d)) for d in data]
        # insert into db
        mongo_collection = self.db[collection_name]
        mongo_collection.insert_many(data)
        
            
class ImportPostgres(object):
    def __init__(self, config_manager):
        self.cfg_util = config_manager

    def connect(self):
        return self.cfg_util.get_pg_connection()

    def cursor(self, conn):
        return conn.cursor(cursor_factory=psycopg2.extras.DictCursor)

    def import_data(self, table_name, input_file, converters={}):
        # parse csv
        header, data = Util.parse_csv(input_file)
        # convert data type/format
        data = Util.convert_fields(header, data, converters)
        # insert into db
        sql = "insert into {table_name} ({fields}) values ({placeholders})".format(
            table_name=table_name, fields=", ".join(header),
            placeholders=", ".join(['%s'] * len(header))
        )
        with self.connect() as conn:
            try:
                with self.cursor(conn) as cur:
                    for row in data:
                        cur.execute(sql, row)
                conn.commit()                
            except:
                conn.rollback()
                logger.exception("Inserting 1 row failed.")


if __name__ == '__main__':
    cfg = ConfigManager('../config.json')
    ImportMongo(cfg).import_data(
        "customer_companies", 
        'test_data/Test task - Mongo - customer_companies.csv',
        converters={"company_id": int}
    )
    ImportMongo(cfg).import_data(
        'customers', 'test_data/Test task - Mongo - customers.csv',
        converters={
            "company_id": int,
            "credit_cards": json.loads
        }
    )
    ImportPostgres(cfg).import_data("orders", 'test_data/Test task - Postgres - orders.csv')
    ImportPostgres(cfg).import_data(
        "order_items", 'test_data/Test task - Postgres - order_items.csv',
        converters={
            "price_per_unit": lambda s: s or None,
            "quantity": lambda s: s or None
        }
    )
    ImportPostgres(cfg).import_data(
        "deliveries", 'test_data/Test task - Postgres - deliveries.csv'
    )
