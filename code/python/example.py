# pip3 install neo4j
# python3 example.py

from neo4j import GraphDatabase, basic_auth

cypher_query = '''
MATCH (s:Stream{name:$streamer})<-[:VIP|:MODERATOR]-(user)
RETURN user.name as moderator_or_vip LIMIT 20
'''

with GraphDatabase.driver(
    "neo4j://<HOST>:<BOLTPORT>",
    auth=("<USERNAME>", "<PASSWORD>")
) as driver:
    result = driver.execute_query(
        cypher_query,
        streamer="ludwig",
        database_="neo4j")
    for record in result.records:
        print(record['moderator_or_vip'])
