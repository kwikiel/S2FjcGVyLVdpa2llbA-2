import requests
import time 
# Test if server is alive and returns empty resource
r = requests.get("http://localhost:8080/api/fetcher")

r.status_code

# Basic case
if r.status_code == 200: 
    print("Server is alive")
else:
    print("General server error")


# Test insertion into in memory database

payload = {'ID': 1, 'URL': 'https://httpbin.org/range/15', 'INTERVAL': 3}

r = requests.post('http://localhost:8080/api/fetcher', 
    json=payload)

if r.status_code == 200: 
    print("Server is alive and POST works")
else:
    print("Error: Cannot insert fetch into database")


# Check if Resource with ID 1 exists
r = requests.get("http://localhost:8080/api/fetcher/1")
if r.json() == payload: 
    print("Server is alive and POST works")
else:
    print("Error: Cannot insert fetch into database")


# Start worker, wait 10 seconds and check if Server is fetching urls 
s = requests.get("http://localhost:8080/worker")
if s.text == "...": 
    print("Worker started")
else:
    print("Error: cannot start the worker")


print("Waiting for 5 seconds")
for x in range(5):
    print("Sleeping "+str(x)+" second...")
    time.sleep(1)

s = requests.get("http://localhost:8080/api/fetcher/1/history")
if (len(s.json()["1"]) > 1): 
    print("Fetching operational")
else:
    print("Error: cannot start the worker")