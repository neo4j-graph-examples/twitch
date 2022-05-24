// npm install --save neo4j-driver
// node example.js
const neo4j = require('neo4j-driver');
const driver = neo4j.driver('neo4j://<HOST>:<BOLTPORT>',
                  neo4j.auth.basic('<USERNAME>', '<PASSWORD>'), 
                  {});

const query =
  `
  MATCH (s:Stream{name:$streamer})<-[:VIP|:MODERATOR]-(user)
  RETURN user.name as moderator_or_vip LIMIT 20
  `;

const params = {"streamer": "ludwig"};

const session = driver.session({database:"neo4j"});

session.run(query, params)
  .then((result) => {
    result.records.forEach((record) => {
        console.log(record.get('moderator_or_vip'));
    });
    session.close();
    driver.close();
  })
  .catch((error) => {
    console.error(error);
  });
