#! /usr/bin/python3
 
from sys import argv
from os import path, mkdir
import urllib3

if len(argv) != 2:
    print("Usage: load_input.py day")
    exit(1)

if not  path.exists(".env"):
    print("No environment file found - exiting")
    exit(1)

token = "not_set"
with open(".env", "r") as f:
    # token=abc
    line = f.readline()
    elems = line.split("=")
    if len(elems) == 2:
        token = elems[1]

if token == "not_set":
    print("No token available - exiting")
    exit(1)


day = int(argv[1])
# day = 13
print(f"Preparing day {day:02}")

date = f"d{day:02}"

if not path.exists(date):
    mkdir(date)

print(f"Downloading input file...")
input_url = f"https://adventofcode.com/2023/day/{day}/input"
http = urllib3.PoolManager()
sessioncookie = f"session={token}"
r = http.request('GET', input_url,
    headers = {
        'Cookie' : sessioncookie
    })
if r.status == 200:
    with open(f'./{date}/input.txt', 'w') as f:
        f.write(r.data.decode())
else:
    print("Error downloading input file - please check manually")
