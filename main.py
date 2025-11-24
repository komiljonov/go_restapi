import requests
import threading
import time

URL = "http://localhost:8080/api/ping"

def hit(i):
    print(f"Request {i} START")
    t0 = time.time()
    resp = requests.get(URL)
    dt = time.time() - t0
    print(f"Request {i} END ({dt:.2f}s)  -> {resp.text}")

# Create two threads that fire at the same time
t1 = threading.Thread(target=hit, args=(1,))
t2 = threading.Thread(target=hit, args=(2,))

t1.start()
t2.start()

t1.join()
t2.join()