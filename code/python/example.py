# pip3 install neo4j
# python3 example.py

from neo4j import GraphDatabase, basic_auth

driver = GraphDatabase.driver(
  "neo4j://<HOST>:<BOLTPORT>",
  auth=basic_auth("<USERNAME>", "<PASSWORD>"))

cypher_query = '''
MATCH (s:Stream{name:$streamer})<-[:VIP|:MODERATOR]-(user)
RETURN user.name as moderator_or_vip LIMIT 20
'''

with driver.session(database="neo4j") as session:
  results = session.read_transaction(
    lambda tx: tx.run(cypher_query,
                      streamer="ludwig").data())
  for record in results:
    print(record['moderator_or_vip'])

driver.close()
