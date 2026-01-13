import pytest
import subprocess
import time
import psycopg2
import requests
import signal
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT
import os

DB_NAMES = ["db1", "db2"]
ADMIN_URL = "postgres://borg:borg@localhost:5432/postgres"

def recreate_database(db_name):
    conn = psycopg2.connect(ADMIN_URL)
    conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT)
    
    with conn.cursor() as cur:
        print(f"--- Recreating database: {db_name} ---")
        cur.execute(f"""
            SELECT pg_terminate_backend(pg_stat_activity.pid)
            FROM pg_stat_activity
            WHERE pg_stat_activity.datname = '{db_name}'
              AND pid <> pg_backend_pid();
        """)
        cur.execute(f"DROP DATABASE IF EXISTS {db_name};")
        cur.execute(f"CREATE DATABASE {db_name};")
    conn.close()

@pytest.fixture(scope="session", autouse=True)
def manage_backend():
    ports = ["8080", "8081"]
    for name in DB_NAMES:
        recreate_database(name)

    cmd = [
        './tests/run_servers.sh', 
        ports[0], DB_NAMES[0], 
        ports[1], DB_NAMES[1]
    ]
    proc = subprocess.Popen(
        cmd,
        stdout=subprocess.PIPE,
        text=True,
        preexec_fn=os.setsid 
    )
    line = proc.stdout.readline().strip()
    pids = [p for p in line.split(',') if p]
    
    for port in ports:
        timeout = 10
        start = time.time()
        while time.time() - start < timeout:
            try:
                if requests.get(f"http://localhost:{port}/health").status_code == 200:
                    break
            except:
                time.sleep(0.5)
        else:
            os.killpg(os.getpgid(proc.pid), signal.SIGTERM)
            pytest.fail(f"Backend setup failed: Server on port {port} did not respond with 200 OK")

    yield

    print("\n--- Shutting down backend servers ---")
    try:
        os.killpg(os.getpgid(proc.pid), signal.SIGTERM)
    except ProcessLookupError:
        pass
