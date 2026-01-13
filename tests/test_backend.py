import requests
import psycopg2
import pytest
from collections import namedtuple

SERVER_1 = "http://localhost:8080"
SERVER_2 = "http://localhost:8081"
SERVERS = [SERVER_1, SERVER_2]

def test_health_check():
    for s in SERVERS:
        resp = requests.get(f"{s}/health")
        assert resp.status_code == 200

USER_1 = "alice"
USER_2 = "bob"
USERS = [USER_1, USER_2]

DB_1 = "db1"
DB_2 = "db2"
DBS = [DB_1, DB_2]

@pytest.fixture
def db_inspector():
    def _query(db_name, sql):
        conn = psycopg2.connect(f"postgres://borg:borg@localhost:5432/{db_name}")
        with conn.cursor() as cur:
            cur.execute(sql)
            rows = cur.fetchall()
            # Get column names from the cursor description
            colnames = [desc[0] for desc in cur.description]
        conn.close()
        
        # Define a NamedTuple type 'Row' with the database column names
        Row = namedtuple('Row', colnames)
        
        # Convert every tuple in 'rows' into a 'Row' object
        return [Row(*r) for r in rows]
    
    return _query

def test_register(db_inspector):
    def register(server, username):
        payload = {
            "username": username,
            "password": "password"
        }
        response = requests.post(f"{server}/auth/register", json=payload)
        return response

    for i in range(2):
        usrs = db_inspector(DBS[i], "SELECT * FROM users")
        assert len(usrs) == 0
        resp1 = register(SERVERS[i], USERS[i])
        assert resp1.status_code == 201
        usrs = db_inspector(DBS[i], "SELECT * FROM users")
        assert len(usrs) == 1

@pytest.fixture(scope="session")
def user_tokens():
    tokens = {}
    for i, name in enumerate(USERS):
        resp = requests.post(
            f"{SERVERS[i]}/auth/login", 
            json={"username": name, "password": "password"}
        )
        assert resp.status_code == 200
        tokens[name] = {"Authorization": resp.headers.get("Authorization")}
        
    return tokens

def test_posting(user_tokens, db_inspector):
    for i in range(2):
        header = user_tokens[USERS[i]]
        
        psts = db_inspector(DBS[i], "SELECT * FROM statuses")
        assert len(psts) == 0
        resp = requests.post(f"{SERVERS[i]}/api/posts", 
                             json={
                                "userID": 1,
                                "content": f"I am {USERS[i]}",
                             }, 
                             headers=header)
        assert resp.status_code == 201
        psts = db_inspector(DBS[i], "SELECT * FROM statuses")
        assert len(psts) == 1

def test_remote_search(db_inspector):
    url = f"{SERVER_1}/api/accounts/lookup"
    url_remote = SERVER_2
    url_remote = url_remote.replace("http://", "").replace("https://", "").rstrip("/")
    params = {
        "acct": f"@{USER_2}@{url_remote}"
    }
    resp = requests.get(url, params=params)
    assert resp.status_code == 200
    remote_acc = db_inspector(DB_1, f"SELECT * FROM accounts where username like '{USER_2}'")
    assert(remote_acc[0].username == "bob")
    assert(remote_acc[0].domain != None)

def test_remote_follow(user_tokens, db_inspector):
    remote_acc = db_inspector(DB_1, f"SELECT * FROM accounts where username like '{USER_2}'")[0]

    header = user_tokens[USER_1]
    resp = requests.post(f"{SERVER_1}/api/accounts/{remote_acc.id}/follow", headers=header)
    assert(resp.status_code == 201)
    # check if  follow appeared locally
    local_follows = db_inspector(DB_1, f"SELECT * FROM follows")
    assert(len(local_follows) == 1)
    # check if  follow appeared  on remote server
    remote_follows = db_inspector(DB_2, f"SELECT * FROM follows")
    assert(len(remote_follows) == 1)

    # check if remote got the follower and if it's domain is correctly set
    remote_follower = db_inspector(DB_2, f"SELECT * FROM accounts where username like '{USER_1}'")[0]
    assert(remote_follower.domain != None)
